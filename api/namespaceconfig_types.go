// this file is prelim, just messing around with CRDs might not even use this
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceConfigSpec defines the desired state of NamespaceConfig
type NamespaceConfigSpec struct {
	Labels map[string]string `json:"labels,omitempty"`
}

// NamespaceConfigStatus defines the observed state
type NamespaceConfigStatus struct {
	Applied bool `json:"applied,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type NamespaceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceConfigSpec   `json:"spec,omitempty"`
	Status NamespaceConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type NamespaceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceConfig `json:"items"`
}
