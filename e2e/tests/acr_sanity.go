package tests

import (
	"os"

	"github.com/alcideio/iskan/e2e/framework"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("[sanity][acr] Cluster Scan Sanity", func() {
	if os.Getenv("E2E_ACR_PULLSECRET") == "" {
		framework.Logf("Missing E2E_ACR_PULLSECRET")
		return
	}

	f, _ := framework.NewDefaultFramework("acr-sanity")
	Context("ACR", func() {
		RegistrySanity(f, "acr")
	})
})
