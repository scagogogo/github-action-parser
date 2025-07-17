# Basic Parsing

This example demonstrates the fundamental parsing capabilities of the GitHub Action Parser library.

## Parse Action File

The most basic operation is parsing a GitHub Action file:

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse action.yml file
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatalf("Failed to parse action: %v", err)
    }

    // Display basic information
    fmt.Printf("Action Name: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
    fmt.Printf("Author: %s\n", action.Author)
}
```

## Access Input Parameters

```go
// Access input parameters
fmt.Println("\nInputs:")
for name, input := range action.Inputs {
    fmt.Printf("  %s:\n", name)
    fmt.Printf("    Description: %s\n", input.Description)
    fmt.Printf("    Required: %t\n", input.Required)
    if input.Default != "" {
        fmt.Printf("    Default: %s\n", input.Default)
    }
    if input.Deprecated {
        fmt.Printf("    Deprecated: true\n")
    }
}
```

## Access Output Parameters

```go
// Access output parameters
fmt.Println("\nOutputs:")
for name, output := range action.Outputs {
    fmt.Printf("  %s:\n", name)
    fmt.Printf("    Description: %s\n", output.Description)
    if output.Value != "" {
        fmt.Printf("    Value: %s\n", output.Value)
    }
}
```

## Check Action Type

```go
// Check what type of action this is
fmt.Printf("\nAction Type: %s\n", action.Runs.Using)

switch action.Runs.Using {
case "composite":
    fmt.Printf("Composite action with %d steps\n", len(action.Runs.Steps))
case "docker":
    fmt.Printf("Docker action using image: %s\n", action.Runs.Image)
case "node16", "node20":
    fmt.Printf("JavaScript action with main: %s\n", action.Runs.Main)
}
```

## Access Branding Information

```go
// Access branding information
if action.Branding.Icon != "" || action.Branding.Color != "" {
    fmt.Println("\nBranding:")
    if action.Branding.Icon != "" {
        fmt.Printf("  Icon: %s\n", action.Branding.Icon)
    }
    if action.Branding.Color != "" {
        fmt.Printf("  Color: %s\n", action.Branding.Color)
    }
}
```

## Complete Example

Here's a complete example that demonstrates all basic parsing features:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <action.yml>")
    }
    
    actionFile := os.Args[1]
    
    // Parse the action file
    action, err := parser.ParseFile(actionFile)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", actionFile, err)
    }
    
    // Display comprehensive information
    displayActionInfo(action)
}

func displayActionInfo(action *parser.ActionFile) {
    fmt.Printf("=== Action Information ===\n")
    fmt.Printf("Name: %s\n", action.Name)
    fmt.Printf("Description: %s\n", action.Description)
    
    if action.Author != "" {
        fmt.Printf("Author: %s\n", action.Author)
    }
    
    // Display inputs
    if len(action.Inputs) > 0 {
        fmt.Printf("\n=== Inputs (%d) ===\n", len(action.Inputs))
        for name, input := range action.Inputs {
            fmt.Printf("• %s", name)
            if input.Required {
                fmt.Printf(" (required)")
            }
            fmt.Printf("\n  %s\n", input.Description)
            if input.Default != "" {
                fmt.Printf("  Default: %s\n", input.Default)
            }
        }
    }
    
    // Display outputs
    if len(action.Outputs) > 0 {
        fmt.Printf("\n=== Outputs (%d) ===\n", len(action.Outputs))
        for name, output := range action.Outputs {
            fmt.Printf("• %s\n", name)
            fmt.Printf("  %s\n", output.Description)
            if output.Value != "" {
                fmt.Printf("  Value: %s\n", output.Value)
            }
        }
    }
    
    // Display runtime information
    fmt.Printf("\n=== Runtime ===\n")
    fmt.Printf("Using: %s\n", action.Runs.Using)
    
    switch action.Runs.Using {
    case "composite":
        fmt.Printf("Steps: %d\n", len(action.Runs.Steps))
        for i, step := range action.Runs.Steps {
            fmt.Printf("  %d. %s\n", i+1, step.Name)
        }
    case "docker":
        fmt.Printf("Image: %s\n", action.Runs.Image)
        if action.Runs.Entrypoint != "" {
            fmt.Printf("Entrypoint: %s\n", action.Runs.Entrypoint)
        }
    case "node16", "node20":
        fmt.Printf("Main: %s\n", action.Runs.Main)
        if action.Runs.Pre != "" {
            fmt.Printf("Pre: %s\n", action.Runs.Pre)
        }
        if action.Runs.Post != "" {
            fmt.Printf("Post: %s\n", action.Runs.Post)
        }
    }
    
    // Display branding
    if action.Branding.Icon != "" || action.Branding.Color != "" {
        fmt.Printf("\n=== Branding ===\n")
        if action.Branding.Icon != "" {
            fmt.Printf("Icon: %s\n", action.Branding.Icon)
        }
        if action.Branding.Color != "" {
            fmt.Printf("Color: %s\n", action.Branding.Color)
        }
    }
}
```

## Error Handling

Always handle errors appropriately when parsing files:

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    // Check for specific error types
    if os.IsNotExist(err) {
        log.Fatal("Action file does not exist")
    } else if strings.Contains(err.Error(), "yaml") {
        log.Fatal("Invalid YAML syntax in action file")
    } else {
        log.Fatalf("Failed to parse action: %v", err)
    }
}
```

## Next Steps

- Learn about [Workflow Parsing](/examples/workflow-parsing)
- Explore [Validation](/examples/validation) features
- Check out [Utility Functions](/examples/utilities)
