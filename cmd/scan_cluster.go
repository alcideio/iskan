package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alcideio/iskan/pkg/report"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

func NewCommandScanCluster() *cobra.Command {
	clusterContext := ""
	cliReportFilter := *(types.NewDefaultPolicy().ReportFilter)
	policy := types.NewDefaultPolicy()
	format := ""
	outfile := "-"
	vulAPIConfig := ""
	reportConfig := ""

	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"scan-cluster"},
		Short:   "Get vulnerabilities information on the presently running containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			config, err := types.LoadVulnProvidersConfig(vulAPIConfig)
			if err != nil {
				return err
			}

			var reportFilter *types.ReportFilter
			if reportConfig != "" {
				reportFilter, err = types.LoadReportFilter(reportConfig)
				if err != nil {
					return err
				}
			} else {
				reportFilter = &cliReportFilter
			}

			policy.ReportFilter = reportFilter
			policy.Init()

			scanner, err := scan.NewClusterScanner(clusterContext, policy, config)
			if err != nil {
				return err
			}

			util.ConsolePrinter(fmt.Sprintf("Scanning Cluster Context '%v'", color.HiBlueString(clusterContext)))
			start := time.Now()

			if policy.ScanScope.NamespaceInclude != "" {
				util.ConsolePrinter(fmt.Sprintf("Scan Scope ==> Including Namespaces '%v'", color.HiGreenString(policy.ScanScope.NamespaceInclude)))
			}

			if policy.ScanScope.NamespaceExclude != "" {
				util.ConsolePrinter(fmt.Sprintf("Scan Scope ==> Excluding Namespaces '%v'", color.HiRedString(policy.ScanScope.NamespaceExclude)))
			}

			if len(policy.ScanScope.RegistryExclusion) > 0 {
				//FIXME - SNOOZE
				for _, e := range policy.ScanScope.RegistryExclusion {
					util.ConsolePrinter(fmt.Sprintf("Scan Scope ==> Excluding Registry '%v'", color.HiRedString(e.Registry)))
				}
			}

			scanReport, err := scanner.Scan()
			if err != nil {
				return err
			}
			delta := time.Now().Sub(start)
			util.ConsolePrinter(fmt.Sprintf("Cluster UID '%v'", color.HiBlueString(scanReport.ClusterId)))

			util.ConsolePrinter(fmt.Sprintf("Scan Completed within '%v'", color.HiBlueString(delta.Round(time.Second).String())))
			util.ConsolePrinter(fmt.Sprintf("There are '%v' Pods - Analyzed '%v' Skipped '%v'",
				scanReport.Summary.AnalyzedPodCount+scanReport.Summary.ExcludedPodCount,
				scanReport.Summary.AnalyzedPodCount, scanReport.Summary.ExcludedPodCount))
			if scanReport.Policy.ReportFilter.Severities != "" {
				util.ConsolePrinter(fmt.Sprintf("Showing Findings with Severity level(s): '%v'", color.HiBlueString(scanReport.Policy.ReportFilter.Severities)))
			}
			if scanReport.Policy.ReportFilter.CvssGreaterThan > 0 {
				util.ConsolePrinter(fmt.Sprintf("Showing Findings with CVSS score higher than '%v'", color.HiBlueString(fmt.Sprint(scanReport.Policy.ReportFilter.CvssGreaterThan))))
			}
			if scanReport.Policy.ReportFilter.FixableOnly {
				util.ConsolePrinter(fmt.Sprintf("Showing Findings with fixable CVEs"))
			}

			util.ConsolePrinter(fmt.Sprintf("Severity Summary (uniques count across the scan scope)\n\n%v", scanReport.Summary.ClusterSeverity.Table()))
			util.ConsolePrinter(fmt.Sprintf("Severity Summary By Namespace (count is not of uniques)\n\n%v", scanReport.Summary.NamespaceSeverity.Table(nil)))
			util.ConsolePrinter(fmt.Sprintf("Severity Summary By Pod\n\n%v", scanReport.Summary.PodSeverity.Table(scanReport.Summary.PodFixableSeverity)))
			util.ConsolePrinter(fmt.Sprintf("Skipped, Failed or Missing Scan Info\n%v\n", strings.Join(scanReport.Summary.FailedOrSkippedPods, "\n")))
			util.ConsolePrinter(fmt.Sprintf("Incomplete Image Scan\n%v\n", strings.Join(scanReport.Summary.FailedOrSkippedImages, "\n")))
			o := os.Stdout
			if outfile != "-" {
				f, err := os.Create(outfile)
				if err != nil {
					util.ConsolePrinter(fmt.Sprintf("Failed to create '%v' - %v", color.HiBlueString(outfile), color.HiRedString(err.Error())))
				} else {
					util.ConsolePrinter(fmt.Sprintf("Saving Scan Report to '%v'", color.HiBlueString(outfile)))
					o = f
				}
			}

			switch format {
			case "json":
				encoder := json.NewEncoder(o)
				encoder.SetIndent("", "\t")
				return encoder.Encode(scanReport)
			case "yaml":
				data, err := yaml.Marshal(scanReport)
				if err != nil {
					return err
				}
				_, err = o.Write(data)
				return err
			case "advisor":
				advisorReport := scanner.GetAdvisorReport()
				data, err := yaml.Marshal(advisorReport)
				if err != nil {
					return err
				}
				_, err = o.Write(data)
				return err
			case "html":
				advisorReport := scanner.GetAdvisorReport()
				r := report.HtmlReport{
					Report: advisorReport,
				}
				data, err := r.Generate()
				if err != nil {
					return err
				}
				_, err = o.Write([]byte(data))
				return err
			default:
				return fmt.Errorf("'%v' is not supported", format)
			}
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&clusterContext, "cluster-context", "", "Cluster Context .use 'kubectl config get-contexts' to list available contexts")
	flags.StringVarP(&format, "format", "f", "json", "Output format. Supported formats: json | yaml | html")
	flags.StringVarP(&outfile, "outfile", "o", "alcide-iskan.report", "Output file name. Use '-' to output to stdout")
	flags.StringVarP(&vulAPIConfig, "api-config", "c", "", "The Vulnerability API configuration file name")
	flags.StringVarP(&reportConfig, "report-config", "r", "", "The Report configuration file name")

	ReportFilterFlags(&cliReportFilter, flags)
	ScanFilterFlags(policy.ScanScope, flags)
	ScanRateLimitFlags(&policy.RateLimit, flags)

	return cmd
}
