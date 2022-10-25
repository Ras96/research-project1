package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Ras96/research-project1/retriever"
	"github.com/Ras96/research-project1/retriever/editdistance"
	"github.com/Ras96/research-project1/retriever/exactmatch"
	"github.com/Ras96/research-project1/retriever/tfidf"
	"github.com/manifoldco/promptui"
)

func main() {
	// 応答選択の基準となる辞書を作成
	dict := makeResponseDictionary()

	// プロンプトから応答選択手法を決定
	r := selectRetrieverMethodInPrompt(dict)

	for {
		// プロンプトから入力文字列を決定
		req := getRequestMessageInPrompt()

		// 選択した応答選択手法で最も適した応答を出力
		res := r.Retrieve(req)
		fmt.Println("response:", res)
		fmt.Println()
	}
}

// JSONデータから必要なものだけを抽出する構造体
type simpleData struct {
	Turns []struct {
		Utterance string `json:"utterance"`
	} `json:"turns"`
}

// JSONファイルを読みこみ、発話に対する返答を記録
func makeResponseDictionary() retriever.Dictionary {
	files, _ := jsonFiles.ReadDir(jsonDirName)
	kv := make(retriever.Dictionary)
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
func selectRetrieverMethodInPrompt(dict retriever.Dictionary) retriever.Retriever {
	const (
		methodExactMatch   = "Exact Match"
		methodEditDistance = "Edit Distance"
		methodTfIdf        = "TF*IDF"
	)

	p := promptui.Select{
		Label: "Which method do you want to use?",
		Items: []string{
			methodExactMatch,
			methodEditDistance,
			methodTfIdf,
		},
		HideHelp: true,
	}

	_, method, _ := p.Run()

	var r retriever.Retriever
	switch method {
	case methodExactMatch:
		r = exactmatch.NewExactMatchRetriever(dict)
	case methodEditDistance:
		r = editdistance.NewEditDistanceRetriever(dict)
	case methodTfIdf:
		r = tfidf.NewTfIdfRetriever(dict)
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

	req, err := p.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrEOF) || errors.Is(err, promptui.ErrInterrupt) {
			fmt.Println("Bye!")
			os.Exit(0)
		} else {
			panic(err)
		}
	}

	return req
}
