package index

import (
	"strings"
)

type Tokenizer = func(string) []string

type Document struct {
	index   int
	rawText string
}

type Index struct {
	word2doc map[string][]int
}

type IndexedLines struct {
	lines     []string
	count     int
	index     Index
	docs      []Document
	tokenizer Tokenizer
}

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

func NewIndexedLines(tokenizer Tokenizer) IndexedLines {
	i := IndexedLines{}
	if i.lines == nil {
		i.lines = []string{}
	}
	if i.index.word2doc == nil {
		i.index = Index{word2doc: map[string][]int{}}
	}
	i.tokenizer = tokenizer
	return i
}

func (i *IndexedLines) AddLine(line string) {
	docId := i.count // docId = index in array
	for _, word := range i.tokenizer(line) {
		toLower := strings.ToLower(word)
		i.index.addWord(toLower, docId)
		for _, nGram := range findNgrams(toLower, 0, 0) {
			i.index.addWord(nGram, docId)
		}
	}
	i.lines = append(i.lines, line)
	i.docs = append(i.docs, Document{rawText: line, index: docId})
	i.count++
}

func findNgrams(toLower string, min int, max int) []string {
	return []string{}
}

func (index *Index) addWord(word string, docId int) {
	if _, ok := index.word2doc[word]; !ok {
		index.word2doc[word] = []int{}
	}
	index.word2doc[word] = append(index.word2doc[word], docId)
}

func (i IndexedLines) FilterEntries(query string) []string {
	if query != "" {
		if docs, ok := i.index.word2doc[strings.ToLower(query)]; ok {
			rawStrings := make([]string, len(docs))
			for j, docId := range docs {
				rawStrings[j] = i.docs[docId].rawText
			}
			return rawStrings
		}
		return []string{}
	}
	rawStrings := make([]string, len(i.docs))
	for j, doc := range i.docs {
		rawStrings[j] = doc.rawText
	}
	return rawStrings
}

func (i IndexedLines) Count() int {
	return i.count
}
