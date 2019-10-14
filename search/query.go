package search

import "strings"

type SubQuery struct {
	Field string
	Query string
}

//Converts a query string to a list of queries
// they should all match (AND)
func ParseQuery(query string) []SubQuery {
	subqueryStrings := strings.Split(query, " ")
	subqueries := []SubQuery{}
	for _, s := range subqueryStrings {
		subQuery := SubQuery{
			Field: "$",
			Query: strings.ToLower(s),
		}
		fieldQuery := strings.Split(s, ":")
		// if query is like field:query
		if len(fieldQuery) > 1 {
			subQuery.Field = fieldQuery[0]
			subQuery.Query = strings.ToLower(fieldQuery[1])
		}
		if subQuery.Query != "" {
			subqueries = append(subqueries, subQuery)
		}
	}
	return subqueries
}
