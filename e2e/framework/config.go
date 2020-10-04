package framework

import (
	"flag"
	"os"
)

type Config struct {
	PullSecrets map[string]*string

	ApiConfigFile string
}

var GlobalConfig Config = Config{
	PullSecrets: map[string]*string{
		"gcr": &gcrPullsecret,
	},
	ApiConfigFile: "",
}

var gcrPullsecret = ""

func RegisterFrameworkFlags() {
	flag.StringVar(&gcrPullsecret, "iskan.gcr-pull-secret", os.Getenv("E2E_GCR_PULLSECRET"), "The pull secret one would place in a Kubernetes Secret Object. Use E2E_GCR_PULLSECRET")
	flag.StringVar(&GlobalConfig.ApiConfigFile, "iskan.api-config", os.Getenv("E2E_API_CONFIG"), "The Vulnerability API configuration - use E2E_API_CONFIG")
}
