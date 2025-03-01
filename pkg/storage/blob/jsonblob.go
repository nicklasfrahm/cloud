package blob

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/server/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	builderrest "sigs.k8s.io/apiserver-runtime/pkg/builder/rest"
)

// NewJSONBLOBStorageProvider returns a storage provider that creates JSON
// files in an object storage bucket.
//
// For namespaced objects, the path is:
// /<group>/<kind>s/namespaces/<namespace>/<name>.json
//
// For cluster-scoped objects, the path is:
// /<group>/<kind>s/<name>.json
func NewJSONBLOBStorageProvider(obj resource.Object, objectStorageConfig []byte) builderrest.ResourceHandlerProvider {
	return func(scheme *runtime.Scheme, getter generic.RESTOptionsGetter) (rest.Storage, error) {
		groupResource := obj.GetGroupVersionResource().GroupResource()
		codec, _, err := storage.NewStorageCodec(storage.StorageCodecConfig{
			StorageMediaType:  runtime.ContentTypeJSON,
			StorageSerializer: serializer.NewCodecFactory(scheme),
			StorageVersion:    scheme.PrioritizedVersionsForGroup(obj.GetGroupVersionResource().Group)[0],
			MemoryVersion:     scheme.PrioritizedVersionsForGroup(obj.GetGroupVersionResource().Group)[0],
			// Not relevant as we are using BLOB storage.
			Config: storagebackend.Config{},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create storage codec: %v", err)
		}

		return NewBLOBREST(
			groupResource,
			codec,
			obj.NamespaceScoped(),
			obj.New,
			obj.NewList,
			objectStorageConfig,
		)
	}
}
