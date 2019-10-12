package index

import (
	"fmt"
	"strings"

	"github.com/txominpelu/fnd/search"
)

// Tokenizers
type Tokenizer = func(string) []string

func CommandLineTokenizer() Tokenizer {
	return func(s string) []string {
		results := []string{}
		builder := strings.Builder{}
		for _, r := range s {
			if r == ' ' || r == rune('/') || r == '.' {
				if builder.Len() > 0 {
					results = append(results, builder.String())
					builder.Reset()
				}
			} else {
				builder.WriteRune(r)
			}
		}
		if builder.Len() > 0 {
			results = append(results, builder.String())
		}
		return results
	}
}

// Glossary
// Word2Doc: mapping from word to docs containing it
// PerFieldWord2Doc: mapping from field to word2doc

type Word2Doc = map[string]map[int]bool

type PerFieldWord2Doc struct {
	perfieldWord2Doc map[string]Word2Doc
}

type IndexedLines struct {
	lines     []string
	count     int
	index     PerFieldWord2Doc
	docs      []search.Document
	tokenizer Tokenizer
	parser    search.Parser
}

func NewIndexedLines(tokenizer Tokenizer, parser search.Parser) IndexedLines {
	i := IndexedLines{}
	if i.lines == nil {
		i.lines = []string{}
	}
	if i.index.perfieldWord2Doc == nil {
		i.index = PerFieldWord2Doc{perfieldWord2Doc: map[string]Word2Doc{}}
	}
	i.tokenizer = tokenizer
	i.parser = parser
	return i
}

func (i *IndexedLines) AddLine(line string) {
	docId := i.count // docId = index in array
	i.lines = append(i.lines, line)
	m := i.parser.Parse()(line)
	parsedLine := map[string]string{}
	for k, interf := range m {
		switch v := interf.(type) {
		case int:
			parsedLine[k] = fmt.Sprintf("%d", v)
		case string:
			parsedLine[k] = v
		default:
			parsedLine[k] = fmt.Sprintf("%v", v)
		}

	}
	parsedLine["$"] = line
	i.docs = append(i.docs, search.Document{
		RawText:    line,
		ID:         docId,
		ParsedLine: parsedLine,
	})
	index(m, &(i.index.perfieldWord2Doc), docId, i.tokenizer, i.parser)
	i.count++
}

// it expects one json object
// for each key it indexes the values:
//   if value is string:
//      indexLine(value)
//   else if value is array:
//      for elem in array:
//        indexElem(elem)
//   else:
//      ignore element -> it ignores null, numbers and nested objects
func index(parsedLine map[string]interface{}, perfield *map[string]Word2Doc, docId int, tokenizer Tokenizer, parser search.Parser) {
	for key, val := range parsedLine {
		switch val.(type) {
		case []interface{}:
			arr := val.([]interface{})
			for _, r := range arr {
				indexElem(perfield, key, r, docId, tokenizer)
				indexElem(perfield, "$", r, docId, tokenizer)
			}
		default:
			indexElem(perfield, key, val, docId, tokenizer)
			indexElem(perfield, "$", val, docId, tokenizer)

		}
	}
}

func indexElem(perfield *map[string]Word2Doc, key string, val interface{}, docId int, tokenizer Tokenizer) {
	switch val.(type) {
	case string:
		indexLine(perfield, key, val.(string), docId, tokenizer)
	}
}

func indexLine(perfield *map[string]Word2Doc, field string, line string, docId int, tokenizer Tokenizer) {
	for _, word := range tokenizer(line) {
		addWord(perfield, field, word, docId)
		//FIXME: add ngrams
		for _, ngram := range findNgrams(word, 0, 0) {
			addWord(perfield, field, ngram, docId)
		}
	}
}

func addWord(perfieldPointer *map[string]Word2Doc, field string, word string, docId int) {
	perfield := *perfieldPointer
	word = strings.ToLower(word)
	if _, ok := perfield[field]; !ok {
		perfield[field] = map[string]map[int]bool{}
	}
	if _, ok2 := perfield[field][word]; !ok2 {
		perfield[field][word] = map[int]bool{}
	}
	perfield[field][word][docId] = true
}

func findNgrams(toLower string, min int, max int) []string {
	return []string{}
}

func (i IndexedLines) Count() int {
	return i.count
}