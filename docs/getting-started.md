# Getting Started

## Installation

Install the GitHub Action Parser library using Go modules:

```bash
go get github.com/scagogogo/github-action-parser
```

## Basic Usage

Import the parser package in your Go code:

```go
import "github.com/scagogogo/github-action-parser/pkg/parser"
```

### Parse an Action File

```go
package main

import (
    "fmt"
    "os"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse action.yml file
    action, err := parser.ParseFile("path/to/action.yml")
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        os.Exit(1)
    }

    // Access action metadata
    fmt.Printf("Action Name: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
    
    // Access input parameters
    for name, input := range action.Inputs {
        fmt.Printf("Input %s: %s (Required: %t)\n", 
            name, input.Description, input.Required)
    }
    
    // Access output parameters
    for name, output := range action.Outputs {
        fmt.Printf("Output %s: %s\n", name, output.Description)
    }
}
```

### Parse a Workflow File

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse workflow.yml file
    workflow, err := parser.ParseFile("path/to/.github/workflows/ci.yml")
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        return
    }
    
    // Access workflow jobs
    for jobId, job := range workflow.Jobs {
        fmt.Printf("Job: %s (%s)\n", jobId, job.Name)
        
        // Access job steps
        for i, step := range job.Steps {
            if step.Name != "" {
                fmt.Printf("  Step %d: %s\n", i+1, step.Name)
            } else if step.Run != "" {
                fmt.Printf("  Step %d: Run command\n", i+1)
            } else if step.Uses != "" {
                fmt.Printf("  Step %d: Uses %s\n", i+1, step.Uses)
            }
        }
    }
}
```

### Validate Files

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse any GitHub Action/Workflow file
    action, err := parser.ParseFile("path/to/file.yml")
    if err != nil {
        fmt.Printf("Error parsing file: %v\n", err)
        return
    }
    
    // Create validator and validate
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) > 0 {
        fmt.Println("Validation errors:")
        for _, err := range errors {
            fmt.Printf("- %s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Println("The file is valid!")
    }
}
```

### Parse Directory

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all workflow files in a directory
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        fmt.Printf("Error parsing workflows: %v\n", err)
        return
    }
    
    for path, workflow := range workflows {
        fmt.Printf("Workflow: %s\n", path)
        fmt.Printf("  Name: %s\n", workflow.Name)
        fmt.Printf("  Jobs: %d\n", len(workflow.Jobs))
    }
}
```

## Next Steps

- Explore the [API Reference](/api/) for detailed documentation
- Check out [Examples](/examples/) for more use cases
- Learn about [validation features](/api/validation)
- Discover [utility functions](/api/utilities) for advanced usage
