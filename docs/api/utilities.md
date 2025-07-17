# Utilities API

The utilities API provides helper functions for type conversion, data processing, and working with reusable workflows.

## Type Conversion Functions

### MapOfStringInterface

Converts a YAML map to `map[string]interface{}`.

```go
func MapOfStringInterface(v interface{}) (map[string]interface{}, error)
```

#### Parameters

- **v** (`interface{}`): Input value to convert

#### Returns

- `map[string]interface{}`: Converted map
- `error`: Error if conversion fails

#### Description

This function handles the common YAML unmarshaling scenario where maps can have different key types. It converts:
- `map[string]interface{}` → returns as-is
- `map[interface{}]interface{}` → converts keys to strings
- `nil` → returns `nil`

#### Usage Example

```go
// Working with workflow triggers
switch on := workflow.On.(type) {
case map[interface{}]interface{}:
    // Convert to string-keyed map
    onMap, err := parser.MapOfStringInterface(on)
    if err != nil {
        log.Fatal(err)
    }
    
    for event, config := range onMap {
        fmt.Printf("Trigger: %s\n", event)
    }
}

// Working with job configuration
if job.With != nil {
    withMap, err := parser.MapOfStringInterface(job.With)
    if err != nil {
        log.Fatal(err)
    }
    
    for key, value := range withMap {
        fmt.Printf("Input %s: %v\n", key, value)
    }
}
```

#### Error Cases

- Returns error if input contains non-string keys that cannot be converted
- Returns error for unsupported input types

### MapOfStringString

Converts a YAML map to `map[string]string`.

```go
func MapOfStringString(v interface{}) (map[string]string, error)
```

#### Parameters

- **v** (`interface{}`): Input value to convert

#### Returns

- `map[string]string`: Converted map with string values
- `error`: Error if conversion fails

#### Description

Converts various map types to a string-only map. Handles:
- `map[string]string` → returns as-is
- `map[string]interface{}` → converts values to strings
- `map[interface{}]interface{}` → converts keys and values to strings
- `nil` → returns `nil`

#### Usage Example

```go
// Working with environment variables
if job.Env != nil {
    envMap, err := parser.MapOfStringString(job.Env)
    if err != nil {
        log.Fatal(err)
    }
    
    for key, value := range envMap {
        fmt.Printf("ENV %s=%s\n", key, value)
    }
}

// Converting step environment variables
if step.Env != nil {
    stepEnv, err := parser.MapOfStringString(step.Env)
    if err != nil {
        log.Printf("Warning: Could not convert step env: %v", err)
    } else {
        for k, v := range stepEnv {
            fmt.Printf("  %s: %s\n", k, v)
        }
    }
}
```

#### Error Cases

- Returns error if values cannot be converted to strings
- Returns error if keys cannot be converted to strings
- Returns error for unsupported input types

## Reusable Workflow Functions

### IsReusableWorkflow

Checks if a workflow is intended to be called by other workflows.

```go
func IsReusableWorkflow(action *ActionFile) bool
```

#### Parameters

- **action** (`*ActionFile`): The workflow to check

#### Returns

- `bool`: True if the workflow is reusable, false otherwise

#### Description

Determines if a workflow is reusable by checking for the `workflow_call` trigger event in the `on` field.

#### Usage Example

```go
workflow, err := parser.ParseFile(".github/workflows/reusable.yml")
if err != nil {
    log.Fatal(err)
}

if parser.IsReusableWorkflow(workflow) {
    fmt.Println("This is a reusable workflow")
    
    // Extract inputs and outputs
    inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
    outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow)
    
    fmt.Printf("Inputs: %d\n", len(inputs))
    fmt.Printf("Outputs: %d\n", len(outputs))
} else {
    fmt.Println("This is a regular workflow")
}

// Batch check workflows
workflows, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatal(err)
}

reusableCount := 0
for path, workflow := range workflows {
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("Reusable: %s\n", path)
        reusableCount++
    }
}

fmt.Printf("Found %d reusable workflows\n", reusableCount)
```

### ExtractInputsFromWorkflowCall

Extracts input definitions from a reusable workflow.

```go
func ExtractInputsFromWorkflowCall(action *ActionFile) (map[string]Input, error)
```

#### Parameters

- **action** (`*ActionFile`): The reusable workflow

#### Returns

- `map[string]Input`: Map of input names to Input definitions
- `error`: Error if extraction fails

#### Description

Extracts input parameter definitions from the `workflow_call` trigger configuration. Returns `nil` if the workflow is not reusable or has no inputs defined.

#### Usage Example

```go
workflow, err := parser.ParseFile("reusable-workflow.yml")
if err != nil {
    log.Fatal(err)
}

if !parser.IsReusableWorkflow(workflow) {
    fmt.Println("Not a reusable workflow")
    return
}

inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("Failed to extract inputs: %v", err)
}

if len(inputs) == 0 {
    fmt.Println("No inputs defined")
} else {
    fmt.Printf("Workflow inputs:\n")
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

Extracts output definitions from a reusable workflow.

```go
func ExtractOutputsFromWorkflowCall(action *ActionFile) (map[string]Output, error)
```

#### Parameters

- **action** (`*ActionFile`): The reusable workflow

#### Returns

- `map[string]Output`: Map of output names to Output definitions
- `error`: Error if extraction fails

#### Description

Extracts output parameter definitions from the `workflow_call` trigger configuration. Returns `nil` if the workflow is not reusable or has no outputs defined.

#### Usage Example

```go
workflow, err := parser.ParseFile("reusable-workflow.yml")
if err != nil {
    log.Fatal(err)
}

outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("Failed to extract outputs: %v", err)
}

if len(outputs) == 0 {
    fmt.Println("No outputs defined")
} else {
    fmt.Printf("Workflow outputs:\n")
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

A utility type for handling YAML fields that can be either a string or an array of strings.

### Methods

#### Contains

```go
func (s *StringOrStringSlice) Contains(value string) bool
```

Checks if the given value is contained in the string or string slice.

#### String

```go
func (s *StringOrStringSlice) String() string
```

Returns a string representation. For single values, returns the value. For multiple values, returns a comma-separated list.

#### Usage Example

```go
// This type is used internally but can be useful for custom processing
var triggers StringOrStringSlice

// Simulate YAML unmarshaling
yaml.Unmarshal([]byte(`["push", "pull_request"]`), &triggers)

if triggers.Contains("push") {
    fmt.Println("Triggered by push events")
}

fmt.Printf("All triggers: %s\n", triggers.String())
// Output: All triggers: push, pull_request
```

## Error Handling Patterns

### Graceful Type Conversion

```go
func safeMapConversion(v interface{}) map[string]string {
    result, err := parser.MapOfStringString(v)
    if err != nil {
        // Fallback to string interface map
        if interfaceMap, err2 := parser.MapOfStringInterface(v); err2 == nil {
            result = make(map[string]string)
            for k, val := range interfaceMap {
                result[k] = fmt.Sprintf("%v", val)
            }
        }
    }
    return result
}
```
