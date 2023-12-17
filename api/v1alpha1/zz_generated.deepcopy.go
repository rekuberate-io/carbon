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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityIssuer) DeepCopyInto(out *CarbonIntensityIssuer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityIssuer.
func (in *CarbonIntensityIssuer) DeepCopy() *CarbonIntensityIssuer {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityIssuer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CarbonIntensityIssuer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityIssuerList) DeepCopyInto(out *CarbonIntensityIssuerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CarbonIntensityIssuer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityIssuerList.
func (in *CarbonIntensityIssuerList) DeepCopy() *CarbonIntensityIssuerList {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityIssuerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CarbonIntensityIssuerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityIssuerSpec) DeepCopyInto(out *CarbonIntensityIssuerSpec) {
	*out = *in
	if in.ProviderRef != nil {
		in, out := &in.ProviderRef, &out.ProviderRef
		*out = new(v1.ObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityIssuerSpec.
func (in *CarbonIntensityIssuerSpec) DeepCopy() *CarbonIntensityIssuerSpec {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityIssuerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CarbonIntensityIssuerStatus) DeepCopyInto(out *CarbonIntensityIssuerStatus) {
	*out = *in
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
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CarbonIntensityIssuerStatus.
func (in *CarbonIntensityIssuerStatus) DeepCopy() *CarbonIntensityIssuerStatus {
	if in == nil {
		return nil
	}
	out := new(CarbonIntensityIssuerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElectricityMaps) DeepCopyInto(out *ElectricityMaps) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElectricityMaps.
func (in *ElectricityMaps) DeepCopy() *ElectricityMaps {
	if in == nil {
		return nil
	}
	out := new(ElectricityMaps)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElectricityMaps) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElectricityMapsList) DeepCopyInto(out *ElectricityMapsList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ElectricityMaps, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElectricityMapsList.
func (in *ElectricityMapsList) DeepCopy() *ElectricityMapsList {
	if in == nil {
		return nil
	}
	out := new(ElectricityMapsList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ElectricityMapsList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElectricityMapsSpec) DeepCopyInto(out *ElectricityMapsSpec) {
	*out = *in
	if in.CommercialTrialEndpoint != nil {
		in, out := &in.CommercialTrialEndpoint, &out.CommercialTrialEndpoint
		*out = new(string)
		**out = **in
	}
	if in.ApiKeyRef != nil {
		in, out := &in.ApiKeyRef, &out.ApiKeyRef
		*out = new(v1.SecretReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElectricityMapsSpec.
func (in *ElectricityMapsSpec) DeepCopy() *ElectricityMapsSpec {
	if in == nil {
		return nil
	}
	out := new(ElectricityMapsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElectricityMapsStatus) DeepCopyInto(out *ElectricityMapsStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElectricityMapsStatus.
func (in *ElectricityMapsStatus) DeepCopy() *ElectricityMapsStatus {
	if in == nil {
		return nil
	}
	out := new(ElectricityMapsStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Simulator) DeepCopyInto(out *Simulator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Simulator.
func (in *Simulator) DeepCopy() *Simulator {
	if in == nil {
		return nil
	}
	out := new(Simulator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Simulator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorList) DeepCopyInto(out *SimulatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Simulator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorList.
func (in *SimulatorList) DeepCopy() *SimulatorList {
	if in == nil {
		return nil
	}
	out := new(SimulatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SimulatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorSpec) DeepCopyInto(out *SimulatorSpec) {
	*out = *in
	if in.Bootstrap != nil {
		in, out := &in.Bootstrap, &out.Bootstrap
		*out = new(bool)
		**out = **in
	}
	if in.Replacement != nil {
		in, out := &in.Replacement, &out.Replacement
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorSpec.
func (in *SimulatorSpec) DeepCopy() *SimulatorSpec {
	if in == nil {
		return nil
	}
	out := new(SimulatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimulatorStatus) DeepCopyInto(out *SimulatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimulatorStatus.
func (in *SimulatorStatus) DeepCopy() *SimulatorStatus {
	if in == nil {
		return nil
	}
	out := new(SimulatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WattTime) DeepCopyInto(out *WattTime) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WattTime.
func (in *WattTime) DeepCopy() *WattTime {
	if in == nil {
		return nil
	}
	out := new(WattTime)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WattTime) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WattTimeList) DeepCopyInto(out *WattTimeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]WattTime, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WattTimeList.
func (in *WattTimeList) DeepCopy() *WattTimeList {
	if in == nil {
		return nil
	}
	out := new(WattTimeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *WattTimeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WattTimeSpec) DeepCopyInto(out *WattTimeSpec) {
	*out = *in
	if in.Password != nil {
		in, out := &in.Password, &out.Password
		*out = new(v1.SecretReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WattTimeSpec.
func (in *WattTimeSpec) DeepCopy() *WattTimeSpec {
	if in == nil {
		return nil
	}
	out := new(WattTimeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WattTimeStatus) DeepCopyInto(out *WattTimeStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WattTimeStatus.
func (in *WattTimeStatus) DeepCopy() *WattTimeStatus {
	if in == nil {
		return nil
	}
	out := new(WattTimeStatus)
	in.DeepCopyInto(out)
	return out
}
