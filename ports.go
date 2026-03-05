package main

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type PortProcess struct {
	Port    int
	PID     int
	Name    string
	Command string
}

func ListDevPorts(showAll bool) ([]PortProcess, error) {
	out, err := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-n", "-P").CombinedOutput()
	if err != nil {
		if len(out) == 0 {
			return nil, nil
		}
	}

	results := parseLsofOutput(string(out))

	minPort := 3000
	maxPort := 9999
	if showAll {
		minPort = 1024
		maxPort = 65535
	}

	var filtered []PortProcess
	for _, pp := range results {
		if pp.Port >= minPort && pp.Port <= maxPort {
			filtered = append(filtered, pp)
		}
	}

	return filtered, nil
}

func parseLsofOutput(output string) []PortProcess {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		return nil
	}

	seen := make(map[int]bool)
	var results []PortProcess

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		command := fields[0]

		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		lastField := fields[len(fields)-1]
		if lastField != "(LISTEN)" {
			continue
		}

		nameCol := fields[len(fields)-2]
		colonIdx := strings.LastIndex(nameCol, ":")
		if colonIdx == -1 {
			continue
		}

		portStr := nameCol[colonIdx+1:]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		if seen[port] {
			continue
		}
		seen[port] = true

		results = append(results, PortProcess{
			Port:    port,
			PID:     pid,
			Name:    command,
			Command: command,
		})
	}

	for i := range results {
		results[i].Command = getFullCommand(results[i].PID)
		if results[i].Command == "" {
			results[i].Command = results[i].Name
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

func getFullCommand(pid int) string {
	out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "args=").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
