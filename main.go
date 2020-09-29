package main

import (
	"bytes"
	goflag "flag"
	"fmt"
	"os"

	"github.com/alcideio/iskan/cmd"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func iSkanGenCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "iskan",
		Short: "iskan",
		Long:  `iskan`,
	}

	var genBashCompletionCmd = &cobra.Command{
		Use:   "bash-completion",
		Short: "Generate bash completion. source < (advisor bash-completion)",
		Long:  "Generate bash completion. source < (advisor bash-completion)",
		Run: func(cmd *cobra.Command, args []string) {
			out := new(bytes.Buffer)
			_ = rootCmd.GenBashCompletion(out)
			println(out.String())
		},
	}

	cmds := []*cobra.Command{
		cmd.NewCommandVersion(),
		cmd.NewCommandScanImage(),
		cmd.NewCommandScanCluster(),
		genBashCompletionCmd,
	}

	flags := rootCmd.PersistentFlags()

	klog.InitFlags(nil)
	flags.AddGoFlagSet(goflag.CommandLine)

	// Hide all klog flags except for -v
	goflag.CommandLine.VisitAll(func(f *goflag.Flag) {
		if f.Name != "v" {
			flags.Lookup(f.Name).Hidden = true
		}
	})

	rootCmd.AddCommand(cmds...)

	return rootCmd
}

func main() {
	rootCmd := iSkanGenCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
