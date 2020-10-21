package tests

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/klog"

	"github.com/alcideio/iskan/e2e/framework"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
)

func RegistrySanity(f *framework.Framework, kind string) {
	Context("API Provider Sanity", func() {
		var err error
		iterator := framework.NewTestImageIterator(framework.FilterByKind(kind))

		for {
			var scanner *scan.ClusterScanner
			var report *types.ClusterScanReport

			info, hasMore := iterator.Next()
			if !hasMore {
				break
			}

			It(fmt.Sprint("Creates Pods with test image '", info.Image, "'"), func() {
				Expect(f).ToNot(BeNil())

				_, pod := f.DeployTestImage(info)

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

					fullPodName := fmt.Sprint(pod.Namespace, "/", pod.Name)
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

				})
			})
		}
	})
}
