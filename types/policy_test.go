package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"testing"
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
