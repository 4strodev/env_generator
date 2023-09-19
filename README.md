# .ENV generator
This is a tool that I use to generate .env files and .ts interfaces for my backend projects

## Installation
```sh
go install github.com/4strodev/env_generator
```

## Usage
```sh
# Explicit
env_generator --schema="env-schema.json" --out-env=".env" --out-ts="env-variables.ts"

# Implicit
env_generator
```
