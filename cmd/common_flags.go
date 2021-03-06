package cmd

import (
	"github.com/alcideio/iskan/pkg/types"
	flag "github.com/spf13/pflag"
)

func ReportFilterFlags(filter *types.ReportFilter, flags *flag.FlagSet) {
	flags.Float32Var(&filter.CvssGreaterThan, "filter-cvss", 0, "Include CVEs with CVSS score greater or equal than the specified number. Valid values: 0.0-10.0")
	flags.BoolVar(&filter.FixableOnly, "filter-fixable-only", false, "Include CVEs with which are fixable")
	flags.StringVar(&filter.Severities, "filter-severity", "", "Select which severities to include. Comma seperated MINIMAL,LOW,MEDIUM,HIGH,CRITICAL")
	//flags.DurationVar(&filter.CveNewerThan, "filter-newer-than", 0, "Only show CVEs newer than the specified duration.")
}

func ScanFilterFlags(filter *types.ScanScope, flags *flag.FlagSet) {
	flags.StringVar(&filter.NamespaceExclude, "namespace-exclude", "kube-system", "Namespaces to exclude from the scan")
	flags.StringVar(&filter.NamespaceInclude, "namespace-include", "*", "Namespaces to include in the scan")
}

func ScanRateLimitFlags(filter *types.ScanRateLimit, flags *flag.FlagSet) {
	flags.Float32Var(&filter.ApiQPS, "scan-api-qps", 30, "Indicates the maximum QPS to the vuln providers")
	flags.Int32Var(&filter.ApiBurst, "scan-api-burst", 100, "Maximum burst for throttle")
}
