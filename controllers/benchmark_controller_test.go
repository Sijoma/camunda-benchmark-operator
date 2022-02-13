package controllers

import (
	"context"
	"time"

	v1 "k8s.io/api/apps/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	benchmarkv1alpha1 "github.com/sijoma/camunda-benchmark-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("Benchmark controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		BenchmarkName      = "test-benchmark"
		BenchmarkNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a Benchmark", func() {
		It("Should start to be in progress running", func() {
			By("By creating a new Benchmark")
			ctx := context.Background()
			benchmark := &benchmarkv1alpha1.Benchmark{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "cloud.camunda.io/v1alpha1",
					Kind:       "Benchmark",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      BenchmarkName,
					Namespace: BenchmarkNamespace,
				},
				Spec: benchmarkv1alpha1.BenchmarkSpec{
					CredentialsSecretName: "test-credentials",
					ProcessStarterRate:    100,
					StarterReplicas:       50,
					WorkerCount:           10,
					Duration:              "2s",
				},
			}
			Expect(k8sClient.Create(ctx, benchmark)).Should(Succeed())

			benchmarkLookupKey := types.NamespacedName{Name: BenchmarkName, Namespace: BenchmarkNamespace}
			createdBenchmark := &benchmarkv1alpha1.Benchmark{}

			// We'll need to retry getting this newly created Benchmark, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, benchmarkLookupKey, createdBenchmark)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdBenchmark.Spec.ProcessStarterRate).Should(Equal(100))

			By("By checking the Benchmark is running")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, benchmarkLookupKey, createdBenchmark)
				if err != nil {
					return "", err
				}
				return createdBenchmark.Status.Progress, nil
			}, duration, interval).Should(Equal("Running"))

			By("By checking that the deployment count is two")
			Eventually(func() (int, error) {
				deploymentList := v1.DeploymentList{
					TypeMeta: v1.Deployment{}.TypeMeta,
				}
				err := k8sClient.List(ctx, &deploymentList)
				if err != nil {
					return 0, err
				}
				return len(deploymentList.Items), err
			}, duration, interval).Should(Equal(2))

			By("By checking that Benchmark is getting finished")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, benchmarkLookupKey, createdBenchmark)
				if err != nil {
					return "", err
				}
				return createdBenchmark.Status.Progress, nil
			}, duration, interval).Should(Equal("Done"))

			By("By checking that the deployment count is cleaned up")
			Eventually(func() (int, error) {
				deploymentList := v1.DeploymentList{
					TypeMeta: v1.Deployment{}.TypeMeta,
				}
				err := k8sClient.List(ctx, &deploymentList)
				if err != nil {
					return 0, err
				}
				return len(deploymentList.Items), err
			}, duration, interval).Should(Equal(0))
		})
	})

})
