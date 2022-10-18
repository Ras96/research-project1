package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
)

const jsonDirName = "corpus/json/init100"

//go:embed corpus/json/init100/*.json
var jsonFiles embed.FS

// JSONデータから必要なものだけを抽出する構造体
type simpleData struct {
	Turns []struct {
		Utterance string `json:"utterance"`
	} `json:"turns"`
}

func main() {
	kv := makeResponseMap()

	if len(os.Args) < 2 {
		fmt.Println("Please input your message.")
		os.Exit(1)
	}

	req := os.Args[1]
	if res, ok := kv[req]; ok {
		fmt.Println(res)
	} else {
		fmt.Println("I don't know.")
	}
}

// JSONファイルを読みこみ、発話に対する返答を記録
func makeResponseMap() map[string]string {
	files, _ := jsonFiles.ReadDir(jsonDirName)
	kv := make(map[string]string)
	for _, f := range files {
		bytes, _ := jsonFiles.ReadFile(jsonDirName + "/" + f.Name())

		var data simpleData
		json.Unmarshal(bytes, &data)

		turns := data.Turns
		for i := 0; i+1 < len(turns); i++ {
			kv[turns[i].Utterance] = turns[i+1].Utterance
		}
	}

	return kv
}
