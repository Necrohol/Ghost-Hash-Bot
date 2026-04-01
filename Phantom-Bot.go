package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SearchPaths defines where we hunt for eclasses, in order of speed.
var SearchPaths = []string{
	"/var/db/repos/gentoo/eclass",    // Local Main
	"/var/db/repos/*/eclass",         // Local Overlays (Wildcard)
}

const CodebergFallback = "https://codeberg.org/gentoo/gentoo/raw/branch/master/eclass/"

func isFetchHeavy(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// We only care about eclasses that define these specific "Heavy" keywords
	keywords := []string{"SRC_URI", "EGIT_REPO_URI", "cargo_src_unpack", "zig-fetch", "go-module_src_unpack"}
	
	for scanner.Scan() {
		line := scanner.Text()
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				return true
			}
		}
	}
	return false
}

func Hunt(targetEclass string) {
	eclassName := targetEclass + ".eclass"
	found := false

	// 1. Try Local Disk First (Fastest)
	for _, pathPattern := range SearchPaths {
		matches, _ := filepath.Glob(filepath.Join(pathPattern, eclassName))
		if len(matches) > 0 {
			if isFetchHeavy(matches[0]) {
				fmt.Printf("🎯 Found Fetch-Heavy Eclass: %s\n", matches[0])
				GenerateToml(matches[0])
				found = true
				break
			}
		}
	}

	// 2. Fallback to Codeberg (Cross-Platform / No-Linux mode)
	if !found {
		fmt.Printf("🌐 Local not found or skipped. Checking Codeberg: %s%s\n", CodebergFallback, eclassName)
		// Logic to HTTP GET from Codeberg and run isFetchHeavy on the response body...
	}
}

func GenerateToml(path string) {
	// Logic to extract SRC_URI regex and build the [eclass.toml]
	name := strings.TrimSuffix(filepath.Base(path), ".eclass")
	fmt.Printf("\n[%s.toml]\n", name)
	fmt.Println("status = \"active\"")
	fmt.Println("fetch_strategy = \"ghost-stream\"")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: phantom-bot <eclass-name>")
		return
	}
	Hunt(os.Args[1])
}
