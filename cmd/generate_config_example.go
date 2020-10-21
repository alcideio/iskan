package cmd

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/vulnprovider"
	"github.com/spf13/cobra"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
	"time"
)

func NewCommandGenerateExample() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate Configuration Example",
	}

	cmd.AddCommand(
		NewCommandGenerateApiConfigExample(),
		NewCommandGenerateReportFilterExample(),
	)

	return cmd
}

func NewCommandGenerateApiConfigExample() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-config",
		Short: "Generate Api Configuration File Example",
		Run: func(cmd *cobra.Command, args []string) {
			c := types.VulnProvidersConfig{
				Providers: []types.VulnProviderConfig{
					{
						Kind:       vulnprovider.ProviderKind_ECR,
						Repository: "666666.dkr.ecr.us-west-2.amazonaws.com/iskan",
						Creds: types.VulnProviderAPICreds{
							ECR: &types.ECR{
								AccessKeyId:     "mykeyid",
								SecretAccessKey: "mysecretkey",
								Region:          "us-west-2",
							},
						},
					},
					{
						Kind:       vulnprovider.ProviderKind_ACR,
						Repository: "myrepo.azurecr.io/iskan",
						Creds: types.VulnProviderAPICreds{
							ACR: &types.Azure{
								TenantId:       "my-tenant-uuid",
								SubscriptionId: "my-subscription-uuid",
								ClientId:       "client-id",
								ClientSecret:   "client-secret",
								CloudName:      "AZUREPUBLICCLOUD",
							},
						},
					},
					{
						Kind:       vulnprovider.ProviderKind_GCR,
						Repository: "gcr.io/myproj",
						Creds: types.VulnProviderAPICreds{
							GCR: `{
  "type": "service_account",
  "project_id": "myproj",
  "private_key_id": "someprivatekey",
  "private_key": "-----BEGIN PRIVATE KEY-----\n\n-----END PRIVATE KEY-----\n",
  "client_email": "imagevulreader@someprivatekey.iam.gserviceaccount.com",
  "client_id": "someclientid",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/imagevulreader%40myproj.iam.gserviceaccount.com"
}
`,
						},
					},
					{
						Kind:       vulnprovider.ProviderKind_TRIVY,
						Repository: "*",
						Creds: types.VulnProviderAPICreds{
							Trivy: &types.TrivyConfig{
								CacheDir:      "/home/iskan/.cache/trivy",
								ReportsDir:    "/home/iskan/.cache/reports",
								DebugMode:     false,
								VulnType:      "os,library",
								Severity:      "UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL",
								IgnoreUnfixed: false,
								SkipUpdate:    false,
								//GitHubToken:   "",
								Insecure: false,
							},
						},
					},
				},
			}

			d, err := yaml.Marshal(&c)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(d))
		},
	}

	cmd.AddCommand()

	return cmd
}

func NewCommandGenerateReportFilterExample() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report-filter",
		Short: "Generate Report Filter File Example",
		Run: func(cmd *cobra.Command, args []string) {
			c := types.ReportFilter{
				Severities:      strings.Join([]string{grafeas.Severity_CRITICAL.String(), grafeas.Severity_HIGH.String()}, ","),
				CvssGreaterThan: 1.0,
				FixableOnly:     false,
				VulnerabilityExclusions: []*types.VulnerabilityExclusion{
					&types.VulnerabilityExclusion{
						CVE:         "CVE-2020-666",
						Reason:      "Nasty CVE",
						SnoozeUntil: time.Now().Unix(),
						SnoozedBy:   "yours-truely",
					},
				},
			}

			d, err := yaml.Marshal(&c)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(d))
		},
	}

	cmd.AddCommand()

	return cmd
}
