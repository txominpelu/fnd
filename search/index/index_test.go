package index

import (
	"reflect"
	"testing"

	"github.com/txominpelu/fnd/search"
)

func TestBasicQuery(t *testing.T) {
	indexedLines := NewIndexedLines(CommandLineTokenizer())
	lines := []string{
		"hello world",
		"this is the best WOrld",
		"this won't match",
	}
	for _, l := range lines {
		indexedLines.AddDocument(search.ParseLine(search.PlainTextParser(), l))
	}
	expected := []string{
		"hello world",
		"this is the best WOrld",
	}
	gotDocs := search.SortDocuments(
		indexedLines.FilterEntries(search.ParseQuery("world")),
		&indexedLines,
	)
	got := make([]string, len(gotDocs))
	for i, d := range gotDocs {
		got[i] = d.RawText
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected: '%v' but got '%v'", expected, got)
	}

}

func TestTokenizer(t *testing.T) {
	tokenizer := CommandLineTokenizer()
	got := tokenizer("/hello.world tal/")
	expected := []string{
		"hello",
		"world",
		"tal",
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected: '%v' but got '%v'", expected, got)
	}

}
