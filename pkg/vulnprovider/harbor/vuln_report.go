package harbor

//Until these data models become part of the official API
//See https://github.com/goharbor/harbor/tree/master/src/pkg/scan/report

const (
	// None - only used to mark the overall severity of the scanned artifacts,
	// means no vulnerabilities attached with the artifacts,
	// (might be bypassed by the CVE whitelist).
	None Severity = "None"
	// Unknown - either a security problem that has not been assigned to a priority yet or
	// a priority that the scanner did not recognize.
	Unknown Severity = "Unknown"
	// Negligible - technically a security problem, but is only theoretical in nature, requires
	// a very special situation, has almost no install base, or does no real damage.
	Negligible Severity = "Negligible"
	// Low - a security problem, but is hard to exploit due to environment, requires a
	// user-assisted attack, a small install base, or does very little damage.
	Low Severity = "Low"
	// Medium - a real security problem, and is exploitable for many people. Includes network
	// daemon denial of service attacks, cross-site scripting, and gaining user privileges.
	Medium Severity = "Medium"
	// High - a real problem, exploitable for many people in a default installation. Includes
	// serious remote denial of service, local root privilege escalations, or data loss.
	High Severity = "High"
	// Critical - a world-burning problem, exploitable for nearly all people in a default installation.
	// Includes remote root privilege escalations, or massive data loss.
	Critical Severity = "Critical"
)

// Severity is a standard scale for measuring the severity of a vulnerability.
type Severity string

// Scanner represents metadata of a Scanner Adapter which allow Harbor to lookup a scanner capable of
// scanning a given Artifact stored in its registry and making sure that it can interpret a
// returned result.
type Scanner struct {
	// The name of the scanner.
	Name string `json:"name"`
	// The name of the scanner's provider.
	Vendor string `json:"vendor"`
	// The version of the scanner.
	Version string `json:"version"`
}

// Report model for vulnerability scan
type Report struct {
	// Time of generating this report
	GeneratedAt string `json:"generated_at"`
	// Scanner of generating this report
	Scanner *Scanner `json:"scanner"`
	// A standard scale for measuring the severity of a vulnerability.
	Severity string `json:"severity"`
	// Vulnerability list
	Vulnerabilities []*VulnerabilityItem `json:"vulnerabilities"`
}

// VulnerabilityItem represents one found vulnerability
type VulnerabilityItem struct {
	// The unique identifier of the vulnerability.
	// e.g: CVE-2017-8283
	ID string `json:"id"`
	// An operating system or software dependency package containing the vulnerability.
	// e.g: dpkg
	Package string `json:"package"`
	// The version of the package containing the vulnerability.
	// e.g: 1.17.27
	Version string `json:"version"`
	// The version of the package containing the fix if available.
	// e.g: 1.18.0
	FixVersion string `json:"fix_version"`
	// A standard scale for measuring the severity of a vulnerability.
	Severity string `json:"severity"`
	// example: dpkg-source in dpkg 1.3.0 through 1.18.23 is able to use a non-GNU patch program
	// and does not offer a protection mechanism for blank-indented diff hunks, which allows remote
	// attackers to conduct directory traversal attacks via a crafted Debian source package, as
	// demonstrated by using of dpkg-source on NetBSD.
	Description string `json:"description"`
	// The list of link to the upstream database with the full description of the vulnerability.
	// Format: URI
	// e.g: List [ "https://security-tracker.debian.org/tracker/CVE-2017-8283" ]
	Links []string `json:"links"`
	// The artifact digest which the vulnerability belonged
	// e.g: sha256@ee1d00c5250b5a886b09be2d5f9506add35dfb557f1ef37a7e4b8f0138f32956
	ArtifactDigest string `json:"artifact_digest"`
}
