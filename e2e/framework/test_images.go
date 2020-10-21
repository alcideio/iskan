package framework

import "k8s.io/apimachinery/pkg/util/sets"

type ImagePropertyChecker func() bool

func ImagePropertyTrue() bool  { return true }
func ImagePropertyFalse() bool { return false }

type TestImageProperties struct {
	Description                string
	HasScanFailures            ImagePropertyChecker
	HasVulnerabilities         ImagePropertyChecker
	HasCriticalVulnerabilities ImagePropertyChecker
	HasFixableVulnerabilities  ImagePropertyChecker
}

type TestImageInfo struct {
	Image string
	TestImageProperties
	PullSecret string
}

var CleanImageBuiltFromScratch = TestImageProperties{
	Description:                "Clean (Vuln. free) Image Built from Scratch",
	HasScanFailures:            ImagePropertyFalse,
	HasVulnerabilities:         ImagePropertyFalse,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var DistrolessImage = TestImageProperties{
	Description:                "Distroless Image",
	HasScanFailures:            ImagePropertyFalse,
	HasVulnerabilities:         ImagePropertyTrue,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var AlpineImage = TestImageProperties{
	Description:                "Alpine Image",
	HasScanFailures:            ImagePropertyFalse,
	HasVulnerabilities:         ImagePropertyTrue,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var CleanImage = TestImageProperties{
	Description:                "Clean Image",
	HasScanFailures:            ImagePropertyFalse,
	HasVulnerabilities:         ImagePropertyFalse,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var FailedScanImage = TestImageProperties{
	Description:                "Image which was failed to scan",
	HasScanFailures:            ImagePropertyTrue,
	HasVulnerabilities:         ImagePropertyFalse,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var ManyVulnsImage = TestImageProperties{
	Description:                "Image with many vulnerabilities",
	HasScanFailures:            ImagePropertyFalse,
	HasVulnerabilities:         ImagePropertyTrue,
	HasCriticalVulnerabilities: ImagePropertyTrue,
	HasFixableVulnerabilities:  ImagePropertyTrue,
}

type TestImageIterator struct {
	Filter func(img *TestImageInfo) bool

	current int
}

func NewTestImageIterator(filter func(img *TestImageInfo) bool) *TestImageIterator {
	return &TestImageIterator{
		Filter:  filter,
		current: 0,
	}
}

func (i *TestImageIterator) Next() (*TestImageInfo, bool) {
	defer func() { i.current++ }()

	for {
		if i.current >= len(TestImages) {
			return nil, false
		}

		if i.Filter == nil || i.Filter(&TestImages[i.current]) {
			return &TestImages[i.current], true
		}
		i.current++
	}
}

func FilterByKind(kind string) func(info *TestImageInfo) bool {
	return func(info *TestImageInfo) bool {
		if info.PullSecret != kind {
			return false
		}

		return true
	}
}

func FilterByKinds(kinds sets.String) func(info *TestImageInfo) bool {
	return func(info *TestImageInfo) bool {
		if !kinds.Has(info.PullSecret) {
			return false
		}

		return true
	}
}

func FilterByPrivateRegistries() func(info *TestImageInfo) bool {
	var images = sets.NewString("ecr", "gcr" /*"acr"*/)
	return func(info *TestImageInfo) bool {
		if !images.Has(info.PullSecret) {
			return false
		}

		return true
	}
}

func FilterByPublicRegistries() func(info *TestImageInfo) bool {
	var images = sets.NewString("")
	return func(info *TestImageInfo) bool {
		if !images.Has(info.PullSecret) {
			return false
		}

		return true
	}
}

var TestImages = []TestImageInfo{
	//GCR
	{
		Image:               "gcr.io/dcvisor-162009/iskan/e2e/zerovuln_scratch:latest",
		TestImageProperties: CleanImageBuiltFromScratch,
		PullSecret:          "gcr",
	},
	{
		Image:               "gcr.io/dcvisor-162009/iskan/e2e/zerovuln_distroless:latest",
		TestImageProperties: DistrolessImage,
		PullSecret:          "gcr",
	},

	//ECR
	{
		Image:               "893825821121.dkr.ecr.us-west-2.amazonaws.com/iskan/zerovuln_distroless:latest",
		TestImageProperties: CleanImage,
		PullSecret:          "ecr",
	},
	{
		Image:               "893825821121.dkr.ecr.us-west-2.amazonaws.com/iskan/zerovuln_scratch:latest",
		TestImageProperties: FailedScanImage, //ECR doesn't like images from scratch
		PullSecret:          "ecr",
	},

	//ACR
	{
		Image:               "alcide.azurecr.io/iskan/zerovuln_distroless:latest",
		TestImageProperties: CleanImage,
		PullSecret:          "acr",
	},
	//{
	//	Image:               "alcide.azurecr.io/iskan/manyvuln:latest",
	//	TestImageProperties: ManyVulnsImage,
	//	PullSecret:          "acr",
	//},
	//{
	//	Image:               "alcide.azurecr.io/iskan/zerovuln_scratch:latest",
	//	TestImageProperties: FailedScanImage,
	//	PullSecret:          "acr",
	//},
	//{
	//	Image:               "alcide.azurecr.io/iskan/vuln_alpine:latest",
	//	TestImageProperties: AlpineImage,
	//	PullSecret:          "acr",
	//},

	//Inline Scan Engine
	{
		Image:               "iskan/zerovuln_distroless:latest",
		TestImageProperties: DistrolessImage,
		PullSecret:          "",
	},
	{
		Image:               "iskan/zerovuln_scratch:latest",
		TestImageProperties: CleanImage,
		PullSecret:          "",
	},
	{
		Image:               "iskan/vuln_alpine:latest",
		TestImageProperties: AlpineImage, //ECR doesn't like images from scratch
		PullSecret:          "",
	},
}
