// Copyright 2019 Grackle Operator authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatastoreSpec defines the desired state of a Grackle datastore.
type DatastoreSpec struct {
	// Host is the hostname/IP address for the datastore.
	Host string `json:"host"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Grackle is the Schema for the grackles API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Grackle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GrackleSpec   `json:"spec,omitempty"`
	Status GrackleStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GrackleList contains a list of Grackle
type GrackleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Grackle `json:"items"`
}

// GracklePhase is a label for the condition of a Grackle at the current time.
type GracklePhase corev1.PodPhase

// GrackleSpec defines the desired state of Grackle
// +k8s:openapi-gen=true
type GrackleSpec struct {
	// Datastore is the specification for the storage of Tweets.
	Datastore DatastoreSpec `json:"datastore,omitempty"`

	// Ingest is the specification for the ingest of Tweets.
	Ingest *IngestSpec `json:"ingest,omitempty"`

	// Web is the specification for the Web UI.
	Web *WebSpec `json:"web,omitempty"`
}

// GrackleStatus defines the observed state of Grackle
// +k8s:openapi-gen=true
type GrackleStatus struct {
	Phase GracklePhase `json:"phase,omitempty"`
}

// IngestSpec defines the desired state of the Grackle Ingest process.
type IngestSpec struct {
	// Track is a slice of search terms to ingest.
	// Each list item is a comma-delimited string of keywords for a seperate Twitter stream.
	Track []string `json:"track,omitempty"`

	// TwitterSecret is the name of a secret containing Twitter API credentials.
	TwitterSecret string `json:"twitterSecret,omitempty"`

	// Version is the Grackle image tag to use for ingest.
	Version string `json:"version,omitempty"`
}

// WebSpec defines the desired state of the Grackle Web UI process.
type WebSpec struct {
	// Replicas is the number of web UI nodes to provision.
	Replicas *int32 `json:"replicas,omitempty"`

	// Version is the Grackle image tag to use for the web UI.
	Version string `json:"version,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Grackle{}, &GrackleList{})
}
