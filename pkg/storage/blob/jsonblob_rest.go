package blob

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
)

// NewBLOBREST creates a new REST storage for the given group resource.
func NewBLOBREST(
	groupResource schema.GroupResource,
	codec runtime.Codec,
	rootPath string,
	isNamespaced bool,
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
) rest.Storage {
	// TODO: Continue implementation.
}
