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
		"gcr": &gcrPullSecret,
		"ecr": &ecrPullSecret,
		"acr": &acrPullSecret,
		"":    &noSecret,
	},
	ApiConfigFile: "",
}

var gcrPullSecret = ""
var ecrPullSecret = ""
var acrPullSecret = ""
var noSecret = ""

func RegisterFrameworkFlags() {
	flag.StringVar(&gcrPullSecret, "iskan.gcr-pull-secret", os.Getenv("E2E_GCR_PULLSECRET"), "The pull secret one would place in a Kubernetes Secret Object. Use E2E_GCR_PULLSECRET")
	flag.StringVar(&ecrPullSecret, "iskan.ecr-pull-secret", os.Getenv("E2E_ECR_PULLSECRET"), "The pull secret one would place in a Kubernetes Secret Object. Use E2E_ECR_PULLSECRET")
	flag.StringVar(&acrPullSecret, "iskan.acr-pull-secret", os.Getenv("E2E_ACR_PULLSECRET"), "The pull secret one would place in a Kubernetes Secret Object. Use E2E_ACR_PULLSECRET")
	flag.StringVar(&GlobalConfig.ApiConfigFile, "iskan.api-config", os.Getenv("E2E_API_CONFIG"), "The Vulnerability API configuration - use E2E_API_CONFIG")
}
