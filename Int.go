package main

import "fmt"

// Fetcher is the K.I.S.S. interface for all eclass helpers
type Fetcher func(pv, pn string, toml map[string]string) string

var Registry = map[string]Fetcher{}

func Register(name string, f Fetcher) {
    Registry[name] = f
}

func GetURL(eclass, pv, pn string, toml map[string]string) (string, error) {
    if f, ok := Registry[eclass]; ok {
        return f(pv, pn, toml), nil
    }
    return "", fmt.Errorf("no helper for eclass: %s", eclass)
}
