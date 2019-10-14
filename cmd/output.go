package cmd

import (
	"bytes"
	"fmt"
	"text/template"
)

type renderOutput = func(map[string]string) string

func getRenderer(outputColumn string, outputTemplate string) renderOutput {
	if outputTemplate == "" || outputColumn != "$" {
		return renderByColumn(outputColumn)
	}
	return renderByTemplate(outputTemplate)
}

func renderByTemplate(outputTemplate string) renderOutput {
	tmpl, err := template.New("test").Parse(outputTemplate)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing output template: %s", err))
	}
	return func(parsedLine map[string]string) string {
		var output bytes.Buffer
		err := tmpl.Execute(&output, parsedLine)
		if err != nil {
			//FIXME: deal with errors
			panic(fmt.Sprintf("Error while executing output template: %s", err))
		}
		return output.String()
	}
}

func renderByColumn(outputColumn string) renderOutput {
	return func(parsedLine map[string]string) string {
		return parsedLine[outputColumn]
	}
}
