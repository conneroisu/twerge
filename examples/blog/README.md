# Blog Example

This example demonstrates a blog UI built with twerge and Tailwind CSS.

## Features
- Article layout with proper typography
- Responsive design for all screen sizes
- Author information section and tag system
- Custom styling for code blocks
- Uses twerge for optimized Tailwind class management

## Setup and Running

1. Generate the templ components:
```sh
templ generate ./views
```

2. Run the code generation:
```sh
go run gen.go
```

3. Run the server:
```sh
go run main.go
```

4. Open your browser and navigate to http://localhost:8080

## Requirements
- Go 1.24+
- templ v0.3.857+
- TailwindCSS CLI

## Note on Code Example Problems

If you encounter issues with code examples in the blog content, particularly with the twerge.If function, this is intentional to show the syntax while avoiding templ parsing errors. In a real application, you would use the actual function with proper syntax.