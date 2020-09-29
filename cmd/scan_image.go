package cmd

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/report"
	"os"

	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/types"
	"github.com/spf13/cobra"
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

			config, err := types.LoadRegistriesConfig(vulAPIConfig)
			if err != nil {
				return err
			}

			regsConfig := map[string]*types.RegistryConfig{}
			for i, r := range config.Registries {
				regsConfig[r.Repository] = &config.Registries[i]
			}

			registryConfig := scan.RegistryConfigForImage(image, regsConfig)

			policy.ReportFilter = &reportFilter
			policy.Init()

			res, err := scan.ScanImage(image, policy, registryConfig)
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
