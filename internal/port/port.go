package port

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	xexec "github.com/zhravan/defenux/internal/exec"
)

type Listener struct {
	Protocol string
	Address  string
	Port     int
	Process  string
	PID      int
}

func List() ([]Listener, error) {
	out, err := xexec.Run("ss", "-H", "-lntp")
	if err != nil {
		return nil, fmt.Errorf("ss: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return parse(out)
}

func parse(data []byte) ([]Listener, error) {
	var listeners []Listener
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 4 { continue }
		address, portText := splitAddress(fields[3])
		port, err := strconv.Atoi(portText)
		if err != nil { continue }
		l := Listener{Protocol: fields[0], Address: address, Port: port}
		if len(fields) > 5 { l.Process, l.PID = parseProcess(fields[5:]) }
		listeners = append(listeners, l)
	}
	return listeners, s.Err()
}

func splitAddress(value string) (string, string) {
	if i := strings.LastIndexByte(value, ':'); i >= 0 {
		return strings.Trim(value[:i], "[]"), value[i+1:]
	}
	return value, ""
}

func parseProcess(fields []string) (string, int) {
	joined := strings.Join(fields, " ")
	if i := strings.Index(joined, "pid="); i >= 0 {
		rest := joined[i+4:]
		if end := strings.IndexAny(rest, ",)"); end >= 0 { rest = rest[:end] }
		pid, _ := strconv.Atoi(rest)
		return joined, pid
	}
	return joined, 0
}
