package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RenameRule defines how a URL template maps to a Manifest filename
type RenameRule struct {
	Eclass  string
	Pattern string
	Target  string
}

// SniffRenameLogic parses an eclass to find '->' rename patterns
func SniffRenameLogic(path string) (*RenameRule, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regex to find the '->' operator within a SRC_URI context
	// Matches: SRC_URI="... -> ${P}.tar.gz"
	re := regexp.MustCompile(`SRC_URI=.*?\s+->\s+["']?([^"'\s]+)["']?`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		
		// If we find the '->' operator, extract the target name
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			return &RenameRule{
				Eclass:  strings.TrimSuffix(filepath.Base(path), ".eclass"),
				Pattern: "standard-rename",
				Target:  matches[1],
			}, nil
		}
	}
	return nil, fmt.Errorf("no rename logic found in %s", path)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: regex-rename <path-to-eclass>")
		return
	}

	path := os.Args[1]
	rule, err := SniffRenameLogic(path)
	if err != nil {
		// If it doesn't rename, it's a "flat" fetch (skip)
		os.Exit(0) 
	}

	// Output as a TOML snippet for the Ghost-Hash bot
	fmt.Printf("[rename_logic.%s]\n", rule.Eclass)
	fmt.Printf("type = \"%s\"\n", rule.Pattern)
	fmt.Printf("template = \"%s\"\n", rule.Target)
}
