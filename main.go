package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
)

type DataType string

type VarDef struct {
	Name     string   `koanf:"name"`
	Type     DataType `koanf:"type"`
	Required bool     `koanf:"required"`
	Default  string   `koanf:"default"`
}

type Schema struct {
	Vars []VarDef `koanf:"vars"`
}

const (
	INT    = "int"
	FLOAT  = "float"
	STRING = "string"
	BOOL   = "bool"
)

var k = koanf.New(".")

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "schema",
				Value: "./env-schema.json",
				Usage: "schema where env variables are defined",
			},
		},
		Action: convertSchemaToEnvFile,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func convertSchemaToEnvFile(ctx *cli.Context) error {
	envFile, err := os.OpenFile(".env", os.O_TRUNC | os.O_WRONLY | os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error opening env file: %s", err)
	}
	defer envFile.Close()

	envInterfaceFile, err := os.OpenFile("env-variables.ts", os.O_TRUNC | os.O_WRONLY | os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error opening typescript interface file: %s", err)
	}
	defer envInterfaceFile.Close()

	schemaPath := ctx.String("schema")
	err = k.Load(file.Provider(schemaPath), json.Parser())
	if err != nil {
		log.Fatalf("Error loading config file: %s\n", err)
	}

	schema := Schema{}
	k.Unmarshal("", &schema)

	envInterfaceFile.WriteString("export interface EnvVariables {\n")
	for _, v := range schema.Vars {
		envFile.WriteString(fmt.Sprintf("%s=\"%s\"\n", v.Name, v.Default))
		dataType, err := getTypescriptType(v.Type)
		if err != nil {
			log.Fatal(err)
		}
		envInterfaceFile.WriteString(fmt.Sprintf("\t%s: %s;\n", v.Name, dataType))
	}
	envInterfaceFile.WriteString("}\n")


	return nil
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
