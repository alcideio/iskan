package cmd

import (
	"fmt"
	"os"

	"github.com/alcideio/iskan/pkg/report"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/vulnprovider/api"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/flowcontrol"
)

func NewCommandScanImage() *cobra.Command {
	image := ""
	format := "json"
	vulAPIConfig := ""
	reportFilter := *(types.NewDefaultPolicy().ReportFilter)
	policy := types.NewDefaultPolicy()

	cmd := &cobra.Command{
		Use:     "image",
		Aliases: []string{"scan-image", "i", "container", "scan-container"},
		Short:   "Get vulnerabilities information for a given container image",
		Example: `iskan image --image="gcr.io/myproj/path/to/myimage:v1.0" --api-config myconfig.yaml -f table --filter-severity CRITICAL,HIGH`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if image == "" {
				return cmd.Usage()
			}

			config, err := api.LoadVulnProvidersConfig(vulAPIConfig)
			if err != nil {
				return err
			}

			//klog.V(10).Info("loaded supplied config", pretty.Sprint(config))
			regsConfig := map[string]*api.VulnProviderConfig{}
			for i, r := range config.Providers {
				regsConfig[r.Repository] = &config.Providers[i]
			}

			registryConfig := scan.RegistryConfigForImage(image, regsConfig)

			policy.ReportFilter = &reportFilter
			policy.Init()

			res, err := scan.ScanImage(image, policy, registryConfig, flowcontrol.NewFakeAlwaysRateLimiter())
			if err != nil {
				fmt.Println(err)
				return nil
			}

			if !res.CompletedOK {
				fmt.Println(res.Reason)
				return nil
			}

			return report.ReportImageScanResult(format, res, os.Stdout)
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&vulAPIConfig, "api-config", "c", "", "The Vulnerability API configuration file name")
	flags.StringVarP(&image, "image", "i", "", "container image for which vulnerabilities information should be obtained")
	flags.StringVarP(&format, "format", "f", "json", "Output format. Supported formats: json | yaml | table")
	ReportFilterFlags(&reportFilter, flags)

	return cmd
}
