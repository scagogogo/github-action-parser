# Examples

This section provides comprehensive examples of using the GitHub Action Parser library in various scenarios.

## Overview

The examples are organized by functionality and complexity:

1. **[Basic Parsing](/examples/basic-parsing)** - Simple action file parsing and data access
2. **[Workflow Parsing](/examples/workflow-parsing)** - Working with GitHub workflow files
3. **[Validation](/examples/validation)** - Validating actions and workflows
4. **[Reusable Workflows](/examples/reusable-workflows)** - Working with reusable workflows
5. **[Utility Functions](/examples/utilities)** - Advanced utility functions and type conversion

## Quick Start Examples

### Parse an Action File

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
}
```

### Validate a Workflow

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    workflow, err := parser.ParseFile(".github/workflows/ci.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    validator := parser.NewValidator()
    errors := validator.Validate(workflow)
    
    if len(errors) == 0 {
        fmt.Println("Workflow is valid!")
    } else {
        fmt.Printf("Found %d validation errors\n", len(errors))
    }
}
```

### Check for Reusable Workflows

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatal(err)
    }
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            fmt.Printf("Reusable workflow: %s\n", path)
        }
    }
}
```

## Example Repository Structure

The examples assume the following repository structure:

```
my-project/
├── action.yml                    # Main action file
├── .github/
│   └── workflows/
│       ├── ci.yml                # CI workflow
│       ├── release.yml           # Release workflow
│       └── reusable-build.yml    # Reusable workflow
├── actions/
│   ├── setup/
│   │   └── action.yml            # Custom setup action
│   └── deploy/
│       └── action.yml            # Custom deploy action
└── main.go                       # Your application
```

## Sample YAML Files

### Sample Action File (action.yml)

```yaml
name: My Custom Action
description: A sample action for demonstration
author: Your Name

inputs:
  environment:
    description: Target environment
    required: true
    default: development
  version:
    description: Version to deploy
    required: false

outputs:
  deployment-url:
    description: URL of the deployed application
    value: ${{ steps.deploy.outputs.url }}

runs:
  using: composite
  steps:
    - name: Setup
      run: echo "Setting up environment ${{ inputs.environment }}"
      shell: bash
    - name: Deploy
      id: deploy
      run: |
        echo "Deploying version ${{ inputs.version }}"
        echo "url=https://app.example.com" >> $GITHUB_OUTPUT
      shell: bash

branding:
  icon: rocket
  color: blue
```

### Sample Workflow File (.github/workflows/ci.yml)

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  NODE_VERSION: 18

jobs:
  test:
    name: Run Tests
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: npm
          
      - name: Install dependencies
        run: npm ci
        
      - name: Run tests
        run: npm test
        
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        if: success()

  build:
    name: Build Application
    runs-on: ubuntu-latest
    needs: test
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Build
        run: |
          echo "Building application..."
          mkdir -p dist
          echo "Built successfully" > dist/app.txt
          
      - name: Upload artifacts
        uses: actions/upload-artifact@v3
        with:
          name: build-artifacts
          path: dist/
```

### Sample Reusable Workflow (.github/workflows/reusable-build.yml)

```yaml
name: Reusable Build Workflow

on:
  workflow_call:
    inputs:
      environment:
        description: Target environment
        required: true
        type: string
      node-version:
        description: Node.js version to use
        required: false
        type: string
        default: '18'
    outputs:
      build-version:
        description: Version of the built application
        value: ${{ jobs.build.outputs.version }}
    secrets:
      NPM_TOKEN:
        description: NPM authentication token
        required: true

jobs:
  build:
    name: Build for ${{ inputs.environment }}
    runs-on: ubuntu-latest
    
    outputs:
      version: ${{ steps.version.outputs.version }}
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        
      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}
          registry-url: https://registry.npmjs.org/
          
      - name: Get version
        id: version
        run: |
          VERSION=$(date +%Y%m%d)-$(git rev-parse --short HEAD)
          echo "version=$VERSION" >> $GITHUB_OUTPUT
          
      - name: Install dependencies
        run: npm ci
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
          
      - name: Build for environment
        run: |
          echo "Building for ${{ inputs.environment }}"
          npm run build:${{ inputs.environment }}
```

## Running the Examples

Each example can be run independently. Make sure you have:

1. Go 1.20 or later installed
2. The GitHub Action Parser library installed:
   ```bash
   go get github.com/scagogogo/github-action-parser
   ```

3. Sample YAML files in your project directory

Then navigate to each example section for detailed instructions and code.
