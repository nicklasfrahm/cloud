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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE! THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required. Any new fields you add
// must have json tags for the fields to be serialized.

// KubernetesClusterSpecCluster defines cluster-wide configuration.
type KubernetesClusterSpecCluster struct {
	// AllowSchedulingOnControlPlanes allows pods to be scheduled on control plane nodes.
	// +kubebuilder:default=false
	AllowSchedulingOnControlPlanes bool `json:"allowSchedulingOnControlPlanes,omitempty"`
}

// KubernetesClusterSpecMachinePools defines machine pool-specific configuration.
type KubernetesClusterSpecMachinePool struct {
	// Name is the name of the machine pool.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
}

// KubernetesClusterSpecInfrastructureControlPlane defines control plane-specific infrastructure configuration.
type KubernetesClusterSpecInfrastructureControlPlane struct {
	// MachinePools defines the machine pools for the control plane.
	// +kubebuilder:validation:MinItems=1
	MachinePools []KubernetesClusterSpecMachinePool `json:"machinePools,omitempty"`
}

// KubernetesClusterSpecInfrastructure defines infrastructure-specific configuration.
type KubernetesClusterSpecInfrastructure struct {
	// ControlPlane defines control plane-specific configuration.
	// +kubebuilder:validation:Required
	ControlPlane KubernetesClusterSpecInfrastructureControlPlane `json:"controlPlane,omitempty"`
}

// KubernetesClusterSpec defines the desired state of KubernetesCluster
type KubernetesClusterSpec struct {
	// Cluster defines cluster-wide configuration.
	// +kubebuilder:validation:Optional
	Cluster KubernetesClusterSpecCluster `json:"cluster,omitempty"`
	// Infrastructure defines infrastructure-specific configuration.
	// +kubebuilder:validation:Required
	Infrastructure KubernetesClusterSpecInfrastructure `json:"infrastructure,omitempty"`
}

// KubernetesClusterStatus defines the observed state of KubernetesCluster
type KubernetesClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KubernetesCluster is the Schema for the kubernetesclusters API
type KubernetesCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubernetesClusterSpec   `json:"spec,omitempty"`
	Status KubernetesClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubernetesClusterList contains a list of KubernetesCluster
type KubernetesClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubernetesCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubernetesCluster{}, &KubernetesClusterList{})
}
