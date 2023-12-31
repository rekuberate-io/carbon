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

// CarbonIntensityIssuerSpec defines the desired state of CarbonIntensityIssuer
type CarbonIntensityIssuerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default=12
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=12
	// +kubebuilder:validation:Maximum=24
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	ForecastRefreshIntervalInHours int32 `json:"forecastRefreshIntervalHours"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=24
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	LiveRefreshIntervalInHours int32 `json:"liveRefreshIntervalHours"`

	// +kubebuilder:validation:Required
	Zone string `json:"zone"`

	// +kubebuilder:validation:Required
	ProviderRef *v1.ObjectReference `json:"providerRef,omitempty"`
}

// CarbonIntensityIssuerStatus defines the observed state of CarbonIntensityIssuer
type CarbonIntensityIssuerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastForecast    *metav1.Time `json:"lastForecast,omitempty"`
	LastUpdate      *metav1.Time `json:"lastUpdate,omitempty"`
	NextUpdate      *metav1.Time `json:"nextUpdate,omitempty"`
	CarbonIntensity *string      `json:"carbonIntensity,omitempty"`

	// Conditions store the status conditions of the Memcached instances
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CarbonIntensityIssuer is the Schema for the carbonintensityissuers API
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.providerRef.name`
// +kubebuilder:printcolumn:name="Zone",type=string,JSONPath=`.spec.zone`
// +kubebuilder:printcolumn:name="Forecast INVL(h)",type=string,JSONPath=`.spec.forecastRefreshIntervalHours`
// +kubebuilder:printcolumn:name="Last Forecast",type=string,JSONPath=`.status.lastForecast`
// +kubebuilder:printcolumn:name="CI (gCO2eq/KWh)",type=string,JSONPath=`.status.carbonIntensity`
// +kubebuilder:printcolumn:name="Last Update",type=string,JSONPath=`.status.lastUpdate`
// +kubebuilder:printcolumn:name="Next Update",type=string,JSONPath=`.status.nextUpdate`
type CarbonIntensityIssuer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CarbonIntensityIssuerSpec   `json:"spec,omitempty"`
	Status CarbonIntensityIssuerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CarbonIntensityIssuerList contains a list of CarbonIntensityIssuer
type CarbonIntensityIssuerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CarbonIntensityIssuer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CarbonIntensityIssuer{}, &CarbonIntensityIssuerList{})
}
