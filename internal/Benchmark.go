package internal

import (
	cloudv1alpha "github.com/sijoma/camunda-benchmark-operator/api/v1alpha1"
	"github.com/sijoma/camunda-benchmark-operator/internal/components"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type Benchmark struct {
	Name      string
	Resources []client.Object
}

func NewBenchmark(benchmarkSpec *cloudv1alpha.Benchmark) *Benchmark {
	benchmark := Benchmark{}

	// These objects currently constitute a "benchmark",
	// we pass no further configuration to them
	componentTemplates := []string{
		"simple-starter-deployment.yaml",
		"worker-deployment.yaml",
	}

	for _, template := range componentTemplates {
		comp := createComponent(benchmarkSpec, template)
		benchmark.Resources = append(benchmark.Resources, comp)
	}

	return &benchmark
}

func createComponent(bm *cloudv1alpha.Benchmark, name string) client.Object {
	benchmarkInfos := map[string]string{
		"credentialsSecretName": bm.Spec.CredentialsSecretName,
		"starterReplicas":       strconv.Itoa(bm.Spec.StarterReplicas),
		"simpleStarterRate":     strconv.Itoa(bm.Spec.ProcessStarterRate),
		"workerReplicas":        strconv.Itoa(bm.Spec.WorkerCount),
	}
	comp, err := components.Get(name, true, benchmarkInfos)
	if err != nil {
		return nil
	}
	comp.SetNamespace(bm.Namespace)
	comp.SetName(bm.Name + "-" + comp.GetName())
	return comp
}
