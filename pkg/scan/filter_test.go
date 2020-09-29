package scan

import (
	"github.com/alcideio/iskan/api"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"testing"
)

var finding = grafeas.Occurrence{
	Name:        "",
	ResourceUri: "",
	NoteName:    "",
	Kind:        grafeas.NoteKind_VULNERABILITY,
	Remediation: "",
	CreateTime:  nil,
	UpdateTime:  nil,
	Details: &grafeas.Occurrence_Vulnerability{
		Vulnerability: &grafeas.VulnerabilityOccurrence{
			Type:              "",
			Severity:          grafeas.Severity_HIGH,
			CvssScore:         0,
			PackageIssue:      nil,
			ShortDescription:  "",
			LongDescription:   "",
			RelatedUrls:       nil,
			EffectiveSeverity: 0,
			FixAvailable:      false,
		},
	},
}

func Test_ResultFilter(t *testing.T) {
	type test struct {
		policy          *api.Policy
		finding         *grafeas.Occurrence
		includeExpected bool
	}

	tests := []test{
		{
			policy: func() *api.Policy {
				p := api.NewDefaultPolicy()
				p.ReportFilter.Severities = "HIGH,CRITICAL"
				return p
			}(),
			finding:         &finding,
			includeExpected: true,
		},
		{
			policy: func() *api.Policy {
				p := api.NewDefaultPolicy()
				p.ReportFilter.Severities = "MEDIUM"
				return p
			}(),
			finding:         &finding,
			includeExpected: false,
		},
	}

	for _, tst := range tests {
		include := RuntimeResultFilter.IncludeResult(tst.policy, tst.finding)
		if tst.includeExpected != include {
			t.Errorf("Expected behvaior failed ['%v' != '%v]'\n%v", tst.includeExpected, include, pretty.Sprint(tst.policy.ReportFilter))
		}
	}
}
