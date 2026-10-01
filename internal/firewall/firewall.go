package firewall

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

type Backend string

const (
	None      Backend = "none"
	UFW       Backend = "ufw"
	Firewalld Backend = "firewalld"
	NFT       Backend = "nftables"
	Iptables  Backend = "iptables"
)

type Info struct {
	Backend Backend
	State   string
}

func Detect() Backend {
	for _, item := range []struct {
		name    string
		backend Backend
	}{
		{"ufw", UFW},
		{"firewall-cmd", Firewalld},
		{"nft", NFT},
		{"iptables", Iptables},
	} {
		if _, err := exec.LookPath(item.name); err == nil {
			return item.backend
		}
	}
	return None
}

func Status() (Info, error) {
	backend := Detect()

	switch backend {
	case None:
		return Info{Backend: None, State: "inactive"}, nil
	case UFW:
		out, err := xexec.Run("ufw", "status")
		if err != nil {
			return Info{Backend: backend, State: "unknown"}, nil
		}
		return Info{Backend: backend, State: parseUFWStatus(out)}, nil
	case Firewalld:
		out, err := xexec.Run("firewall-cmd", "--state")
		if err != nil {
			return Info{Backend: backend, State: "inactive"}, nil
		}
		return Info{Backend: backend, State: parseState(out)}, nil
	case NFT:
		out, err := xexec.Run("nft", "list", "ruleset")
		if err != nil {
			return Info{Backend: backend, State: "unknown"}, nil
		}
		return Info{Backend: backend, State: rulesState(out)}, nil
	case Iptables:
		out, err := xexec.Run("iptables", "-S")
		if err != nil {
			return Info{Backend: backend, State: "unknown"}, nil
		}
		return Info{Backend: backend, State: rulesState(out)}, nil
	default:
		return Info{}, fmt.Errorf("unsupported firewall backend: %s", backend)
	}
}

func List() (string, error) {
	switch backend := Detect(); backend {
	case None:
		return fmt.Sprintln("no firewall detected"), nil
	case UFW:
		return run("ufw", "status", "verbose")
	case Firewalld:
		return run("firewall-cmd", "--list-all-zones")
	case NFT:
		return run("nft", "list", "ruleset")
	case Iptables:
		return run("iptables", "-S")
	default:
		return "", fmt.Errorf("unsupported firewall backend: %s", backend)
	}
}

func Allow(port int, protocol string) error {
	return rule("allow", port, protocol)
}

func Deny(port int, protocol string) error {
	return rule("deny", port, protocol)
}

func Enable() error {
	switch Detect() {
	case UFW:
		_, err := xexec.Run("ufw", "--force", "enable")
		return err
	case Firewalld:
		_, err := xexec.Run("systemctl", "enable", "--now", "firewalld")
		return err
	default:
		return fmt.Errorf("firewall enable is unsupported for backend %s", Detect())
	}
}

func Disable() error {
	switch Detect() {
	case UFW:
		_, err := xexec.Run("ufw", "disable")
		return err
	case Firewalld:
		_, err := xexec.Run("systemctl", "disable", "--now", "firewalld")
		return err
	default:
		return fmt.Errorf("firewall disable is unsupported for backend %s", Detect())
	}
}

func ParseSpec(value string) (int, string, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid rule: %s (use PORT/PROTOCOL)", value)
	}

	port, err := strconv.Atoi(parts[0])
	if err != nil || port < 1 || port > 65535 {
		return 0, "", fmt.Errorf("invalid port: %s", parts[0])
	}

	protocol := strings.ToLower(parts[1])
	if protocol != "tcp" && protocol != "udp" {
		return 0, "", fmt.Errorf("invalid protocol: %s", protocol)
	}

	return port, protocol, nil
}

func rule(action string, port int, protocol string) error {
	switch backend := Detect(); backend {
	case UFW:
		_, err := xexec.Run("ufw", action, fmt.Sprintf("%d/%s", port, protocol))
		return err
	case Firewalld:
		spec := fmt.Sprintf("%d/%s", port, protocol)
		args := []string{"--permanent", "--add-port", spec}

		if action == "deny" {
			spec = fmt.Sprintf("rule port port=%q protocol=%q drop", strconv.Itoa(port), protocol)
			args = []string{"--permanent", "--add-rich-rule", spec}
		}

		if _, err := xexec.Run("firewall-cmd", args...); err != nil {
			return err
		}

		_, err := xexec.Run("firewall-cmd", "--reload")
		return err
	case Iptables:
		target := "ACCEPT"
		if action == "deny" {
			target = "DROP"
		}

		_, err := xexec.Run(
			"iptables",
			"-I", "INPUT",
			"-p", protocol,
			"--dport", strconv.Itoa(port),
			"-j", target,
		)
		return err
	default:
		return fmt.Errorf("firewall rule control is unsupported for backend %s", backend)
	}
}

func run(name string, args ...string) (string, error) {
	out, err := xexec.Run(name, args...)
	if err != nil {
		return "", fmt.Errorf("%s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func parseUFWStatus(data []byte) string {
	status := strings.ToLower(string(data))
	if strings.Contains(status, "status: active") {
		return "active"
	}
	if strings.Contains(status, "status: inactive") {
		return "inactive"
	}
	return "unknown"
}

func parseState(data []byte) string {
	if strings.EqualFold(strings.TrimSpace(string(data)), "running") {
		return "active"
	}
	return "inactive"
}

func rulesState(data []byte) string {
	if strings.TrimSpace(string(data)) == "" {
		return "inactive"
	}
	return "active"
}
