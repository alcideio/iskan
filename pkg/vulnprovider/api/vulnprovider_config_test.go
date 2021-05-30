package api

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func Test_ConfigLoader(t *testing.T) {
	c := VulnProvidersConfig{
		Providers: []VulnProviderConfig{
			{
				Kind:       "gcr",
				Repository: "us.gcr.io/k8s-artifacts-prod/external-dns",
				Creds: VulnProviderAPICreds{
					GCR: "someserviceaccount",
				},
			},
		},
	}

	f, err := ioutil.TempFile("/tmp/", "iskan-reg-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create file - %v", err)
	}
	fname := f.Name()
	defer os.Remove(fname)

	d, err := yaml.Marshal(&c)
	if err != nil {
		t.Fatalf("Failed to marshal - %v", err)
	}
	f.Write(d)
	f.Sync()
	f.Close()

	rc, err := LoadVulnProvidersConfig(fname)
	if err != nil {
		t.Fatalf("Failed to load - %v", err)
	}

	assertions := assert.New(t)
	assertions.Equal(&c, rc, "NOT EQUAL")
	assertions.Len(rc.Providers, 1, "Incorrect length")
}

func Test_ConfigLoaderFromBuffer(t *testing.T) {
	config := `
providers:
  - kind: "gcr"
    repository: "gcr.io/yourproject"
    creds:
      gcr: |
        {
          "type": "service_account",
          "project_id": "yourproject",
          "private_key_id": "XXX",
          "private_key": "",
          "client_email": "imagevulreader@yourproject.iam.gserviceaccount.com",
          "client_id": "666",
          "auth_uri": "https://accounts.google.com/o/oauth2/auth",
          "token_uri": "https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
          "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/imagevulreader%40yourproject.iam.gserviceaccount.com"
        }
  - kind: "ecr"
    repository: "yourawsaccount.dkr.ecr.us-west-2.amazonaws.com/iskan"
    creds:
      ecr:
        accessKeyId: AWSKEY
        secretAccessKey: AWSSECRET
        region: us-west-2
  - kind: "acr"
    repository: "alcide.azurecr.io/iskan"
    creds:
      acr:
        tenantId: mytenantid
        subscriptionId: subscrrptionId
        clientId: clientId
        clientSecret: clientsecret
        cloudName: "AZUREPUBLICCLOUD"
  - kind: "trivy"
    # Use "*" for a capture all images
    repository: "*"
    creds:
      trivy:
        debugMode: false
`

	rc, err := LoadVulnProvidersConfigFromBuffer([]byte(config))
	if err != nil {
		t.Fatalf("Failed to load - %v", err)
	}

	assertions := assert.New(t)
	assertions.Len(rc.Providers, 4, "Incorrect length")
}
