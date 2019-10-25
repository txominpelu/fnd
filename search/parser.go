package search

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/txominpelu/fnd/log"
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

func TabularParser(headers []string, delimiter rune) Parser {
	parse := func(line string) map[string]interface{} {
		columns := strings.FieldsFunc(line, func(r rune) bool { return r == delimiter })
		result := map[string]interface{}{}
		for i := 0; i < len(headers) && i < len(columns); i++ {
			if i == len(headers)-1 {
				//for the last column take everything that's left
				result[headers[i]] = strings.TrimSpace(strings.Join(columns[i:], string([]rune{delimiter})))
			} else {
				result[headers[i]] = strings.TrimSpace(columns[i])
			}
		}
		return result
	}
	return Parser{
		headers: headers,
		parse:   parse,
	}
}

func FormatNameToParser(format string, firstline string, headers []string, logger *log.StandardLogger, delimiter rune) Parser {
	var p Parser
	switch format {
	case "plain":
		p = PlainTextParser()
	case "json":
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(firstline), &m)
		logger.CheckError(
			err,
			fmt.Sprintf("when parsing first line '' as json", firstline),
		)
		headers := []string{}
		for k, _ := range m {
			headers = append(headers, k)
		}
		sort.StringSlice(headers).Sort()
		p = JsonParser(headers, logger)
	case "tabular":
		headers := strings.FieldsFunc(firstline, func(r rune) bool { return r == delimiter })
		trimmedHeaders := []string{}
		for _, h := range headers {
			trimmedHeaders = append(trimmedHeaders, strings.TrimSpace(h))
		}
		p = TabularParser(trimmedHeaders, delimiter)
	default:
		err := fmt.Errorf("pass invalid --line_format '%s' should be one of (plain/tabular/json) \n", format)
		logger.CheckError(err, "")
	}
	if len(headers) > 0 {
		p.headers = headers
	}
	return p
}

func PlainTextParser() Parser {
	parse := func(line string) map[string]interface{} { return map[string]interface{}{"$": line} }
	return Parser{
		headers: []string{"$"},
		parse:   parse,
	}
}

func JsonParser(headers []string, logger *log.StandardLogger) Parser {
	parse := func(line string) map[string]interface{} {
		m := map[string]interface{}{}
		err := json.Unmarshal([]byte(line), &m)
		logger.WarnIfErr(err, "when parsing line as json")
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
	loweredParsed["$"] = strings.ToLower(line)
	return Document{
		RawText: line,
		//only one level key, value
		ParsedLine:    parsedLine,    //for display and return of command
		LoweredParsed: loweredParsed, //for case insensitive search
	}
}
