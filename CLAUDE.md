# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Run tests: `go test ./...`
- Run single test: `go test -run TestName ./path/to/package`
- Format code: `go fmt ./...`
- Generate code: `go generate ./...`
- Run TailwindCSS: `tailwindcss -i input.css -o _static/dist/styles.css`

## Code Style Guidelines
- Imports: Group standard library first, then external, then internal packages
- Error handling: Always check errors and use descriptive error messages
- Types: Use meaningful type names and favor composition over inheritance
- Naming: Follow Go conventions (CamelCase for exported, camelCase for unexported)
- Documentation: Document all exported functions, types, and constants
- File structure: Follow the project's existing organization pattern
- Testing: Write tests for all new functionality