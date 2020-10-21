package util

import (
	"fmt"
	"github.com/fatih/color"
)

var iskanPrefix = color.New(color.FgHiBlue).SprintFunc()
var lineMsg = color.New(color.FgHiWhite).SprintFunc()
var TitleSprint = color.New(color.FgHiWhite).SprintFunc()

func ConsolePrinter(msg string) {
	fmt.Println(iskanPrefix("[Alcide iSkan]"), lineMsg(msg))
}
