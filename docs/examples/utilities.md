# Utility Functions

This example demonstrates how to use the utility functions provided by the GitHub Action Parser library for type conversion and data processing.

## Map Conversion Utilities

### MapOfStringInterface

Convert various map types to `map[string]interface{}`:

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Example map with interface{} keys
    interfaceMap := map[interface{}]interface{}{
        "key1": "value1",
        "key2": 42,
        "key3": true,
        "key4": []string{"a", "b", "c"},
    }
    
    // Convert to map[string]interface{}
    stringMap, err := parser.MapOfStringInterface(interfaceMap)
    if err != nil {
        log.Fatalf("Failed to convert map: %v", err)
    }
    
    fmt.Println("Converted map:")
    for key, value := range stringMap {
        fmt.Printf("  %s: %v (type: %T)\n", key, value, value)
    }
}
```

### MapOfStringString

Convert various map types to `map[string]string`:

```go
// Example map with interface{} values
interfaceValueMap := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}

// Convert to map[string]string
stringStringMap, err := parser.MapOfStringString(interfaceValueMap)
if err != nil {
    log.Fatalf("Failed to convert map: %v", err)
}

fmt.Println("\nConverted string map:")
for key, value := range stringStringMap {
    fmt.Printf("  %s: %s\n", key, value)
}

// This will fail because not all values are strings
mixedMap := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
}

_, err = parser.MapOfStringString(mixedMap)
if err != nil {
    fmt.Printf("\nExpected error: %v\n", err)
}
```

## Working with Reusable Workflows

### IsReusableWorkflow

Check if a workflow is reusable:

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse workflow file
    workflow, err := parser.ParseFile(".github/workflows/ci.yml")
    if err != nil {
        log.Fatalf("Failed to parse workflow: %v", err)
    }
    
    // Check if it's a reusable workflow
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("This is a reusable workflow")
    } else {
        fmt.Println("This is a regular workflow")
    }
}
```

### ExtractInputsFromWorkflowCall

Extract input parameters from a reusable workflow:

```go
// Check if workflow is reusable and extract inputs
if parser.IsReusableWorkflow(workflow) {
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        log.Fatalf("Failed to extract inputs: %v", err)
    }
    
    fmt.Printf("Found %d inputs:\n", len(inputs))
    for name, input := range inputs {
        fmt.Printf("  %s:\n", name)
        fmt.Printf("    Description: %s\n", input.Description)
        fmt.Printf("    Required: %t\n", input.Required)
        if input.Default != "" {
            fmt.Printf("    Default: %s\n", input.Default)
        }
    }
}
```

### ExtractOutputsFromWorkflowCall

Extract output parameters from a reusable workflow:

```go
// Extract outputs from reusable workflow
if parser.IsReusableWorkflow(workflow) {
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        log.Fatalf("Failed to extract outputs: %v", err)
    }
    
    fmt.Printf("Found %d outputs:\n", len(outputs))
    for name, output := range outputs {
        fmt.Printf("  %s:\n", name)
        fmt.Printf("    Description: %s\n", output.Description)
        if output.Value != "" {
            fmt.Printf("    Value: %s\n", output.Value)
        }
    }
}
```

## StringOrStringSlice Type

Working with the `StringOrStringSlice` type:

```go
package main

import (
    "fmt"
    "log"
    "gopkg.in/yaml.v3"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Example YAML with string or string array fields
    yamlData := `
singleString: value1
stringArray:
  - item1
  - item2
  - item3
`
    
    // Parse YAML
    var data struct {
        SingleString parser.StringOrStringSlice `yaml:"singleString"`
        StringArray  parser.StringOrStringSlice `yaml:"stringArray"`
    }
    
    err := yaml.Unmarshal([]byte(yamlData), &data)
    if err != nil {
        log.Fatalf("Failed to parse YAML: %v", err)
    }
    
    // Access single string
    fmt.Printf("SingleString.Value: %s\n", data.SingleString.Value)
    fmt.Printf("SingleString.Values: %v\n", data.SingleString.Values)
    
    // Access string array
    fmt.Printf("StringArray.Value: %s\n", data.StringArray.Value)
    fmt.Printf("StringArray.Values: %v\n", data.StringArray.Values)
    
    // Check if contains a value
    if data.StringArray.Contains("item2") {
        fmt.Println("StringArray contains 'item2'")
    }
    
    // String representation
    fmt.Printf("StringArray as string: %s\n", data.StringArray.String())
}
```

## Error Handling Patterns

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Example of safe map conversion with error handling
    interfaceMap := map[interface{}]interface{}{
        "key1": "value1",
        "key2": 42,
        "key3": true,
    }
    
    // Try to convert to string-string map (will fail)
    stringMap, err := safeMapConversion(interfaceMap)
    if err != nil {
        fmt.Printf("Warning: %v\n", err)
        fmt.Printf("Fallback map: %v\n", stringMap)
    } else {
        fmt.Printf("Converted map: %v\n", stringMap)
    }
}

// Safe conversion with fallback
func safeMapConversion(input interface{}) (map[string]string, error) {
    // Try direct conversion first
    result, err := parser.MapOfStringString(input)
    if err == nil {
        return result, nil
    }
    
    // If that fails, try to convert to map[string]interface{} first
    interfaceMap, err := parser.MapOfStringInterface(input)
    if err != nil {
        return nil, fmt.Errorf("failed to convert map: %w", err)
    }
    
    // Then manually convert values to strings
    result = make(map[string]string)
    var conversionErrors []string
    
    for key, value := range interfaceMap {
        switch v := value.(type) {
        case string:
            result[key] = v
        case int, int64, float64, bool:
            result[key] = fmt.Sprintf("%v", v)
        default:
            conversionErrors = append(conversionErrors, 
                fmt.Sprintf("cannot convert %s: %v (%T) to string", key, v, v))
        }
    }
    
    if len(conversionErrors) > 0 {
        return result, fmt.Errorf("partial conversion with %d errors", len(conversionErrors))
    }
    
    return result, nil
}
```

## Batch Processing with Utilities

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all workflows
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    // Process each workflow
    for path, workflow := range workflows {
        fmt.Printf("=== %s ===\n", path)
        
        // Process workflow triggers
        processTriggers(workflow.On)
        
        // Process environment variables
        processEnvVars(workflow.Env)
        
        // Process jobs
        for jobID, job := range workflow.Jobs {
            fmt.Printf("Job: %s\n", jobID)
            
            // Process job environment variables
            processEnvVars(job.Env)
            
            // Process steps
            for _, step := range job.Steps {
                if step.With != nil {
                    processStepInputs(step.With)
                }
            }
        }
        
        fmt.Println()
    }
}

func processTriggers(on interface{}) {
    fmt.Println("Triggers:")
    
    switch v := on.(type) {
    case string:
        fmt.Printf("  %s\n", v)
    case []interface{}:
        for _, trigger := range v {
            fmt.Printf("  %v\n", trigger)
        }
    case map[string]interface{}:
        stringMap, err := parser.MapOfStringInterface(v)
        if err != nil {
            fmt.Printf("  Error converting triggers: %v\n", err)
            return
        }
        
        for event, config := range stringMap {
            fmt.Printf("  %s: %v\n", event, config)
        }
    default:
        fmt.Printf("  Unknown trigger type: %T\n", on)
    }
}

func processEnvVars(env interface{}) {
    if env == nil {
        return
    }
    
    fmt.Println("Environment Variables:")
    
    envVars, err := parser.MapOfStringString(env)
    if err != nil {
        fmt.Printf("  Error converting env vars: %v\n", err)
        return
    }
    
    for key, value := range envVars {
        fmt.Printf("  %s: %s\n", key, value)
    }
}

func processStepInputs(with interface{}) {
    fmt.Println("  Step Inputs:")
    
    inputs, err := parser.MapOfStringInterface(with)
    if err != nil {
        fmt.Printf("    Error converting inputs: %v\n", err)
        return
    }
    
    for key, value := range inputs {
        fmt.Printf("    %s: %v\n", key, value)
    }
}
```

## Complete Utility Example

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
        log.Fatal("Usage: go run main.go <workflow.yml>")
    }
    
    workflowFile := os.Args[1]
    
    // Parse workflow file
    workflow, err := parser.ParseFile(workflowFile)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", workflowFile, err)
    }
    
    // Comprehensive workflow analysis using utility functions
    fmt.Printf("=== Analyzing %s ===\n\n", workflowFile)
    
    // Basic info
    fmt.Printf("Name: %s\n", workflow.Name)
    
    // Process triggers with type conversion
    fmt.Println("\n=== Triggers ===")
    processTriggers(workflow.On)
    
    // Process environment variables
    if workflow.Env != nil {
        fmt.Println("\n=== Environment Variables ===")
        envVars, err := parser.MapOfStringString(workflow.Env)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            for key, value := range envVars {
                fmt.Printf("%s: %s\n", key, value)
            }
        }
    }
    
    // Check if reusable
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("\n=== Reusable Workflow ===")
        
        // Extract inputs
        inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
        if err != nil {
            fmt.Printf("Failed to extract inputs: %v\n", err)
        } else {
            fmt.Printf("Inputs: %d\n", len(inputs))
            for name, input := range inputs {
                fmt.Printf("  %s: %s (required: %t)\n", 
                    name, input.Description, input.Required)
            }
        }
        
        // Extract outputs
        outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
        if err != nil {
            fmt.Printf("Failed to extract outputs: %v\n", err)
        } else {
            fmt.Printf("Outputs: %d\n", len(outputs))
            for name, output := range outputs {
                fmt.Printf("  %s: %s\n", name, output.Description)
            }
        }
    }
    
    // Process jobs
    fmt.Printf("\n=== Jobs (%d) ===\n", len(workflow.Jobs))
    for jobID, job := range workflow.Jobs {
        fmt.Printf("\nJob: %s\n", jobID)
        
        // Process runs-on
        if job.RunsOn != nil {
            fmt.Printf("Runs on: ")
            switch runsOn := job.RunsOn.(type) {
            case string:
                fmt.Printf("%s\n", runsOn)
            case []interface{}:
                fmt.Printf("%v\n", runsOn)
            default:
                fmt.Printf("%v (type: %T)\n", runsOn, runsOn)
            }
        }
        
        // Process steps
        fmt.Printf("Steps: %d\n", len(job.Steps))
    }
}

func processTriggers(on interface{}) {
    switch v := on.(type) {
    case string:
        fmt.Printf("Single event: %s\n", v)
    case []interface{}:
        fmt.Println("Multiple events:")
        for _, event := range v {
            fmt.Printf("  - %v\n", event)
        }
    case map[interface{}]interface{}:
        fmt.Println("Complex events:")
        events, err := parser.MapOfStringInterface(v)
        if err != nil {
            fmt.Printf("Error converting events: %v\n", err)
            return
        }
        
        for event, config := range events {
            fmt.Printf("  %s: %v\n", event, config)
        }
    case map[string]interface{}:
        fmt.Println("Complex events:")
        for event, config := range v {
            fmt.Printf("  %s: %v\n", event, config)
        }
    default:
        fmt.Printf("Unknown trigger type: %T\n", on)
    }
}
```

## Next Steps

- Check out the [API Reference](/api/utilities) for detailed documentation
- Learn about [Validation](/examples/validation) features
- Explore [Reusable Workflows](/examples/reusable-workflows) in detail
