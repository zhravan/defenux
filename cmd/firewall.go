package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhravan/defenux/internal/firewall"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Inspect and control the firewall",
}

var firewallStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show firewall status",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := firewall.Status()
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "backend: %s\nstate: %s\n", status.Backend, status.State)
		return nil
	},
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List firewall rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := firewall.List()
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.OutOrStdout(), rules)
		return nil
	},
}

var firewallAllowCmd = &cobra.Command{
	Use:   "allow <port/protocol>",
	Short: "Allow inbound traffic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, protocol, err := firewall.ParseSpec(args[0])
		if err != nil {
			return err
		}
		return firewall.Allow(port, protocol)
	},
}

var firewallDenyCmd = &cobra.Command{
	Use:   "deny <port/protocol>",
	Short: "Deny inbound traffic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		port, protocol, err := firewall.ParseSpec(args[0])
		if err != nil {
			return err
		}
		return firewall.Deny(port, protocol)
	},
}

var firewallEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable the firewall",
	RunE: func(cmd *cobra.Command, args []string) error {
		return firewall.Enable()
	},
}

var firewallDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable the firewall",
	RunE: func(cmd *cobra.Command, args []string) error {
		return firewall.Disable()
	},
}

func init() {
	firewallCmd.AddCommand(
		firewallStatusCmd,
		firewallListCmd,
		firewallAllowCmd,
		firewallDenyCmd,
		firewallEnableCmd,
		firewallDisableCmd,
	)
	rootCmd.AddCommand(firewallCmd)
}
