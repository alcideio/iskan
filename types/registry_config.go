package types

import (
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

type RegistryAPICreds struct {
	GCR string
}

type RegistryConfig struct {
	//Repo Kind
	Kind string

	//Repo FQDN
	Repository string

	//API Access Credentials
	Creds RegistryAPICreds
}

type RegistriesConfig struct {
	Registries []RegistryConfig
}

func LoadRegistriesConfig(fname string) (*RegistriesConfig, error) {
	rc := &RegistriesConfig{
		Registries: []RegistryConfig{},
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

func LoadRegistriesConfigFromBuffer(data []byte) (*RegistriesConfig, error) {
	rc := &RegistriesConfig{
		Registries: []RegistryConfig{},
	}

	err := yaml.Unmarshal(data, rc)
	if err != nil {
		return nil, err
	}

	return rc, err
}
