package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// NamespaceConfigSpec defines the desired state of NamespaceConfig
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NamespaceConfigSpec struct {
	Labels map[string]string `json:"labels,omitempty"`
}

// NamespaceConfigStatus defines the observed state
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NamespaceConfigStatus struct {
	Applied bool `json:"applied,omitempty"`
}

// NamespaceConfig is the Schema for the namespaceconfigs API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NamespaceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceConfigSpec   `json:"spec,omitempty"`
	Status NamespaceConfigStatus `json:"status,omitempty"`
}

// NamespaceConfigList contains a list of NamespaceConfig
// +kubebuilder:object:root=true
type NamespaceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceConfig `json:"items"`
}

// DeepCopyObject implements runtime.Object for NamespaceConfig
func (in *NamespaceConfig) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(NamespaceConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject implements runtime.Object for NamespaceConfigList
func (in *NamespaceConfigList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(NamespaceConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copies all properties of this object into another object of the same type.
func (in *NamespaceConfig) DeepCopyInto(out *NamespaceConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Spec.Labels != nil {
		out.Spec.Labels = make(map[string]string, len(in.Spec.Labels))
		for key, val := range in.Spec.Labels {
			out.Spec.Labels[key] = val
		}
	}
	out.Status = in.Status
}

// DeepCopyInto copies all properties of this list object into another list object of the same type.
func (in *NamespaceConfigList) DeepCopyInto(out *NamespaceConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		out.Items = make([]NamespaceConfig, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
}
