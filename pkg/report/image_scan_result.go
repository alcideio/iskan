package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/types"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
)

func reportImageScanResultAsTable(res *types.ImageScanResult, w io.WriteCloser) error {

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"SCAN INFO", ""})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	headerRows := [][]string{
		{util.TitleSprint("Image"), res.Image},
		{util.TitleSprint("Completed"), fmt.Sprint(res.CompletedOK)},
		{util.TitleSprint("Reason"), res.Reason},
		{util.TitleSprint("Findings"), fmt.Sprint(len(res.Findings))},
		{util.TitleSprint("Excluded"), fmt.Sprint(res.ExcludeCount)},
		{util.TitleSprint("Summary"), res.Summary.String()},
	}
	table.AppendBulk(headerRows)
	table.Render()

	if len(res.Findings) == 0 {
		return nil
	}

	rows := [][]string{}
	for _, o := range res.Findings {
		vul := o.GetVulnerability()
		if vul == nil {
			continue
		}

		urls := []string{}
		for _, url := range vul.RelatedUrls {
			urls = append(urls, url.Url)
		}

		for _, pkg := range vul.PackageIssue {

			row := []string{
				o.NoteName,
				pkg.AffectedPackage,
				pkg.AffectedVersion.FullName,
				vul.Severity.String(),
				vul.EffectiveSeverity.String(),
				fmt.Sprint(vul.CvssScore),
				fmt.Sprint(vul.FixAvailable),
				vul.LongDescription,
				strings.Join(urls, " "),
			}
			rows = append(rows, row)

		}

	}

	sort.Slice(rows, func(i, j int) bool {
		a := grafeas.Severity_value[rows[i][3]]
		b := grafeas.Severity_value[rows[j][3]]

		return a > b
	})

	table = tablewriter.NewWriter(w)
	table.SetHeader([]string{"CVE", "PACKAGE", "VERSION", "SEVERITY", "EFFECTIVE", "CVSS SCORE", "FIX AVAIL.", "DESCRIPTION", "MORE INFO"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetFooter([]string{util.TitleSprint("Summary"), res.Summary.String(), "", "", "", "", "", "", ""})

	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	//table.SetAutoMergeCells(true)

	table.AppendBulk(rows)
	table.Render()

	return nil
}

func ReportImageScanResult(format string, res *types.ImageScanResult, w io.WriteCloser) error {
	switch format {
	case "table":
		return reportImageScanResultAsTable(res, w)
	case "json":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "\t")
		return encoder.Encode(res)
	case "yaml":
		data, err := yaml.Marshal(res)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	default:
		return fmt.Errorf("'%v' is not supported", format)
	}
}
