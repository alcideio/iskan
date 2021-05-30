package api

import (
	"context"
	"google.golang.org/genproto/googleapis/grafeas/v1"
)

type ImageVulnerabilitiesFinder interface {
	//Which Registry Platform it supports
	Type() string

	ListOccurrences(ctx context.Context, containerImage string) (*ImageScanResult, error)
}

type ImageScanResult struct {
	Image string

	CompletedOK bool
	Reason      string

	//If completed ok - this value should be populated with findings (if there are any)
	Findings []*grafeas.Occurrence

	//Stats
	Summary      SeveritySummary
	Fixable      SeveritySummary
	ExcludeCount uint32
}
