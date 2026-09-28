package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zhravan/defenux/internal/port"
)

var portCmd = &cobra.Command{Use: "port", Short: "Inspect listening ports"}

var portListCmd = &cobra.Command{
	Use: "list", Short: "List listening ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		listeners, err := port.List()
		if err != nil {
			return err
		}
		for _, l := range listeners {
			printListener(cmd, l)
		}
		return nil
	},
}

var portStatusCmd = &cobra.Command{
	Use:   "status <port>",
	Short: "Show listeners for a port",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		portNumber, err := strconv.Atoi(args[0])
		if err != nil || portNumber < 1 || portNumber > 65535 {
			return fmt.Errorf("invalid port: %s", args[0])
		}
		listeners, err := port.Status(portNumber)
		if err != nil {
			return err
		}
		if len(listeners) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "port %d: not listening\n", portNumber)
			return nil
		}
		for _, l := range listeners {
			printListener(cmd, l)
		}
		return nil
	},
}

func printListener(cmd *cobra.Command, l port.Listener) {
	fmt.Fprintf(cmd.OutOrStdout(), "%s %s:%d", l.Protocol, l.Address, l.Port)
	if l.PID > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), " pid=%d", l.PID)
	}
	if l.Process != "" {
		fmt.Fprintf(cmd.OutOrStdout(), " %s", l.Process)
	}
	fmt.Fprintln(cmd.OutOrStdout())
}

func init() {
	portCmd.AddCommand(portListCmd, portStatusCmd)
	rootCmd.AddCommand(portCmd)
}
