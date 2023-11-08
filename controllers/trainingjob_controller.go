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

package controllers

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	tj "github.com/kuizhiqing/trainingjob-operator/api/v1beta1"
)

type JobManager interface {
	GetObject() client.Object
	NewJob() client.Object
	GetStatus() *tj.TrainingJobStatus
	GetTJob() *tj.TrainingJob
}

type JobManagerFunc func(tjob *tj.TrainingJob) JobManager
type GetObjectFunc func() client.Object

var (
	ctrlRefKey        = ".metadata.controller"
	apiGVStr          = tj.GroupVersion.String()
	JobManagerBuilder = map[string]JobManagerFunc{}
	RegisteredSchemes = map[string]GetObjectFunc{}
)

// TrainingJobReconciler reconciles a TrainingJob object
type TrainingJobReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=kubeflow.org,resources=trainingjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeflow.org,resources=trainingjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeflow.org,resources=trainingjobs/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *TrainingJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logr := log.FromContext(ctx)

	var tjob tj.TrainingJob
	if err := r.Get(ctx, req.NamespacedName, &tjob); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logr.Info("Reconcile", "version", tjob.ResourceVersion, "ddl", tjob.ObjectMeta.DeletionTimestamp)

	if tjob.Status.JobStatus.CompletionTime != nil {
		return ctrl.Result{}, nil
	}

	if mfc, ok := JobManagerBuilder[tjob.Spec.Type]; ok {
		m := mfc(&tjob)
		return r.doReconcile(ctx, req, m)
	} else {
		logr.Info("JobTypeError", "JobType", tjob.Spec.Type)
	}
	// TODO(kuizhiqing) handle run/clean policy

	return ctrl.Result{}, nil
}

func (r *TrainingJobReconciler) doReconcile(ctx context.Context, req ctrl.Request, jobManager JobManager) (ctrl.Result, error) {
	obj := jobManager.GetObject()
	tjob := jobManager.GetTJob()
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			job := jobManager.NewJob()
			if err = ctrl.SetControllerReference(tjob, job, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			err = r.createResource(ctx, tjob, job)
		} else {
			// TODO(kuizhiqing): sync spec
		}
		return ctrl.Result{}, err
	} else {
		status := jobManager.GetStatus()
		if !equality.Semantic.DeepEqual(status, tjob.Status) {
			tjob.Status = *status
			err = r.Status().Update(ctx, tjob)
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *TrainingJobReconciler) deleteResource(ctx context.Context, tjob *tj.TrainingJob, obj client.Object) error {
	if obj.GetDeletionTimestamp() != nil {
		return nil
	}
	tp := obj.GetObjectKind().GroupVersionKind().Kind
	if err := r.Delete(ctx, obj, client.PropagationPolicy(metav1.DeletePropagationBackground)); (err) != nil {
		r.Recorder.Event(tjob, corev1.EventTypeWarning, "Delete", fmt.Sprintf("delete failed %s %s", tp, obj.GetName()))
		return err
	}
	r.Recorder.Event(tjob, corev1.EventTypeNormal, "Deleted", fmt.Sprintf("deleted %s %s", tp, obj.GetName()))
	return nil
}

func (r *TrainingJobReconciler) createResource(ctx context.Context, tjob *tj.TrainingJob, obj client.Object) error {
	tp := obj.GetObjectKind().GroupVersionKind().Kind
	if err := r.Create(ctx, obj); err != nil {
		r.Recorder.Event(tjob, corev1.EventTypeWarning, "Create", fmt.Sprintf("create failed %s %s", tp, obj.GetName()))
		return err
	}
	r.Recorder.Event(tjob, corev1.EventTypeNormal, "Created", fmt.Sprintf("created %s %s", tp, obj.GetName()))
	return nil

}

func indexerFunc(rawObj client.Object) []string {
	owner := metav1.GetControllerOf(rawObj)
	if owner == nil {
		return nil
	}
	if owner.APIVersion != apiGVStr || owner.Kind != tj.KIND {
		return nil
	}

	// ...and if so, return it
	return []string{owner.Name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TrainingJobReconciler) SetupWithManager(mgr ctrl.Manager, enabledSchemes string) error {
	logr := log.FromContext(context.Background())

	builder := ctrl.NewControllerManagedBy(mgr).For(&tj.TrainingJob{})

	schemes := strings.Split(enabledSchemes, ",")
	for k := range RegisteredSchemes {
		if len(schemes) == 0 || contains(schemes, k) {
			logr.Info("Setup", "enable", k)
			if err := mgr.GetFieldIndexer().
				IndexField(context.Background(), RegisteredSchemes[k](), ctrlRefKey, indexerFunc); err != nil {
				return err
			}
			builder = builder.Owns(RegisteredSchemes[k]())
		} else {
			delete(RegisteredSchemes, k)
			delete(JobManagerBuilder, k)
		}
	}

	return builder.Complete(r)
}
