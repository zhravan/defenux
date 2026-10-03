package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	xuser "github.com/zhravan/defenux/internal/user"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Inspect local user accounts",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List local user accounts",
	RunE: func(cmd *cobra.Command, args []string) error {
		users, err := xuser.List()
		if err != nil {
			return err
		}

		for _, item := range users {
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"%s\t%d\t%d\t%s\t%s\t%s\t%s\n",
				item.Name,
				item.UID,
				item.GID,
				item.Class,
				item.Admin,
				item.Home,
				item.Shell,
			)
		}
		return nil
	},
}

var userStatusCmd = &cobra.Command{
	Use:   "status <name>",
	Short: "Show user account status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := xuser.Status(args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "name: %s\n", item.Name)
		fmt.Fprintf(out, "uid: %d\n", item.UID)
		fmt.Fprintf(out, "gid: %d\n", item.GID)
		fmt.Fprintf(out, "class: %s\n", item.Class)
		fmt.Fprintf(out, "home: %s\n", item.Home)
		fmt.Fprintf(out, "shell: %s\n", item.Shell)
		fmt.Fprintf(out, "admin: %s\n", item.Admin)
		fmt.Fprintf(out, "lock: %s\n", item.Lock)
		fmt.Fprintf(out, "password-changed: %s\n", item.PasswordChanged)
		fmt.Fprintf(out, "password-expires: %s\n", item.PasswordExpires)
		fmt.Fprintf(out, "account-expires: %s\n", item.AccountExpires)
		fmt.Fprintf(out, "last-login: %s\n", item.LastLogin)
		return nil
	},
}

func init() {
	userCmd.AddCommand(userListCmd, userStatusCmd)
	rootCmd.AddCommand(userCmd)
}
