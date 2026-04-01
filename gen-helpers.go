package main

import (
	"os"
	"text/template"
)

type EclassData struct {
	EclassName  string
	Pattern     string
	URLTemplate string
}

func main() {
	// Example data from your TOML-ifyer
	data := EclassData{
		EclassName:  "zig-utils",
		Pattern:     "standard-rename",
		URLTemplate: "https://github.com/zigtools/zls/archive/refs/tags/%s.tar.gz",
	}

	tmpl, _ := template.ParseFiles("helpers.tmpl")
	f, _ := os.Create("zig_helpers.go")
	defer f.Close()

	tmpl.Execute(f, data)
}
