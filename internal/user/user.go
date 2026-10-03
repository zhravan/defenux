package user

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

type User struct {
	Name            string
	UID             int
	GID             int
	Home            string
	Shell           string
	Class           string
	Admin           string
	Lock            string
	PasswordChanged string
	PasswordExpires string
	AccountExpires  string
	LastLogin       string
}

func List() ([]User, error) {
	passwd, err := xexec.Run("getent", "passwd")
	if err != nil {
		return nil, fmt.Errorf("getent passwd: %w: %s", err, strings.TrimSpace(string(passwd)))
	}

	admin, err := adminMembers()
	if err != nil {
		return nil, err
	}

	users := parsePasswd(passwd, admin)
	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})
	return users, nil
}

func Status(name string) (User, error) {
	out, err := xexec.Run("getent", "passwd", name)
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return User{}, fmt.Errorf("user not found: %s", name)
	}

	admin, err := adminMembers()
	if err != nil {
		return User{}, err
	}

	users := parsePasswd(out, admin)
	if len(users) == 0 {
		return User{}, fmt.Errorf("user not found: %s", name)
	}

	item := users[0]
	item.Lock = passwordState(name)
	item.PasswordChanged, item.PasswordExpires, item.AccountExpires = passwordAging(name)
	item.LastLogin = lastLogin(name)
	return item, nil
}

func parsePasswd(data []byte, admin map[string]bool) []User {
	var users []User
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ":", 7)
		if len(fields) != 7 {
			continue
		}

		uid, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		gid, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		name := fields[0]
		users = append(users, User{
			Name:  name,
			UID:   uid,
			GID:   gid,
			Home:  fields[5],
			Shell: fields[6],
			Class: classify(uid),
			Admin: yesNo(admin[name]),
			Lock:  "unknown",
		})
	}

	return users
}

func classify(uid int) string {
	switch {
	case uid == 0:
		return "root"
	case uid < 1000:
		return "system"
	default:
		return "human"
	}
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func adminMembers() (map[string]bool, error) {
	out, err := xexec.Run("getent", "group")
	if err != nil {
		return nil, fmt.Errorf("getent group: %w: %s", err, strings.TrimSpace(string(out)))
	}

	admin := map[string]bool{}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), ":", 4)
		if len(fields) != 4 || (fields[0] != "sudo" && fields[0] != "wheel") {
			continue
		}

		for _, member := range strings.Split(fields[3], ",") {
			if member != "" {
				admin[member] = true
			}
		}
	}

	return admin, nil
}

func passwordState(name string) string {
	out, err := xexec.Run("passwd", "-S", name)
	if err != nil {
		return "unknown"
	}
	return parsePasswordState(out)
}

func parsePasswordState(data []byte) string {
	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return "unknown"
	}

	switch fields[1] {
	case "L", "LK":
		return "locked"
	case "NP":
		return "no-password"
	case "P":
		return "active"
	default:
		return fields[1]
	}
}

func passwordAging(name string) (string, string, string) {
	out, err := xexec.Run("chage", "-l", name)
	if err != nil {
		return "unknown", "unknown", "unknown"
	}
	return parsePasswordAging(out)
}

func parsePasswordAging(data []byte) (string, string, string) {
	var changed, expires, account string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), ":")
		if !ok {
			continue
		}

		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "Last password change":
			changed = value
		case "Password expires":
			expires = value
		case "Account expires":
			account = value
		}
	}

	return agingValue(changed), agingValue(expires), agingValue(account)
}

func agingValue(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func lastLogin(name string) string {
	out, err := xexec.Run("lastlog", "-u", name)
	if err != nil {
		return "unknown"
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return "unknown"
	}

	value := strings.TrimSpace(lines[len(lines)-1])
	if value == "" || strings.Contains(value, "**Never logged in**") {
		return "never"
	}
	return value
}
