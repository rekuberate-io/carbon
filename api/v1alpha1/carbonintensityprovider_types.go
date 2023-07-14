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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CarbonIntensityProviderSpec defines the desired state of CarbonIntensityProvider
type CarbonIntensityProviderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default=12
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=12
	// +kubebuilder:validation:Maximum=24
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	RefreshIntervalInHours *int32 `json:"interval"`

	// +kubebuilder:validation:Enum=watttime;electricitymaps
	// +kubebuilder:default:=electricitymaps
	Provider string `json:"provider"`

	// +kubebuilder:validation:Enum=average;marginal
	// +kubebuilder:default:=average
	Signal string `json:"signal"`

	Config *v1.SecretReference `json:"config,omitempty"`
}

// CarbonIntensityProviderStatus defines the observed state of CarbonIntensityProvider
type CarbonIntensityProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastUpdate *metav1.Time `json:"lastUpdate,omitempty"`
	NextUpdate *metav1.Time `json:"nextUpdate,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CarbonIntensityProvider is the Schema for the carbonintensityproviders API
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Signal",type=string,JSONPath=`.spec.signal`
// +kubebuilder:printcolumn:name="Interval(Hours)",type=string,JSONPath=`.spec.interval`
// +kubebuilder:printcolumn:name="Last Update",type=string,JSONPath=`.status.lastUpdate`
// +kubebuilder:printcolumn:name="Next Update",type=string,JSONPath=`.status.nextUpdate`
type CarbonIntensityProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CarbonIntensityProviderSpec   `json:"spec,omitempty"`
	Status CarbonIntensityProviderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CarbonIntensityProviderList contains a list of CarbonIntensityProvider
type CarbonIntensityProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CarbonIntensityProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CarbonIntensityProvider{}, &CarbonIntensityProviderList{})
}
