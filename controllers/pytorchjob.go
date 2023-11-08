package controllers

import (
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	tj "github.com/kuizhiqing/trainingjob-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	pyTorchJobType = "pytorchjob"
	nprocKey       = "nprocPerNode"
)

type PyTorchJobManager struct {
	TJob *tj.TrainingJob
	Job  *kubeflowv1.PyTorchJob
}

func NewPyTorchJobManager(tjob *tj.TrainingJob) *PyTorchJobManager {
	return &PyTorchJobManager{
		TJob: tjob,
		Job:  &kubeflowv1.PyTorchJob{},
	}
}

func (m *PyTorchJobManager) GetTJob() *tj.TrainingJob {
	return m.TJob
}

func (m *PyTorchJobManager) GetObject() client.Object {
	return m.Job
}

func (m *PyTorchJobManager) NewJob() client.Object {
	// runPolicy := kubeflowv1.RunPolicy{}
	// convert(m.TJob.Spec.RunPolicy, &runPolicy)

	// replicaSpecs := make(map[kubeflowv1.ReplicaType]*kubeflowv1.ReplicaSpec)
	// convert(m.TJob.Spec.ReplicaSpecs, &replicaSpecs)

	runPolicy := m.TJob.Spec.RunPolicy

	replicaSpecs := m.TJob.Spec.ReplicaSpecs

	slots := pickStrFromOptionalMap(m.TJob.Spec.OptionalMap, nprocKey, "1")

	pytorchjobSpec := kubeflowv1.PyTorchJobSpec{
		RunPolicy:           runPolicy,
		PyTorchReplicaSpecs: replicaSpecs,
		NprocPerNode:        &slots,
	}

	pytorchjob := &kubeflowv1.PyTorchJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   m.TJob.Namespace,
			Name:        m.TJob.Name,
			Labels:      m.TJob.Labels,
			Annotations: m.TJob.Annotations,
		},
		Spec: pytorchjobSpec,
	}
	return pytorchjob
}

func (m *PyTorchJobManager) GetStatus() *tj.TrainingJobStatus {
	status := tj.TrainingJobStatus{}
	if m.Job != nil {
		// convert(m.Job.Status, &status.JobStatus)
		status.JobStatus = m.Job.Status
	}

	status.State = getJobState(&status.JobStatus)

	return &status
}

func init() {
	JobManagerBuilder[pyTorchJobType] = func(tjob *tj.TrainingJob) JobManager { return NewPyTorchJobManager(tjob) }
	RegisteredSchemes[pyTorchJobType] = func() client.Object { return &kubeflowv1.PyTorchJob{} }
}
