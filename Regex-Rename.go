package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RenameRule struct {
	Eclass  string
	URL     string
	Target  string
}

func SniffRenameLogic(path string) (*RenameRule, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Improved Regex: Finds "URL -> Rename" even with spacing/quotes
	// Groups: 1=URL, 2=Target
	re := regexp.MustCompile(`SRC_URI=["']?([^"'\s]+)\s+->\s+([^"'\s]+)["']?`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			return &RenameRule{
				Eclass: strings.TrimSuffix(filepath.Base(path), ".eclass"),
				URL:    matches[1],
				Target: matches[2],
			}, nil
		}
	}
	return nil, fmt.Errorf("no complex SRC_URI logic in %s", path)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: regex-rename <path-to-eclass>")
		return
	}

	path := os.Args[1]
	rule, err := SniffRenameLogic(path)
	if err != nil {
		// Just output a flat mapping if no rename is found
		eName := strings.TrimSuffix(filepath.Base(path), ".eclass")
		fmt.Printf("[%s]\nrename = false\n", eName)
		return
	}

	fmt.Printf("[%s]\n", rule.Eclass)
	fmt.Printf("url_template = \"%s\"\n", rule.URL)
	fmt.Printf("rename_template = \"%s\"\n", rule.Target)
}
