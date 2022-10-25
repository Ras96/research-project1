//go:build !debug

package main

import "embed"

const (
	jsonDirName = "corpus/json/rest1046"
	isDebug     = false
)

//go:embed corpus/json/rest1046/*.json
var jsonFiles embed.FS
