package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

// ManifestEntry holds the metadata required for a Gentoo Manifest line
type ManifestEntry struct {
	Filename string
	Size     int64
	Blake2b  string
	Sha512   string
}

// GhostHash streams a URL and calculates hashes without writing to disk
func GhostHash(url string, distName string) (*ManifestEntry, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Initialize Gentoo-standard hashes
	b2, _ := blake2b.New512(nil)
	s512 := sha512.New()

	// MultiWriter pipes the network stream to both hashers AND /dev/null (io.Discard)
	// This ensures we never fill up /tmp or RAM.
	mw := io.MultiWriter(b2, s512, io.Discard)

	// Stream the data
	size, err := io.Copy(mw, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("streaming failed: %v", err)
	}

	return &ManifestEntry{
		Filename: distName,
		Size:     size,
		Blake2b:  fmt.Sprintf("%x", b2.Sum(nil)),
		Sha512:   fmt.Sprintf("%x", s512.Sum(nil)),
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ghost-hashbot <URL> [distname]")
		os.Exit(1)
	}

	url := os.Args[1]
	distName := filepath.Base(url)
	if len(os.Args) > 2 {
		distName = os.Args[2]
	}

	entry, err := GhostHash(url, distName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output in standard Gentoo Manifest format
	fmt.Printf("DIST %s %d BLAKE2B %s SHA512 %s\n",
		entry.Filename, entry.Size, entry.Blake2b, entry.Sha512)
}
