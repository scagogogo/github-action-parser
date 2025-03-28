package parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Test parsing a regular action file
	action, err := ParseFile("testdata/action.yml")
	if err != nil {
		t.Fatalf("Failed to parse action file: %v", err)
	}

	// Verify basic fields
	if action.Name != "Example GitHub Action" {
		t.Errorf("Expected action name to be 'Example GitHub Action', got '%s'", action.Name)
	}

	if action.Description != "An example GitHub Action for testing the parser" {
		t.Errorf("Expected action description to be 'An example GitHub Action for testing the parser', got '%s'", action.Description)
	}

	if action.Author != "GitHub Action Parser" {
		t.Errorf("Expected action author to be 'GitHub Action Parser', got '%s'", action.Author)
	}

	// Verify branding
	if action.Branding.Icon != "code" {
		t.Errorf("Expected branding icon to be 'code', got '%s'", action.Branding.Icon)
	}

	if action.Branding.Color != "blue" {
		t.Errorf("Expected branding color to be 'blue', got '%s'", action.Branding.Color)
	}

	// Verify inputs
	if len(action.Inputs) != 3 {
		t.Errorf("Expected 3 inputs, got %d", len(action.Inputs))
	}

	filePathInput, ok := action.Inputs["file-path"]
	if !ok {
		t.Errorf("Expected 'file-path' input to be defined")
	} else {
		if !filePathInput.Required {
			t.Errorf("Expected 'file-path' input to be required")
		}
		if filePathInput.Description != "Path to the file to process" {
			t.Errorf("Expected 'file-path' description to be 'Path to the file to process', got '%s'", filePathInput.Description)
		}
	}

	// Verify outputs
	if len(action.Outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(action.Outputs))
	}

	resultOutput, ok := action.Outputs["result"]
	if !ok {
		t.Errorf("Expected 'result' output to be defined")
	} else {
		if resultOutput.Description != "The result of the action" {
			t.Errorf("Expected 'result' description to be 'The result of the action', got '%s'", resultOutput.Description)
		}
	}

	// Verify runs configuration
	if action.Runs.Using != "composite" {
		t.Errorf("Expected 'using' to be 'composite', got '%s'", action.Runs.Using)
	}

	if len(action.Runs.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(action.Runs.Steps))
	}
}

// TestParseNonExistentFile tests parsing a file that doesn't exist
func TestParseNonExistentFile(t *testing.T) {
	_, err := ParseFile("testdata/non-existent-file.yml")
	if err == nil {
		t.Errorf("Expected error when parsing non-existent file, got nil")
	}
}

// TestParseInvalidYAML tests parsing an invalid YAML file
func TestParseInvalidYAML(t *testing.T) {
	invalidYAML := `
name: Invalid YAML
description: This is not valid YAML
inputs:
  this-is-invalid:
    - missing colon
    "unclosed quote
`
	_, err := Parse(strings.NewReader(invalidYAML))
	if err == nil {
		t.Errorf("Expected error when parsing invalid YAML, got nil")
	}
}

func TestParseWorkflow(t *testing.T) {
	// Test parsing a workflow file
	workflow, err := ParseFile("testdata/workflow.yml")
	if err != nil {
		t.Fatalf("Failed to parse workflow file: %v", err)
	}

	// Verify basic fields
	if workflow.Name != "CI/CD Workflow" {
		t.Errorf("Expected workflow name to be 'CI/CD Workflow', got '%s'", workflow.Name)
	}

	// Verify jobs
	if len(workflow.Jobs) != 4 {
		t.Errorf("Expected 4 jobs, got %d", len(workflow.Jobs))
	}

	lintJob, ok := workflow.Jobs["lint"]
	if !ok {
		t.Errorf("Expected 'lint' job to be defined")
	} else {
		if lintJob.Name != "Lint Code" {
			t.Errorf("Expected 'lint' job name to be 'Lint Code', got '%s'", lintJob.Name)
		}
	}

	testJob, ok := workflow.Jobs["test"]
	if !ok {
		t.Errorf("Expected 'test' job to be defined")
	} else {
		if testJob.Needs == nil {
			t.Errorf("Expected 'test' job to have dependencies")
		}
	}

	deployJob, ok := workflow.Jobs["deploy"]
	if !ok {
		t.Errorf("Expected 'deploy' job to be defined")
	} else {
		if deployJob.If == "" {
			t.Errorf("Expected 'deploy' job to have an 'if' condition")
		}
	}

	// Test workflow environment variables
	if len(workflow.Env) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(workflow.Env))
	}

	nodeVersion, ok := workflow.Env["NODE_VERSION"]
	if !ok || nodeVersion != "16" {
		t.Errorf("Expected NODE_VERSION to be '16', got '%s'", nodeVersion)
	}
}

func TestParseReusableWorkflow(t *testing.T) {
	// Test parsing a reusable workflow file
	workflow, err := ParseFile("testdata/reusable-workflow.yml")
	if err != nil {
		t.Fatalf("Failed to parse reusable workflow file: %v", err)
	}

	// Check if it's a reusable workflow
	if !IsReusableWorkflow(workflow) {
		t.Errorf("Expected workflow to be a reusable workflow")
	}

	// Verify inputs
	inputs, err := ExtractInputsFromWorkflowCall(workflow)
	if err != nil {
		t.Fatalf("Failed to extract inputs: %v", err)
	}

	if len(inputs) != 4 {
		t.Errorf("Expected 4 inputs, got %d", len(inputs))
	}

	artifactNameInput, ok := inputs["artifact-name"]
	if !ok {
		t.Errorf("Expected 'artifact-name' input to be defined")
	} else {
		if !artifactNameInput.Required {
			t.Errorf("Expected 'artifact-name' input to be required")
		}
	}

	// Verify default values
	nodeVersionInput, ok := inputs["node-version"]
	if !ok {
		t.Errorf("Expected 'node-version' input to be defined")
	} else {
		if nodeVersionInput.Default != "16" {
			t.Errorf("Expected 'node-version' default to be '16', got '%s'", nodeVersionInput.Default)
		}
	}

	// Verify outputs
	outputs, err := ExtractOutputsFromWorkflowCall(workflow)
	if err != nil {
		t.Fatalf("Failed to extract outputs: %v", err)
	}

	if len(outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(outputs))
	}

	buildTimeOutput, ok := outputs["build-time"]
	if !ok {
		t.Errorf("Expected 'build-time' output to be defined")
	} else {
		if buildTimeOutput.Description != "Time taken to build the project" {
			t.Errorf("Expected 'build-time' description to be 'Time taken to build the project', got '%s'", buildTimeOutput.Description)
		}
		if !strings.Contains(buildTimeOutput.Value, "jobs.build.outputs.build-time") {
			t.Errorf("Expected 'build-time' value to reference job output, got '%s'", buildTimeOutput.Value)
		}
	}
}

// TestParseDir tests parsing a directory of action files
func TestParseDir(t *testing.T) {
	actions, err := ParseDir("testdata")
	if err != nil {
		t.Fatalf("Failed to parse directory: %v", err)
	}

	// Should have parsed 3 files
	expectedFiles := 3
	if len(actions) != expectedFiles {
		t.Errorf("Expected %d actions, got %d", expectedFiles, len(actions))
	}

	// Verify we have the expected files
	fileNames := []string{"action.yml", "workflow.yml", "reusable-workflow.yml"}
	for _, fileName := range fileNames {
		_, exists := actions[fileName]
		if !exists {
			t.Errorf("Expected to find %s in parsed files, but it was not found", fileName)
		}
	}
}

// TestParseDirNonExistent tests parsing a non-existent directory
func TestParseDirNonExistent(t *testing.T) {
	_, err := ParseDir("non-existent-dir")
	if err == nil {
		t.Errorf("Expected error when parsing non-existent directory, got nil")
	}
}

func TestParse(t *testing.T) {
	// Test parsing from a reader
	file, err := os.Open("testdata/action.yml")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	action, err := Parse(file)
	if err != nil {
		t.Fatalf("Failed to parse from reader: %v", err)
	}

	if action.Name != "Example GitHub Action" {
		t.Errorf("Expected action name to be 'Example GitHub Action', got '%s'", action.Name)
	}
}

// TestParseEmptyReader tests parsing from an empty reader
func TestParseEmptyReader(t *testing.T) {
	action, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Errorf("Unexpected error when parsing empty reader: %v", err)
	}

	// An empty reader should result in an empty ActionFile, not an error
	if action == nil {
		t.Errorf("Expected non-nil ActionFile from empty reader")
	}
}

// TestParseErrorReader tests parsing from a reader that returns an error
func TestParseErrorReader(t *testing.T) {
	errorReader := &ErrorReader{Err: io.ErrUnexpectedEOF}
	_, err := Parse(errorReader)
	if err == nil {
		t.Errorf("Expected error when parsing from error reader, got nil")
	}
}

// ErrorReader is a reader that always returns an error
type ErrorReader struct {
	Err error
}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, r.Err
}

func TestValidateAction(t *testing.T) {
	// Test validating an action file
	action, err := ParseFile("testdata/action.yml")
	if err != nil {
		t.Fatalf("Failed to parse action file: %v", err)
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) > 0 {
		t.Errorf("Expected no validation errors, got %d: %v", len(errors), errors)
	}

	// Test with an invalid action (by modifying a valid one)
	action.Name = ""
	action.Description = ""

	errors = validator.Validate(action)
	if len(errors) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(errors))
	}

	// Check the specific validation errors
	foundNameError := false
	foundDescError := false
	for _, err := range errors {
		if err.Field == "name" {
			foundNameError = true
		}
		if err.Field == "description" {
			foundDescError = true
		}
	}

	if !foundNameError {
		t.Errorf("Expected validation error for missing name")
	}
	if !foundDescError {
		t.Errorf("Expected validation error for missing description")
	}
}

// TestValidateJavaScriptAction tests validation of JavaScript actions
func TestValidateJavaScriptAction(t *testing.T) {
	// Create a JavaScript action for testing
	action := &ActionFile{
		Name:        "JavaScript Action",
		Description: "A JavaScript Action",
		Runs: RunsConfig{
			Using: "node16",
			// Missing Main field
		},
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(errors))
	} else if errors[0].Field != "runs.main" {
		t.Errorf("Expected validation error for runs.main, got %s", errors[0].Field)
	}

	// Fix the action by adding the main field
	action.Runs.Main = "index.js"
	errors = validator.Validate(action)

	if len(errors) > 0 {
		t.Errorf("Expected no validation errors after fixing, got %d: %v", len(errors), errors)
	}
}

// TestValidateDockerAction tests validation of Docker actions
func TestValidateDockerAction(t *testing.T) {
	// Create a Docker action for testing
	action := &ActionFile{
		Name:        "Docker Action",
		Description: "A Docker Action",
		Runs: RunsConfig{
			Using: "docker",
			// Missing Image field
		},
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(errors))
	} else if errors[0].Field != "runs.image" {
		t.Errorf("Expected validation error for runs.image, got %s", errors[0].Field)
	}

	// Fix the action by adding the image field
	action.Runs.Image = "Dockerfile"
	errors = validator.Validate(action)

	if len(errors) > 0 {
		t.Errorf("Expected no validation errors after fixing, got %d: %v", len(errors), errors)
	}
}

// TestValidateWorkflow tests validation of workflows
func TestValidateWorkflow(t *testing.T) {
	// Create a workflow with empty jobs map
	workflow := &ActionFile{
		Name: "Invalid Workflow",
		Jobs: map[string]Job{}, // Empty jobs map
	}

	validator := NewValidator()
	errors := validator.Validate(workflow)

	// Should have error for missing 'on' field and empty jobs
	if len(errors) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(errors))
		for i, e := range errors {
			t.Logf("Error %d: %s - %s", i, e.Field, e.Message)
		}
	}

	// Check for specific errors
	hasOnError := false
	hasJobsError := false
	for _, err := range errors {
		if err.Field == "on" {
			hasOnError = true
		}
		if err.Field == "jobs" {
			hasJobsError = true
		}
	}

	if !hasOnError {
		t.Errorf("Expected error for missing 'on' trigger")
	}
	if !hasJobsError {
		t.Errorf("Expected error for empty jobs")
	}

	// Add a basic 'on' trigger and a job without runs-on or uses
	workflow.On = map[string]interface{}{"push": nil}
	workflow.Jobs = map[string]Job{
		"test": {
			Name:  "Test Job",
			Steps: []Step{}, // Empty steps
		},
	}

	errors = validator.Validate(workflow)

	// Should have error for job without runs-on or uses, and for empty steps
	expectedErrors := 2
	if len(errors) != expectedErrors {
		t.Errorf("Expected %d validation errors, got %d", expectedErrors, len(errors))
		for i, e := range errors {
			t.Logf("Error %d: %s - %s", i, e.Field, e.Message)
		}
	}

	// Fix the job by adding runs-on and steps
	workflow.Jobs["test"] = Job{
		Name:   "Test Job",
		RunsOn: "ubuntu-latest",
		Steps: []Step{
			{
				Name: "Test Step",
				Run:  "echo 'Hello World'",
			},
		},
	}

	errors = validator.Validate(workflow)
	if len(errors) > 0 {
		t.Errorf("Expected no validation errors after fixing, got %d: %v", len(errors), errors)
	}
}

// TestValidateInvalidSteps tests validation of invalid steps
func TestValidateInvalidSteps(t *testing.T) {
	// Create a workflow with an invalid step
	workflow := &ActionFile{
		Name: "Workflow with Invalid Step",
		On:   map[string]interface{}{"push": nil},
		Jobs: map[string]Job{
			"test": {
				Name:   "Test Job",
				RunsOn: "ubuntu-latest",
				Steps: []Step{
					{
						Name: "Invalid Step",
						// Missing both 'uses' and 'run'
					},
				},
			},
		},
	}

	validator := NewValidator()
	errors := validator.Validate(workflow)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(errors))
	} else if !strings.Contains(errors[0].Field, "steps[0]") {
		t.Errorf("Expected validation error for invalid step, got field %s", errors[0].Field)
	}
}

// TestIsValidMethod tests the IsValid method
func TestIsValidMethod(t *testing.T) {
	validator := NewValidator()

	// Empty validator should be valid
	if !validator.IsValid() {
		t.Errorf("Expected new validator to be valid")
	}

	// Add an error and check again
	validator.addError("test", "test error")
	if validator.IsValid() {
		t.Errorf("Expected validator with errors to be invalid")
	}
}

// TestStringOrStringSlice tests the StringOrStringSlice utilities
func TestStringOrStringSlice(t *testing.T) {
	// Test Contains method
	sss := StringOrStringSlice{
		Values: []string{"a", "b", "c"},
	}

	if !sss.Contains("a") {
		t.Errorf("Expected StringOrStringSlice to contain 'a'")
	}

	if sss.Contains("d") {
		t.Errorf("Expected StringOrStringSlice to not contain 'd'")
	}

	// Test String method
	sss1 := StringOrStringSlice{
		Value:  "single",
		Values: []string{"single"},
	}
	if sss1.String() != "single" {
		t.Errorf("Expected String() to return 'single', got '%s'", sss1.String())
	}

	sss2 := StringOrStringSlice{
		Values: []string{"a", "b", "c"},
	}
	if sss2.String() != "a, b, c" {
		t.Errorf("Expected String() to return 'a, b, c', got '%s'", sss2.String())
	}
}

// TestMapUtilities tests the map conversion utilities
func TestMapUtilities(t *testing.T) {
	// Test MapOfStringInterface with map[string]interface{}
	input1 := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	result1, err := MapOfStringInterface(input1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result1["key1"] != "value1" || result1["key2"] != 42 {
		t.Errorf("MapOfStringInterface did not preserve values")
	}

	// Test MapOfStringInterface with map[interface{}]interface{}
	input2 := map[interface{}]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	result2, err := MapOfStringInterface(input2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result2["key1"] != "value1" || result2["key2"] != 42 {
		t.Errorf("MapOfStringInterface did not convert keys correctly")
	}

	// Test MapOfStringInterface with nil
	result3, err := MapOfStringInterface(nil)
	if err != nil {
		t.Errorf("Unexpected error with nil input: %v", err)
	}
	if result3 != nil {
		t.Errorf("Expected nil result for nil input")
	}

	// Test MapOfStringInterface with invalid type
	_, err = MapOfStringInterface(42)
	if err == nil {
		t.Errorf("Expected error for invalid input type")
	}

	// Test MapOfStringString
	stringMap := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	result4, err := MapOfStringString(stringMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result4["key1"] != "value1" || result4["key2"] != "value2" {
		t.Errorf("MapOfStringString did not preserve values")
	}
}
