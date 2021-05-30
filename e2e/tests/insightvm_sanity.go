package tests

import (
	"fmt"
	. "github.com/onsi/ginkgo"

	"github.com/alcideio/iskan/e2e/framework"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"
)

var _ = Describe("[regression][insightvm] Test InsightVM vulnerability providers", func() {
	Context("InsightVM", func() {
		var config *types.VulnProvidersConfig
		tstConfig := types.VulnProvidersConfig{
			Providers: []types.VulnProviderConfig{
				{Kind: "insightvm", Repository: "*", Creds: types.VulnProviderAPICreds{
					InsightVM: &types.InsightVM{
						ApiKey: "",
						Region: "",
					},
				}},
			},
		}

		f, err := framework.NewDefaultFramework("insightvm")

		It(fmt.Sprint("Init test successfully"), func() {
			config, err = types.LoadVulnProvidersConfigFromBuffer([]byte(framework.GlobalConfig.ApiConfigFile))
			framework.ExpectNoError(err)
			//framework.Logf("loaded supplied config - \n%v\n - \n%v\n", framework.GlobalConfig.ApiConfigFile, pretty.Sprint(*config))

			//Copy relevant API-Key
			for _, p := range config.Providers {
				if p.Kind == "insightvm" {
					//tstConfig.Providers[0].Repository = p.Repository
					tstConfig.Providers[0].Creds.InsightVM.ApiKey = p.Creds.InsightVM.ApiKey
					tstConfig.Providers[0].Creds.InsightVM.Region = p.Creds.InsightVM.Region
				}
			}

		})

		ImageScanSanity(f, &tstConfig, []string{"insightvm"})

	})
})
