package tests

import (
	"github.com/alcideio/iskan/e2e/framework"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("[sanity][gcr] Cluster Scan Sanity", func() {
	f, _ := framework.NewDefaultFramework("gcr-sanity")
	Context("GCR", func() {
		RegistrySanity(f, "gcr")
	})
})
