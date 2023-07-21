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
	ForecastRefreshIntervalInHours *int32 `json:"forecastRefreshInterval"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=24
	// +kubebuilder:validation:ExclusiveMinimum=false
	// +kubebuilder:validation:ExclusiveMaximum=false
	LiveRefreshIntervalInHours *int32 `json:"liveRefreshInterval"`

	// +kubebuilder:validation:Enum=watttime;electricitymaps;simulator
	// +kubebuilder:default:=electricitymaps
	Provider string `json:"provider"`

	// +kubebuilder:validation:Enum=average;marginal
	// +kubebuilder:default:=average
	EmissionsType string `json:"emissionsType"`

	WattTimeConfiguration        *WattTimeConfigurationSpec        `json:"watttime,omitempty"`
	ElectricityMapsConfiguration *ElectricityMapsConfigurationSpec `json:"electricitymaps,omitempty"`
	SimulatorConfiguration       *SimulatorConfigurationSpec       `json:"simulator,omitempty"`
}

// CarbonIntensityProviderStatus defines the observed state of CarbonIntensityProvider
type CarbonIntensityProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Zone            *string      `json:"zone,omitempty"`
	Provider        *string      `json:"provider,omitempty"`
	LastForecast    *metav1.Time `json:"lastForecast,omitempty"`
	LastUpdate      *metav1.Time `json:"lastUpdate,omitempty"`
	NextUpdate      *metav1.Time `json:"nextUpdate,omitempty"`
	CarbonIntensity *string      `json:"carbonIntensity,omitempty"`
}

type WattTimeConfigurationSpec struct {
	Username string              `json:"username"`
	Region   string              `json:"region"`
	Password *v1.SecretReference `json:"password"`
}

type ElectricityMapsConfigurationSpec struct {
	// +kubebuilder:validation:Enum=commercial;commercial_trial;free_tier
	// +kubebuilder:default:=free_tier
	Subscription            string              `json:"subscription,omitempty"`
	CommercialTrialEndpoint *string             `json:"commercialTrialEndpoint,omitempty"`
	Zone                    *string             `json:"zone,omitempty"`
	ApiKey                  *v1.SecretReference `json:"apiKey"`
}

type SimulatorConfigurationSpec struct {

	// +kubebuilder:default:=false
	Randomize *bool `json:"randomize,omitempty"`

	// +kubebuilder:default:=SIM-1
	Zone *string `json:"zone,omitempty"`
}

type GeolocationSpec struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CarbonIntensityProvider is the Schema for the carbonintensityproviders API
// +kubebuilder:printcolumn:name="Provider",type=string,JSONPath=`.spec.provider`
// +kubebuilder:printcolumn:name="Zone",type=string,JSONPath=`.status.zone`
// +kubebuilder:printcolumn:name="Forecast INVL(h)",type=string,JSONPath=`.spec.forecastRefreshInterval`
// +kubebuilder:printcolumn:name="Last Forecast",type=string,JSONPath=`.status.lastForecast`
// +kubebuilder:printcolumn:name="CI (gCO2eq/KWh)",type=string,JSONPath=`.status.carbonIntensity`
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
