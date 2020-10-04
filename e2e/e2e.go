package e2e

import (
	"k8s.io/apimachinery/pkg/util/rand"

	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	runtimeutils "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	// ensure auth plugins are loaded
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	//The tests
	_ "github.com/alcideio/iskan/e2e/tests"
)

const (
	// namespaceCleanupTimeout is how long to wait for the namespace to be deleted.
	// If there are any orphaned namespaces to clean up, this test is running
	// on a long lived cluster. A long wait here is preferably to spurious test
	// failures caused by leaked resources from a previous test run.
	namespaceCleanupTimeout = 15 * time.Minute
)

var _ = ginkgo.SynchronizedBeforeSuite(
	func() []byte {
		// Reference common test to make the import valid.
		//setupSuite()
		return nil
	},
	func(data []byte) {
		// Run on all Ginkgo nodes
		//setupSuitePerGinkgoNode()
	})

var _ = ginkgo.SynchronizedAfterSuite(
	func() {
		//CleanupSuite()
	},
	func() {
		//AfterSuiteActions()
	})

// RunE2ETests checks configuration parameters (specified through flags) and then runs
// E2E tests using the Ginkgo runner.
// If a "report directory" is specified, one or more JUnit test reports will be
// generated in this directory, and cluster logs will also be saved.
// This function is called on each Ginkgo node in parallel mode.
func RunE2ETests(t *testing.T) {
	runtimeutils.ReallyCrash = true
	logs.InitLogs()
	defer logs.FlushLogs()

	// Disable skipped tests unless they are explicitly requested.
	if config.GinkgoConfig.FocusString == "" && config.GinkgoConfig.SkipString == "" {
		config.GinkgoConfig.SkipString = `\[Flaky\]|\[Feature:.+\]`
	}

	// Run tests through the Ginkgo runner with output to console + JUnit for Jenkins
	var r []ginkgo.Reporter

	klog.Infof("Starting e2e run %q on Ginkgo node %d", rand.String(5), config.GinkgoConfig.ParallelNode)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Kubernetes e2e suite", r)
}
