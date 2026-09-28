package cmd

import "github.com/spf13/cobra"

const version = "0.0.1"

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version)
		},
	})
}
