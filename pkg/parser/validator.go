package parser

import (
	"fmt"
)

// ValidationError represents an error found during validation
type ValidationError struct {
	Field   string
	Message string
}

// Validator validates an ActionFile to ensure it meets GitHub's requirements
type Validator struct {
	errors []ValidationError
}

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]ValidationError, 0),
	}
}

// Validate checks if an ActionFile is valid according to GitHub's requirements
func (v *Validator) Validate(action *ActionFile) []ValidationError {
	v.errors = make([]ValidationError, 0)

	// Check action metadata for composite or Docker actions
	if action.Runs.Using != "" {
		v.validateActionMetadata(action)
	}

	// Check workflow for workflow files - jobs can be empty map but still considered a workflow
	if action.Jobs != nil {
		v.validateWorkflow(action)
	}

	return v.errors
}

// validateActionMetadata validates action.yml/action.yaml files
func (v *Validator) validateActionMetadata(action *ActionFile) {
	// Name is required
	if action.Name == "" {
		v.addError("name", "Action name is required")
	}

	// Description is required
	if action.Description == "" {
		v.addError("description", "Action description is required")
	}

	// Validate runs configuration
	if action.Runs.Using == "" {
		v.addError("runs.using", "Action must specify 'using' field")
	} else {
		switch action.Runs.Using {
		case "node16", "node20":
			if action.Runs.Main == "" {
				v.addError("runs.main", "JavaScript actions require a 'main' entry point")
			}
		case "docker":
			if action.Runs.Image == "" && action.Runs.Using == "docker" {
				v.addError("runs.image", "Docker actions require an 'image' to use")
			}
		case "composite":
			if len(action.Runs.Steps) == 0 {
				v.addError("runs.steps", "Composite actions require at least one step")
			}
		default:
			v.addError("runs.using", fmt.Sprintf("Unsupported action type: %s", action.Runs.Using))
		}
	}
}

// validateWorkflow validates workflow files
func (v *Validator) validateWorkflow(action *ActionFile) {
	// On trigger is required
	if action.On == nil {
		v.addError("on", "Workflow must have at least one trigger")
	}

	// Validate jobs
	if len(action.Jobs) == 0 {
		v.addError("jobs", "Workflow must have at least one job")
	}

	for jobID, job := range action.Jobs {
		// Either 'runs-on' or 'uses' is required for a job
		if job.RunsOn == nil && job.Uses == "" {
			v.addError(fmt.Sprintf("jobs.%s", jobID), "Job must specify either 'runs-on' or 'uses'")
		}

		// Validate steps if defined
		if job.Steps != nil && len(job.Steps) == 0 {
			v.addError(fmt.Sprintf("jobs.%s.steps", jobID), "Job must have at least one step if steps are defined")
		}

		// Validate steps
		for i, step := range job.Steps {
			if step.Uses == "" && step.Run == "" {
				v.addError(fmt.Sprintf("jobs.%s.steps[%d]", jobID, i), "Step must have either 'uses' or 'run'")
			}
		}
	}
}

// addError adds a validation error to the list
func (v *Validator) addError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsValid returns true if there are no validation errors
func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}
