---
layout: home

hero:
  name: "GitHub Action Parser"
  text: "Go Library for GitHub Actions"
  tagline: "Parse, validate and process GitHub Action YAML files with ease"
  image:
    src: /logo.svg
    alt: GitHub Action Parser
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: API Reference
      link: /api/
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/github-action-parser

features:
  - icon: 📄
    title: Parse Action Files
    details: Parse GitHub Action YAML files (action.yml/action.yaml) with full type safety and validation.
  - icon: ⚙️
    title: Workflow Support
    details: Parse GitHub Workflow files (.github/workflows/*.yml) and extract job definitions, steps, and triggers.
  - icon: ✅
    title: Validation
    details: Validate actions and workflows according to GitHub's specifications with detailed error reporting.
  - icon: 🔄
    title: Reusable Workflows
    details: Detect and process reusable workflows with input/output parameter extraction.
  - icon: 🛠️
    title: Utility Functions
    details: Type conversion and data processing utilities for working with YAML data structures.
  - icon: 📁
    title: Batch Processing
    details: Parse all Action and Workflow files in a directory recursively with a single function call.
---

## Quick Start

Install the library:

```bash
go get github.com/scagogogo/github-action-parser
```

Parse an action file:

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
}
```

## Supported GitHub Action Features

- **Action Metadata**: Name, description, author information
- **Input Parameters**: With validation requirements and default values
- **Output Parameters**: With descriptions and values
- **Docker Actions**: Container-based actions
- **JavaScript Actions**: Node.js 16/20 actions
- **Composite Actions**: Multi-step composite actions
- **Workflow Jobs**: Job definitions and dependencies
- **Workflow Triggers**: Event-based triggers
- **Reusable Workflows**: Callable workflows with parameters
- **Secrets Handling**: Reusable workflow secrets processing

## Why GitHub Action Parser?

- **Type Safe**: Full Go type definitions for all GitHub Action structures
- **Comprehensive**: Supports all GitHub Action and Workflow features
- **Validated**: Built-in validation according to GitHub specifications
- **Well Tested**: 98.9% test coverage with comprehensive test suite
- **Easy to Use**: Simple API with clear documentation and examples
- **Performance**: Optimized for parsing large numbers of files efficiently
