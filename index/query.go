package index

import "strings"

// Query
func (i IndexedLines) FilterEntries(query string) []string {
	if query != "" {
		subQueries := parseQuery(query)
		results := i.index.perfieldWord2Doc[subQueries[0].field][strings.ToLower(subQueries[0].query)]
		for _, sQ := range subQueries[1:] {
			if docs, ok := i.index.perfieldWord2Doc[sQ.field][strings.ToLower(sQ.query)]; ok {
				results = intersection(results, docs)
			} else {
				return []string{}
			}
		}
		rawStrings := make([]string, len(results))
		{
			j := 0
			for docId := range results {
				rawStrings[j] = i.docs[docId].rawText
				j++
			}
		}
		return rawStrings
	}
	// otherwise if no query all docs match
	rawStrings := make([]string, len(i.docs))
	for j, doc := range i.docs {
		rawStrings[j] = doc.rawText
	}
	return rawStrings
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

type SubQuery struct {
	field string
	query string
}

//Converts a query string to a list of queries
// they should all match (AND)
func parseQuery(query string) []SubQuery {
	subqueryStrings := strings.Split(query, " ")
	subqueries := make([]SubQuery, len(subqueryStrings))
	for i, s := range subqueryStrings {
		subQuery := SubQuery{
			field: "$",
			query: strings.ToLower(s),
		}
		fieldQuery := strings.Split(s, ":")
		// if query is like field:query
		if len(fieldQuery) > 1 {
			subQuery.field = fieldQuery[0]
			subQuery.query = strings.ToLower(fieldQuery[1])
		}
		subqueries[i] = subQuery
	}
	return subqueries
}
