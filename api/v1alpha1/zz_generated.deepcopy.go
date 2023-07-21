//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityProvider) DeepCopyInto(out *CarbonIntensityProvider) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityProvider.
func (in *CarbonIntensityProvider) DeepCopy() *CarbonIntensityProvider {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityProvider)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CarbonIntensityProvider) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityProviderList) DeepCopyInto(out *CarbonIntensityProviderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CarbonIntensityProvider, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityProviderList.
func (in *CarbonIntensityProviderList) DeepCopy() *CarbonIntensityProviderList {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityProviderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CarbonIntensityProviderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityProviderSpec) DeepCopyInto(out *CarbonIntensityProviderSpec) {
	*out = *in
	if in.WattTimeConfiguration != nil {
		in, out := &in.WattTimeConfiguration, &out.WattTimeConfiguration
		*out = new(WattTimeConfigurationSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ElectricityMapsConfiguration != nil {
		in, out := &in.ElectricityMapsConfiguration, &out.ElectricityMapsConfiguration
		*out = new(ElectricityMapsConfigurationSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.SimulatorConfiguration != nil {
		in, out := &in.SimulatorConfiguration, &out.SimulatorConfiguration
		*out = new(SimulatorConfigurationSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityProviderSpec.
func (in *CarbonIntensityProviderSpec) DeepCopy() *CarbonIntensityProviderSpec {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityProviderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityProviderStatus) DeepCopyInto(out *CarbonIntensityProviderStatus) {
	*out = *in
	if in.Zone != nil {
		in, out := &in.Zone, &out.Zone
		*out = new(string)
		**out = **in
	}
	if in.Provider != nil {
		in, out := &in.Provider, &out.Provider
		*out = new(string)
		**out = **in
	}
	if in.LastForecast != nil {
		in, out := &in.LastForecast, &out.LastForecast
		*out = (*in).DeepCopy()
	}
	if in.LastUpdate != nil {
		in, out := &in.LastUpdate, &out.LastUpdate
		*out = (*in).DeepCopy()
	}
	if in.NextUpdate != nil {
		in, out := &in.NextUpdate, &out.NextUpdate
		*out = (*in).DeepCopy()
	}
	if in.CarbonIntensity != nil {
		in, out := &in.CarbonIntensity, &out.CarbonIntensity
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityProviderStatus.
func (in *CarbonIntensityProviderStatus) DeepCopy() *CarbonIntensityProviderStatus {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityProviderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElectricityMapsConfigurationSpec) DeepCopyInto(out *ElectricityMapsConfigurationSpec) {
	*out = *in
	if in.CommercialTrialEndpoint != nil {
		in, out := &in.CommercialTrialEndpoint, &out.CommercialTrialEndpoint
		*out = new(string)
		**out = **in
	}
	if in.ApiKey != nil {
		in, out := &in.ApiKey, &out.ApiKey
		*out = new(v1.SecretReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElectricityMapsConfigurationSpec.
func (in *ElectricityMapsConfigurationSpec) DeepCopy() *ElectricityMapsConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(ElectricityMapsConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeolocationSpec) DeepCopyInto(out *GeolocationSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeolocationSpec.
func (in *GeolocationSpec) DeepCopy() *GeolocationSpec {
	if in == nil {
		return nil
	}
	out := new(GeolocationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorConfigurationSpec) DeepCopyInto(out *SimulatorConfigurationSpec) {
	*out = *in
	if in.Randomize != nil {
		in, out := &in.Randomize, &out.Randomize
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorConfigurationSpec.
func (in *SimulatorConfigurationSpec) DeepCopy() *SimulatorConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(SimulatorConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WattTimeConfigurationSpec) DeepCopyInto(out *WattTimeConfigurationSpec) {
	*out = *in
	if in.Password != nil {
		in, out := &in.Password, &out.Password
		*out = new(v1.SecretReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WattTimeConfigurationSpec.
func (in *WattTimeConfigurationSpec) DeepCopy() *WattTimeConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(WattTimeConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}
