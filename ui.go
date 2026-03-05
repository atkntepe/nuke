package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func readKey() string {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil || n == 0 {
		return ""
	}

	if n == 1 {
		switch buf[0] {
		case 27:
			return "esc"
		case 3:
			return "ctrl-c"
		case 13:
			return "enter"
		case ' ':
			return "space"
		case 'q':
			return "q"
		case 'a':
			return "a"
		case 'j':
			return "down"
		case 'k':
			return "up"
		}
		return ""
	}

	if n >= 3 && buf[0] == 27 && buf[1] == '[' {
		switch buf[2] {
		case 'A':
			return "up"
		case 'B':
			return "down"
		}
	}

	return ""
}

func SelectItems(items []string) ([]int, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	defer term.Restore(fd, oldState)

	if len(items) == 0 {
		return nil, nil
	}

	cursor := 0
	selected := make([]bool, len(items))

	render := func() {
		var sb strings.Builder
		sb.WriteString("\033[H\033[J")
		sb.WriteString("Select items (\033[1m\xe2\x86\x91\xe2\x86\x93\033[0m move, \033[1mspace\033[0m select, \033[1ma\033[0m all, \033[1menter\033[0m confirm, \033[1mq/esc\033[0m quit)\r\n\r\n")
		for i, item := range items {
			check := "[ ]"
			if selected[i] {
				check = "[x]"
			}
			if i == cursor {
				if selected[i] {
					sb.WriteString(fmt.Sprintf("\033[7m\033[1m  %s %s\033[0m\r\n", check, item))
				} else {
					sb.WriteString(fmt.Sprintf("\033[7m  %s %s\033[0m\r\n", check, item))
				}
			} else if selected[i] {
				sb.WriteString(fmt.Sprintf("\033[1m  %s %s\033[0m\r\n", check, item))
			} else {
				sb.WriteString(fmt.Sprintf("  %s %s\r\n", check, item))
			}
		}
		os.Stdout.WriteString(sb.String())
	}

	render()

	for {
		key := readKey()

		switch key {
		case "q", "esc":
			fmt.Print("\033[H\033[J")
			return nil, nil
		case "space":
			selected[cursor] = !selected[cursor]
		case "enter":
			fmt.Print("\033[H\033[J")
			var indices []int
			for i, s := range selected {
				if s {
					indices = append(indices, i)
				}
			}
			return indices, nil
		case "a":
			allSelected := true
			for _, s := range selected {
				if !s {
					allSelected = false
					break
				}
			}
			for i := range selected {
				selected[i] = !allSelected
			}
		case "up":
			if cursor > 0 {
				cursor--
			}
		case "down":
			if cursor < len(items)-1 {
				cursor++
			}
		case "ctrl-c":
			term.Restore(fd, oldState)
			fmt.Print("\033[H\033[J")
			os.Exit(0)
		}

		render()
	}
}

func SelectPorts(ports []PortProcess) ([]PortProcess, error) {
	if len(ports) == 0 {
		return nil, nil
	}

	maxPortWidth := 4
	maxNameWidth := 4
	for _, p := range ports {
		pw := len(fmt.Sprintf("%d", p.Port))
		if pw > maxPortWidth {
			maxPortWidth = pw
		}
		if len(p.Name) > maxNameWidth {
			maxNameWidth = len(p.Name)
		}
	}

	items := make([]string, len(ports))
	for i, p := range ports {
		items[i] = fmt.Sprintf("%-*d   %-*s   %s", maxPortWidth, p.Port, maxNameWidth, p.Name, p.Command)
	}

	indices, err := SelectItems(items)
	if err != nil {
		return nil, err
	}

	var result []PortProcess
	for _, idx := range indices {
		result = append(result, ports[idx])
	}
	return result, nil
}
