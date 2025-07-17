# GitHub Action Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/github-action-parser.svg)](https://pkg.go.dev/github.com/scagogogo/github-action-parser)
[![Go CI](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml)
[![Documentation](https://github.com/scagogogo/github-action-parser/actions/workflows/docs.yml/badge.svg)](https://scagogogo.github.io/github-action-parser/)
[![Coverage](https://img.shields.io/badge/coverage-98.9%25-brightgreen)](https://github.com/scagogogo/github-action-parser)

A Go library for parsing, validating and processing GitHub Action YAML files.

---

## 📚 Documentation

**🌐 [Complete Documentation Website](https://scagogogo.github.io/github-action-parser/)**

- 📖 [English Documentation](https://scagogogo.github.io/github-action-parser/)
- 🇨🇳 [Chinese Documentation](https://scagogogo.github.io/github-action-parser/zh/)

The documentation includes:
- 🚀 **Getting Started Guide** - Quick setup and basic usage
- 📋 **API Reference** - Complete API documentation with examples
- 💡 **Examples** - Practical code examples and use cases
- ✅ **Validation Guide** - How to validate actions and workflows
- 🔄 **Reusable Workflows** - Working with reusable workflows

---

## Features

- Parse GitHub Action YAML files (`action.yml`/`action.yaml`)
- Parse GitHub Workflow files (`.github/workflows/*.yml`)
- Validate actions and workflows according to GitHub's specifications
- Support for composite, Docker, and JavaScript actions
- Extract metadata, inputs, outputs, jobs, and step information
- Detect and process reusable workflows
- Type conversion and data processing utilities
- Batch parsing of all Action and Workflow files in directories

## Installation

```bash
go get github.com/scagogogo/github-action-parser
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse an action file
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
    
    // Validate the action
    validator := parser.NewValidator()
    if errors := validator.Validate(action); len(errors) > 0 {
        fmt.Printf("Validation errors: %v\n", errors)
    } else {
        fmt.Println("Action is valid!")
    }
}
```

## Documentation

- **📖 [Full Documentation](https://scagogogo.github.io/github-action-parser/)** - Complete API reference and guides
- **🚀 [Getting Started](https://scagogogo.github.io/github-action-parser/getting-started)** - Quick start guide
- **📚 [API Reference](https://scagogogo.github.io/github-action-parser/api/)** - Detailed API documentation
- **💡 [Examples](https://scagogogo.github.io/github-action-parser/examples/)** - Code examples and use cases

### Chinese Documentation

- **📖 [Complete Documentation](https://scagogogo.github.io/github-action-parser/zh/)** - Full API reference and guides
- **🚀 [Getting Started](https://scagogogo.github.io/github-action-parser/zh/getting-started)** - Quick start guide
- **📚 [API Reference](https://scagogogo.github.io/github-action-parser/zh/api/)** - Detailed API documentation
- **💡 [Examples](https://scagogogo.github.io/github-action-parser/zh/examples/)** - Code examples and use cases

## Examples

### Parse Action File

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Action: %s\n", action.Name)
for name, input := range action.Inputs {
    fmt.Printf("Input %s: required=%t\n", name, input.Required)
}
```

### Parse Workflow File

```go
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
if err != nil {
    log.Fatal(err)
}

for jobID, job := range workflow.Jobs {
    fmt.Printf("Job %s has %d steps\n", jobID, len(job.Steps))
}
```

### Validate Files

```go
validator := parser.NewValidator()
errors := validator.Validate(action)

if len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("Error in %s: %s\n", err.Field, err.Message)
    }
}
```

### Parse Directory

```go
actions, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d workflow files\n", len(actions))
```

### Check Reusable Workflows

```go
if parser.IsReusableWorkflow(workflow) {
    inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
    fmt.Printf("Reusable workflow with %d inputs\n", len(inputs))
}
```

## Supported GitHub Action Features

- ✅ Action metadata (name, description, author)
- ✅ Input parameters with validation requirements
- ✅ Output parameters with descriptions and values
- ✅ Docker container actions
- ✅ JavaScript actions (Node.js 16/20)
- ✅ Composite actions
- ✅ Workflow job definitions
- ✅ Workflow triggers (events)
- ✅ Reusable workflows
- ✅ Job and step dependencies
- ✅ Secrets handling for reusable workflows

## Testing

The library has comprehensive test coverage (98.9%) and includes:

- Unit tests for all functions
- Integration tests with real GitHub Action files
- Validation tests for GitHub specifications
- Performance benchmarks

```bash
go test ./pkg/parser/
go test -bench=. ./pkg/parser/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Links

- **Documentation**: https://scagogogo.github.io/github-action-parser/
- **Go Package**: https://pkg.go.dev/github.com/scagogogo/github-action-parser
- **GitHub Repository**: https://github.com/scagogogo/github-action-parser
- **Issues**: https://github.com/scagogogo/github-action-parser/issues
