/*
Copyright 2024.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=0

	// The version of the cluster. I'm not exactly sure how cluster
	// versions work so this will always be "v1".
	Version string `json:"version"`

	//+kubebuilder:validation:MinLength=0

	// The name of the cluster. Again, not sure if clusters have names but seems useful to have in a scan.
	Name string `json:"name,omitempty"`

	//+kubebuilder:validation:MinItems=0

	// list of nodes (create a new type for this shit)
	Nodes []Node `json:"nodes"`
}

// Struct representing a node.
type Node struct {
	//+kubebuilder:validation:MinLength=0
	Name         string     `json:"name"`
	UID          int32      `json:"uid"`
	NumberOfPods int32      `json:"numberOfPods"`
	Master       bool       `json:"master"`
	Status       NodeStatus `json:"status"`
}

// NodeStatus describes the status of a node. Only one of the given statuses can be specified.
// +kubebuilder:validation:Enum=Active;Inactive;Error
type NodeStatus string

const (
	// The node is active.
	ActiveStatus NodeStatus = "Active"

	// The node is inactive.
	InactiveStatus NodeStatus = "Inactive"

	// The node has encountered an error of some sort.
	ErrorStatus NodeStatus = "Error"
)

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Note: Not sure what observed state I would record here since I just log the data collected in the ClusterScan.
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterScan is the Schema for the clusterscans API
type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}
