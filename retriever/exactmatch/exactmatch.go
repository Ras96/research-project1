package exactmatch

import "github.com/Ras96/research-project1/retriever"

type exactmatchRetriever struct{
	dict retriever.Dictionary
}

func NewExactMatchRetriever(dict retriever.Dictionary) retriever.Retriever {
	return &exactmatchRetriever{dict}
}

func (r *exactmatchRetriever) Retrieve(req string) string {
	if res, ok := r.dict[req]; ok {
		return res
	} else {
		return "I don't know."
	}
}
