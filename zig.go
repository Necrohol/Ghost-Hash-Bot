package main

import "fmt"

func init() {
    Register("zig-utils", ZigFetch)
}

func ZigFetch(pv, pn string, toml map[string]string) string {
    // Template logic: logic sniffed by regex-rename.go
    // Using the TOML definition as the source of truth
    base := toml["primary_url"] 
    // Return the resolved URL
    return fmt.Sprintf("%s/archive/%s.tar.gz", base, pv)
}
