package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func renderInterface(schema Schema) (string, error) {
	resultBuffer := bytes.Buffer{}
	result := ""
	templateContent, err := templatesFS.ReadFile("templates/interface-ts.tmpl")
	if err != nil {
		return "", fmt.Errorf("Error loading typescript template: %s", err)
	}

	funcMap := template.FuncMap{
		"getTypescriptType": getTypescriptType,
	}
	tsTemplate := template.New("ts interface").Funcs(funcMap)
	tmpl, err := tsTemplate.Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %s", err)
	}

	tmpl.Execute(&resultBuffer, schema)

	result = resultBuffer.String()
	return result, nil
}

func getTypescriptType(dataType DataType) (string, error) {
	if (dataType == INT) {
		return "number", nil
	}
	if (dataType == FLOAT) {
		return "number", nil
	}
	if (dataType == BOOL) {
		return "boolean", nil
	}
	if (dataType == STRING) {
		return "string", nil
	}

	return "", fmt.Errorf("Invalid data type")
}

