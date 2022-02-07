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

package controllers

import (
	"context"

	"github.com/sijoma/camunda-benchmark-operator/internal"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudv1alpha "github.com/sijoma/camunda-benchmark-operator/api/v1alpha1"
)

// BenchmarkReconciler reconciles a Benchmark object
type BenchmarkReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// CRUD core: secrets
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;

// CRUD core: events, services and configmaps
// +kubebuilder:rbac:groups=core,resources=events;services;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status;configmaps/status,verbs=get

// CRUD apps: deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

//+kubebuilder:rbac:groups=cloud.camunda.io,resources=benchmarks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.camunda.io,resources=benchmarks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.camunda.io,resources=benchmarks/finalizers,verbs=update

func (r *BenchmarkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var benchmark = new(cloudv1alpha.Benchmark)
	if err := r.Get(ctx, req.NamespacedName, benchmark); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	bench := internal.NewBenchmark(benchmark)
	for _, resource := range bench.Resources {
		// Set owner reference.
		err := controllerutil.SetControllerReference(benchmark, resource, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Apply the object data.
		force := true
		err = r.Client.Patch(ctx, resource, client.Apply, &client.PatchOptions{Force: &force, FieldManager: "benchmark-reconciler"})
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BenchmarkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha.Benchmark{}).
		Owns(&core.Service{}).
		Owns(&apps.Deployment{}).
		Complete(r)
}
