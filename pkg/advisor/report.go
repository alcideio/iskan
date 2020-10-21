package advisor

import (
	"bytes"
	"fmt"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/pkg/version"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
	"time"
)

func GenerateAdvisorReport(res *types.ScanTaskResult) (*AdvisorClusterReport, error) {
	advisorReport := &AdvisorClusterReport{
		AdvisorReportHeader: AdvisorReportHeader{
			CreationTimeStamp: time.Now().Format(time.RFC3339),
			ReportUID:         rand.String(10),
			Info:              "Kubernetes Native Container Image Scanner - By Alcide",
			ScannerVersion:    fmt.Sprintf("%v-%v", version.Version, version.Commit),
			MSTimeStamp:       0,
			ProfileID:         "iskan",
		},
		Reports: map[string]*AdvisorReportData{},
	}

	reportData := &AdvisorReportData{
		ResourceName:      "Pod Vulnerability Scan",
		ResourceNamespace: "KubeAdvisor",
		ResourceUID:       "iskan",
		ResourceKind:      "Kubernetes Image Vulnerabilities Scan",
		Results:           []CheckResult{},
	}

	var missingImageInfo, failedImages sets.String
	var check *CheckResult

	for _, pod := range res.ScannedPods {
		missingImageInfo = sets.NewString()
		failedImages = sets.NewString()

		podSummary := types.NewSeveritySummary()

		podContainers := [][]v1.ContainerStatus{
			pod.Status.InitContainerStatuses,
			pod.Status.EphemeralContainerStatuses,
			pod.Status.ContainerStatuses,
		}

		findings := &AdditionalFindings{
			SeverityCount: nil,
			Severity:      "",
			Findings:      []*AdditionalFinding{},
		}

		for _, l := range podContainers {
			for _, c := range l {
				var finding *AdditionalFinding

				image := util.GetImageId(c.Image, c.ImageID)

				imageFindings, exist := res.Findings[image]
				if !exist {
					missingImageInfo.Insert(image)
					continue
				}

				if !imageFindings.CompletedOK {
					failedImages.Insert(imageFindings.Image)
					continue
				}

				for _, occurence := range imageFindings.Findings {
					finding = new(AdditionalFinding)

					vul := occurence.GetVulnerability()

					finding.Id = occurence.NoteName
					finding.References = make([]string, len(vul.RelatedUrls))
					for i, url := range vul.RelatedUrls {
						finding.References[i] = url.Url
					}
					finding.Title = vul.ShortDescription
					finding.Description = vul.LongDescription
					finding.Severity = vul.Severity.String()
					finding.Info = map[string]string{}
					finding.Info["Image"] = imageFindings.Image
					finding.Info["FixAvailable"] = fmt.Sprint(vul.FixAvailable)
					finding.Info["EffectiveSeverity"] = vul.EffectiveSeverity.String()
					finding.Info["CvssScore"] = fmt.Sprint(vul.CvssScore)
					for _, pkg := range vul.PackageIssue {
						finding.Info["AffectedPackage"] = pkg.AffectedPackage
						if pkg.AffectedVersion != nil {
							finding.Info["AffectedVersion"] = pkg.AffectedVersion.FullName
						}
						finding.Info["FixedPackage"] = pkg.FixedPackage
						if pkg.FixedVersion != nil {
							finding.Info["FixedVersion"] = pkg.FixedVersion.FullName
						}
					}

					findings.Findings = append(findings.Findings, finding)
				}

				podSummary.Add(imageFindings.Summary)
			}
		}

		check = new(CheckResult)
		check.Resource.Name = pod.Name
		check.Resource.Namespace = pod.Namespace
		check.Resource.Kind = "Pod"
		check.Resource.Group = ""
		check.Resource.Version = "v1"
		check.ResultUID = strings.ToLower(fmt.Sprintf("iskan.pod@%v@%v@%v", pod.Namespace, pod.Name, pod.UID))
		check.Check.ModuleId = "iskan.1"
		check.Check.ModuleTitle = "Kubernetes Image Vulnerabilities Scan"
		check.Check.GroupId = "1"
		check.Check.GroupTitle = "Pod Vulnerabilities Scan"
		check.Check.CheckId = "1"
		check.Check.CheckTitle = "Container Image Vulnerability Scan"
		check.CheckId = fmt.Sprint("iskan.1.1.1.", pod.UID)

		msgBuf := bytes.NewBufferString("")
		fmt.Fprintf(msgBuf, "Total number of identified vulnerabilties in '%v/%v' by severity: %v\n", pod.Namespace, pod.Name, podSummary.String())

		if missingImageInfo.Len() > 0 {
			fmt.Fprintf(msgBuf,
				"Note: The following images did not have vulnerability information: '%v'\n",
				strings.Join(missingImageInfo.List(), ","))
		}

		if failedImages.Len() > 0 {
			fmt.Fprintf(msgBuf,
				"Note: Extracting vulnerability information for the following images were not successful: '%v'\n",
				strings.Join(failedImages.List(), ","))
		}

		check.Action = AdmissionAction_Alert.String()
		check.Message = msgBuf.String()
		check.Recommendation = `"Distroless" images contain only your application and its runtime dependencies. 
They do not contain package managers, shells or any other programs you would expect to find in a standard Linux distribution.
Restricting what's in your runtime container to precisely what's necessary for your app is the best practice to employ.
It improves the signal to noise of CVE scanners and reduces the burden of establishing provenance to just what you need.
`
		severity, _ := podSummary.Max()
		switch severity {
		case grafeas.Severity_CRITICAL.String():
			check.Severity = CheckSeverity_Critical.String()
		case grafeas.Severity_HIGH.String():
			check.Severity = CheckSeverity_High.String()
		case grafeas.Severity_MEDIUM.String():
			check.Severity = CheckSeverity_Medium.String()
		case grafeas.Severity_LOW.String():
			check.Severity = CheckSeverity_Low.String()
		case grafeas.Severity_MINIMAL.String():
			check.Severity = CheckSeverity_Low.String()
		case "":
			check.Severity = CheckSeverity_Pass.String()
		}

		check.References = []string{
			"https://github.com/GoogleContainerTools/distroless",
		}

		check.AdditionalFindings = map[string]*AdditionalFindings{}
		check.AdditionalFindings["iskan"] = findings

		reportData.Results = append(reportData.Results, *check)
	}

	advisorReport.Reports["Kubernetes Image Scan"] = reportData

	return advisorReport, nil
}
