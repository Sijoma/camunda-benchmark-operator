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
	"time"

	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudv1alpha "github.com/sijoma/camunda-benchmark-operator/api/v1alpha1"
	"github.com/sijoma/camunda-benchmark-operator/internal"
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
		log.Error(err, "unable to fetch benchmark")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	bench := internal.NewBenchmark(benchmark)
	if benchmark.Status.Progress == "Done" {
		for _, resource := range bench.Resources {
			err := r.Delete(ctx, resource)
			if err != nil {
				log.Error(err, "unable to delete resource")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	for _, resource := range bench.Resources {
		controllerutil.CreateOrUpdate(ctx, r.Client, resource, func() error {
			return controllerutil.SetControllerReference(benchmark, resource, r.Scheme)
		})
	}

	if benchmark.Status.StartTime == nil {
		benchmark.Status.StartTime = &metav1.Time{Time: time.Now()}
	}

	// Reconcile function will rerun every 10 seconds
	result := ctrl.Result{
		Requeue:      true,
		RequeueAfter: 10,
	}
	currentTime := time.Now()
	benchmarkDuration, _ := time.ParseDuration(benchmark.Spec.Duration)
	if currentTime.After(benchmark.Status.StartTime.Time.Add(benchmarkDuration)) {
		result.Requeue = false
		benchmark.Status.Progress = "Done"
	} else {
		benchmark.Status.Progress = "Running"
	}
	err := r.Status().Update(ctx, benchmark)
	if err != nil {
		log.Error(err, "unable to update benchmark status")
		return result, err
	}

	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BenchmarkReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudv1alpha.Benchmark{}).
		Owns(&apps.Deployment{}).
		Complete(r)
}
