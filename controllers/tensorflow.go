package controllers

import (
	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	tj "github.com/kuizhiqing/trainingjob-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	tfJobType              = "tfjob"
	successPolicyKey       = "successPolicy"
	enableDynamicWorkerKey = "enableDynamicWorker"
)

type TFJobManager struct {
	TJob *tj.TrainingJob
	Job  *kubeflowv1.TFJob
}

func NewTFJobManager(tjob *tj.TrainingJob) *TFJobManager {
	return &TFJobManager{
		TJob: tjob,
		Job:  &kubeflowv1.TFJob{},
	}
}

func (m *TFJobManager) GetTJob() *tj.TrainingJob {
	return m.TJob
}

func (m *TFJobManager) GetObject() client.Object {
	return m.Job
}

func (m *TFJobManager) NewJob() client.Object {
	// runPolicy := kubeflowv1.RunPolicy{}
	// convert(m.TJob.Spec.RunPolicy, &runPolicy)

	// replicaSpecs := make(map[kubeflowv1.ReplicaType]*kubeflowv1.ReplicaSpec)
	// convert(m.TJob.Spec.ReplicaSpecs, &replicaSpecs)

	successPolicy := kubeflowv1.SuccessPolicy(pickStrFromOptionalMap(m.TJob.Spec.OptionalMap, successPolicyKey, ""))
	enableDynamicWorker := pickBoolFromOptionalMap(m.TJob.Spec.OptionalMap, enableDynamicWorkerKey, false)

	runPolicy := m.TJob.Spec.RunPolicy

	replicaSpecs := m.TJob.Spec.ReplicaSpecs

	tfjobSpec := kubeflowv1.TFJobSpec{
		RunPolicy:           runPolicy,
		TFReplicaSpecs:      replicaSpecs,
		SuccessPolicy:       &successPolicy,
		EnableDynamicWorker: enableDynamicWorker,
	}

	tfjob := &kubeflowv1.TFJob{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   m.TJob.Namespace,
			Name:        m.TJob.Name,
			Labels:      m.TJob.Labels,
			Annotations: m.TJob.Annotations,
		},
		Spec: tfjobSpec,
	}
	return tfjob
}

func (m *TFJobManager) GetStatus() *tj.TrainingJobStatus {
	status := tj.TrainingJobStatus{}
	if m.Job != nil {
		// convert(m.Job.Status, &status.JobStatus)
		status.JobStatus = m.Job.Status
	}

	status.State = getJobState(&status.JobStatus)

	return &status
}

func init() {
	JobManagerBuilder[tfJobType] = func(tjob *tj.TrainingJob) JobManager { return NewTFJobManager(tjob) }
	RegisteredSchemes[tfJobType] = func() client.Object { return &kubeflowv1.TFJob{} }
}
