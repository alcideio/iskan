package types

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
)

func Test_ConfigLoader(t *testing.T) {
	c := RegistriesConfig{
		Registries: []RegistryConfig{
			{
				Kind:       "gcr",
				Repository: "us.gcr.io/k8s-artifacts-prod/external-dns",
				Creds: RegistryAPICreds{
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

	rc, err := LoadRegistriesConfig(fname)
	if err != nil {
		t.Fatalf("Failed to load - %v", err)
	}

	assertions := assert.New(t)
	assertions.Equal(&c, rc, "NOT EQUAL")
}
