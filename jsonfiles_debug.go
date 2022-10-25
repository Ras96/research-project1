//go:build debug

// -tags debugをつけて実行するとinit100のみが使用される

package main

import "embed"

const (
	jsonDirName = "corpus/json/init100"
	isDebug     = true
)

//go:embed corpus/json/init100/*.json
var jsonFiles embed.FS
