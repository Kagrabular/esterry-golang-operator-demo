package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion group version used to register these objects
var GroupVersion = schema.GroupVersion{Group: "example.com", Version: "v1"}

// registers types
var (
	SchemeBuilder = runtime.NewSchemeBuilder(
		addKnownTypes,
	)
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		GroupVersion,
		&NamespaceConfig{},
		&NamespaceConfigList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
