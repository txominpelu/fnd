package search

import (
	"encoding/json"
	"fmt"
	"strings"
)

//Parser takes a line and creates a table (map from field to value)
type Parser struct {
	headers []string
	parse   func(string) map[string]interface{}
}

func (p Parser) Headers() []string {
	return p.headers
}

func (p Parser) Parse() func(string) map[string]interface{} {
	return p.parse
}

func TabularParser(headers []string) Parser {
	parse := func(line string) map[string]interface{} {
		columns := strings.Fields(line)
		result := map[string]interface{}{}
		for i := 0; i < len(headers) && i < len(columns); i++ {
			if i == len(headers)-1 {
				//for the last column take everything that's left
				result[headers[i]] = strings.Join(columns[i:], " ")
			} else {
				result[headers[i]] = columns[i]
			}
		}
		return result
	}
	return Parser{
		headers: headers,
		parse:   parse,
	}
}

func FormatNameToParser(format string, firstline string) Parser {
	switch format {
	case "plain":
		return PlainTextParser()
	case "json":
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(firstline), &m)
		//FIXME: error handle
		if err != nil {
			panic("Failed parsing line as json")
		}
		headers := []string{}
		for k, _ := range m {
			headers = append(headers, k)
		}
		return JsonParser(headers)
	case "tabular":
		headers := strings.Fields(firstline)
		return TabularParser(headers)
	default:
		//FIXME: decide once and for all how to handle errors
		panic(fmt.Sprintf("format '%s' is not valid", format))
	}

}

func PlainTextParser() Parser {
	parse := func(line string) map[string]interface{} { return map[string]interface{}{"$": line} }
	return Parser{
		headers: []string{"$"},
		parse:   parse,
	}
}

func JsonParser(headers []string) Parser {
	parse := func(line string) map[string]interface{} {
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(line), &m)
		//FIXME: if line cannot be parsed, just ignore, maybe log
		if err != nil {
			panic("Failed parsing line as json")
		}
		return m
	}
	return Parser{
		headers: headers,
		parse:   parse,
	}
}

func ParseLine(parser Parser, line string) Document {
	m := parser.Parse()(line)
	parsedLine := map[string]string{}
	loweredParsed := map[string]string{}
	for k, interf := range m {
		switch v := interf.(type) {
		case int:
			parsedLine[k] = fmt.Sprintf("%d", v)
		case string:
			parsedLine[k] = v
		default:
			parsedLine[k] = fmt.Sprintf("%v", v)
		}
		loweredParsed[k] = strings.ToLower(parsedLine[k])
	}
	parsedLine["$"] = line
	return Document{
		RawText: line,
		//only one level key, value
		ParsedLine:    parsedLine,    //for display and return of command
		LoweredParsed: loweredParsed, //for case insensitive search
	}
}
