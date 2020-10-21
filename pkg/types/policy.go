package types

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/util"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

type VulnerabilityExclusion struct {
	CVE         string
	Reason      string
	SnoozeUntil int64
	SnoozedBy   string
}

//Evaluate Each Occurence Against filter to exclude occurences from report.
// The conditions are ANDed - if any if the conditions
type ReportFilter struct {
	//Empty Severity --> Include All Severity
	Severities string

	//CVSS Score is greater than the specified threshold
	CvssGreaterThan float32

	//Include only fixable vulnerabilities
	FixableOnly bool

	//Skip CVEs Older than X Days
	//CveNewerThan time.Duration

	//Specific CVEs
	VulnerabilityExclusions []*VulnerabilityExclusion
}

func LoadReportFilter(fname string) (*ReportFilter, error) {
	rc := &ReportFilter{
		Severities:              "",
		CvssGreaterThan:         0,
		FixableOnly:             false,
		VulnerabilityExclusions: []*VulnerabilityExclusion{},
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

type RegistryExclusion struct {
	Registry string
	Reason   string

	SnoozeBy string
	// Unix returns t as a Unix time, the number of seconds elapsed
	// since January 1, 1970 UTC.
	SnoozeUntil int64
}

type ScanScope struct {
	NamespaceExclude string
	NamespaceInclude string

	RegistryExclusion []*RegistryExclusion

	nsInclude sets.String
	nsExclude sets.String
}

func (f *ScanScope) Init() {
	f.nsInclude, f.nsExclude = util.GetNamespaceSets(f.NamespaceInclude, f.NamespaceExclude)
}

func (f *ScanScope) IsNamespaceIncluded(ns string) bool {
	return util.IsNamespaceIncluded(ns, f.nsInclude, f.nsExclude)
}

func (f *ScanScope) ShouldScanImage(image string) (bool, string) {
	for _, e := range f.RegistryExclusion {
		if strings.HasPrefix(strings.ToLower(image), strings.ToLower(e.Registry)) {
			if e.SnoozeUntil == 0 {
				return false, fmt.Sprintf("Image '%v' excluded from vulnerability reporting", image)
			}

			if e.SnoozeUntil > 0 && time.Now().Unix() < e.SnoozeUntil {
				return false, fmt.Sprintf("Image '%v' excluded from vulnerability reporting. Exclusion active until %v", image, time.Unix(e.SnoozeUntil, 0).Format(time.RFC3339))
			}
		}
	}

	return true, ""
}

type ScanRateLimit struct {
	ApiQPS   float32 `json:"apiQPS,omitempty" envconfig:"ISKAN_API_QPS"    default:"50.0"  doc:"indicates the maximum QPS to the vuln providers"`
	ApiBurst int32   `json:"apiBurst,omitempty" envconfig:"ISKAN_API_BURST"  default:"100"   doc:"Maximum burst for throttle"`
}

type Policy struct {
	ScanScope *ScanScope

	ReportFilter *ReportFilter

	RateLimit ScanRateLimit
}

func (p *Policy) Init() {
	if p.ScanScope != nil {
		p.ScanScope.Init()
	}
}

func NewDefaultPolicy() *Policy {
	return &Policy{
		ScanScope: &ScanScope{
			NamespaceExclude:  "kube-system",
			NamespaceInclude:  "",
			RegistryExclusion: []*RegistryExclusion{},
		},

		ReportFilter: &ReportFilter{
			Severities:              "",
			CvssGreaterThan:         0,
			FixableOnly:             false,
			VulnerabilityExclusions: []*VulnerabilityExclusion{},
		},

		RateLimit: ScanRateLimit{
			ApiQPS:   30,
			ApiBurst: 100,
		},
	}
}
