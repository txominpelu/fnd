package fuzzy

import (
	"github.com/txominpelu/fnd/search"
)

type FuzzySearcher struct {
	docs   []search.Document
	docIds []int
}

func NewFuzzySearcher() FuzzySearcher {
	return FuzzySearcher{
		docs:   []search.Document{},
		docIds: []int{},
	}
}

func (f *FuzzySearcher) AddDocument(d search.Document) {
	f.docIds = append(f.docIds, len(f.docs))
	f.docs = append(f.docs, d)
}

func (f *FuzzySearcher) FilterEntries(subQueries []search.SubQuery) []int {
	if len(subQueries) > 0 {
		results := make([]int, len(f.docIds))
		for i, dId := range f.docIds {
			if len(results) > i {
				results[i] = dId
			}
		}
		for _, subQ := range subQueries {
			results = f.filter(results, subQ)
		}
		return results
	}
	// otherwise if no query all docs match
	return f.docIds

}

func (f *FuzzySearcher) Count() int {
	return len(f.docs)
}

func (f *FuzzySearcher) GetDocById(docId int) search.Document {
	return f.docs[docId]
}

func (f *FuzzySearcher) filter(docIds []int, subQuery search.SubQuery) []int {
	result := []int{}
	for _, docId := range docIds {
		v := f.GetDocById(docId).LoweredParsed[subQuery.Field]
		if matchesFuzzy(v, subQuery.Query) {
			result = append(result, docId)
		}
	}
	return result
}

func matchesFuzzy(text string, fuzzy string) bool {
	i := len(fuzzy)
	if i > 0 {
		fuzzyChan := func() chan rune {
			out := make(chan rune, len(fuzzy))
			for _, r := range fuzzy {
				out <- r
			}
			return out
		}()
		charToSearch := <-fuzzyChan
		for _, ch := range text {
			if ch == charToSearch {
				i--
				if i > 0 {
					charToSearch = <-fuzzyChan
				} else {
					break
				}
			}
		}
		return i == 0
	} else {
		return true
	}
}
