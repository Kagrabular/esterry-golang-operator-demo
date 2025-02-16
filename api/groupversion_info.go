// this file is prelim, just messing around with CRDs might not even use this, but might need it for custom resources I'll have to see what I'm trying to target
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Define API version
var (
	GroupVersion  = schema.GroupVersion{Group: "example.com", Version: "v1"}
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}
)
