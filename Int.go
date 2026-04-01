package main

import "errors"

// Registry for eclass-specific URI builders
var EclassRegistry = make(map[string]func(pv, pn string) string)

func ResolveURL(eclass, pv, pn string) (string, error) {
	if fn, ok := EclassRegistry[eclass]; ok {
		return fn(pv, pn), nil
	}
	return "", errors.New("eclass not supported in registry")
}
