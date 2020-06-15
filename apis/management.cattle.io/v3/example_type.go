package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ExampleSpec `json:"spec,omitempty"`
}

type ExampleSpec struct {
	DisplayName string `json:"displayName"`
}
