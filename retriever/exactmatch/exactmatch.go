package exactmatch

import "github.com/Ras96/research-project1/retriever"

type exactmatchRetriever struct{}

func NewExactMatchRetriever() retriever.Retriever {
	return &exactmatchRetriever{}
}

func (r *exactmatchRetriever) Retrieve(dict map[string]string, req string) string {
	if res, ok := dict[req]; ok {
		return res
	} else {
		return "I don't know."
	}
}
