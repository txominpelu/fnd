package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/txominpelu/fnd/log"
)

type renderOutput = func(map[string]string) string

func getRenderer(outputColumn string, outputTemplate string, logger *log.StandardLogger) renderOutput {
	if outputTemplate == "" || outputColumn != "$" {
		return renderByColumn(outputColumn)
	}
	return renderByTemplate(outputTemplate, logger)
}

func renderByTemplate(outputTemplate string, logger *log.StandardLogger) renderOutput {
	tmpl, err := template.New("test").Parse(outputTemplate)
	logger.CheckError(err, fmt.Sprintf("while parsing output template: %s"))
	return func(parsedLine map[string]string) string {
		var output bytes.Buffer
		err := tmpl.Execute(&output, parsedLine)
		logger.CheckError(err, fmt.Sprintf("while executing output template: %s"))
		return output.String()
	}
}

func renderByColumn(outputColumn string) renderOutput {
	return func(parsedLine map[string]string) string {
		return parsedLine[outputColumn]
	}
}
