package tests

import (
	"fmt"
	"github.com/alcideio/iskan/e2e/framework"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/klog"
)

var _ = Describe("[regression][scan-scope] Scan Cluster with multiple vulnerability providers", func() {
	Context("Test Scenario", func() {
		f, err := framework.NewDefaultFramework("scan-scope")
		nTestImages := 0
		It("Creates Pods with test images", func() {
			Expect(f).ToNot(BeNil())
			framework.ExpectNoError(err)

			By(fmt.Sprintf("Deploying all test images into namespace '%v' ", f.Namespace), func() {
				iterator := framework.NewTestImageIterator(framework.FilterByPrivateRegistries())

				for {
					info, hasMore := iterator.Next()
					if !hasMore {
						break
					}
					nTestImages++

					framework.ExpectNoError(err)
					Expect(f).ToNot(BeNil())

					_, _ = f.DeployTestImage(info)
				}
			})

			var scanner *scan.ClusterScanner
			var report *types.ClusterScanReport
			policy := types.NewDefaultPolicy()
			policy.ScanScope.NamespaceInclude = f.Namespace
			policy.ScanScope.RegistryExclusion = []*types.RegistryExclusion{}
			iterator := framework.NewTestImageIterator(nil)
			for {
				info, hasMore := iterator.Next()
				if !hasMore {
					break
				}

				repo, _, _, _ := util.ParseImageName(info.Image)

				var exclusion *types.RegistryExclusion
				exclusion = &types.RegistryExclusion{
					Registry:    repo,
					Reason:      "",
					SnoozeBy:    "",
					SnoozeUntil: 0,
				}

				policy.ScanScope.RegistryExclusion = append(policy.ScanScope.RegistryExclusion, exclusion)
			}

			scanner = f.NewClusterScanner(policy)

			By(fmt.Sprintf("scanning the namespace '%v' within the cluster", f.Namespace), func() {
				Expect(scanner).NotTo(BeNil())
				report, err = scanner.Scan()
				framework.ExpectNoError(err)
				klog.V(5).Infof("%v", pretty.Sprint(report))
			})

			By("verifying that all images were excluded", func() {
				framework.Logf("scan report - %v", pretty.Sprint(report))
				Expect(report.Summary.ExcludedPodCount).To(BeNumerically(">", uint32(nTestImages)))
				Expect(report.Summary.AnalyzedPodCount).To(BeNumerically("==", uint32(0)))
				Expect(len(report.Summary.FailedOrSkippedPods)).To(BeZero())
			})
		})
	})
})
