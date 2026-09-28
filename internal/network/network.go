package network

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

type Interface struct {
	Name string
	State string
	MAC  string
}

func Interfaces() ([]Interface, error) {
	out, err := xexec.Run("ip", "-o", "link", "show")
	if err != nil {
		return nil, fmt.Errorf("ip link: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return parseInterfaces(out), nil
}

func parseInterfaces(data []byte) []Interface {
	var result []Interface
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 3 {
			continue
		}
		name := strings.TrimSuffix(fields[1], ":")
		state := ""
		mac := ""
		for i, field := range fields {
			if field == "state" && i+1 < len(fields) {
				state = fields[i+1]
			}
			if field == "link/ether" && i+1 < len(fields) {
				mac = fields[i+1]
			}
		}
		result = append(result, Interface{Name: name, State: state, MAC: mac})
	}
	return result
}

type Route struct {
	Destination string
	Gateway     string
	Device      string
}

func Routes() ([]Route, error) {
	out, err := xexec.Run("ip", "-o", "route", "show")
	if err != nil {
		return nil, fmt.Errorf("ip route: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return parseRoutes(out), nil
}

func parseRoutes(data []byte) []Route {
	var result []Route
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 1 {
			continue
		}
		route := Route{Destination: fields[0]}
		for i, field := range fields {
			if field == "via" && i+1 < len(fields) {
				route.Gateway = fields[i+1]
			}
			if field == "dev" && i+1 < len(fields) {
				route.Device = fields[i+1]
			}
		}
		result = append(result, route)
	}
	return result
}
