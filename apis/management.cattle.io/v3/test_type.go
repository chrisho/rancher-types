package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TestSpec `json:"spec,omitempty"`
}

type TestSpec struct {
	DisplayName string `json:"displayName"`
}
