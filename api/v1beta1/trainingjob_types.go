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

package v1beta1

import (
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const KIND = "TrainingJob"

type JobState string

const (
	Created    JobState = "Created"
	Restarting JobState = "Restarting"
	Pending    JobState = "Pending"
	Running    JobState = "Running"
	Succeeded  JobState = "Succeeded"
	Failed     JobState = "Failed"
	Unknown    JobState = "Unknown"
)

// TrainingJobSpec defines the desired state of TrainingJob
type TrainingJobSpec struct {
	// Type indicates the training framework or workload type
	Type string `json:"type,omitempty"`

	// RunPolicy encapsulates various runtime policies of the distributed training
	// job, for example how to clean up resources and how long the job can stay
	// active.
	//+kubebuilder:validation:Optional
	RunPolicy kubeflowv1.RunPolicy `json:"runPolicy"`

	// ElasticPolicy aim to define global elastic configuration
	ElasticPolicy *kubeflowv1.ElasticPolicy `json:"elasticPolicy,omitempty"`

	// A map of ReplicaType (type) to ReplicaSpec (value).
	// For example,
	//   {
	//     "Master": PyTorchReplicaSpec,
	//     "Worker": PyTorchReplicaSpec,
	//   }
	ReplicaSpecs map[kubeflowv1.ReplicaType]*kubeflowv1.ReplicaSpec `json:"replicaSpecs"`

	// OptionalMap store user defined data.
	OptionalMap map[string]string `json:"optionalMap,omitempty"`
}

// TrainingJobStatus defines the observed state of TrainingJob
type TrainingJobStatus struct {
	kubeflowv1.JobStatus `json:"status,omitempty"`

	// State is a simple, high-level summary of where the Job is in its lifecycle
	// +optional
	State JobState `json:"state,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=tjob
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=`.spec.type`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=`.status.state`

// TrainingJob is the Schema for the trainingjobs API
type TrainingJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrainingJobSpec   `json:"spec,omitempty"`
	Status TrainingJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TrainingJobList contains a list of TrainingJob
type TrainingJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrainingJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrainingJob{}, &TrainingJobList{})
}
