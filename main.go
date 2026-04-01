package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Simple CLI for now (RAD mode)
	if len(os.Args) < 4 {
		fmt.Println("Usage: ghost <eclass> <pn> <pv>")
		os.Exit(1)
	}

	eclass, pn, pv := os.Args[1], os.Args[2], os.Args[3]

	// 1. Get the URL from Int.go / Registry
	url, err := ResolveURL(eclass, pv, pn)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// 2. Perform the Ghost-Hash
	distName := filepath.Base(url) // Or handle -> renames here
	entry, err := GhostHash(url, distName)
	if err != nil {
		log.Fatalf("Ghosting failed: %v", err)
	}

	// 3. Print the Gentoo-ready line
	fmt.Println(entry.ToManifestLine())
}
