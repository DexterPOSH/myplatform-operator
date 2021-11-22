/*
Copyright 2021.

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

const (
	Organization = "dexterposh"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InhouseAppSpec defines the desired state of InhouseApp
type InhouseAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AppId uniquely identifies an app on MyPlatform
	AppId string `json:"appId"`

	// Language mentions the programming language for the app on the platform
	// +kubebuilder:validation:Enum=csharp;python;go
	Language string `json:"language"`

	// OS specifies the type of Operating System
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=windows;linux
	// +kubebuilder:default:=linux
	OS string `json:"os"`

	// InstanceSize is the T-Shirt size for the deployment
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=small;medium;large
	// +kubebuilder:default:=small
	InstanceSize string `json:"instanceSize"`

	// EnvironmenType specifies the type of environment
	// +kubebuilder:validation:Enum=dev;test;prod
	EnvironmentType string `json:"environmentType"`

	// Replicas indicate the replicas to mantain
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=1
	Replicas int32 `json:"replicas"`
}

// InhouseAppStatus defines the observed state of InhouseApp
type InhouseAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Pods are the name of the Pods hosting the App
	Pods []string `json:"pods"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// InhouseApp is the Schema for the inhouseapps API
type InhouseApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InhouseAppSpec   `json:"spec,omitempty"`
	Status InhouseAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InhouseAppList contains a list of InhouseApp
type InhouseAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InhouseApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InhouseApp{}, &InhouseAppList{})
}
