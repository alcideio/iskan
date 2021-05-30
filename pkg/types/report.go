package types

import (
	"time"

	"github.com/alcideio/iskan/pkg/vulnprovider/api"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type PodSpecSummary struct {
	Name      string
	Namespace string

	Spec *v1.PodSpec

	Severity api.SeveritySummary
	Fixable  api.SeveritySummary

	ScanFailures uint32
}

type ClusterScanReportSummary struct {
	// ns/podName --> Severity Vector
	PodSummary map[string]PodSpecSummary

	ClusterSeverity api.SeveritySummary

	NamespaceSeverity api.SeveritySummaryMap

	// ns/podName --> Severity Vector
	PodSeverity        api.SeveritySummaryMap
	PodFixableSeverity api.SeveritySummaryMap

	FailedOrSkippedPods   []string
	FailedOrSkippedImages []string

	AnalyzedPodCount uint32
	ExcludedPodCount uint32
}

type ScanTaskResult struct {
	Findings map[string]*api.ImageScanResult

	ScannedPods []*v1.Pod
	SkippedPods []*v1.Pod
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
	Findings map[string]*api.ImageScanResult

	// High level stats about this report
	Summary ClusterScanReportSummary
}

func NewClusterScanReportSummary() *ClusterScanReportSummary {
	return &ClusterScanReportSummary{
		PodSummary:            map[string]PodSpecSummary{},
		ClusterSeverity:       api.NewSeveritySummary(),
		NamespaceSeverity:     map[string]api.SeveritySummary{},
		PodSeverity:           map[string]api.SeveritySummary{},
		PodFixableSeverity:    map[string]api.SeveritySummary{},
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
		Findings:          map[string]*api.ImageScanResult{},
		Summary:           ClusterScanReportSummary{},
	}
}
