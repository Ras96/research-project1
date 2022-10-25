package editdistance

import (
	"fmt"

	"github.com/Ras96/research-project1/retriever"
	"github.com/agnivade/levenshtein"
)

type editDistanceRetriever struct{
	dict retriever.Dictionary
}

// *editDistanceRetrieverはretriever.Retrieverインターフェイスを満たす
func NewEditDistanceRetriever(dict retriever.Dictionary) retriever.Retriever {
	return &editDistanceRetriever{dict}
}

func (r *editDistanceRetriever) Retrieve(req string) string {
	var (
		minDist int = 1e9
		bestRes string
		ref     string
	)

	for k, v := range r.dict {
		d := levenshtein.ComputeDistance(k, req)
		if d < minDist {
			minDist = d
			bestRes = v
			ref = k
		}
	}

	fmt.Println("request      :", req)
	fmt.Println("best response:", bestRes)
	fmt.Println("edit distance:", minDist)
	fmt.Println("reference    :", ref)

	return bestRes
}
