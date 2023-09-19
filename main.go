package main

import (
	"embed"
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

var (
	k = koanf.New(".")
	//go:embed templates
	templatesFS embed.FS
)

func main() {
	app := &cli.App{
		Name:        "env generator",
		Description: "A simple tool that generates .env files and .ts interfaces based on json schema",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "schema",
				Value: "./env-schema.json",
				Usage: "schema where env variables are defined",
			},
			&cli.StringFlag{
				Name:  "out-env",
				Value: ".env",
				Usage: "path to create .env file",
			},
			&cli.StringFlag{
				Name:  "out-ts",
				Value: "env-variables.ts",
				Usage: "path to create ts interface file",
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
	envFile, err := os.OpenFile(ctx.String("out-env"), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("Error opening env file: %s", err)
	}
	defer envFile.Close()

	envInterfaceFile, err := os.OpenFile(ctx.String("out-ts"), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, os.ModePerm)
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

	result, err := renderInterface(schema)
	if err != nil {
		log.Fatal(err)
	}
	envInterfaceFile.WriteString(result)

	result, err = renderEnvFile(schema)
	if err != nil {
		log.Fatal(err)
	}
	envFile.WriteString(result)

	return nil
}
