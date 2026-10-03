package service

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

type Service struct {
	Name        string
	Description string
	State       string
	SubState    string
	Enabled     string
	MainPID     int
	Type        string
}

func List() ([]Service, error) {
	enabled, err := unitFileStates()
	if err != nil {
		return nil, err
	}

	out, err := xexec.Run(
		"systemctl",
		"list-units",
		"--type=service",
		"--all",
		"--no-legend",
		"--no-pager",
	)
	if err != nil {
		return nil, fmt.Errorf("systemctl list-units: %w: %s", err, strings.TrimSpace(string(out)))
	}

	services := parseListUnits(out, enabled)
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})
	return services, nil
}

func Status(name string) (Service, error) {
	out, err := xexec.Run(
		"systemctl",
		"show",
		name,
		"--no-pager",
		"-p", "Id",
		"-p", "Description",
		"-p", "ActiveState",
		"-p", "SubState",
		"-p", "UnitFileState",
		"-p", "MainPID",
		"-p", "Type",
	)
	if err != nil {
		return Service{}, fmt.Errorf("systemctl show %s: %w: %s", name, err, strings.TrimSpace(string(out)))
	}

	item := parseShow(out)
	if item.Name == "" {
		item.Name = name
	}
	return item, nil
}

func unitFileStates() (map[string]string, error) {
	out, err := xexec.Run(
		"systemctl",
		"list-unit-files",
		"--type=service",
		"--no-legend",
		"--no-pager",
	)
	if err != nil {
		return nil, fmt.Errorf("systemctl list-unit-files: %w: %s", err, strings.TrimSpace(string(out)))
	}

	states := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 {
			states[fields[0]] = fields[1]
		}
	}
	return states, nil
}

func parseListUnits(data []byte, enabled map[string]string) []Service {
	var services []Service

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}

		name := fields[0]
		description := strings.Join(fields[4:], " ")
		services = append(services, Service{
			Name:        name,
			State:       fields[2],
			SubState:    fields[3],
			Enabled:     enabled[name],
			Description: description,
		})
	}
	return services
}

func parseShow(data []byte) Service {
	item := Service{}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		key, value, ok := strings.Cut(scanner.Text(), "=")
		if !ok {
			continue
		}

		switch key {
		case "Id":
			item.Name = value
		case "Description":
			item.Description = value
		case "ActiveState":
			item.State = value
		case "SubState":
			item.SubState = value
		case "UnitFileState":
			item.Enabled = value
		case "MainPID":
			item.MainPID, _ = strconv.Atoi(value)
		case "Type":
			item.Type = value
		}
	}
	return item
}
