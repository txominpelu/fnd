package index

import (
	"reflect"
	"testing"
)

func TestBasicQuery(t *testing.T) {
	lines := NewIndexedLines(CommandLineTokenizer(), PlainTextParser())
	lines.AddLine("hello world")
	lines.AddLine("this is the best WOrld")
	lines.AddLine("this won't match")
	expected := []string{
		"hello world",
		"this is the best WOrld",
	}

	gotDocs := lines.FilterEntries("world")
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
