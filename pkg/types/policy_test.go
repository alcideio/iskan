package types

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func Test_ReportFilter_Loader(t *testing.T) {
	filter := ReportFilter{
		Severities:      "CRITICAL,HIGH",
		CvssGreaterThan: 3.0,
		FixableOnly:     true,
		VulnerabilityExclusions: []*VulnerabilityExclusion{
			&VulnerabilityExclusion{
				CVE:         "CVE-2020-666666",
				Reason:      "Number of the beast",
				SnoozeUntil: 1,
				SnoozedBy:   "Devil",
			},
		},
	}

	f, err := ioutil.TempFile("/tmp/", "iskan-report-filter-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create file - %v", err)
	}
	fname := f.Name()
	defer os.Remove(fname)

	d, err := yaml.Marshal(&filter)
	if err != nil {
		t.Fatalf("Failed to marshal - %v", err)
	}
	fmt.Println(string(d))
	f.Write(d)
	f.Sync()
	f.Close()

	reportFilter, err := LoadReportFilter(fname)
	if err != nil {
		t.Fatalf("Failed to load - %v", err)
	}

	assertions := assert.New(t)
	assertions.Equal(&filter, reportFilter, "NOT EQUAL")
}

func Test_ImageRegistry_Filter(t *testing.T) {
	filter := ScanScope{
		NamespaceExclude: "excludeme",
		NamespaceInclude: "includeme",
		RegistryExclusion: []*RegistryExclusion{
			&RegistryExclusion{
				Registry:    "myreg.com/excluded",
				Reason:      "taking a break",
				SnoozeBy:    "devil",
				SnoozeUntil: time.Now().Add(time.Hour).Unix(),
			},
			&RegistryExclusion{
				Registry:    "myreg.com/expired",
				Reason:      "expired exception",
				SnoozeBy:    "devil",
				SnoozeUntil: time.Now().Add(-1 * time.Hour).Unix(),
			},
			&RegistryExclusion{
				Registry:    "myreg.com/forever",
				Reason:      "forever exclude",
				SnoozeBy:    "devil",
				SnoozeUntil: 0,
			},
		},
	}

	filter.Init()

	assertions := assert.New(t)
	assertions.Equal(filter.IsNamespaceIncluded("includeme"), true)
	assertions.Equal(filter.IsNamespaceIncluded("excludeme"), false)
	assertions.Equal(filter.IsNamespaceIncluded("other"), false)
	assertions.Equal(filter.IsNamespaceIncluded("inc"), false)

	shouldScan, reason := filter.ShouldScanImage("myreg.com/excluded/myimage:v1.0")
	assertions.Equal(shouldScan, false, reason)

	shouldScan, reason = filter.ShouldScanImage("myreg.com/iskan/zz:v1")
	assertions.Equal(shouldScan, true, reason)

	shouldScan, reason = filter.ShouldScanImage("myreg.com/expired/zz:v1")
	assertions.Equal(shouldScan, true, reason)

	shouldScan, reason = filter.ShouldScanImage("myreg.com/forever/zz:v1")
	assertions.Equal(shouldScan, false, reason)
}
