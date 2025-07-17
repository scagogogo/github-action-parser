# Validation

This example demonstrates how to validate GitHub Actions and Workflows using the built-in validation features.

## Basic Validation

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse action file
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatalf("Failed to parse action: %v", err)
    }
    
    // Create validator and validate
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) == 0 {
        fmt.Println("✅ Action is valid!")
    } else {
        fmt.Printf("❌ Found %d validation errors:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

## Validate Multiple Files

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all YAML files in current directory
    files, err := parser.ParseDir(".")
    if err != nil {
        log.Fatalf("Failed to parse directory: %v", err)
    }
    
    validator := parser.NewValidator()
    totalErrors := 0
    
    for path, action := range files {
        fmt.Printf("\n=== Validating %s ===\n", filepath.Base(path))
        
        errors := validator.Validate(action)
        if len(errors) == 0 {
            fmt.Println("✅ Valid")
        } else {
            fmt.Printf("❌ %d errors:\n", len(errors))
            for _, err := range errors {
                fmt.Printf("  - %s: %s\n", err.Field, err.Message)
            }
            totalErrors += len(errors)
        }
    }
    
    fmt.Printf("\n=== Summary ===\n")
    fmt.Printf("Files checked: %d\n", len(files))
    fmt.Printf("Total errors: %d\n", totalErrors)
}
```

## Validation with Detailed Reporting

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
        log.Fatal("Usage: go run main.go <file.yml>")
    }
    
    filename := os.Args[1]
    
    // Parse file
    action, err := parser.ParseFile(filename)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", filename, err)
    }
    
    // Validate with detailed reporting
    validateWithDetails(filename, action)
}

func validateWithDetails(filename string, action *parser.ActionFile) {
    fmt.Printf("=== Validating %s ===\n\n", filename)
    
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) == 0 {
        fmt.Println("✅ File is valid!")
        displayFileInfo(action)
        return
    }
    
    // Group errors by category
    fieldErrors := make(map[string][]parser.ValidationError)
    for _, err := range errors {
        fieldErrors[err.Field] = append(fieldErrors[err.Field], err)
    }
    
    fmt.Printf("❌ Found %d validation errors:\n\n", len(errors))
    
    // Display errors by category
    for field, errs := range fieldErrors {
        fmt.Printf("Field: %s\n", field)
        for _, err := range errs {
            fmt.Printf("  ❌ %s\n", err.Message)
            provideSuggestion(err)
        }
        fmt.Println()
    }
    
    // Provide general suggestions
    fmt.Println("💡 General Suggestions:")
    provideGeneralSuggestions(action, errors)
}

func provideSuggestion(err parser.ValidationError) {
    suggestions := map[string]string{
        "name":        "Add a descriptive name for your action",
        "description": "Add a clear description explaining what your action does",
        "runs.using":  "Specify a supported runtime: node16, node20, docker, or composite",
        "runs.main":   "For JavaScript actions, specify the main entry point file",
        "runs.image":  "For Docker actions, specify the Docker image or Dockerfile",
        "runs.steps":  "For composite actions, add at least one step",
        "on":          "Specify at least one trigger event for workflows",
        "jobs":        "Add at least one job to your workflow",
    }
    
    if suggestion, exists := suggestions[err.Field]; exists {
        fmt.Printf("    💡 %s\n", suggestion)
    }
}

func provideGeneralSuggestions(action *parser.ActionFile, errors []parser.ValidationError) {
    // Suggest based on action type
    if action.Runs.Using != "" {
        switch action.Runs.Using {
        case "composite":
            fmt.Println("  - For composite actions, ensure each step has either 'uses' or 'run'")
            fmt.Println("  - Consider adding step names for better readability")
        case "docker":
            fmt.Println("  - For Docker actions, ensure your Dockerfile exists")
            fmt.Println("  - Consider specifying an entrypoint if needed")
        case "node16", "node20":
            fmt.Println("  - For JavaScript actions, ensure your main file exists")
            fmt.Println("  - Consider adding pre/post scripts if needed")
        }
    }
    
    // Suggest based on error patterns
    hasRequiredFieldErrors := false
    for _, err := range errors {
        if err.Field == "name" || err.Field == "description" {
            hasRequiredFieldErrors = true
            break
        }
    }
    
    if hasRequiredFieldErrors {
        fmt.Println("  - Required fields (name, description) are essential for GitHub Actions")
        fmt.Println("  - These help users understand what your action does")
    }
    
    // Workflow-specific suggestions
    if len(action.Jobs) > 0 {
        fmt.Println("  - For workflows, ensure each job has 'runs-on' or 'uses'")
        fmt.Println("  - Check that all referenced actions exist and are accessible")
    }
}

func displayFileInfo(action *parser.ActionFile) {
    fmt.Println("\n📋 File Information:")
    
    if action.Name != "" {
        fmt.Printf("  Name: %s\n", action.Name)
    }
    
    if action.Description != "" {
        fmt.Printf("  Description: %s\n", action.Description)
    }
    
    if action.Runs.Using != "" {
        fmt.Printf("  Type: %s action\n", action.Runs.Using)
    }
    
    if len(action.Jobs) > 0 {
        fmt.Printf("  Type: Workflow with %d jobs\n", len(action.Jobs))
    }
    
    if len(action.Inputs) > 0 {
        fmt.Printf("  Inputs: %d\n", len(action.Inputs))
    }
    
    if len(action.Outputs) > 0 {
        fmt.Printf("  Outputs: %d\n", len(action.Outputs))
    }
}
```

## Batch Validation with Summary

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Validate all workflows in .github/workflows
    validateDirectory(".github/workflows", "Workflows")
    
    // Validate action files in current directory
    validateDirectory(".", "Actions")
}

func validateDirectory(dir, category string) {
    fmt.Printf("\n=== Validating %s in %s ===\n", category, dir)
    
    files, err := parser.ParseDir(dir)
    if err != nil {
        log.Printf("Failed to parse %s: %v", dir, err)
        return
    }
    
    if len(files) == 0 {
        fmt.Printf("No YAML files found in %s\n", dir)
        return
    }
    
    validator := parser.NewValidator()
    
    validFiles := 0
    totalErrors := 0
    errorsByType := make(map[string]int)
    
    for path, action := range files {
        errors := validator.Validate(action)
        
        filename := filepath.Base(path)
        if len(errors) == 0 {
            fmt.Printf("✅ %s\n", filename)
            validFiles++
        } else {
            fmt.Printf("❌ %s (%d errors)\n", filename, len(errors))
            totalErrors += len(errors)
            
            // Count error types
            for _, err := range errors {
                errorsByType[err.Field]++
            }
        }
    }
    
    // Display summary
    fmt.Printf("\n📊 %s Summary:\n", category)
    fmt.Printf("  Total files: %d\n", len(files))
    fmt.Printf("  Valid files: %d\n", validFiles)
    fmt.Printf("  Files with errors: %d\n", len(files)-validFiles)
    fmt.Printf("  Total errors: %d\n", totalErrors)
    
    if len(errorsByType) > 0 {
        fmt.Printf("\n🔍 Most common errors:\n")
        for field, count := range errorsByType {
            fmt.Printf("  %s: %d occurrences\n", field, count)
        }
    }
}
```

## Custom Validation Rules

```go
package main

import (
    "fmt"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

// Custom validator with additional rules
type CustomValidator struct {
    *parser.Validator
}

func NewCustomValidator() *CustomValidator {
    return &CustomValidator{
        Validator: parser.NewValidator(),
    }
}

func (cv *CustomValidator) ValidateWithCustomRules(action *parser.ActionFile) []parser.ValidationError {
    // Run standard validation first
    errors := cv.Validator.Validate(action)
    
    // Add custom validation rules
    errors = append(errors, cv.validateNaming(action)...)
    errors = append(errors, cv.validateSecurity(action)...)
    errors = append(errors, cv.validateBestPractices(action)...)
    
    return errors
}

func (cv *CustomValidator) validateNaming(action *parser.ActionFile) []parser.ValidationError {
    var errors []parser.ValidationError
    
    // Check action name follows conventions
    if action.Name != "" {
        if !strings.Contains(strings.ToLower(action.Name), "action") && len(action.Jobs) == 0 {
            errors = append(errors, parser.ValidationError{
                Field:   "name",
                Message: "Action name should contain 'Action' for clarity",
            })
        }
        
        if len(action.Name) > 50 {
            errors = append(errors, parser.ValidationError{
                Field:   "name",
                Message: "Action name should be 50 characters or less",
            })
        }
    }
    
    return errors
}

func (cv *CustomValidator) validateSecurity(action *parser.ActionFile) []parser.ValidationError {
    var errors []parser.ValidationError
    
    // Check for hardcoded secrets in steps
    for jobID, job := range action.Jobs {
        for i, step := range job.Steps {
            if step.Run != "" {
                if strings.Contains(strings.ToLower(step.Run), "password") ||
                   strings.Contains(strings.ToLower(step.Run), "token") {
                    errors = append(errors, parser.ValidationError{
                        Field:   fmt.Sprintf("jobs.%s.steps[%d].run", jobID, i),
                        Message: "Avoid hardcoding secrets in run commands",
                    })
                }
            }
        }
    }
    
    return errors
}

func (cv *CustomValidator) validateBestPractices(action *parser.ActionFile) []parser.ValidationError {
    var errors []parser.ValidationError
    
    // Check for step names
    for jobID, job := range action.Jobs {
        unnamedSteps := 0
        for _, step := range job.Steps {
            if step.Name == "" {
                unnamedSteps++
            }
        }
        
        if unnamedSteps > len(job.Steps)/2 {
            errors = append(errors, parser.ValidationError{
                Field:   fmt.Sprintf("jobs.%s.steps", jobID),
                Message: "Consider adding names to steps for better readability",
            })
        }
    }
    
    // Check for action versioning
    for jobID, job := range action.Jobs {
        for i, step := range job.Steps {
            if step.Uses != "" && strings.Contains(step.Uses, "actions/") {
                if !strings.Contains(step.Uses, "@v") && !strings.Contains(step.Uses, "@main") {
                    errors = append(errors, parser.ValidationError{
                        Field:   fmt.Sprintf("jobs.%s.steps[%d].uses", jobID, i),
                        Message: "Consider pinning action to a specific version",
                    })
                }
            }
        }
    }
    
    return errors
}

func main() {
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Use custom validator
    validator := NewCustomValidator()
    errors := validator.ValidateWithCustomRules(action)
    
    if len(errors) == 0 {
        fmt.Println("✅ Action passes all validation rules!")
    } else {
        fmt.Printf("Found %d issues:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

## Validation in CI/CD

Here's an example of how to use validation in a CI/CD pipeline:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    exitCode := 0
    
    // Validate all workflow files
    if err := validateWorkflows(); err != nil {
        fmt.Printf("❌ Workflow validation failed: %v\n", err)
        exitCode = 1
    }
    
    // Validate action files
    if err := validateActions(); err != nil {
        fmt.Printf("❌ Action validation failed: %v\n", err)
        exitCode = 1
    }
    
    if exitCode == 0 {
        fmt.Println("✅ All validations passed!")
    }
    
    os.Exit(exitCode)
}

func validateWorkflows() error {
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        return err
    }
    
    validator := parser.NewValidator()
    hasErrors := false
    
    for path, workflow := range workflows {
        errors := validator.Validate(workflow)
        if len(errors) > 0 {
            fmt.Printf("❌ %s:\n", filepath.Base(path))
            for _, err := range errors {
                fmt.Printf("  - %s: %s\n", err.Field, err.Message)
            }
            hasErrors = true
        }
    }
    
    if hasErrors {
        return fmt.Errorf("workflow validation failed")
    }
    
    return nil
}

func validateActions() error {
    // Look for action.yml or action.yaml files
    actionFiles := []string{"action.yml", "action.yaml"}
    
    validator := parser.NewValidator()
    hasErrors := false
    
    for _, filename := range actionFiles {
        if _, err := os.Stat(filename); os.IsNotExist(err) {
            continue
        }
        
        action, err := parser.ParseFile(filename)
        if err != nil {
            return fmt.Errorf("failed to parse %s: %w", filename, err)
        }
        
        errors := validator.Validate(action)
        if len(errors) > 0 {
            fmt.Printf("❌ %s:\n", filename)
            for _, err := range errors {
                fmt.Printf("  - %s: %s\n", err.Field, err.Message)
            }
            hasErrors = true
        }
    }
    
    if hasErrors {
        return fmt.Errorf("action validation failed")
    }
    
    return nil
}
```

## Next Steps

- Learn about [Reusable Workflows](/examples/reusable-workflows)
- Explore [Utility Functions](/examples/utilities) for advanced processing
- Check out the [API Reference](/api/validation) for detailed validation documentation
