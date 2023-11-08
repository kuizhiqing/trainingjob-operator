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
// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	kubeflowv1beta1 "github.com/kuizhiqing/trainingjob-operator/api/v1beta1"
	versioned "github.com/kuizhiqing/trainingjob-operator/client/clientset/versioned"
	internalinterfaces "github.com/kuizhiqing/trainingjob-operator/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/kuizhiqing/trainingjob-operator/client/listers/kubeflow/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TrainingJobInformer provides access to a shared informer and lister for
// TrainingJobs.
type TrainingJobInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.TrainingJobLister
}

type trainingJobInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTrainingJobInformer constructs a new informer for TrainingJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTrainingJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTrainingJobInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTrainingJobInformer constructs a new informer for TrainingJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTrainingJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1beta1().TrainingJobs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1beta1().TrainingJobs(namespace).Watch(context.TODO(), options)
			},
		},
		&kubeflowv1beta1.TrainingJob{},
		resyncPeriod,
		indexers,
	)
}

func (f *trainingJobInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTrainingJobInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *trainingJobInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kubeflowv1beta1.TrainingJob{}, f.defaultInformer)
}

func (f *trainingJobInformer) Lister() v1beta1.TrainingJobLister {
	return v1beta1.NewTrainingJobLister(f.Informer().GetIndexer())
}
