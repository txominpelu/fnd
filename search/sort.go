package search

import "sort"

type Compare = func(int, int) bool

type docSorter struct {
	docIds []int
	by     func(d1, d2 int) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *docSorter) Len() int {
	return len(s.docIds)
}

// Swap is part of sort.Interface.
func (s *docSorter) Swap(i, j int) {
	s.docIds[i], s.docIds[j] = s.docIds[j], s.docIds[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *docSorter) Less(i, j int) bool {
	return s.by(s.docIds[i], s.docIds[j])
}

func Sort(docIds []int, by func(int, int) bool) {
	dSorter := &docSorter{
		docIds: docIds,
		by:     by,
	}
	sort.Sort(dSorter)
}

func SortDocuments(docIds []int, searcher TextSearcher, by func(int, int) bool) []Document {
	Sort(docIds, by)
	docs := make([]Document, len(docIds))
	for j, docId := range docIds {
		docs[j] = searcher.GetDocById(docId)
	}
	return docs
}
