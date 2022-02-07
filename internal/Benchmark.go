package internal

import (
	cloudv1 "github.com/sijoma/camunda-benchmark-operator/api/v1alpha1"
	"github.com/sijoma/camunda-benchmark-operator/internal/components"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Benchmark struct {
	Name      string
	Resources []client.Object
}

func NewBenchmark(benchmarkSpec *cloudv1.Benchmark) *Benchmark {
	benchmark := Benchmark{}

	// These objects currently constitute a "benchmark",
	// we pass no further configuration to them
	componentTemplates := []string{
		"simple-starter-deployment.yaml", "simple-service.yaml",
		"worker-deployment.yaml", "worker-service.yaml",
	}

	for _, template := range componentTemplates {
		comp := createComponent(benchmarkSpec.Namespace, template)
		benchmark.Resources = append(benchmark.Resources, comp)
	}

	return &benchmark
}

func createComponent(ns string, name string) client.Object {
	comp, err := components.Get(name, true, "")
	if err != nil {
		return nil
	}
	comp.SetNamespace(ns)
	return comp
}
