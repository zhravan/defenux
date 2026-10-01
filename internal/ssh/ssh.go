package ssh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

var configPaths = []string{
	"/etc/ssh/sshd_config",
}

type StatusInfo struct {
	Service    string
	State      string
	ConfigPath string
	Ports      []int
}

type ConfigInfo struct {
	Path                 string
	Ports                []int
	PermitRootLogin      string
	PasswordAuthentication string
	PubkeyAuthentication string
	PermitEmptyPasswords string
	X11Forwarding         string
	MaxAuthTries          string
	AllowUsers            []string
	AllowGroups           []string
}

func Status() (StatusInfo, error) {
	service, state := serviceStatus()
	path := configPath()
	var ports []int

	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			ports = parseConfig(data).Ports
		}
	}

	return StatusInfo{
		Service:    service,
		State:      state,
		ConfigPath: path,
		Ports:      ports,
	}, nil
}

func Config() (ConfigInfo, error) {
	path := configPath()

	if _, err := exec.LookPath("sshd"); err == nil {
		args := []string{"-T"}
		if path != "" {
			args = append(args, "-f", path)
		}

		if out, err := xexec.Run("sshd", args...); err == nil {
			config := parseConfig(out)
			config.Path = path
			return config, nil
		}
	}

	if path == "" {
		return ConfigInfo{}, fmt.Errorf("sshd config not found")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ConfigInfo{}, fmt.Errorf("read sshd config: %w", err)
	}

	config := parseConfig(data)
	config.Path = path
	return config, nil
}

func serviceStatus() (string, string) {
	if _, err := exec.LookPath("systemctl"); err != nil {
		return "unknown", "unknown"
	}

	for _, service := range []string{"sshd", "ssh"} {
		out, _ := xexec.Run("systemctl", "is-active", service)
		switch state := strings.TrimSpace(string(out)); state {
		case "active", "inactive", "failed", "activating", "deactivating":
			return service, state
		}
	}

	return "unknown", "unknown"
}

func configPath() string {
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func parseConfig(data []byte) ConfigInfo {
	config := ConfigInfo{}

	s := bufio.NewScanner(strings.NewReader(string(data)))
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 2 || strings.HasPrefix(fields[0], "#") {
			continue
		}

		key := strings.ToLower(fields[0])
		values := fields[1:]
		value := strings.Join(values, " ")

		switch key {
		case "port":
			if port, err := strconv.Atoi(values[0]); err == nil && port > 0 && port <= 65535 {
				config.Ports = append(config.Ports, port)
			}
		case "permitrootlogin":
			config.PermitRootLogin = value
		case "passwordauthentication":
			config.PasswordAuthentication = value
		case "pubkeyauthentication":
			config.PubkeyAuthentication = value
		case "permitemptypasswords":
			config.PermitEmptyPasswords = value
		case "x11forwarding":
			config.X11Forwarding = value
		case "maxauthtries":
			config.MaxAuthTries = value
		case "allowusers":
			config.AllowUsers = values
		case "allowgroups":
			config.AllowGroups = values
		}
	}

	if len(config.Ports) == 0 {
		config.Ports = []int{22}
	}

	return config
}
