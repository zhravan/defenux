package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhravan/defenux/internal/port"
)

var portCmd = &cobra.Command{Use: "port", Short: "Inspect listening ports"}

var portListCmd = &cobra.Command{
	Use: "list", Short: "List listening ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		listeners, err := port.List()
		if err != nil { return err }
		for _, l := range listeners {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s:%d %s", l.Protocol, l.Address, l.Port, l.Process)
			if l.PID > 0 { fmt.Fprintf(cmd.OutOrStdout(), " (pid=%d)", l.PID) }
			fmt.Fprintln(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	portCmd.AddCommand(portListCmd)
	rootCmd.AddCommand(portCmd)
}
