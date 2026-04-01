package main

import (
	"fmt"
	"strings"
	// The core logic for Gentoo Go-module hashing
	egovendor "github.com/williamh/get-ego-vendor/pkg/egovendor"
)

func init() {
	EclassRegistry["go-module"] = func(pv, pn string) ([]string, []string) {
		var urls []string
		var filenames []string

		// Use get-ego-vendor to parse a go.sum (streamed into RAM)
		// For RAD, we assume 'data' is the []byte from the ghost-fetched go.sum
		data := []byte("github.com/pkg/errors v0.8.1 h1:...") 
		
		entries, err := egovendor.ParseGoSum(data)
		if err != nil {
			return nil, nil
		}

		for _, entry := range entries {
			// Construct the Proxy URL: https://proxy.golang.org/<mod>/@v/<ver>.zip
			url := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.zip", entry.Path, entry.Version)
			
			// Standard Gentoo naming: go-module-github-com-pkg-errors-v0.8.1.zip
			safePath := strings.ReplaceAll(entry.Path, "/", "-")
			distname := fmt.Sprintf("go-module-%s-%s.zip", safePath, entry.Version)
			
			urls = append(urls, url)
			filenames = append(filenames, distname)
		}

		return urls, filenames
	}
}
