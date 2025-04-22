# twerge Examples

This directory contains examples demonstrating the use of twerge with Tailwind CSS and Go templ.

## Available Examples

### Simple Example
A basic example showing a landing page with header, main content, and footer.

### Dashboard Example
A modern analytics dashboard UI with metrics cards and data tables.

### Blog Example
A blog interface with article display and typography.

## Setup Requirements

All examples require:
- Go 1.24+
- templ v0.3.857+
- TailwindCSS CLI

## Running an Example

1. Navigate to the example directory:
```sh
cd [example-name]
```

2. Generate the templ components:
```sh
templ generate ./views
```

3. Run the code generation:
```sh
go run gen.go
```

4. Run the server:
```sh
go run main.go
```

5. Open your browser and navigate to http://localhost:8080

## Notes

- Each example follows the same structure with `views/`, `classes/`, and `_static/` directories
- The `gen.go` file handles twerge code generation and TailwindCSS processing
- `classes/classes.go` contains the generated class mappings
- `input.css` and `tailwind.config.js` manage TailwindCSS configuration