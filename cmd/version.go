package cmd

import (
	_ "embed"
	"strings"

	"github.com/spf13/cobra"
)

var (
	//go:embed version.txt
	versionFile string
	version     = strings.TrimSpace(versionFile)
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version)
		},
	})
}
