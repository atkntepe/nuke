package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	showAll := false
	force := false
	var positional []string

	for _, a := range args {
		switch a {
		case "--all":
			showAll = true
		case "--force":
			force = true
		default:
			positional = append(positional, a)
		}
	}

	if len(positional) == 0 {
		runInteractive(showAll, force)
		return
	}

	switch positional[0] {
	case "list":
		runList(showAll)
	case "cache":
		runCache(showAll)
	default:
		runKillPort(positional[0], force)
	}
}

func runInteractive(showAll, force bool) {
	ports, err := ListDevPorts(showAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning ports: %s\n", err)
		os.Exit(1)
	}

	if len(ports) == 0 {
		if showAll {
			fmt.Println("No processes found on ports (1024-65535).")
		} else {
			fmt.Println("No processes found on dev ports (3000-9999).")
			fmt.Println("Run with --all to scan all ports.")
		}
		return
	}

	selected, err := SelectPorts(ports)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if len(selected) == 0 {
		return
	}

	KillProcesses(selected, force)
}

func runList(showAll bool) {
	ports, err := ListDevPorts(showAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning ports: %s\n", err)
		os.Exit(1)
	}

	if len(ports) == 0 {
		if showAll {
			fmt.Println("No processes found on ports (1024-65535).")
		} else {
			fmt.Println("No processes found on dev ports (3000-9999).")
			fmt.Println("Run with --all to scan all ports.")
		}
		return
	}

	maxPortW := 4
	maxPidW := 3
	maxNameW := 4
	for _, p := range ports {
		pw := len(strconv.Itoa(p.Port))
		if pw > maxPortW {
			maxPortW = pw
		}
		pidW := len(strconv.Itoa(p.PID))
		if pidW > maxPidW {
			maxPidW = pidW
		}
		if len(p.Name) > maxNameW {
			maxNameW = len(p.Name)
		}
	}

	fmt.Printf("%-*s   %-*s   %-*s   %s\n", maxPortW, "PORT", maxPidW, "PID", maxNameW, "NAME", "COMMAND")
	for _, p := range ports {
		fmt.Printf("%-*d   %-*d   %-*s   %s\n", maxPortW, p.Port, maxPidW, p.PID, maxNameW, p.Name, p.Command)
	}
}

func runKillPort(portArg string, force bool) {
	port, err := strconv.Atoi(portArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid port: %s\n", portArg)
		os.Exit(1)
	}

	ports, err := ListDevPorts(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning ports: %s\n", err)
		os.Exit(1)
	}

	var target []PortProcess
	for _, p := range ports {
		if p.Port == port {
			target = append(target, p)
			break
		}
	}

	if len(target) == 0 {
		fmt.Printf("no process found on port %d\n", port)
		os.Exit(1)
	}

	KillProcesses(target, force)
}

func runCache(clearAll bool) {
	caches := DetectAvailableCaches()

	if len(caches) == 0 {
		fmt.Println("No cache targets detected.")
		return
	}

	if clearAll {
		ClearCaches(caches)
		return
	}

	items := make([]string, len(caches))
	for i, c := range caches {
		items[i] = c.Name
	}

	indices, err := SelectItems(items)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if len(indices) == 0 {
		return
	}

	var selected []CacheTarget
	for _, idx := range indices {
		selected = append(selected, caches[idx])
	}

	ClearCaches(selected)
}
