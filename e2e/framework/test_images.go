package framework

type ImagePropertyChecker func() bool

func ImagePropertyTrue() bool  { return true }
func ImagePropertyFalse() bool { return false }

type TestImageProperties struct {
	Description                string
	HasVulnerabilities         ImagePropertyChecker
	HasCriticalVulnerabilities ImagePropertyChecker
	HasFixableVulnerabilities  ImagePropertyChecker
}

type TestImageInfo struct {
	TestImageProperties
	PullSecret string
}

var CleanImageBuiltFromScratch = TestImageProperties{
	Description:                "Clean (Vuln. free) Image Built from Scratch",
	HasVulnerabilities:         ImagePropertyFalse,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var DistrolessImage = TestImageProperties{
	Description:                "Distroless Image",
	HasVulnerabilities:         ImagePropertyTrue,
	HasCriticalVulnerabilities: ImagePropertyFalse,
	HasFixableVulnerabilities:  ImagePropertyFalse,
}

var TestImages = map[string]TestImageInfo{
	"gcr.io/dcvisor-162009/iskan/e2e/zerovuln_scratch:latest": {
		TestImageProperties: CleanImageBuiltFromScratch,
		PullSecret:          "gcr",
	},
	"gcr.io/dcvisor-162009/iskan/e2e/zerovuln_distroless:latest": {
		TestImageProperties: DistrolessImage,
		PullSecret:          "gcr",
	},
}
