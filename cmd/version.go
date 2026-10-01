package cmd

import "github.com/spf13/cobra"

var version = "0.0.3"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version)
		},
	})
}
