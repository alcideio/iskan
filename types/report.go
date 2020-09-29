package types

import (
	"time"

	"google.golang.org/genproto/googleapis/grafeas/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

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
