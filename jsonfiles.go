//go:build !debug

package main

import "embed"

const jsonDirName = "corpus/json/rest1046"

//go:embed corpus/json/rest1046/*.json
var jsonFiles embed.FS
