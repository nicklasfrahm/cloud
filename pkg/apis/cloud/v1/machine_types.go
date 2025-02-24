
/*
Copyright 2025.

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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
 	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Machine
// +k8s:openapi-gen=true
type Machine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineSpec   `json:"spec,omitempty"`
	Status MachineStatus `json:"status,omitempty"`
}

// MachineList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MachineList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Machine `json:"items"`
}

// MachineSpec defines the desired state of Machine
type MachineSpec struct {
}

var _ resource.Object = &Machine{}
var _ resourcestrategy.Validater = &Machine{}

func (in *Machine) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *Machine) NamespaceScoped() bool {
	return false
}

func (in *Machine) New() runtime.Object {
	return &Machine{}
}

func (in *Machine) NewList() runtime.Object {
	return &MachineList{}
}

func (in *Machine) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "cloud.nicklasfrahm.dev",
		Version:  "v1",
		Resource: "machines",
	}
}

func (in *Machine) IsStorageVersion() bool {
	return true
}

func (in *Machine) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &MachineList{}

func (in *MachineList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}
// MachineStatus defines the observed state of Machine
type MachineStatus struct {
}

func (in MachineStatus) SubResourceName() string {
	return "status"
}

// Machine implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &Machine{}

func (in *Machine) GetStatus() resource.StatusSubResource {
	return in.Status
}

// MachineStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &MachineStatus{}

func (in MachineStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*Machine).Status = in
}
