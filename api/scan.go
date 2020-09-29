package api

import (
	"context"
)

type ImageVulnerabilitiesFinder interface {
	//Which Registry Platform it supports
	Type() string

	ListOccurrences(ctx context.Context, containerImage string) (*ImageScanResult, error)
}
