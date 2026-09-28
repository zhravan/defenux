package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhravan/defenux/internal/network"
)

var networkCmd = &cobra.Command{Use: "network", Short: "Inspect network state"}

var networkInterfacesCmd = &cobra.Command{
	Use: "interfaces", Short: "List network interfaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := network.Interfaces()
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s", item.Name, item.State)
			if item.MAC != "" {
				fmt.Fprintf(cmd.OutOrStdout(), " %s", item.MAC)
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}
		return nil
	},
}

var networkRoutesCmd = &cobra.Command{
	Use: "routes", Short: "List network routes",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := network.Routes()
		if err != nil {
			return err
		}
		for _, item := range items {
			fmt.Fprintf(cmd.OutOrStdout(), "%s", item.Destination)
			if item.Gateway != "" {
				fmt.Fprintf(cmd.OutOrStdout(), " via %s", item.Gateway)
			}
			if item.Device != "" {
				fmt.Fprintf(cmd.OutOrStdout(), " dev %s", item.Device)
			}
			fmt.Fprintln(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	networkCmd.AddCommand(networkInterfacesCmd, networkRoutesCmd)
	rootCmd.AddCommand(networkCmd)
}
