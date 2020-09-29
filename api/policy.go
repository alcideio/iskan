package api

import (
	"github.com/alcideio/iskan/pkg/util"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"
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
	Registry    string
	Reason      string
	SnoozeUntil int64
	SnoozeBy    string
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

type Policy struct {
	ScanScope *ScanScope

	ReportFilter *ReportFilter
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
	}
}
