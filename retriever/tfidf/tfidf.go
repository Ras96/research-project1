package tfidf

import (
	"fmt"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/ras0q/research-project1/retriever" // mecab-ipadic-2.7.0-20070801
	"github.com/wilcosheh/tfidf"
	"github.com/wilcosheh/tfidf/seg"
	"github.com/wilcosheh/tfidf/similarity"
)

type tfIdfRetriever struct {
	f       *tfidf.TFIDF
	dict    retriever.Dictionary
	isDebug bool
}

func NewTfIdfRetriever(dict retriever.Dictionary, isDebug bool) retriever.Retriever {
	t := NewMyTokenizer()
	f := tfidf.NewTokenizer(t)

	for doc := range dict {
		f.AddDocs(doc)
	}

	return &tfIdfRetriever{f, dict, isDebug}
}

func (r *tfIdfRetriever) Retrieve(req string) string {
	reqW := r.f.Cal(req)

	maxScore := 0.0
	maxDoc := ""
	for doc := range r.dict {
		docW := r.f.Cal(doc)
		score := similarity.Cosine(docW, reqW)
		if score > maxScore {
			maxScore = score
			maxDoc = doc
		}
	}

	if r.isDebug {
		fmt.Println("maxScore      :", maxScore)
		fmt.Println("maxDoc        :", maxDoc)
		fmt.Println("maxDocResponse:", r.dict[maxDoc])
	}

	return fmt.Sprintf("%s (score=%f, ref=%s)", r.dict[maxDoc], maxScore, maxDoc)
}

type myTokenizer struct {
	t *tokenizer.Tokenizer
}

func NewMyTokenizer() seg.Tokenizer {
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	return &myTokenizer{t}
}

func (t *myTokenizer) Seg(text string) []string {
	tokens := t.t.Tokenize(text)
	words := make([]string, 0, len(tokens))
	for _, token := range tokens {
		words = append(words, token.Surface)
	}

	return words
}

func (t *myTokenizer) Free() {}
