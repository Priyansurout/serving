/*
Copyright 2020 The Knative Authors

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

package domainmapping

import (
	"context"

	"k8s.io/client-go/tools/cache"
	netclient "knative.dev/networking/pkg/client/injection/client"
	ingressinformer "knative.dev/networking/pkg/client/injection/informers/networking/v1alpha1/ingress"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/serving/pkg/apis/serving/v1alpha1"
	"knative.dev/serving/pkg/client/injection/informers/serving/v1alpha1/domainmapping"
	kindreconciler "knative.dev/serving/pkg/client/injection/reconciler/serving/v1alpha1/domainmapping"
)

// NewController creates a new DomainMapping controller.
func NewController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)
	domainmappingInformer := domainmapping.Get(ctx)
	ingressInformer := ingressinformer.Get(ctx)

	r := &Reconciler{
		ingressLister: ingressInformer.Lister(),
		netclient:     netclient.Get(ctx),
	}

	impl := kindreconciler.NewImpl(ctx, r)

	logger.Info("Setting up event handlers")
	domainmappingInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	handleControllerOf := cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(v1alpha1.Kind("DomainMapping")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	}
	ingressInformer.Informer().AddEventHandler(handleControllerOf)

	return impl
}