# Workflow Parsing

This example demonstrates how to parse GitHub Workflow files and extract job information, steps, and triggers.

## Parse Workflow File

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

    fmt.Printf("Workflow: %s\n", workflow.Name)
    fmt.Printf("Jobs: %d\n", len(workflow.Jobs))
}
```

## Access Workflow Triggers

```go
// Access workflow triggers (on events)
fmt.Println("\nTriggers:")
switch on := workflow.On.(type) {
case string:
    fmt.Printf("  Single trigger: %s\n", on)
case []interface{}:
    fmt.Printf("  Multiple triggers:\n")
    for _, trigger := range on {
        fmt.Printf("    - %v\n", trigger)
    }
case map[string]interface{}:
    fmt.Printf("  Complex triggers:\n")
    for event, config := range on {
        fmt.Printf("    - %s: %v\n", event, config)
    }
}
```

## Analyze Jobs

```go
// Analyze each job
fmt.Println("\nJobs:")
for jobID, job := range workflow.Jobs {
    fmt.Printf("  %s:\n", jobID)
    if job.Name != "" {
        fmt.Printf("    Name: %s\n", job.Name)
    }
    
    // Check runs-on
    switch runsOn := job.RunsOn.(type) {
    case string:
        fmt.Printf("    Runs on: %s\n", runsOn)
    case []interface{}:
        fmt.Printf("    Runs on: %v\n", runsOn)
    }
    
    // Check dependencies
    if job.Needs != nil {
        switch needs := job.Needs.(type) {
        case string:
            fmt.Printf("    Needs: %s\n", needs)
        case []interface{}:
            fmt.Printf("    Needs: %v\n", needs)
        }
    }
    
    fmt.Printf("    Steps: %d\n", len(job.Steps))
}
```

## Analyze Job Steps

```go
// Detailed step analysis
for jobID, job := range workflow.Jobs {
    fmt.Printf("\n=== Job: %s ===\n", jobID)
    
    for i, step := range job.Steps {
        fmt.Printf("Step %d:\n", i+1)
        
        if step.Name != "" {
            fmt.Printf("  Name: %s\n", step.Name)
        }
        
        if step.ID != "" {
            fmt.Printf("  ID: %s\n", step.ID)
        }
        
        if step.Uses != "" {
            fmt.Printf("  Uses: %s\n", step.Uses)
            
            // Display step inputs
            if len(step.With) > 0 {
                fmt.Printf("  With:\n")
                for key, value := range step.With {
                    fmt.Printf("    %s: %v\n", key, value)
                }
            }
        }
        
        if step.Run != "" {
            fmt.Printf("  Run: %s\n", step.Run)
            if step.Shell != "" {
                fmt.Printf("  Shell: %s\n", step.Shell)
            }
        }
        
        if step.If != "" {
            fmt.Printf("  If: %s\n", step.If)
        }
        
        // Display step environment variables
        if len(step.Env) > 0 {
            fmt.Printf("  Env:\n")
            for key, value := range step.Env {
                fmt.Printf("    %s: %s\n", key, value)
            }
        }
        
        fmt.Println()
    }
}
```

## Check for Reusable Workflows

```go
// Check if any jobs use reusable workflows
fmt.Println("\nReusable Workflow Jobs:")
for jobID, job := range workflow.Jobs {
    if job.Uses != "" {
        fmt.Printf("  %s uses: %s\n", jobID, job.Uses)
        
        // Display inputs passed to reusable workflow
        if len(job.With) > 0 {
            fmt.Printf("    Inputs:\n")
            for key, value := range job.With {
                fmt.Printf("      %s: %v\n", key, value)
            }
        }
        
        // Display secrets passed to reusable workflow
        if job.Secrets != nil {
            fmt.Printf("    Secrets: %v\n", job.Secrets)
        }
    }
}
```

## Parse Multiple Workflows

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // Parse all workflows in .github/workflows directory
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("Failed to parse workflows: %v", err)
    }
    
    fmt.Printf("Found %d workflow files\n\n", len(workflows))
    
    for path, workflow := range workflows {
        analyzeWorkflow(path, workflow)
    }
}

func analyzeWorkflow(path string, workflow *parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filepath.Base(path))
    
    if workflow.Name != "" {
        fmt.Printf("Name: %s\n", workflow.Name)
    }
    
    // Count different types of triggers
    triggerCount := countTriggers(workflow.On)
    fmt.Printf("Triggers: %d\n", triggerCount)
    
    // Analyze jobs
    fmt.Printf("Jobs: %d\n", len(workflow.Jobs))
    
    totalSteps := 0
    reusableJobs := 0
    
    for _, job := range workflow.Jobs {
        totalSteps += len(job.Steps)
        if job.Uses != "" {
            reusableJobs++
        }
    }
    
    fmt.Printf("Total Steps: %d\n", totalSteps)
    if reusableJobs > 0 {
        fmt.Printf("Reusable Workflow Jobs: %d\n", reusableJobs)
    }
    
    // Check if this is a reusable workflow
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("Type: Reusable Workflow\n")
        
        inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
        outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow)
        
        fmt.Printf("Inputs: %d\n", len(inputs))
        fmt.Printf("Outputs: %d\n", len(outputs))
    }
    
    fmt.Println()
}

func countTriggers(on interface{}) int {
    switch triggers := on.(type) {
    case string:
        return 1
    case []interface{}:
        return len(triggers)
    case map[string]interface{}:
        return len(triggers)
    default:
        return 0
    }
}
```

## Environment Variables and Secrets

```go
// Access workflow-level environment variables
if len(workflow.Env) > 0 {
    fmt.Println("\nWorkflow Environment Variables:")
    for key, value := range workflow.Env {
        fmt.Printf("  %s: %s\n", key, value)
    }
}

// Access job-level environment variables
for jobID, job := range workflow.Jobs {
    if len(job.Env) > 0 {
        fmt.Printf("\nJob %s Environment Variables:\n", jobID)
        for key, value := range job.Env {
            fmt.Printf("  %s: %s\n", key, value)
        }
    }
}
```

## Complete Workflow Analysis Example

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
    
    workflow, err := parser.ParseFile(workflowFile)
    if err != nil {
        log.Fatalf("Failed to parse %s: %v", workflowFile, err)
    }
    
    analyzeCompleteWorkflow(workflow)
}

func analyzeCompleteWorkflow(workflow *parser.ActionFile) {
    fmt.Printf("=== Workflow Analysis ===\n")
    
    if workflow.Name != "" {
        fmt.Printf("Name: %s\n", workflow.Name)
    }
    
    // Analyze triggers
    fmt.Printf("\nTriggers:\n")
    analyzeTriggers(workflow.On)
    
    // Analyze jobs
    fmt.Printf("\nJobs (%d):\n", len(workflow.Jobs))
    for jobID, job := range workflow.Jobs {
        analyzeJob(jobID, job)
    }
    
    // Check if reusable
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("\n=== Reusable Workflow ===\n")
        analyzeReusableWorkflow(workflow)
    }
}

func analyzeTriggers(on interface{}) {
    switch triggers := on.(type) {
    case string:
        fmt.Printf("  - %s\n", triggers)
    case []interface{}:
        for _, trigger := range triggers {
            fmt.Printf("  - %v\n", trigger)
        }
    case map[string]interface{}:
        for event, config := range triggers {
            fmt.Printf("  - %s:\n", event)
            if configMap, ok := config.(map[string]interface{}); ok {
                for key, value := range configMap {
                    fmt.Printf("      %s: %v\n", key, value)
                }
            }
        }
    }
}

func analyzeJob(jobID string, job parser.Job) {
    fmt.Printf("\n  %s:\n", jobID)
    
    if job.Name != "" {
        fmt.Printf("    Name: %s\n", job.Name)
    }
    
    if job.Uses != "" {
        fmt.Printf("    Uses: %s (Reusable Workflow)\n", job.Uses)
    } else {
        fmt.Printf("    Steps: %d\n", len(job.Steps))
        
        // Show first few steps
        for i, step := range job.Steps {
            if i >= 3 { // Limit to first 3 steps
                fmt.Printf("    ... and %d more steps\n", len(job.Steps)-3)
                break
            }
            
            stepDesc := fmt.Sprintf("Step %d", i+1)
            if step.Name != "" {
                stepDesc = step.Name
            } else if step.Uses != "" {
                stepDesc = fmt.Sprintf("Uses %s", step.Uses)
            } else if step.Run != "" {
                stepDesc = "Run command"
            }
            
            fmt.Printf("      %d. %s\n", i+1, stepDesc)
        }
    }
}

func analyzeReusableWorkflow(workflow *parser.ActionFile) {
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err == nil && len(inputs) > 0 {
        fmt.Printf("Inputs (%d):\n", len(inputs))
        for name, input := range inputs {
            required := ""
            if input.Required {
                required = " (required)"
            }
            fmt.Printf("  - %s%s: %s\n", name, required, input.Description)
        }
    }
    
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err == nil && len(outputs) > 0 {
        fmt.Printf("Outputs (%d):\n", len(outputs))
        for name, output := range outputs {
            fmt.Printf("  - %s: %s\n", name, output.Description)
        }
    }
}
```

## Next Steps

- Learn about [Validation](/examples/validation) of workflows
- Explore [Reusable Workflows](/examples/reusable-workflows) in detail
- Check out [Utility Functions](/examples/utilities) for advanced processing
