package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zhravan/defenux/internal/service"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Inspect system services",
}

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List system services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := service.List()
		if err != nil {
			return err
		}

		for _, item := range services {
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"%s\t%s\t%s\t%s\t%s\n",
				item.Name,
				item.State,
				item.SubState,
				item.Enabled,
				item.Description,
			)
		}
		return nil
	},
}

var serviceStatusCmd = &cobra.Command{
	Use:   "status <name>",
	Short: "Show service status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := service.Status(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "name: %s\n", item.Name)
		fmt.Fprintf(out, "description: %s\n", item.Description)
		fmt.Fprintf(out, "state: %s\n", item.State)
		fmt.Fprintf(out, "substate: %s\n", item.SubState)
		fmt.Fprintf(out, "enabled: %s\n", item.Enabled)
		fmt.Fprintf(out, "pid: %d\n", item.MainPID)
		fmt.Fprintf(out, "type: %s\n", item.Type)
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(serviceListCmd, serviceStatusCmd)
	rootCmd.AddCommand(serviceCmd)
}
