package types

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/genproto/googleapis/grafeas/v1"
	"sort"
	"strings"
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

func (sm SeveritySummaryMap) Table(aux SeveritySummaryMap) string {
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
		var row []string

		if aux != nil {
			row = []string{
				color.HiWhiteString(k),
				color.HiRedString(fmt.Sprintf("%v/%v", aux[k][grafeas.Severity_CRITICAL], s[grafeas.Severity_CRITICAL])),
				color.RedString(fmt.Sprintf("%v/%v", aux[k][grafeas.Severity_HIGH], s[grafeas.Severity_HIGH])),
				color.HiYellowString(fmt.Sprintf("%v/%v", aux[k][grafeas.Severity_MEDIUM], s[grafeas.Severity_MEDIUM])),
				color.YellowString(fmt.Sprintf("%v/%v", aux[k][grafeas.Severity_LOW], s[grafeas.Severity_LOW])),
				color.BlueString(fmt.Sprintf("%v/%v", aux[k][grafeas.Severity_MINIMAL], s[grafeas.Severity_MINIMAL])),
			}
		} else {
			row = []string{
				color.HiWhiteString(k),
				color.HiRedString(fmt.Sprint(s[grafeas.Severity_CRITICAL])),
				color.RedString(fmt.Sprint(s[grafeas.Severity_HIGH])),
				color.HiYellowString(fmt.Sprint(s[grafeas.Severity_MEDIUM])),
				color.YellowString(fmt.Sprint(s[grafeas.Severity_LOW])),
				color.BlueString(fmt.Sprint(s[grafeas.Severity_MINIMAL])),
			}
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
