package tests

import (
	"fmt"
	"github.com/alcideio/iskan/e2e/framework"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
	"os"
)

var _ = Describe("[regression][local-scanner] Test local scanner vulnerability providers", func() {
	Context("Test Scenario", func() {
		f, err := framework.NewDefaultFramework("local-scanner-trivy")
		nTestImages := 0
		testImageInfo := map[string]*framework.TestImageInfo{}
		It("Creates Pods with test images", func() {
			Expect(f).ToNot(BeNil())
			framework.ExpectNoError(err)

			By(fmt.Sprintf("Deploying all test images into namespace '%v' ", f.Namespace), func() {
				iterator := framework.NewTestImageIterator(framework.FilterByPublicRegistries())
				for {
					info, hasMore := iterator.Next()
					if !hasMore {
						break
					}
					nTestImages++

					framework.ExpectNoError(err)
					Expect(f).ToNot(BeNil())

					_, pod := f.DeployTestImage(info)
					testImageInfo[fmt.Sprintf("%v/%v", pod.Namespace, pod.Name)] = info
				}
			})

			var scanner *scan.ClusterScanner
			var report *types.ClusterScanReport
			policy := types.NewDefaultPolicy()
			policy.ScanScope.NamespaceInclude = f.Namespace

			trivyConfig := types.DefaultTrivyConfig()

			name, err := ioutil.TempDir("/tmp", fmt.Sprintf("%s-", "iskan-local-skan"))
			Expect(err).To(BeNil())
			defer os.RemoveAll(name)

			trivyConfig.CacheDir = name
			trivyConfig.ReportsDir = name

			config := types.VulnProvidersConfig{
				Providers: []types.VulnProviderConfig{
					{Kind: "trivy", Repository: "*", Creds: types.VulnProviderAPICreds{
						Trivy: trivyConfig,
					}},
				},
			}

			scanner = f.NewClusterScannerWithConfig(policy, &config)

			By(fmt.Sprintf("scanning the namespace '%v' within the cluster", f.Namespace), func() {
				Expect(scanner).NotTo(BeNil())
				report, err = scanner.Scan()
				framework.ExpectNoError(err)
				klog.V(5).Infof("%v", pretty.Sprint(report))
			})

			By("verifying that scan results match the image properties", func() {
				framework.Logf("scan report - %v", pretty.Sprint(report))
				Expect(report.Summary.AnalyzedPodCount).To(BeNumerically("==", uint32(nTestImages)))
				for fullPodName, info := range testImageInfo {
					if info.HasScanFailures() {
						framework.Logf("%v - %v", fullPodName, pretty.Sprint(report.Summary))
						failedPods := sets.NewString(report.Summary.FailedOrSkippedPods...)
						Expect(failedPods.Has(fullPodName)).To(BeTrue())
					} else {
						podSummary, exist := report.Summary.PodSummary[fullPodName]
						framework.Logf("%v - %v", fullPodName, pretty.Sprint(podSummary))
						Expect(podSummary.ScanFailures).To(BeZero())
						Expect(exist).To(BeTrue())

						if info.HasVulnerabilities() {
							var count uint32 = 0
							for _, n := range podSummary.Severity {
								count = count + n
							}
							Expect(count).To(BeNumerically(">", uint32(0)))
						}

						if info.HasCriticalVulnerabilities() {
							var count uint32 = podSummary.Severity[grafeas.Severity_CRITICAL]
							Expect(count).To(BeNumerically(">", 0))
						}
					}
				}

			})
		})
	})
})
