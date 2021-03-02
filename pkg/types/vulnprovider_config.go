package types

import (
	"io/ioutil"
	"sigs.k8s.io/yaml"

	"github.com/kelseyhightower/envconfig"
)

type VulnProviderAPICreds struct {
	GCR string `json:"gcr,omitempty"`

	ECR *ECR `json:"ecr,omitempty"`

	ACR *Azure `json:"acr,omitempty"`

	Trivy *TrivyConfig `json:"trivy,omitempty"`

	Harbor *HarborConfig `json:"harbor,omitempty"`

	InsightVM *InsightVM `json:"insightvm,omitempty"`
}

type InsightVM struct {
	//https://help.rapid7.com/insightvm/en-us/api/api.html#section/Overview
	ApiKey string `json:"apiKey,omitempty" envconfig:"INSIGHTVM_APIKEY" default:""`
	Region string `json:"region,omitempty" envconfig:"INSIGHTVM_REGION" default:"us"`
}

type HarborConfig struct {
	Insecure bool   `json:"insecure,omitempty" envconfig:"HARBOR_INSECURE" default:"false"`
	Host     string `json:"host,omitempty" envconfig:"HARBOR_HOST" default:"localhost"`
	Username string `json:"username,omitempty" envconfig:"HARBOR_USERNAME" default:""`
	Password string `json:"password,omitempty" envconfig:"HARBOR_PASSWORD" default:""`
}

type TrivyConfig struct {
	CacheDir      string `json:"cacheDir,omitempty" envconfig:"SCANNER_TRIVY_CACHE_DIR" default:"/home/iskan/.cache/trivy"`
	ReportsDir    string `json:"reportsDir,omitempty" envconfig:"SCANNER_TRIVY_REPORTS_DIR" default:"/home/iskan/.cache/reports"`
	DebugMode     bool   `json:"debugMode,omitempty" envconfig:"SCANNER_TRIVY_DEBUG_MODE" default:"false"`
	VulnType      string `json:"vulnType,omitempty" envconfig:"SCANNER_TRIVY_VULN_TYPE" default:"os,library"`
	Severity      string `json:"severity,omitempty" envconfig:"SCANNER_TRIVY_SEVERITY" default:"UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL"`
	IgnoreUnfixed bool   `json:"ignoreUnfixed,omitempty" envconfig:"SCANNER_TRIVY_IGNORE_UNFIXED" default:"false"`
	SkipUpdate    bool   `json:"skipUpdate,omitempty" envconfig:"SCANNER_TRIVY_SKIP_UPDATE" default:"false"`
	GitHubToken   string `json:"githubToken,omitempty" envconfig:"SCANNER_TRIVY_GITHUB_TOKEN"`
	Insecure      bool   `json:"insecure,omitempty" envconfig:"SCANNER_TRIVY_INSECURE" default:"false"`
}

func DefaultTrivyConfig() *TrivyConfig {
	s := TrivyConfig{}
	envconfig.Process("SCANNER_TRIVY", &s)

	return &s
}

type ECR struct {
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
	Region          string `json:"region,omitempty"`
}

type Azure struct {
	// Azure Related
	// https://docs.microsoft.com/en-us/cli/azure/ad/sp?view=azure-cli-latest#az_ad_sp_reset_credentials
	TenantId       string `json:"tenantId,omitempty" envconfig:"AZURE_TENANT_ID"       default:""                  doc:"Tenant Id"`
	SubscriptionId string `json:"subscriptionId,omitempty" envconfig:"AZURE_SUBSCRIPTION_ID" default:""                  doc:"Subscription Id"`
	ClientId       string `json:"clientId,omitempty" envconfig:"AZURE_CLIENT_ID"       default:""                  doc:"Client Id"`
	ClientSecret   string `json:"clientSecret,omitempty" envconfig:"AZURE_CLIENT_SECRET"   default:""                  doc:"Client Secret"`
	//ResourceGroupName string `json:"resourceGroupName,omitempty" envconfig:"AZURE_RG_NAME"         default:""                  doc:"Resource Group Name"`
	CloudName string `json:"cloudName,omitempty" envconfig:"AZURE_CLOUD_NAME"      default:"AZUREPUBLICCLOUD"  doc:"AZUREPUBLICCLOUD, AZURECHINACLOUD, AZUREUSGOVERNMENTCLOUD, AZUREGERMANCLOUD"`
}

type VulnProviderConfig struct {
	//Repo Kind
	Kind string

	//Repo FQDN
	Repository string

	//API Access Credentials
	Creds VulnProviderAPICreds
}

type VulnProvidersConfig struct {
	Providers []VulnProviderConfig
}

func LoadVulnProvidersConfig(fname string) (*VulnProvidersConfig, error) {
	rc := &VulnProvidersConfig{
		Providers: []VulnProviderConfig{},
	}
	if fname == "" {
		return rc, nil
	}

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, rc)
	if err != nil {
		return nil, err
	}

	return rc, err
}

func LoadVulnProvidersConfigFromBuffer(data []byte) (*VulnProvidersConfig, error) {
	rc := &VulnProvidersConfig{
		Providers: []VulnProviderConfig{},
	}

	if len(data) == 0 {
		return rc, nil
	}

	err := yaml.Unmarshal(data, rc)
	if err != nil {
		return nil, err
	}

	return rc, err
}
