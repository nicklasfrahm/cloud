package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/thanos-io/objstore"
	"github.com/thanos-io/objstore/client"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

var (
	// ErrNamespaceNotExists is returned when the namespace does not exist.
	ErrNamespaceNotExists = fmt.Errorf("namespace does not exist")
	// ErrObjectExists is returned when the object already exists.
	ErrObjectExists = fmt.Errorf("object already exists")
)

// Ensure that blobREST implements the interfaces.
var _ rest.StandardStorage = &blobREST{}
var _ rest.Scoper = &blobREST{}
var _ rest.Storage = &blobREST{}

// NewBLOBREST creates a new REST storage for the given group resource.
func NewBLOBREST(
	groupResource schema.GroupResource,
	codec runtime.Codec,
	isNamespaced bool,
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	objectStorageConfig []byte,
) (rest.Storage, error) {
	// TODO: Replace logger with compatible OTEL logger.
	bucket, err := client.NewBucket(log.NewJSONLogger(os.Stdout), objectStorageConfig, groupResource.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create object storage client: %v", err)
	}

	rest := &blobREST{
		TableConvertor: rest.NewDefaultTableConvertor(groupResource),
		codec:          codec,
		isNamespaced:   isNamespaced,
		newFunc:        newFunc,
		newListFunc:    newListFunc,
		watchers:       make(map[int]*blobWatch, 10),
		bucket:         bucket,
	}

	return rest, nil
}

type blobREST struct {
	rest.TableConvertor
	codec        runtime.Codec
	isNamespaced bool

	muWriter sync.RWMutex
	watchers map[int]*blobWatch

	newFunc     func() runtime.Object
	newListFunc func() runtime.Object

	bucket objstore.Bucket
}

func (b *blobREST) notifyWatchers(event watch.Event) {
	b.muWriter.RLock()
	defer b.muWriter.RUnlock()
	for _, w := range b.watchers {
		w.ch <- event
	}
}

func (b *blobREST) New() runtime.Object {
	return b.newFunc()
}

// Destroy is a no-op for BLOB storage.
func (b *blobREST) Destroy() {}

func (b *blobREST) NewList() runtime.Object {
	return b.newListFunc()
}

func (b *blobREST) NamespaceScoped() bool {
	return b.isNamespaced
}

func (b *blobREST) Get(
	ctx context.Context,
	name string,
	options *metav1.GetOptions,
) (runtime.Object, error) {
	objName, err := b.objectFileName(ctx, name)
	if err != nil {
		gvk := b.newFunc().GetObjectKind().GroupVersionKind()

		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    gvk.Group,
			Resource: gvk.Kind,
		}, name)
	}

	return b.read(ctx, objName)
}

func (b *blobREST) List(
	ctx context.Context,
	options *metainternalversion.ListOptions,
) (runtime.Object, error) {
	newListObj := b.newListFunc()

	val, err := getListPtr(newListObj)
	if err != nil {
		return nil, fmt.Errorf("failed to get list ptr: %v", err)
	}

	dirName, err := b.objectDirName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object directory name: %v", err)
	}

	if err := b.bucket.Iter(ctx, dirName, func(objName string) error {
		if !strings.HasSuffix(objName, ".json") {
			return nil
		}

		obj, err := b.read(ctx, objName)
		if err != nil {
			return fmt.Errorf("failed to get object: %v", err)
		}

		appendItem(val, obj)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to iterate objects in bucket: %v", err)
	}

	return newListObj, nil
}

func (b *blobREST) Create(
	ctx context.Context,
	obj runtime.Object,
	createValidation rest.ValidateObjectFunc,
	options *metav1.CreateOptions,
) (runtime.Object, error) {
	if createValidation != nil {
		if err := createValidation(ctx, obj); err != nil {
			return nil, fmt.Errorf("failed to validate object: %v", err)
		}
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to get object accessor: %v", err)
	}

	objName, err := b.objectFileName(ctx, accessor.GetName())
	if err != nil {
		return nil, fmt.Errorf("failed to get object file name: %v", err)
	}

	objExists, err := b.exists(ctx, objName)
	if err != nil {
		return nil, fmt.Errorf("failed to check for object existence: %v", err)
	}

	if objExists {
		gvk := obj.GetObjectKind().GroupVersionKind()

		return nil, errors.NewAlreadyExists(schema.GroupResource{
			Group:    gvk.Group,
			Resource: gvk.Kind,
		}, accessor.GetName())
	}

	if err := b.write(ctx, objName, obj); err != nil {
		return nil, fmt.Errorf("failed to write object: %v", err)
	}

	b.notifyWatchers(watch.Event{
		Type:   watch.Added,
		Object: obj,
	})

	return obj, nil
}

func (b *blobREST) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *metav1.UpdateOptions,
) (runtime.Object, bool, error) {
	var isCreate bool

	oldObj, err := b.Get(ctx, name, nil)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, false, fmt.Errorf("failed to get object: %v", err)
		}

		if !forceAllowCreate {
			gvk := b.newFunc().GetObjectKind().GroupVersionKind()

			return nil, false, errors.NewNotFound(schema.GroupResource{
				Group:    gvk.Group,
				Resource: gvk.Kind,
			}, name)
		}

		isCreate = true
	}

	updatedObj, err := objInfo.UpdatedObject(ctx, oldObj)
	if err != nil {
		return nil, false, fmt.Errorf("failed to update object: %v", err)
	}

	objName, err := b.objectFileName(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get object file name: %v", err)
	}

	if isCreate {
		if createValidation != nil {
			if err := createValidation(ctx, updatedObj); err != nil {
				return nil, false, fmt.Errorf("failed to validate object: %v", err)
			}

			if err := b.write(ctx, objName, updatedObj); err != nil {
				return nil, false, fmt.Errorf("failed to write object: %v", err)
			}

			b.notifyWatchers(watch.Event{
				Type:   watch.Added,
				Object: updatedObj,
			})

			return updatedObj, true, nil
		}
	}

	if updateValidation != nil {
		if err := updateValidation(ctx, oldObj, updatedObj); err != nil {
			return nil, false, fmt.Errorf("failed to validate object update: %v", err)
		}
	}

	if err := b.write(ctx, objName, updatedObj); err != nil {
		return nil, false, fmt.Errorf("failed to write object: %v", err)
	}

	b.notifyWatchers(watch.Event{
		Type:   watch.Modified,
		Object: updatedObj,
	})

	return updatedObj, false, nil
}

func (b *blobREST) Delete(
	ctx context.Context,
	name string,
	deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions,
) (runtime.Object, bool, error) {
	objName, err := b.objectFileName(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get object file name: %v", err)
	}

	oldObj, err := b.Get(ctx, name, nil)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get object: %v", err)
	}

	if deleteValidation != nil {
		if err := deleteValidation(ctx, oldObj); err != nil {
			return nil, false, fmt.Errorf("failed to validate object deletion: %v", err)
		}
	}

	if err := b.bucket.Delete(ctx, objName); err != nil {
		return nil, false, fmt.Errorf("failed to delete object: %v", err)
	}

	b.notifyWatchers(watch.Event{
		Type:   watch.Deleted,
		Object: oldObj,
	})

	return oldObj, true, nil
}

func (b *blobREST) DeleteCollection(
	ctx context.Context,
	deleteValidation rest.ValidateObjectFunc,
	options *metav1.DeleteOptions,
	listOptions *metainternalversion.ListOptions,
) (runtime.Object, error) {
	newListObj := b.newListFunc()

	val, err := getListPtr(newListObj)
	if err != nil {
		return nil, fmt.Errorf("failed to get list ptr: %v", err)
	}

	dirName, err := b.objectDirName(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object directory name: %v", err)
	}

	if err := b.bucket.Iter(ctx, dirName, func(objName string) error {
		if !strings.HasSuffix(objName, ".json") {
			return nil
		}

		obj, err := b.read(ctx, objName)
		if err != nil {
			return fmt.Errorf("failed to get object: %v", err)
		}

		if deleteValidation != nil {
			if err := deleteValidation(ctx, obj); err != nil {
				return fmt.Errorf("failed to validate object deletion: %v", err)
			}
		}

		// TODO: Add logging if deleting an object fails.
		_ = b.bucket.Delete(ctx, objName)

		appendItem(val, obj)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to iterate objects in bucket: %v", err)
	}

	return newListObj, nil
}

func (b *blobREST) Watch(
	ctx context.Context,
	options *metainternalversion.ListOptions,
) (watch.Interface, error) {
	bw := &blobWatch{
		blob: b,
		id:   len(b.watchers),
		ch:   make(chan watch.Event, 10),
	}

	// On initial watch, send all the existing objects.
	list, err := b.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	danger := reflect.ValueOf(list).Elem()
	items := danger.FieldByName("Items")

	for i := 0; i < items.Len(); i++ {
		obj := items.Index(i).Addr().Interface().(runtime.Object)
		bw.ch <- watch.Event{
			Type:   watch.Added,
			Object: obj,
		}
	}

	b.muWriter.Lock()
	defer b.muWriter.Unlock()
	b.watchers[bw.id] = bw

	return bw, nil
}

func (b *blobREST) objectDirName(ctx context.Context) (string, error) {
	gvk := b.newFunc().GetObjectKind().GroupVersionKind()
	plural := strings.ToLower(gvk.Kind) + "s"

	dirName := gvk.Group + "/" + plural

	if b.isNamespaced {
		ns, ok := genericapirequest.NamespaceFrom(ctx)
		if !ok {
			return "", ErrNamespaceNotExists
		}

		return filepath.Join(dirName, "namespaces", ns), nil
	}

	return dirName, nil
}

func (b *blobREST) objectFileName(ctx context.Context, name string) (string, error) {
	dirName, err := b.objectDirName(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get object directory name: %v", err)
	}

	return filepath.Join(dirName, name+".json"), nil
}

func (b *blobREST) read(
	ctx context.Context,
	name string,
) (runtime.Object, error) {
	blob, err := b.bucket.Get(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from bucket: %v", err)
	}

	content, err := io.ReadAll(blob)
	if err != nil {
		return nil, fmt.Errorf("failed to get content from blob: %v", err)
	}

	newObj := b.newFunc()
	decodedObj, _, err := b.codec.Decode(content, nil, newObj)
	if err != nil {
		return nil, fmt.Errorf("failed to decode object: %v", err)
	}

	return decodedObj, nil
}

func (b *blobREST) exists(ctx context.Context, name string) (bool, error) {
	exists, err := b.bucket.Exists(ctx, name)
	if err != nil {
		return false, fmt.Errorf("failed to check for object existence: %v", err)
	}

	return exists, nil
}

func (b *blobREST) write(
	ctx context.Context,
	name string,
	obj runtime.Object,
) error {
	content, err := runtime.Encode(b.codec, obj)
	if err != nil {
		return fmt.Errorf("failed to encode object: %v", err)
	}

	if err := b.bucket.Upload(ctx, name, bytes.NewBuffer(content)); err != nil {
		return fmt.Errorf("failed to write object to bucket: %v", err)
	}

	return nil
}

type blobWatch struct {
	blob *blobREST
	id   int
	ch   chan watch.Event
}

func (w *blobWatch) Stop() {
	w.blob.muWriter.Lock()
	defer w.blob.muWriter.Unlock()
	delete(w.blob.watchers, w.id)
	close(w.ch)
}

func (w *blobWatch) ResultChan() <-chan watch.Event {
	return w.ch
}

func appendItem(value reflect.Value, obj runtime.Object) {
	value.Set(reflect.Append(value, reflect.ValueOf(obj).Elem()))
}

func getListPtr(obj runtime.Object) (reflect.Value, error) {
	listPtr, err := meta.GetItemsPtr(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	val, err := conversion.EnforcePtr(listPtr)
	if err != nil || val.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("need ptr to slice: %v", err)
	}

	return val, nil
}
