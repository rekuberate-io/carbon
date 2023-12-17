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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SimulatorSpec defines the desired state of Simulator
type SimulatorSpec struct {
	// +kubebuilder:default:=false
	// +kubebuilder:validation:Type=boolean
	Bootstrap *bool `json:"bootstrap,omitempty"`

	// +kubebuilder:default:=false
	// +kubebuilder:validation:Type=boolean
	Replacement *bool `json:"replacement,omitempty"`

	// +kubebuilder:default:=true
	// +kubebuilder:validation:Type=boolean
	StableSample *bool `json:"stableSample,omitempty"`
}

// SimulatorStatus defines the observed state of Simulator
type SimulatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Simulator is the Schema for the simulators API
// +kubebuilder:printcolumn:name="Bootstrap",type=string,JSONPath=`.spec.bootstrap`
// +kubebuilder:printcolumn:name="Replacement",type=string,JSONPath=`.spec.replacement`
// +kubebuilder:printcolumn:name="Stable Sample",type=string,JSONPath=`.spec.stableSample`
type Simulator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SimulatorSpec   `json:"spec,omitempty"`
	Status SimulatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SimulatorList contains a list of Simulator
type SimulatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Simulator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Simulator{}, &SimulatorList{})
}
