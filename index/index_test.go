package index

import (
	"reflect"
	"testing"
)

func TestBasicQuery(t *testing.T) {
	lines := NewIndexedLines(CommandLineTokenizer())
	lines.AddLine("hello world")
	lines.AddLine("this is the best WOrld")
	lines.AddLine("this won't match")
	expected := []string{
		"hello world",
		"this is the best WOrld",
	}

	got := lines.FilterEntries("world")
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