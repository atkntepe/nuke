package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

func KillProcesses(ports []PortProcess, force bool) error {
	signal := "-15"
	if force {
		signal = "-9"
	}

	for _, pp := range ports {
		pidStr := strconv.Itoa(pp.PID)
		err := exec.Command("kill", signal, pidStr).Run()
		if err != nil {
			fmt.Printf("  failed  %d  %s  (PID %d) — %s\n", pp.Port, pp.Name, pp.PID, err)
			continue
		}

		time.Sleep(500 * time.Millisecond)

		out, _ := exec.Command("lsof", "-i", ":"+strconv.Itoa(pp.Port), "-sTCP:LISTEN", "-n", "-P").Output()
		if len(out) == 0 {
			fmt.Printf("  killed  %d  %s  (PID %d)\n", pp.Port, pp.Name, pp.PID)
		} else {
			fmt.Printf("  failed  %d  %s  (PID %d) — try with --force\n", pp.Port, pp.Name, pp.PID)
		}
	}

	return nil
}
