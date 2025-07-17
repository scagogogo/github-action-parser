# Reusable Workflows

This example demonstrates how to work with reusable workflows, including detection, input/output extraction, and analysis.

## Detect Reusable Workflows

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse workflow file
    workflow, err := parser.ParseFile(".github/workflows/reusable.yml")
    if err != nil {
        log.Fatalf("Failed to parse workflow: %v", err)
    }
    
    // Check if it's a reusable workflow
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("✅ This is a reusable workflow")
        analyzeReusableWorkflow(workflow)
    } else {
        fmt.Println("❌ This is not a reusable workflow")
    }
}

func analyzeReusableWorkflow(workflow *parser.ActionFile) {
    fmt.Printf("Name: %s\n", workflow.Name)
    
    // Extract inputs
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        log.Printf("Failed to extract inputs: %v", err)
    } else {
        fmt.Printf("Inputs: %d\n", len(inputs))
    }
    
    // Extract outputs
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        log.Printf("Failed to extract outputs: %v", err)
    } else {
        fmt.Printf("Outputs: %d\n", len(outputs))
    }
}
```

## Extract and Display Inputs

```go
// Extract and display detailed input information
inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("Failed to extract inputs: %v", err)
}

if len(inputs) > 0 {
    fmt.Printf("\n=== Inputs (%d) ===\n", len(inputs))
    for name, input := range inputs {
        fmt.Printf("• %s", name)
        
        if input.Required {
            fmt.Printf(" (required)")
        } else {
            fmt.Printf(" (optional)")
        }
        
        fmt.Printf("\n  Description: %s\n", input.Description)
        
        if input.Default != "" {
            fmt.Printf("  Default: %s\n", input.Default)
        }
        
        if input.Deprecated {
            fmt.Printf("  ⚠️  Deprecated\n")
        }
        
        fmt.Println()
    }
} else {
    fmt.Println("No inputs defined")
}
```

## Extract and Display Outputs

```go
// Extract and display detailed output information
outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("Failed to extract outputs: %v", err)
}

if len(outputs) > 0 {
    fmt.Printf("\n=== Outputs (%d) ===\n", len(outputs))
    for name, output := range outputs {
        fmt.Printf("• %s\n", name)
        fmt.Printf("  Description: %s\n", output.Description)
        
        if output.Value != "" {
            fmt.Printf("  Value: %s\n", output.Value)
        }
        
        fmt.Println()
    }
} else {
    fmt.Println("No outputs defined")
}
```

## Find All Reusable Workflows

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all workflows
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    fmt.Printf("Scanning %d workflow files...\n\n", len(workflows))
    
    reusableCount := 0
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableCount++
            analyzeReusableWorkflowDetailed(filepath.Base(path), workflow)
        }
    }
    
    fmt.Printf("\n=== Summary ===\n")
    fmt.Printf("Total workflows: %d\n", len(workflows))
    fmt.Printf("Reusable workflows: %d\n", reusableCount)
    fmt.Printf("Regular workflows: %d\n", len(workflows)-reusableCount)
}

func analyzeReusableWorkflowDetailed(filename string, workflow *parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filename)
    
    if workflow.Name != "" {
        fmt.Printf("Name: %s\n", workflow.Name)
    }
    
    // Analyze inputs
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ Failed to extract inputs: %v\n", err)
    } else {
        fmt.Printf("Inputs: %d", len(inputs))
        if len(inputs) > 0 {
            requiredCount := 0
            for _, input := range inputs {
                if input.Required {
                    requiredCount++
                }
            }
            fmt.Printf(" (%d required, %d optional)", requiredCount, len(inputs)-requiredCount)
        }
        fmt.Println()
    }
    
    // Analyze outputs
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ Failed to extract outputs: %v\n", err)
    } else {
        fmt.Printf("Outputs: %d\n", len(outputs))
    }
    
    // Analyze jobs
    fmt.Printf("Jobs: %d\n", len(workflow.Jobs))
    
    fmt.Println()
}
```

## Validate Reusable Workflow Usage

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all workflows to find reusable workflow usage
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    // Find reusable workflows and their callers
    reusableWorkflows := make(map[string]*parser.ActionFile)
    callerWorkflows := make(map[string]*parser.ActionFile)
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableWorkflows[path] = workflow
        } else {
            // Check if this workflow calls any reusable workflows
            for _, job := range workflow.Jobs {
                if job.Uses != "" {
                    callerWorkflows[path] = workflow
                    break
                }
            }
        }
    }
    
    fmt.Printf("Found %d reusable workflows and %d caller workflows\n\n", 
        len(reusableWorkflows), len(callerWorkflows))
    
    // Analyze usage
    for callerPath, callerWorkflow := range callerWorkflows {
        analyzeReusableWorkflowUsage(callerPath, callerWorkflow, reusableWorkflows)
    }
}

func analyzeReusableWorkflowUsage(callerPath string, caller *parser.ActionFile, reusableWorkflows map[string]*parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filepath.Base(callerPath))
    
    for jobID, job := range caller.Jobs {
        if job.Uses == "" {
            continue
        }
        
        fmt.Printf("Job '%s' uses: %s\n", jobID, job.Uses)
        
        // Check if it's a local reusable workflow
        if strings.HasPrefix(job.Uses, "./") {
            localPath := strings.TrimPrefix(job.Uses, "./")
            if reusableWorkflow, exists := reusableWorkflows[localPath]; exists {
                validateReusableWorkflowCall(jobID, job, reusableWorkflow)
            } else {
                fmt.Printf("  ⚠️  Local reusable workflow not found: %s\n", localPath)
            }
        }
        
        // Display inputs passed to reusable workflow
        if len(job.With) > 0 {
            fmt.Printf("  Inputs passed:\n")
            for key, value := range job.With {
                fmt.Printf("    %s: %v\n", key, value)
            }
        }
        
        // Display secrets passed to reusable workflow
        if job.Secrets != nil {
            fmt.Printf("  Secrets: %v\n", job.Secrets)
        }
        
        fmt.Println()
    }
}

func validateReusableWorkflowCall(jobID string, job parser.Job, reusableWorkflow *parser.ActionFile) {
    // Extract expected inputs from reusable workflow
    expectedInputs, err := parser.ExtractInputsFromWorkflowCall(reusableWorkflow)
    if err != nil {
        fmt.Printf("  ❌ Failed to extract expected inputs: %v\n", err)
        return
    }
    
    // Check if all required inputs are provided
    providedInputs := make(map[string]bool)
    for key := range job.With {
        providedInputs[key] = true
    }
    
    missingRequired := []string{}
    extraInputs := []string{}
    
    // Check for missing required inputs
    for name, input := range expectedInputs {
        if input.Required && !providedInputs[name] {
            missingRequired = append(missingRequired, name)
        }
    }
    
    // Check for extra inputs
    for name := range providedInputs {
        if _, exists := expectedInputs[name]; !exists {
            extraInputs = append(extraInputs, name)
        }
    }
    
    // Report validation results
    if len(missingRequired) == 0 && len(extraInputs) == 0 {
        fmt.Printf("  ✅ Input validation passed\n")
    } else {
        if len(missingRequired) > 0 {
            fmt.Printf("  ❌ Missing required inputs: %v\n", missingRequired)
        }
        if len(extraInputs) > 0 {
            fmt.Printf("  ⚠️  Extra inputs (not defined): %v\n", extraInputs)
        }
    }
}
```

## Generate Reusable Workflow Documentation

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Find all reusable workflows
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    reusableWorkflows := make(map[string]*parser.ActionFile)
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableWorkflows[path] = workflow
        }
    }
    
    if len(reusableWorkflows) == 0 {
        fmt.Println("No reusable workflows found")
        return
    }
    
    // Generate documentation
    generateReusableWorkflowDocs(reusableWorkflows)
}

func generateReusableWorkflowDocs(workflows map[string]*parser.ActionFile) {
    fmt.Println("# Reusable Workflows Documentation\n")
    
    // Sort workflows by name
    var paths []string
    for path := range workflows {
        paths = append(paths, path)
    }
    sort.Strings(paths)
    
    for _, path := range paths {
        workflow := workflows[path]
        generateWorkflowDoc(filepath.Base(path), workflow)
    }
}

func generateWorkflowDoc(filename string, workflow *parser.ActionFile) {
    fmt.Printf("## %s\n\n", strings.TrimSuffix(filename, filepath.Ext(filename)))
    
    if workflow.Name != "" {
        fmt.Printf("**Name:** %s\n\n", workflow.Name)
    }
    
    // Extract inputs
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ Failed to extract inputs: %v\n\n", err)
    } else if len(inputs) > 0 {
        fmt.Printf("### Inputs\n\n")
        fmt.Printf("| Name | Description | Required | Default |\n")
        fmt.Printf("|------|-------------|----------|----------|\n")
        
        // Sort inputs by name
        var inputNames []string
        for name := range inputs {
            inputNames = append(inputNames, name)
        }
        sort.Strings(inputNames)
        
        for _, name := range inputNames {
            input := inputs[name]
            required := "No"
            if input.Required {
                required = "Yes"
            }
            
            defaultValue := input.Default
            if defaultValue == "" {
                defaultValue = "-"
            }
            
            fmt.Printf("| `%s` | %s | %s | `%s` |\n", 
                name, input.Description, required, defaultValue)
        }
        fmt.Println()
    }
    
    // Extract outputs
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ Failed to extract outputs: %v\n\n", err)
    } else if len(outputs) > 0 {
        fmt.Printf("### Outputs\n\n")
        fmt.Printf("| Name | Description |\n")
        fmt.Printf("|------|-------------|\n")
        
        // Sort outputs by name
        var outputNames []string
        for name := range outputs {
            outputNames = append(outputNames, name)
        }
        sort.Strings(outputNames)
        
        for _, name := range outputNames {
            output := outputs[name]
            fmt.Printf("| `%s` | %s |\n", name, output.Description)
        }
        fmt.Println()
    }
    
    // Usage example
    fmt.Printf("### Usage Example\n\n")
    fmt.Printf("```yaml\n")
    fmt.Printf("jobs:\n")
    fmt.Printf("  call-reusable-workflow:\n")
    fmt.Printf("    uses: ./.github/workflows/%s\n", filename)
    
    if len(inputs) > 0 {
        fmt.Printf("    with:\n")
        for name, input := range inputs {
            if input.Required {
                fmt.Printf("      %s: # Required - %s\n", name, input.Description)
            }
        }
    }
    
    fmt.Printf("```\n\n")
    
    fmt.Println("---\n")
}
```

## Complete Reusable Workflow Analysis Tool

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  go run main.go analyze    - Analyze all reusable workflows")
        fmt.Println("  go run main.go validate   - Validate reusable workflow usage")
        fmt.Println("  go run main.go docs       - Generate documentation")
        os.Exit(1)
    }
    
    command := os.Args[1]
    
    switch command {
    case "analyze":
        analyzeAllReusableWorkflows()
    case "validate":
        validateReusableWorkflowUsage()
    case "docs":
        generateDocumentation()
    default:
        fmt.Printf("Unknown command: %s\n", command)
        os.Exit(1)
    }
}

func analyzeAllReusableWorkflows() {
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    fmt.Println("=== Reusable Workflow Analysis ===\n")
    
    reusableCount := 0
    totalInputs := 0
    totalOutputs := 0
    
    for path, workflow := range workflows {
        if !parser.IsReusableWorkflow(workflow) {
            continue
        }
        
        reusableCount++
        filename := filepath.Base(path)
        
        fmt.Printf("📄 %s\n", filename)
        if workflow.Name != "" {
            fmt.Printf("   Name: %s\n", workflow.Name)
        }
        
        inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
        outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow)
        
        fmt.Printf("   Inputs: %d, Outputs: %d, Jobs: %d\n", 
            len(inputs), len(outputs), len(workflow.Jobs))
        
        totalInputs += len(inputs)
        totalOutputs += len(outputs)
        
        fmt.Println()
    }
    
    fmt.Printf("=== Summary ===\n")
    fmt.Printf("Reusable workflows: %d\n", reusableCount)
    fmt.Printf("Total inputs: %d\n", totalInputs)
    fmt.Printf("Total outputs: %d\n", totalOutputs)
    
    if reusableCount > 0 {
        fmt.Printf("Average inputs per workflow: %.1f\n", float64(totalInputs)/float64(reusableCount))
        fmt.Printf("Average outputs per workflow: %.1f\n", float64(totalOutputs)/float64(reusableCount))
    }
}

func validateReusableWorkflowUsage() {
    // Implementation similar to previous validation example
    fmt.Println("=== Validating Reusable Workflow Usage ===")
    // ... (validation logic)
}

func generateDocumentation() {
    // Implementation similar to previous documentation example
    fmt.Println("=== Generating Reusable Workflow Documentation ===")
    // ... (documentation generation logic)
}
```

## Next Steps

- Learn about [Utility Functions](/examples/utilities) for advanced processing
- Check out the [API Reference](/api/utilities) for reusable workflow functions
- Explore [Validation](/examples/validation) for comprehensive workflow validation
