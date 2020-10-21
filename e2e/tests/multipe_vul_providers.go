package tests

import (
	"fmt"
	"github.com/alcideio/iskan/e2e/framework"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/klog"
)

var _ = Describe("[regression][multi-vul-providers] Scan Cluster with multiple vulnerability providers", func() {
	Context("Test Scenario", func() {
		f, err := framework.NewDefaultFramework("multi-vul-providers")
		It("Creates Pods with test images", func() {
			nTestImages := 0

			Expect(f).ToNot(BeNil())
			framework.ExpectNoError(err)

			By(fmt.Sprintf("Deploying all test images into namespace '%v' ", f.Namespace), func() {
				iterator := framework.NewTestImageIterator(framework.FilterByPrivateRegistries())
				for {
					info, hasMore := iterator.Next()
					if !hasMore {
						break
					}

					framework.ExpectNoError(err)
					Expect(f).ToNot(BeNil())

					_, _ = f.DeployTestImage(info)
					nTestImages++
				}
			})

			var scanner *scan.ClusterScanner
			var report *types.ClusterScanReport
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
				framework.Logf("scan report - %v", pretty.Sprint(report))
				Expect(report.Summary.AnalyzedPodCount).To(BeNumerically("==", uint32(nTestImages)))

			})
		})
	})
})
