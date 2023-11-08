package controllers

import (
	"strconv"

	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	tj "github.com/kuizhiqing/trainingjob-operator/api/v1beta1"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

func contains(list []string, key string) bool {
	for _, k := range list {
		if k == key {
			return true
		}
	}
	return false
}

func convert(in interface{}, out interface{}) error {
	inn, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(inn, out)
}

func pickIntFromOptionalMap(opts map[string]string, key string, def int) int {
	if value, ok := opts[key]; ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return def
}

func pickStrFromOptionalMap(opts map[string]string, key string, def string) string {
	if value, ok := opts[key]; ok {
		return value
	}
	return def
}

func pickBoolFromOptionalMap(opts map[string]string, key string, def bool) bool {
	if value, ok := opts[key]; ok {
		if value == "true" || value == "True" || value == "TRUE" {
			return true
		}
		return false
	}
	return def
}

func hasConditionType(status *kubeflowv1.JobStatus, tp kubeflowv1.JobConditionType) bool {
	for _, condition := range status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			if condition.Type == tp {
				return true
			}
		}
	}
	return false
}

func getJobState(status *kubeflowv1.JobStatus) tj.JobState {
	if status.CompletionTime == nil {
		if hasConditionType(status, kubeflowv1.JobRunning) {
			return tj.Running
		} else if hasConditionType(status, kubeflowv1.JobRestarting) {
			return tj.Restarting
		} else if hasConditionType(status, kubeflowv1.JobCreated) {
			return tj.Created
		} else {
			return tj.Pending
		}
	} else {
		if hasConditionType(status, kubeflowv1.JobSucceeded) {
			return tj.Succeeded
		} else if hasConditionType(status, kubeflowv1.JobFailed) {
			return tj.Failed
		}
	}
	return tj.Unknown
}
