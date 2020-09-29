package api

import (
	"bytes"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type SeveritySummary map[grafeas.Severity]uint32

func NewSeveritySummary() SeveritySummary {
	return map[grafeas.Severity]uint32{
		grafeas.Severity_SEVERITY_UNSPECIFIED: 0,
		grafeas.Severity_MINIMAL:              0,
		grafeas.Severity_LOW:                  0,
		grafeas.Severity_MEDIUM:               0,
		grafeas.Severity_HIGH:                 0,
		grafeas.Severity_CRITICAL:             0,
	}
}

func (s SeveritySummary) Add(b SeveritySummary) {
	for severity, count := range b {
		s[severity] = s[severity] + count
	}
}

func (s SeveritySummary) String() string {
	return fmt.Sprint(
		grafeas.Severity_CRITICAL.String(), ": ", s[grafeas.Severity_CRITICAL], ",",
		grafeas.Severity_HIGH.String(), ": ", s[grafeas.Severity_HIGH], ",",
		grafeas.Severity_MEDIUM.String(), ": ", s[grafeas.Severity_MEDIUM], ",",
		grafeas.Severity_LOW.String(), ": ", s[grafeas.Severity_LOW], ",",
		grafeas.Severity_MINIMAL.String(), ": ", s[grafeas.Severity_MINIMAL], ",",
	)
}

func (s SeveritySummary) Max() (string, uint32) {
	if s[grafeas.Severity_CRITICAL] > 0 {
		return grafeas.Severity_CRITICAL.String(), s[grafeas.Severity_CRITICAL]
	}
	if s[grafeas.Severity_HIGH] > 0 {
		return grafeas.Severity_HIGH.String(), s[grafeas.Severity_HIGH]
	}
	if s[grafeas.Severity_MEDIUM] > 0 {
		return grafeas.Severity_MEDIUM.String(), s[grafeas.Severity_MEDIUM]
	}
	if s[grafeas.Severity_LOW] > 0 {
		return grafeas.Severity_LOW.String(), s[grafeas.Severity_LOW]
	}
	if s[grafeas.Severity_MINIMAL] > 0 {
		return grafeas.Severity_MINIMAL.String(), s[grafeas.Severity_MINIMAL]
	}
	return "", 0
}

func (s SeveritySummary) Table() string {
	w := bytes.NewBufferString("")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"SEVERITY", "COUNT"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	headerRows := [][]string{
		{color.HiRedString(grafeas.Severity_CRITICAL.String()), color.HiRedString(fmt.Sprint(s[grafeas.Severity_CRITICAL]))},
		{color.RedString(grafeas.Severity_HIGH.String()), color.RedString(fmt.Sprint(s[grafeas.Severity_HIGH]))},
		{color.HiYellowString(grafeas.Severity_MEDIUM.String()), color.HiYellowString(fmt.Sprint(s[grafeas.Severity_MEDIUM]))},
		{color.YellowString(grafeas.Severity_LOW.String()), color.YellowString(fmt.Sprint(s[grafeas.Severity_LOW]))},
		{color.BlueString(grafeas.Severity_MINIMAL.String()), color.BlueString(fmt.Sprint(s[grafeas.Severity_MINIMAL]))},
	}
	table.AppendBulk(headerRows)
	table.Render()

	return w.String()
}

type SeveritySummaryMap map[string]SeveritySummary

func (sm SeveritySummaryMap) Table() string {
	w := bytes.NewBufferString("")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{
		"",
		grafeas.Severity_CRITICAL.String(),
		grafeas.Severity_HIGH.String(),
		grafeas.Severity_MEDIUM.String(),
		grafeas.Severity_LOW.String(),
		grafeas.Severity_MINIMAL.String(),
	})

	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	rows := [][]string{}
	for k, s := range sm {
		row := []string{
			color.HiWhiteString(k),
			color.HiRedString(fmt.Sprint(s[grafeas.Severity_CRITICAL])),
			color.RedString(fmt.Sprint(s[grafeas.Severity_HIGH])),
			color.HiYellowString(fmt.Sprint(s[grafeas.Severity_MEDIUM])),
			color.YellowString(fmt.Sprint(s[grafeas.Severity_LOW])),
			color.BlueString(fmt.Sprint(s[grafeas.Severity_MINIMAL])),
		}

		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) < 0
	})

	table.AppendBulk(rows)
	table.Render()

	return w.String()
}

type PodSpecSummary struct {
	Name      string
	Namespace string

	Spec *v1.PodSpec

	Severity SeveritySummary

	ScanFailures uint32
}

type ClusterScanReportSummary struct {
	// ns/podName --> Severity Vector
	PodSummary map[string]PodSpecSummary

	ClusterSeverity SeveritySummary

	NamespaceSeverity SeveritySummaryMap

	// ns/podName --> Severity Vector
	PodSeverity SeveritySummaryMap

	FailedOrSkippedPods   []string
	FailedOrSkippedImages []string

	AnalyzedPodCount uint32
	ExcludedPodCount uint32
}

type ImageScanResult struct {
	Image string

	CompletedOK bool
	Reason      string

	//If completed ok - this value should be populated with findings (if there are any)
	Findings []*grafeas.Occurrence

	//Stats
	Summary      SeveritySummary
	ExcludeCount uint32
}

type ClusterScanReport struct {
	//A Cluster UUID Identifier ... (namespace uid of kube-system ;P)
	ClusterId string
	// CreationTimestamp is a timestamp representing the time when this report was
	// created. It is represented in RFC3339 form and is in UTC.
	CreationTimeStamp string `json:"CreationTimeStamp,omitempty"`

	//Unique Report UUID
	ReportUUID string

	//The Policy with the report was generated with
	Policy Policy

	// Report Findings
	Findings map[string]*ImageScanResult

	// High level stats about this report
	Summary ClusterScanReportSummary
}

func NewClusterScanReportSummary() *ClusterScanReportSummary {
	return &ClusterScanReportSummary{
		PodSummary:            map[string]PodSpecSummary{},
		ClusterSeverity:       NewSeveritySummary(),
		NamespaceSeverity:     map[string]SeveritySummary{},
		PodSeverity:           map[string]SeveritySummary{},
		AnalyzedPodCount:      0,
		ExcludedPodCount:      0,
		FailedOrSkippedPods:   []string{},
		FailedOrSkippedImages: []string{},
	}
}

func NewClusterScanReport() *ClusterScanReport {
	return &ClusterScanReport{
		ClusterId:         "",
		CreationTimeStamp: time.Now().Format(time.RFC3339),
		ReportUUID:        rand.String(10),
		Policy:            Policy{},
		Findings:          map[string]*ImageScanResult{},
		Summary:           ClusterScanReportSummary{},
	}
}
