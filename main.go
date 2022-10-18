package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ras96/research-project1/retriever"
	"github.com/Ras96/research-project1/retriever/editdistance"
	"github.com/Ras96/research-project1/retriever/exactmatch"
	"github.com/manifoldco/promptui"
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
	// プロンプトから応答選択手法と入力文字列を決定
	r := selectRetrieverMethodInPrompt()
	req := getRequestMessageInPrompt()

	// 応答選択の基準となる辞書を作成
	kv := makeResponseMap()

	// 選択した応答選択手法で最も適した応答を出力
	res := r.Retrieve(kv, req)
	fmt.Println("response:", res)
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

// 応答選択手法をインタラクティブに選択
func selectRetrieverMethodInPrompt() retriever.Retriever {
	const (
		methodExactMatch   = "Exact Match"
		methodEditDistance = "Edit Distance"
	)

	p := promptui.Select{
		Label: "Which method do you want to use?",
		Items: []string{
			methodExactMatch,
			methodEditDistance,
		},
	}

	_, method, _ := p.Run()

	var r retriever.Retriever
	switch method {
	case methodExactMatch:
		r = exactmatch.NewExactMatchRetriever()
	case methodEditDistance:
		r = editdistance.NewEditDistanceRetriever()
	default:
		fmt.Println("Select an method")
		os.Exit(1)
	}

	return r
}

func getRequestMessageInPrompt() string {
	p := promptui.Prompt{
		Label: "Input your message",
		Validate: func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("Input your message")
			}

			return nil
		},
	}

	req, _ := p.Run()

	return req
}
