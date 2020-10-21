package tests

import (
	"github.com/alcideio/iskan/e2e/framework"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("[sanity][ecr] Cluster Scan Sanity", func() {
	f, _ := framework.NewDefaultFramework("ecr-sanity")
	Context("ECR", func() {
		RegistrySanity(f, "ecr")
	})
})
