package search

type TextSearcher interface {
	AddDocument(document Document)
	// Given a bunch of subqueries returns the docIds that match them
	FilterEntries(subQueries []SubQuery) []int
	GetDocById(docId int) Document
	Count() int
}

type Document struct {
	RawText    string
	ParsedLine map[string]string
}
