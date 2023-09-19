package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func renderEnvFile(schema Schema) (string, error) {
	resultBuffer := bytes.Buffer{}
	result := ""
	templateContent, err := templatesFS.ReadFile("templates/env-file.tmpl")
	if err != nil {
		return "", fmt.Errorf("Error loading env file template: %s", err)
	}

	tsTemplate := template.New("ts interface")
	tmpl, err := tsTemplate.Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %s", err)
	}

	tmpl.Execute(&resultBuffer, schema)

	result = resultBuffer.String()
	return result, nil
}

