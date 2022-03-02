package managedresource_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sacv1alpha1 "github.com/rabbitmq/single-active-consumer-operator/api/v1alpha1"
	"github.com/rabbitmq/single-active-consumer-operator/internal/managedresource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ = Describe("SuperstreamExchange", func() {
	var (
		superStreamConsumer           *sacv1alpha1.SuperStreamConsumer
		superStreamConsumerPodBuilder *managedresource.SuperStreamConsumerPodBuilder
		pod                           *corev1.Pod
		podSpec                       corev1.PodSpec
		scheme                        *runtime.Scheme
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(sacv1alpha1.AddToScheme(scheme)).To(Succeed())
		superStreamConsumer = &sacv1alpha1.SuperStreamConsumer{}
		superStreamConsumer.Name = "parent-set"
		superStreamConsumer.Namespace = "parent-namespace"
		superStreamConsumer.Spec.SuperStreamReference = sacv1alpha1.SuperStreamReference{
			Name:      "super-stream-1",
			Namespace: "parent-namespace",
		}

		podSpec = corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "example-container",
					Image: "example-image",
				},
			},
		}

		superStreamConsumerPodBuilder = managedresource.SuperStreamConsumerPod(superStreamConsumer, scheme, podSpec, "super-stream-1", "sample-partition")
		obj, _ := superStreamConsumerPodBuilder.Build()
		pod = obj.(*corev1.Pod)
	})

	Context("Build", func() {
		It("generates a pod object with the correct name", func() {
			Expect(pod.GenerateName).To(Equal("parent-set-sample-partition-"))
		})

		It("generates a pod object with the correct namespace", func() {
			Expect(pod.Namespace).To(Equal(superStreamConsumer.Namespace))
		})
		It("sets expected labels on the Pod", func() {
			Expect(pod.ObjectMeta.Labels).To(HaveKeyWithValue("rabbitmq.com/super-stream", "super-stream-1"))
			Expect(pod.ObjectMeta.Labels).To(HaveKeyWithValue("rabbitmq.com/super-stream-partition", "sample-partition"))
			Expect(pod.ObjectMeta.Labels).To(HaveKeyWithValue("rabbitmq.com/consumer-pod-spec-hash", "5963d9e83cb18c41"))
		})
		It("sets the podSpec", func() {
			Expect(pod.Spec).To(Equal(podSpec))
		})
	})

	Context("Update", func() {
		BeforeEach(func() {
			Expect(superStreamConsumerPodBuilder.Update(pod)).To(Succeed())
		})
		It("sets owner reference", func() {
			Expect(pod.OwnerReferences[0].Name).To(Equal(superStreamConsumer.Name))
		})
	})
})
