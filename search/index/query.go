package index

import (
	"strings"

	"github.com/txominpelu/fnd/search"
)

// Query. Return docIds
func (i IndexedLines) FilterEntries(subQueries []search.SubQuery) []int {
	if len(subQueries) > 0 {
		results := i.index.perfieldWord2Doc[subQueries[0].Field][strings.ToLower(subQueries[0].Query)]
		for _, sQ := range subQueries[1:] {
			if docs, ok := i.index.perfieldWord2Doc[sQ.Field][strings.ToLower(sQ.Query)]; ok {
				results = intersection(results, docs)
			} else {
				return []int{}
			}
		}
		docIds := make([]int, len(results))
		{
			j := 0
			for dId := range results {
				docIds[j] = dId
				j++
			}
		}
		return docIds
	}
	// otherwise if no query all docs match
	return i.docIds
}

func intersection(s1 map[int]bool, s2 map[int]bool) map[int]bool {
	result := map[int]bool{}
	for k := range s1 {
		if _, ok := s2[k]; ok {
			result[k] = true
		}
	}
	return result
}
