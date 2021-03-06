package gcr

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"net/url"
	"strings"

	containeranalysis "cloud.google.com/go/containeranalysis/apiv1"
	grafeas "cloud.google.com/go/grafeas/apiv1"
	"github.com/alcideio/iskan/pkg/util"
	types "github.com/alcideio/iskan/pkg/vulnprovider/api"
	"github.com/docker/distribution/reference"
	"google.golang.org/api/iterator"
	grafeaspb "google.golang.org/genproto/googleapis/grafeas/v1"
)

type imageVulnerabilitiesFinder struct {
	client *containeranalysis.Client
}

func NewImageVulnerabilitiesFinder(cred *types.VulnProviderAPICreds) (types.ImageVulnerabilitiesFinder, error) {
	var client *containeranalysis.Client
	var err error
	ctx := context.Background()

	if cred != nil && cred.GCR != "" {
		//Load the registry specific creds
		client, err = containeranalysis.NewClient(ctx, option.WithCredentialsJSON([]byte(cred.GCR)))
		if err != nil {
			return nil, fmt.Errorf("NewClient: %v", err)
		}
	} else {
		//Auto detect GCP credentials from envs
		client, err = containeranalysis.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("NewClient: %v", err)
		}
	}

	return &imageVulnerabilitiesFinder{
		client: client,
	}, nil

}

//Which Registry Platform it supports
func (i *imageVulnerabilitiesFinder) Type() string {
	return "gcr"
}

func (i *imageVulnerabilitiesFinder) ListOccurrences(ctx context.Context, containerImage string) (*types.ImageScanResult, error) {
	repo, _, digest, err := util.ParseImageName(containerImage)
	if err != nil {
		return nil, err
	}

	repoUrl, err := url.Parse(repo)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(repoUrl.Path, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("failed to extract project from '%v'", repoUrl.Path)
	}
	ctx = context.WithValue(ctx, "project_id", parts[1])

	findings, err := findVulnerabilityOccurrencesForImage(ctx, i.client.GetGrafeasClient(), "https://"+repo+"@"+digest)
	if err != nil {
		return &types.ImageScanResult{
			Findings:    findings,
			CompletedOK: false,
			Reason:      err.Error(),
		}, nil
	}

	return &types.ImageScanResult{
		Findings:    findings,
		CompletedOK: true,
		Reason:      "",
	}, nil
}

// findVulnerabilityOccurrencesForImage retrieves all vulnerability Occurrences associated with a resource.
func findVulnerabilityOccurrencesForImage(ctx context.Context, client *grafeas.Client, resourceURL string) ([]*grafeaspb.Occurrence, error) {
	proj := ctx.Value("project_id")
	req := &grafeaspb.ListOccurrencesRequest{
		Parent:   fmt.Sprintf("projects/%s", proj),
		Filter:   fmt.Sprintf("resourceUrl=%q kind=%q", resourceURL, grafeaspb.NoteKind_VULNERABILITY.String()),
		PageSize: 100,
	}

	//klog.V(5).Infof("[%v][%v]", pretty.Sprint(req), proj)

	var occurrenceList []*grafeaspb.Occurrence
	it := client.ListOccurrences(ctx, req)
	for {
		occ, err := it.Next()
		//klog.V(5).Infof("[%+v][%v]", occ, err)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("occurrence iteration error: %v", err)
		}
		occ.NoteName = strings.TrimPrefix(occ.NoteName, "projects/goog-vulnz/notes/")
		occurrenceList = append(occurrenceList, occ)
	}

	//klog.V(8).Infof("%v", pretty.Sprint(occurrenceList))

	return occurrenceList, nil
}

func resourceURL(reference reference.Reference) string {
	return "resourceUrl=\"https://" + reference.String() + "\""
}

func kindFilterStr(reference reference.Reference, kind grafeaspb.NoteKind) string {
	return resourceURL(reference) + " AND kind=\"" + kind.String() + "\""
}
