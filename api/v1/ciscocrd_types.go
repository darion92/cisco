/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CiscoCRDSpec defines the desired state of CiscoCRD
type CiscoCRDSpec struct {
	Replicas       *int32 `json:"replicas"`
	ContainerImage string `json:"containerImage"`
	Port           int32  `json:"port"`
	Host           string `json:"host"`
}

// CiscoCRDStatus defines the observed state of CiscoCRD
type CiscoCRDStatus struct {
	Workload        CrossNamespaceObjectReference `json:"workload"`
	DeployedService bool                          `json:"deployed_service,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CiscoCRD is the Schema for the ciscocrds API
type CiscoCRD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiscoCRDSpec   `json:"spec,omitempty"`
	Status CiscoCRDStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CiscoCRDList contains a list of CiscoCRD
type CiscoCRDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CiscoCRD `json:"items"`
}

type CrossNamespaceObjectReference struct {
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the referent.
	// +required
	Kind string `json:"kind"`

	// Name of the referent.
	// +required
	Name string `json:"name"`

	// Namespace of the referent, defaults to the namespace of the Kubernetes resource object that contains the reference.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

func init() {
	SchemeBuilder.Register(&CiscoCRD{}, &CiscoCRDList{})
}
