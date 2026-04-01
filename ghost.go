package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/blake2b"
)

// ToManifestLine formats the entry for a Gentoo Manifest file
func (m *ManifestEntry) ToManifestLine() string {
	return fmt.Sprintf("DIST %s %d BLAKE2B %s SHA512 %s",
		m.Filename, m.Size, m.Blake2b, m.Sha512)
}

// GhostHash remains the core streaming function
func GhostHash(url string, distName string) (*ManifestEntry, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	b2, _ := blake2b.New512(nil)
	s512 := sha512.New()

	// The "Ghost" magic: Bytes go into hashes and then vanish
	mw := io.MultiWriter(b2, s512, io.Discard)

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
