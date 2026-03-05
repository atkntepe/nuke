package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CacheTarget struct {
	Name    string
	Command string
	Paths   []string
}

func ListCacheTargets() []CacheTarget {
	home, _ := os.UserHomeDir()

	return []CacheTarget{
		{Name: "npm", Command: "npm cache clean --force"},
		{Name: "yarn", Command: "yarn cache clean"},
		{Name: "pnpm", Command: "pnpm store prune"},
		{Name: "bun", Command: "bun pm cache rm"},
		{Name: "vite", Paths: []string{"node_modules/.vite"}},
		{Name: "next", Paths: []string{".next/cache"}},
		{Name: "turbo", Paths: []string{".turbo"}},
		{Name: "go", Command: "go clean -cache"},
		{Name: "gradle", Paths: []string{filepath.Join(home, ".gradle", "caches")}},
		{Name: "maven", Paths: []string{filepath.Join(home, ".m2", "repository")}},
		{Name: "docker", Command: "docker system prune -f"},
	}
}

func DetectAvailableCaches() []CacheTarget {
	all := ListCacheTargets()
	var available []CacheTarget

	for _, t := range all {
		if t.Command != "" {
			bin := strings.Fields(t.Command)[0]
			if _, err := exec.LookPath(bin); err == nil {
				available = append(available, t)
			}
		} else if len(t.Paths) > 0 {
			for _, p := range t.Paths {
				if _, err := os.Stat(p); err == nil {
					available = append(available, t)
					break
				}
			}
		}
	}

	return available
}

func ClearCaches(targets []CacheTarget) error {
	for _, t := range targets {
		if t.Command != "" {
			parts := strings.Fields(t.Command)
			cmd := exec.Command(parts[0], parts[1:]...)
			if err := cmd.Run(); err != nil {
				fmt.Printf("  failed   %s cache — %s\n", t.Name, err)
			} else {
				fmt.Printf("  cleared  %s cache\n", t.Name)
			}
		} else if len(t.Paths) > 0 {
			failed := false
			for _, p := range t.Paths {
				if err := os.RemoveAll(p); err != nil {
					fmt.Printf("  failed   %s cache — %s\n", t.Name, err)
					failed = true
					break
				}
			}
			if !failed {
				fmt.Printf("  cleared  %s cache\n", t.Name)
			}
		}
	}
	return nil
}
