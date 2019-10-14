package fuzzy

import (
	"reflect"
	"testing"

	"github.com/txominpelu/fnd/search"
)

func TestFuzzy(t *testing.T) {
	fuzzySearcher := NewFuzzySearcher()
	lines := []string{
		"abcd",
		"aaaacbddddd",
		"tttttbxxxxxc",
	}
	for _, l := range lines {
		fuzzySearcher.AddDocument(search.ParseLine(search.PlainTextParser(), l))
	}
	expected := []string{
		"abcd",
		"tttttbxxxxxc",
	}
	gotDocs := search.SortDocuments(
		fuzzySearcher.FilterEntries(search.ParseQuery("bc")),
		&fuzzySearcher,
	)
	got := make([]string, len(gotDocs))
	for i, d := range gotDocs {
		got[i] = d.RawText
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected: '%v' but got '%v'", expected, got)
	}
}

func TestUpperCase(t *testing.T) {
	fuzzySearcher := NewFuzzySearcher()
	lines := []string{
		"aBCd",
		"aaaacbddddd",
		"tttttbxxxxxc",
	}
	for _, l := range lines {
		fuzzySearcher.AddDocument(search.ParseLine(search.PlainTextParser(), l))
	}
	expected := []string{
		"aBCd",
		"tttttbxxxxxc",
	}
	gotDocs := search.SortDocuments(
		fuzzySearcher.FilterEntries(search.ParseQuery("bc")),
		&fuzzySearcher,
	)
	got := make([]string, len(gotDocs))
	for i, d := range gotDocs {
		got[i] = d.RawText
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected: '%v' but got '%v'", expected, got)
	}
}
