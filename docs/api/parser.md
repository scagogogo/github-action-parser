# Parser Functions

This page documents the core parsing functions provided by the GitHub Action Parser library.

## ParseFile

Parses a GitHub Action YAML file at the specified path.

```go
func ParseFile(path string) (*ActionFile, error)
```

### Parameters

- **path** (`string`): File path to the YAML file to parse

### Returns

- `*ActionFile`: Parsed action/workflow structure
- `error`: Error if parsing fails

### Description

`ParseFile` opens and parses a YAML file from the filesystem. It supports both action files (`action.yml`, `action.yaml`) and workflow files (`.github/workflows/*.yml`).

### Usage Example

```go
// Parse an action file
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatalf("Failed to parse action: %v", err)
}

fmt.Printf("Action: %s\n", action.Name)
fmt.Printf("Description: %s\n", action.Description)

// Parse a workflow file
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
if err != nil {
    log.Fatalf("Failed to parse workflow: %v", err)
}

fmt.Printf("Workflow: %s\n", workflow.Name)
fmt.Printf("Jobs: %d\n", len(workflow.Jobs))
```

### Error Handling

The function returns an error in the following cases:
- File does not exist or cannot be opened
- File contains invalid YAML syntax
- File cannot be read due to permissions

```go
action, err := parser.ParseFile("nonexistent.yml")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("File does not exist")
    } else {
        fmt.Printf("Parse error: %v\n", err)
    }
}
```

## Parse

Parses a GitHub Action YAML from an io.Reader.

```go
func Parse(r io.Reader) (*ActionFile, error)
```

### Parameters

- **r** (`io.Reader`): Reader containing YAML data

### Returns

- `*ActionFile`: Parsed action/workflow structure
- `error`: Error if parsing fails

### Description

`Parse` reads YAML data from any `io.Reader` and parses it into an `ActionFile` structure. This is useful for parsing YAML content from various sources like HTTP responses, embedded files, or in-memory data.

### Usage Example

```go
// Parse from a string
yamlContent := `
name: My Action
description: A sample action
runs:
  using: composite
  steps:
    - name: Hello
      run: echo "Hello World"
`

action, err := parser.Parse(strings.NewReader(yamlContent))
if err != nil {
    log.Fatalf("Failed to parse: %v", err)
}

fmt.Printf("Action: %s\n", action.Name)

// Parse from HTTP response
resp, err := http.Get("https://example.com/action.yml")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

action, err = parser.Parse(resp.Body)
if err != nil {
    log.Fatalf("Failed to parse remote action: %v", err)
}
```

### Error Handling

The function returns an error in the following cases:
- Reader returns an error when reading data
- YAML content is malformed
- Data cannot be unmarshaled into ActionFile structure

```go
// Handle different error types
action, err := parser.Parse(reader)
if err != nil {
    if strings.Contains(err.Error(), "unmarshal") {
        fmt.Println("Invalid YAML structure")
    } else if strings.Contains(err.Error(), "read") {
        fmt.Println("Failed to read data")
    } else {
        fmt.Printf("Parse error: %v\n", err)
    }
}
```

## ParseDir

Parses all GitHub Action YAML files in a directory recursively.

```go
func ParseDir(dir string) (map[string]*ActionFile, error)
```

### Parameters

- **dir** (`string`): Directory path to scan for YAML files

### Returns

- `map[string]*ActionFile`: Map of relative file paths to parsed structures
- `error`: Error if parsing fails

### Description

`ParseDir` recursively walks through a directory and parses all YAML files (`.yml` and `.yaml` extensions). It returns a map where keys are relative file paths and values are the parsed `ActionFile` structures.

### Usage Example

```go
// Parse all workflows in .github/workflows
workflows, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatalf("Failed to parse workflows: %v", err)
}

for path, workflow := range workflows {
    fmt.Printf("File: %s\n", path)
    fmt.Printf("  Name: %s\n", workflow.Name)
    fmt.Printf("  Jobs: %d\n", len(workflow.Jobs))
    
    // Check if it's a reusable workflow
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("  Type: Reusable Workflow\n")
    }
}

// Parse all action files in a repository
actions, err := parser.ParseDir(".")
if err != nil {
    log.Fatalf("Failed to parse repository: %v", err)
}

fmt.Printf("Found %d YAML files\n", len(actions))
```

### File Filtering

`ParseDir` automatically filters files based on extension:
- Includes: `.yml`, `.yaml`
- Excludes: All other file types, directories

### Error Handling

The function returns an error in the following cases:
- Directory does not exist or cannot be accessed
- Permission denied when reading directory or files
- Any individual file fails to parse (stops processing)

```go
actions, err := parser.ParseDir("workflows")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("Directory does not exist")
    } else if strings.Contains(err.Error(), "permission denied") {
        fmt.Println("Permission denied")
    } else {
        fmt.Printf("Parse error: %v\n", err)
    }
}
```

### Performance Considerations

- The function processes files sequentially
- Large directories with many files may take time to process
- Memory usage scales with the number and size of files
- Consider using goroutines for parallel processing if needed

```go
// Example: Processing large directories
start := time.Now()
actions, err := parser.ParseDir("large-repo")
if err != nil {
    log.Fatal(err)
}
duration := time.Since(start)

fmt.Printf("Parsed %d files in %v\n", len(actions), duration)
```

## Best Practices

### File Path Handling

Always use proper file path handling for cross-platform compatibility:

```go
import "path/filepath"

// Good: Use filepath.Join for cross-platform paths
actionPath := filepath.Join("actions", "my-action", "action.yml")
action, err := parser.ParseFile(actionPath)

// Good: Use relative paths when possible
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
```

### Error Handling Patterns

Implement comprehensive error handling:

```go
func parseActionSafely(path string) (*parser.ActionFile, error) {
    action, err := parser.ParseFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to parse %s: %w", path, err)
    }

    // Additional validation
    if action.Name == "" {
        return nil, fmt.Errorf("action in %s has no name", path)
    }

    return action, nil
}
```

### Batch Processing

When processing multiple files, consider error handling strategies:

```go
func parseAllActions(dir string) (map[string]*parser.ActionFile, []error) {
    var errors []error
    results := make(map[string]*parser.ActionFile)

    actions, err := parser.ParseDir(dir)
    if err != nil {
        return nil, []error{err}
    }

    for path, action := range actions {
        // Additional validation per file
        validator := parser.NewValidator()
        if validationErrors := validator.Validate(action); len(validationErrors) > 0 {
            for _, ve := range validationErrors {
                errors = append(errors, fmt.Errorf("%s: %s - %s", path, ve.Field, ve.Message))
            }
        } else {
            results[path] = action
        }
    }

    return results, errors
}
```

### Memory Management

For large-scale processing, consider memory usage:

```go
// Process files one at a time to reduce memory usage
func processLargeRepository(dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
            action, err := parser.ParseFile(path)
            if err != nil {
                return fmt.Errorf("failed to parse %s: %w", path, err)
            }

            // Process action immediately, don't store in memory
            processAction(path, action)
        }

        return nil
    })
}
```
