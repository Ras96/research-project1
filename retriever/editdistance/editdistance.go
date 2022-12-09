package editdistance

import (
	"fmt"

	"github.com/agnivade/levenshtein"
	"github.com/ras0q/research-project1/retriever"
)

type editDistanceRetriever struct {
	dict    retriever.Dictionary
	isDebug bool
}

// *editDistanceRetrieverはretriever.Retrieverインターフェイスを満たす
func NewEditDistanceRetriever(dict retriever.Dictionary, isDebug bool) retriever.Retriever {
	return &editDistanceRetriever{dict, isDebug}
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

	if r.isDebug {
		fmt.Println("request      :", req)
		fmt.Println("best response:", bestRes)
		fmt.Println("edit distance:", minDist)
		fmt.Println("reference    :", ref)
	}

	return fmt.Sprintf("%s (dist=%d, ref=%s)", bestRes, minDist, ref)
}
