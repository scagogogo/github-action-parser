# API Reference

The GitHub Action Parser library provides a comprehensive set of types and functions for parsing, validating, and processing GitHub Action and Workflow YAML files.

## Package Overview

```go
import "github.com/scagogogo/github-action-parser/pkg/parser"
```

The `parser` package contains all the functionality needed to work with GitHub Actions and Workflows:

- **Core Types**: Data structures representing GitHub Action and Workflow components
- **Parser Functions**: Functions to parse YAML files and directories
- **Validation**: Tools to validate parsed files according to GitHub specifications
- **Utilities**: Helper functions for type conversion and data processing

## Quick Reference

### Main Functions

| Function | Description |
|----------|-------------|
| [`ParseFile(path string)`](/api/parser#parsefile) | Parse a single YAML file |
| [`Parse(r io.Reader)`](/api/parser#parse) | Parse from an io.Reader |
| [`ParseDir(dir string)`](/api/parser#parsedir) | Parse all YAML files in a directory |
| [`NewValidator()`](/api/validation#newvalidator) | Create a new validator instance |

### Core Types

| Type | Description |
|------|-------------|
| [`ActionFile`](/api/types#actionfile) | Main structure representing an action or workflow |
| [`Input`](/api/types#input) | Input parameter definition |
| [`Output`](/api/types#output) | Output parameter definition |
| [`Job`](/api/types#job) | Workflow job definition |
| [`Step`](/api/types#step) | Individual step in a job |
| [`RunsConfig`](/api/types#runsconfig) | Action execution configuration |

### Validation Types

| Type | Description |
|------|-------------|
| [`Validator`](/api/validation#validator) | Validator for GitHub Action specifications |
| [`ValidationError`](/api/validation#validationerror) | Validation error information |

### Utility Types

| Type | Description |
|------|-------------|
| [`StringOrStringSlice`](/api/utilities#stringorstringslice) | Flexible string/array type for YAML |

## Error Handling

All parsing functions return errors that provide detailed information about parsing failures:

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    // Handle parsing error
    fmt.Printf("Failed to parse: %v\n", err)
    return
}
```

Validation errors are returned as a slice of `ValidationError` structs:

```go
validator := parser.NewValidator()
errors := validator.Validate(action)
for _, err := range errors {
    fmt.Printf("Field %s: %s\n", err.Field, err.Message)
}
```

## Type Safety

The library provides full type safety for all GitHub Action and Workflow structures. All fields are properly typed according to the GitHub Actions specification, with appropriate use of pointers for optional fields and interfaces for flexible data types.

## Performance

The parser is optimized for performance and can efficiently handle:
- Large action and workflow files
- Batch processing of multiple files
- Recursive directory parsing
- Memory-efficient processing of large repositories

## Next Steps

- [Types Reference](/api/types) - Detailed documentation of all data structures
- [Parser Functions](/api/parser) - Complete parsing API documentation
- [Validation](/api/validation) - Validation features and error handling
- [Utilities](/api/utilities) - Helper functions and utilities
