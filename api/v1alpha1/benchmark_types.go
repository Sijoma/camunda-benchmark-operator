/*
Copyright 2022.

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

// BenchmarkSpec defines the desired state of Benchmark
type BenchmarkSpec struct {
	// This has to be the name of a secret that contains clientId, clientSecret,
	// zeebeAddress and authServer. It defaults to "cloud-credentials"
	//+optional
	CredentialsSecretName string `json:"credentialsSecretName"`
	// Starter rate value on the starter deployment
	ProcessStarterRate int `json:"processStarterRate"`
	// Number of replicas for the starter
	StarterReplicas int `json:"starterReplicas"`
	// Number of workers
	WorkerCount int `json:"workerCount"`
	// How long the benchmark should run
	Duration string `json:"duration"`
}

// BenchmarkStatus defines the observed state of Benchmark
type BenchmarkStatus struct {
	StartTime *metav1.Time `json:"startTime"`
	Progress  string       `json:"progress"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced,path=benchmarks,shortName=bench
//+kubebuilder:printcolumn:name="# Workers",type=string,JSONPath=`.spec.workerCount`
//+kubebuilder:printcolumn:name="# Starters",type=string,JSONPath=`.spec.starterReplicas`
//+kubebuilder:printcolumn:name="Starter Rate",type=string,JSONPath=`.spec.processStarterRate`
//+kubebuilder:printcolumn:name="Duration",type=string,JSONPath=`.spec.duration`
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
//+kubebuilder:printcolumn:name="Progress",type=string,JSONPath=`.status.progress`

// Benchmark is the Schema for the benchmarks API
type Benchmark struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BenchmarkSpec   `json:"spec,omitempty"`
	Status BenchmarkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BenchmarkList contains a list of Benchmark
type BenchmarkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Benchmark `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Benchmark{}, &BenchmarkList{})
}
