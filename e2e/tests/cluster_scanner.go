package tests

import (
	"fmt"

	"github.com/alcideio/iskan/e2e/framework"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/types"
	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/klog"
)

var _ = Describe("Cluster Scan Sanity", func() {
	Describe("Verify containers scan with default policy", func() {
		f, err := framework.NewDefaultFramework("scratch")
		Context("Scan Test Images", func() {
			for img, info := range framework.TestImages {
				var secret *v1.Secret
				var pod *v1.Pod
				var scanner *scan.ClusterScanner
				var report *types.ClusterScanReport

				It("Creates Pods with test images", func() {
					framework.ExpectNoError(err)
					Expect(f).ToNot(BeNil())

					framework.Logf("Testing '%v'", info.Description)

					secretName := fmt.Sprintf("secret-%v", rand.String(4))
					By(fmt.Sprintf("creating image pull secret '%v'", secretName), func() {
						pullSecret, exist := framework.GlobalConfig.PullSecrets[info.PullSecret]
						Expect(exist).NotTo(BeFalse())
						Expect(pullSecret).NotTo(BeNil())
						secret = f.CreateImagePullSecret(secretName, *pullSecret)
						Expect(secret).NotTo(BeNil())
					})

					podName := fmt.Sprintf("pod-%v", rand.String(4))
					By(fmt.Sprintf("creating pod '%v' that use the image '%v' ", podName, img), func() {
						pod = f.CreatePodWithContainerImage(podName, img, secret.Name)
						Expect(pod).NotTo(BeNil())
					})

					policy := types.NewDefaultPolicy()
					policy.ScanScope.NamespaceInclude = f.Namespace
					scanner = f.NewClusterScanner(policy)
					By(fmt.Sprintf("scanning the namespace '%v' within the cluster", f.Namespace), func() {
						Expect(scanner).NotTo(BeNil())
						report, err = scanner.Scan()
						framework.ExpectNoError(err)
						klog.V(5).Infof("%v", pretty.Sprint(report))
					})

					By("verifying that scan results match the image properties", func() {
						Expect(report.Summary.AnalyzedPodCount).To(Equal(uint32(1)))
						podSummary, exist := report.Summary.PodSummary[fmt.Sprint(pod.Namespace, "/", pod.Name)]
						Expect(exist).To(BeTrue())
						Expect(podSummary.ScanFailures).To(BeZero())
						//Expect(pod).NotTo(BeNil())
					})
				})
			}
		})
	})
})
