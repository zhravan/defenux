package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	xssh "github.com/zhravan/defenux/internal/ssh"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Inspect SSH server configuration",
}

var sshStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show SSH server status",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, err := xssh.Status()
		if err != nil {
			return err
		}

		fmt.Fprintf(
			cmd.OutOrStdout(),
			"service: %s\nstate: %s\nconfig: %s\nports: %s\n",
			status.Service,
			status.State,
			status.ConfigPath,
			ports(status.Ports),
		)
		return nil
	},
}

var sshConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show effective SSH configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := xssh.Config()
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "path: %s\n", config.Path)
		fmt.Fprintf(out, "port: %s\n", ports(config.Ports))
		fmt.Fprintf(out, "permitrootlogin: %s\n", value(config.PermitRootLogin))
		fmt.Fprintf(out, "passwordauthentication: %s\n", value(config.PasswordAuthentication))
		fmt.Fprintf(out, "pubkeyauthentication: %s\n", value(config.PubkeyAuthentication))
		fmt.Fprintf(out, "permitemptypasswords: %s\n", value(config.PermitEmptyPasswords))
		fmt.Fprintf(out, "x11forwarding: %s\n", value(config.X11Forwarding))
		fmt.Fprintf(out, "maxauthtries: %s\n", value(config.MaxAuthTries))
		fmt.Fprintf(out, "allowusers: %s\n", value(strings.Join(config.AllowUsers, " ")))
		fmt.Fprintf(out, "allowgroups: %s\n", value(strings.Join(config.AllowGroups, " ")))
		return nil
	},
}

func value(s string) string {
	if s == "" {
		return "unset"
	}
	return s
}

func ports(values []int) string {
	if len(values) == 0 {
		return "none"
	}

	var out []string
	for _, port := range values {
		out = append(out, fmt.Sprint(port))
	}
	return strings.Join(out, ",")
}

func init() {
	sshCmd.AddCommand(sshStatusCmd, sshConfigCmd)
	rootCmd.AddCommand(sshCmd)
}
