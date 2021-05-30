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
	api "github.com/alcideio/iskan/pkg/vulnprovider/api"
)

func ImageScanSanity(f *framework.Framework, config *api.VulnProvidersConfig, tags []string) {
	Context("API Provider Sanity", func() {
		var err error
		iterator := framework.NewTestImageIterator(framework.FilterByHasAnyTag(sets.NewString(tags...)))

		for {
			var scanner *scan.ImageScanner
			var report *api.ImageScanResult

			info, hasMore := iterator.Next()
			if !hasMore {
				break
			}

			It(fmt.Sprint("Scan test image '", info.Image, "'"), func() {
				Expect(f).ToNot(BeNil())

				policy := types.NewDefaultPolicy()

				scanner = f.NewImageScannerWithConfig(policy, config)
				By(fmt.Sprintf("scanning image '%v' within the cluster", info.Image), func() {
					Expect(scanner).NotTo(BeNil())
					report, err = scanner.Scan(info.Image)
					framework.ExpectNoError(err)
					klog.V(5).Infof("%v", pretty.Sprint(report))
				})

				By("verifying that scan results match the image properties", func() {
					//Expect(report.Summary.AnalyzedPodCount).To(Equal(uint32(1)))

					//fullPodName := fmt.Sprint(pod.Namespace, "/", pod.Name)
					if info.HasScanFailures() {
						framework.Logf("%v - %v", info.Image, pretty.Sprint(report.Summary))
						Expect(report.CompletedOK).To(BeFalse())
					} else {
						framework.Logf("%v - %v", info.Image, pretty.Sprint(report.Summary))

						if info.HasVulnerabilities() {
							Expect(len(report.Findings)).To(BeNumerically(">", 0))
						}

						if info.HasCriticalVulnerabilities() {
							var count uint32 = report.Summary[grafeas.Severity_CRITICAL]
							Expect(count).To(BeNumerically(">", 0))
						}
					}

				})
			})
		}
	})
}
