package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_action_file <path_to_action_yaml>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	action, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Display basic action information
	fmt.Printf("Action Name: %s\n", action.Name)
	fmt.Printf("Description: %s\n", action.Description)

	// Validate the action
	validator := parser.NewValidator()
	errors := validator.Validate(action)

	if len(errors) > 0 {
		fmt.Println("\nValidation Errors:")
		for _, err := range errors {
			fmt.Printf("- %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("\nAction is valid.")
	}

	// Display inputs if available
	if len(action.Inputs) > 0 {
		fmt.Println("\nInputs:")
		for name, input := range action.Inputs {
			required := "optional"
			if input.Required {
				required = "required"
			}

			defaultValue := "none"
			if input.Default != "" {
				defaultValue = input.Default
			}

			fmt.Printf("- %s (%s, default: %s): %s\n",
				name, required, defaultValue, input.Description)
		}
	}

	// Display jobs if available (for workflow files)
	if len(action.Jobs) > 0 {
		fmt.Println("\nJobs:")
		for jobID, job := range action.Jobs {
			fmt.Printf("- %s: %s\n", jobID, job.Name)

			if len(job.Steps) > 0 {
				fmt.Println("  Steps:")
				for i, step := range job.Steps {
					if step.Name != "" {
						fmt.Printf("  %d. %s\n", i+1, step.Name)
					} else if step.Run != "" {
						fmt.Printf("  %d. Run: %s...\n", i+1, truncate(step.Run, 50))
					} else if step.Uses != "" {
						fmt.Printf("  %d. Uses: %s\n", i+1, step.Uses)
					} else {
						fmt.Printf("  %d. <unnamed step>\n", i+1)
					}
				}
			}
		}
	}

	// Check if it's a reusable workflow
	if parser.IsReusableWorkflow(action) {
		fmt.Println("\nThis is a reusable workflow.")

		inputs, err := parser.ExtractInputsFromWorkflowCall(action)
		if err != nil {
			fmt.Printf("Error extracting workflow inputs: %v\n", err)
		} else if len(inputs) > 0 {
			fmt.Println("\nWorkflow Inputs:")
			for name, input := range inputs {
				required := "optional"
				if input.Required {
					required = "required"
				}

				fmt.Printf("- %s (%s): %s\n", name, required, input.Description)
			}
		}

		outputs, err := parser.ExtractOutputsFromWorkflowCall(action)
		if err != nil {
			fmt.Printf("Error extracting workflow outputs: %v\n", err)
		} else if len(outputs) > 0 {
			fmt.Println("\nWorkflow Outputs:")
			for name, output := range outputs {
				fmt.Printf("- %s: %s\n", name, output.Description)
			}
		}
	}
}

// truncate truncates a string to the specified length and adds "..." if it was truncated
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
