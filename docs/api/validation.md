# Validation API

The validation API provides tools to validate GitHub Action and Workflow files according to GitHub's specifications.

## Validator

The main validator struct for validating ActionFile structures.

```go
type Validator struct {
    errors []ValidationError
}
```

### Description

`Validator` checks ActionFile structures against GitHub's requirements and specifications. It validates both GitHub Actions and Workflows, providing detailed error information for any violations found.

## NewValidator

Creates a new Validator instance.

```go
func NewValidator() *Validator
```

### Returns

- `*Validator`: A new validator instance ready for use

### Usage Example

```go
validator := parser.NewValidator()
```

## Validate

Validates an ActionFile according to GitHub's requirements.

```go
func (v *Validator) Validate(action *ActionFile) []ValidationError
```

### Parameters

- **action** (`*ActionFile`): The action or workflow to validate

### Returns

- `[]ValidationError`: Slice of validation errors found (empty if valid)

### Description

`Validate` performs comprehensive validation of the ActionFile structure, checking:

- **Action Validation**: For files with `runs` configuration
  - Required fields (name, description)
  - Runtime-specific requirements (main script for Node.js, image for Docker, steps for composite)
  - Supported runtime types

- **Workflow Validation**: For files with `jobs` configuration
  - Required trigger events (`on` field)
  - Job requirements (runs-on or uses)
  - Step validation (uses or run required)

### Usage Example

```go
// Basic validation
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

validator := parser.NewValidator()
errors := validator.Validate(action)

if len(errors) > 0 {
    fmt.Println("Validation errors found:")
    for _, err := range errors {
        fmt.Printf("- %s: %s\n", err.Field, err.Message)
    }
} else {
    fmt.Println("Action is valid!")
}

// Validate multiple files
files := []string{"action.yml", "workflow.yml"}
for _, file := range files {
    action, err := parser.ParseFile(file)
    if err != nil {
        fmt.Printf("Failed to parse %s: %v\n", file, err)
        continue
    }
    
    errors := validator.Validate(action)
    if len(errors) > 0 {
        fmt.Printf("%s has %d validation errors:\n", file, len(errors))
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Printf("%s is valid\n", file)
    }
}
```

## IsValid

Checks if the validator has no errors after validation.

```go
func (v *Validator) IsValid() bool
```

### Returns

- `bool`: True if no validation errors exist, false otherwise

### Usage Example

```go
validator := parser.NewValidator()
validator.Validate(action)

if validator.IsValid() {
    fmt.Println("No validation errors")
} else {
    fmt.Println("Validation errors found")
}
```

## ValidationError

Represents a validation error with field and message information.

```go
type ValidationError struct {
    Field   string
    Message string
}
```

### Fields

- **Field** (`string`): The field path where the error occurred
- **Message** (`string`): Human-readable error message

### Usage Example

```go
errors := validator.Validate(action)
for _, err := range errors {
    fmt.Printf("Field: %s\n", err.Field)
    fmt.Printf("Error: %s\n", err.Message)
    fmt.Println("---")
}
```

## Validation Rules

### Action Validation Rules

#### Required Fields
- `name`: Action must have a name
- `description`: Action must have a description
- `runs.using`: Action must specify a runtime

#### Runtime-Specific Rules

**Node.js Actions** (`using: "node16"` or `using: "node20"`):
- `runs.main`: Must specify main entry point script

**Docker Actions** (`using: "docker"`):
- `runs.image`: Must specify Docker image

**Composite Actions** (`using: "composite"`):
- `runs.steps`: Must have at least one step

#### Supported Runtimes
- `node16`: Node.js 16 runtime
- `node20`: Node.js 20 runtime  
- `docker`: Docker container runtime
- `composite`: Composite action runtime

### Workflow Validation Rules

#### Required Fields
- `on`: Workflow must have at least one trigger event
- `jobs`: Workflow must have at least one job

#### Job Rules
- Each job must specify either `runs-on` or `uses`
- If `steps` is defined, it must contain at least one step

#### Step Rules
- Each step must have either `uses` or `run`

### Example Validation Scenarios

#### Valid Action Examples

```yaml
# Valid Node.js Action
name: My Node Action
description: A Node.js action
runs:
  using: node20
  main: index.js

# Valid Docker Action  
name: My Docker Action
description: A Docker action
runs:
  using: docker
  image: Dockerfile

# Valid Composite Action
name: My Composite Action
description: A composite action
runs:
  using: composite
  steps:
    - name: Hello
      run: echo "Hello"
```

#### Valid Workflow Examples

```yaml
# Valid Workflow
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: npm test

# Valid Reusable Workflow
name: Reusable Workflow
on:
  workflow_call:
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Building"
```

## Common Validation Errors

### Action Errors

| Error Field | Message | Solution |
|-------------|---------|----------|
| `name` | Action name is required | Add a `name` field to your action |
| `description` | Action description is required | Add a `description` field to your action |
| `runs.using` | Action must specify 'using' field | Add `runs.using` with a supported runtime |
| `runs.main` | JavaScript actions require a 'main' entry point | Add `runs.main` for Node.js actions |
| `runs.image` | Docker actions require an 'image' to use | Add `runs.image` for Docker actions |
| `runs.steps` | Composite actions require at least one step | Add steps to `runs.steps` for composite actions |
| `runs.using` | Unsupported action type: {type} | Use a supported runtime (node16, node20, docker, composite) |

### Workflow Errors

| Error Field | Message | Solution |
|-------------|---------|----------|
| `on` | Workflow must have at least one trigger | Add trigger events to the `on` field |
| `jobs` | Workflow must have at least one job | Add at least one job to the `jobs` section |
| `jobs.{job-id}` | Job must specify either 'runs-on' or 'uses' | Add `runs-on` or `uses` to the job |
| `jobs.{job-id}.steps` | Job must have at least one step if steps are defined | Add steps or remove empty `steps` array |
| `jobs.{job-id}.steps[{index}]` | Step must have either 'uses' or 'run' | Add `uses` or `run` to the step |

## Advanced Validation Patterns

### Custom Validation Wrapper

```go
type ExtendedValidator struct {
    *parser.Validator
    customRules []ValidationRule
}

type ValidationRule func(*parser.ActionFile) []parser.ValidationError

func NewExtendedValidator() *ExtendedValidator {
    return &ExtendedValidator{
        Validator: parser.NewValidator(),
        customRules: []ValidationRule{
            validateActionNaming,
            validateSecurityPractices,
        },
    }
}

func (ev *ExtendedValidator) ValidateWithCustomRules(action *parser.ActionFile) []parser.ValidationError {
    // Run standard validation
    errors := ev.Validator.Validate(action)

    // Run custom rules
    for _, rule := range ev.customRules {
        errors = append(errors, rule(action)...)
    }

    return errors
}

func validateActionNaming(action *parser.ActionFile) []parser.ValidationError {
    var errors []parser.ValidationError

    if action.Name != "" && !strings.Contains(action.Name, "Action") {
        errors = append(errors, parser.ValidationError{
            Field:   "name",
            Message: "Action name should contain 'Action'",
        })
    }

    return errors
}
```

### Batch Validation

```go
func ValidateRepository(repoPath string) (map[string][]parser.ValidationError, error) {
    results := make(map[string][]parser.ValidationError)
    validator := parser.NewValidator()

    // Parse all YAML files
    files, err := parser.ParseDir(repoPath)
    if err != nil {
        return nil, fmt.Errorf("failed to parse repository: %w", err)
    }

    // Validate each file
    for path, action := range files {
        errors := validator.Validate(action)
        if len(errors) > 0 {
            results[path] = errors
        }
    }

    return results, nil
}

// Usage
validationResults, err := ValidateRepository(".")
if err != nil {
    log.Fatal(err)
}

if len(validationResults) == 0 {
    fmt.Println("All files are valid!")
} else {
    fmt.Printf("Found validation errors in %d files:\n", len(validationResults))
    for path, errors := range validationResults {
        fmt.Printf("\n%s:\n", path)
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

### Validation with Suggestions

```go
func ValidateWithSuggestions(action *parser.ActionFile) {
    validator := parser.NewValidator()
    errors := validator.Validate(action)

    if len(errors) == 0 {
        fmt.Println("✅ Action is valid!")
        return
    }

    fmt.Printf("❌ Found %d validation errors:\n\n", len(errors))

    for _, err := range errors {
        fmt.Printf("Field: %s\n", err.Field)
        fmt.Printf("Error: %s\n", err.Message)

        // Provide suggestions based on error type
        switch err.Field {
        case "name":
            fmt.Println("💡 Suggestion: Add a descriptive name like 'My Awesome Action'")
        case "description":
            fmt.Println("💡 Suggestion: Add a description explaining what your action does")
        case "runs.main":
            fmt.Println("💡 Suggestion: Add 'main: index.js' or your entry point file")
        case "runs.image":
            fmt.Println("💡 Suggestion: Add 'image: Dockerfile' or a Docker image reference")
        case "runs.steps":
            fmt.Println("💡 Suggestion: Add at least one step with 'run' or 'uses'")
        }
        fmt.Println()
    }
}
```
