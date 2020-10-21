package cmd

import (
	"fmt"

	"github.com/alcideio/iskan/pkg/version"
	"github.com/spf13/cobra"
)

var (
	Version = ""
	Commit  = ""
)

func NewCommandVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print iskan version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: " + version.Version + "\nCommit: " + version.Commit)
		},
	}
}
