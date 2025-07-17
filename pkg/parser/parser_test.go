package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	// Test MapOfStringInterface with non-string key
	input3 := map[interface{}]interface{}{
		42:     "value1",
		"key2": "value2",
	}
	_, err = MapOfStringInterface(input3)
	if err == nil {
		t.Errorf("Expected error for non-string key")
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

	// Test MapOfStringString with map[string]interface{} containing strings
	interfaceMap := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	result5, err := MapOfStringString(interfaceMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result5["key1"] != "value1" || result5["key2"] != "value2" {
		t.Errorf("MapOfStringString did not convert interface map correctly")
	}

	// Test MapOfStringString with map[string]interface{} containing non-strings
	interfaceMapWithNonString := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	_, err = MapOfStringString(interfaceMapWithNonString)
	if err == nil {
		t.Errorf("Expected error for non-string value")
	}

	// Test MapOfStringString with map[interface{}]interface{}
	interfaceKeyMap := map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	result6, err := MapOfStringString(interfaceKeyMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result6["key1"] != "value1" || result6["key2"] != "value2" {
		t.Errorf("MapOfStringString did not convert interface key map correctly")
	}

	// Test MapOfStringString with map[interface{}]interface{} containing non-string key
	interfaceKeyMapWithNonStringKey := map[interface{}]interface{}{
		42:     "value1",
		"key2": "value2",
	}
	_, err = MapOfStringString(interfaceKeyMapWithNonStringKey)
	if err == nil {
		t.Errorf("Expected error for non-string key")
	}

	// Test MapOfStringString with map[interface{}]interface{} containing non-string value
	interfaceKeyMapWithNonStringValue := map[interface{}]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	_, err = MapOfStringString(interfaceKeyMapWithNonStringValue)
	if err == nil {
		t.Errorf("Expected error for non-string value")
	}

	// Test MapOfStringString with nil
	result7, err := MapOfStringString(nil)
	if err != nil {
		t.Errorf("Unexpected error with nil input: %v", err)
	}
	if result7 != nil {
		t.Errorf("Expected nil result for nil input")
	}

	// Test MapOfStringString with invalid type
	_, err = MapOfStringString(42)
	if err == nil {
		t.Errorf("Expected error for invalid input type")
	}
}

// TestStringOrStringSliceUnmarshalYAML tests the UnmarshalYAML method
func TestStringOrStringSliceUnmarshalYAML(t *testing.T) {
	// Test unmarshaling a string
	var sss1 StringOrStringSlice
	err := sss1.UnmarshalYAML(func(v interface{}) error {
		// First try will be for string
		if str, ok := v.(*string); ok {
			*str = "test-string"
			return nil
		}
		// If not string, return error to trigger slice attempt
		return fmt.Errorf("not a string")
	})
	if err != nil {
		t.Errorf("Unexpected error unmarshaling string: %v", err)
	}
	if sss1.Value != "test-string" {
		t.Errorf("Expected Value to be 'test-string', got '%s'", sss1.Value)
	}
	if len(sss1.Values) != 1 || sss1.Values[0] != "test-string" {
		t.Errorf("Expected Values to be ['test-string'], got %v", sss1.Values)
	}

	// Test unmarshaling a slice of strings
	var sss2 StringOrStringSlice
	callCount := 0
	err = sss2.UnmarshalYAML(func(v interface{}) error {
		callCount++
		if callCount == 1 {
			// First call tries string, should fail
			return fmt.Errorf("not a string")
		}
		// Second call tries slice
		if slice, ok := v.(*[]string); ok {
			*slice = []string{"item1", "item2", "item3"}
			return nil
		}
		return fmt.Errorf("not a slice")
	})
	if err != nil {
		t.Errorf("Unexpected error unmarshaling slice: %v", err)
	}
	if sss2.Value != "item1" {
		t.Errorf("Expected Value to be 'item1', got '%s'", sss2.Value)
	}
	if len(sss2.Values) != 3 || sss2.Values[0] != "item1" || sss2.Values[1] != "item2" || sss2.Values[2] != "item3" {
		t.Errorf("Expected Values to be ['item1', 'item2', 'item3'], got %v", sss2.Values)
	}

	// Test unmarshaling an empty slice
	var sss3 StringOrStringSlice
	callCount3 := 0
	err = sss3.UnmarshalYAML(func(v interface{}) error {
		callCount3++
		if callCount3 == 1 {
			// First call tries string, should fail
			return fmt.Errorf("not a string")
		}
		// Second call tries slice
		if slice, ok := v.(*[]string); ok {
			*slice = []string{}
			return nil
		}
		return fmt.Errorf("not a slice")
	})
	if err != nil {
		t.Errorf("Unexpected error unmarshaling empty slice: %v", err)
	}
	if sss3.Value != "" {
		t.Errorf("Expected Value to be empty for empty slice, got '%s'", sss3.Value)
	}
	if len(sss3.Values) != 0 {
		t.Errorf("Expected Values to be empty, got %v", sss3.Values)
	}

	// Test unmarshaling invalid type (should fail)
	var sss4 StringOrStringSlice
	err = sss4.UnmarshalYAML(func(v interface{}) error {
		return fmt.Errorf("unmarshal error")
	})
	if err == nil {
		t.Errorf("Expected error when unmarshaling fails for both string and slice")
	}
}

// TestIsReusableWorkflow tests the IsReusableWorkflow function
func TestIsReusableWorkflow(t *testing.T) {
	// Test with workflow_call in map[string]interface{}
	workflow1 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{},
			"push":          nil,
		},
	}
	if !IsReusableWorkflow(workflow1) {
		t.Errorf("Expected workflow with workflow_call to be reusable")
	}

	// Test with workflow_call in map[interface{}]interface{}
	workflow2 := &ActionFile{
		On: map[interface{}]interface{}{
			"workflow_call": map[string]interface{}{},
			"push":          nil,
		},
	}
	if !IsReusableWorkflow(workflow2) {
		t.Errorf("Expected workflow with workflow_call to be reusable")
	}

	// Test without workflow_call
	workflow3 := &ActionFile{
		On: map[string]interface{}{
			"push": nil,
		},
	}
	if IsReusableWorkflow(workflow3) {
		t.Errorf("Expected workflow without workflow_call to not be reusable")
	}

	// Test with map[interface{}]interface{} but non-string key
	workflow4 := &ActionFile{
		On: map[interface{}]interface{}{
			42:     "value",
			"push": nil,
		},
	}
	if IsReusableWorkflow(workflow4) {
		t.Errorf("Expected workflow with non-string key to not be reusable")
	}

	// Test with nil On
	workflow5 := &ActionFile{
		On: nil,
	}
	if IsReusableWorkflow(workflow5) {
		t.Errorf("Expected workflow with nil On to not be reusable")
	}

	// Test with unsupported On type
	workflow6 := &ActionFile{
		On: "push",
	}
	if IsReusableWorkflow(workflow6) {
		t.Errorf("Expected workflow with string On to not be reusable")
	}
}

// TestExtractInputsFromWorkflowCall tests the ExtractInputsFromWorkflowCall function
func TestExtractInputsFromWorkflowCall(t *testing.T) {
	// Test with valid workflow_call inputs
	workflow1 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"inputs": map[string]interface{}{
					"input1": map[string]interface{}{
						"description": "First input",
						"required":    true,
						"default":     "default1",
					},
					"input2": map[string]interface{}{
						"description": "Second input",
						"required":    false,
					},
				},
			},
		},
	}

	inputs, err := ExtractInputsFromWorkflowCall(workflow1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(inputs))
	}
	if inputs["input1"].Description != "First input" {
		t.Errorf("Expected input1 description to be 'First input', got '%s'", inputs["input1"].Description)
	}
	if !inputs["input1"].Required {
		t.Errorf("Expected input1 to be required")
	}
	if inputs["input1"].Default != "default1" {
		t.Errorf("Expected input1 default to be 'default1', got '%s'", inputs["input1"].Default)
	}

	// Test without workflow_call
	workflow2 := &ActionFile{
		On: map[string]interface{}{
			"push": nil,
		},
	}
	inputs2, err := ExtractInputsFromWorkflowCall(workflow2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if inputs2 != nil {
		t.Errorf("Expected nil inputs for non-reusable workflow")
	}

	// Test without inputs in workflow_call
	workflow3 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{},
		},
	}
	inputs3, err := ExtractInputsFromWorkflowCall(workflow3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if inputs3 != nil {
		t.Errorf("Expected nil inputs when no inputs defined")
	}

	// Test with invalid workflow_call type
	workflow4 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": "invalid",
		},
	}
	_, err = ExtractInputsFromWorkflowCall(workflow4)
	if err == nil {
		t.Errorf("Expected error for invalid workflow_call type")
	}

	// Test with invalid inputs type
	workflow5 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"inputs": "invalid",
			},
		},
	}
	_, err = ExtractInputsFromWorkflowCall(workflow5)
	if err == nil {
		t.Errorf("Expected error for invalid inputs type")
	}

	// Test with invalid input definition type
	workflow6 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"inputs": map[string]interface{}{
					"input1": "invalid",
				},
			},
		},
	}
	_, err = ExtractInputsFromWorkflowCall(workflow6)
	if err == nil {
		t.Errorf("Expected error for invalid input definition type")
	}
}

// TestExtractOutputsFromWorkflowCall tests the ExtractOutputsFromWorkflowCall function
func TestExtractOutputsFromWorkflowCall(t *testing.T) {
	// Test with valid workflow_call outputs
	workflow1 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"outputs": map[string]interface{}{
					"output1": map[string]interface{}{
						"description": "First output",
						"value":       "${{ jobs.build.outputs.result }}",
					},
					"output2": map[string]interface{}{
						"description": "Second output",
					},
				},
			},
		},
	}

	outputs, err := ExtractOutputsFromWorkflowCall(workflow1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(outputs) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(outputs))
	}
	if outputs["output1"].Description != "First output" {
		t.Errorf("Expected output1 description to be 'First output', got '%s'", outputs["output1"].Description)
	}
	if outputs["output1"].Value != "${{ jobs.build.outputs.result }}" {
		t.Errorf("Expected output1 value to be '${{ jobs.build.outputs.result }}', got '%s'", outputs["output1"].Value)
	}

	// Test without workflow_call
	workflow2 := &ActionFile{
		On: map[string]interface{}{
			"push": nil,
		},
	}
	outputs2, err := ExtractOutputsFromWorkflowCall(workflow2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if outputs2 != nil {
		t.Errorf("Expected nil outputs for non-reusable workflow")
	}

	// Test without outputs in workflow_call
	workflow3 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{},
		},
	}
	outputs3, err := ExtractOutputsFromWorkflowCall(workflow3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if outputs3 != nil {
		t.Errorf("Expected nil outputs when no outputs defined")
	}

	// Test with invalid workflow_call type
	workflow4 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": "invalid",
		},
	}
	_, err = ExtractOutputsFromWorkflowCall(workflow4)
	if err == nil {
		t.Errorf("Expected error for invalid workflow_call type")
	}

	// Test with invalid outputs type
	workflow5 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"outputs": "invalid",
			},
		},
	}
	_, err = ExtractOutputsFromWorkflowCall(workflow5)
	if err == nil {
		t.Errorf("Expected error for invalid outputs type")
	}

	// Test with invalid output definition type
	workflow6 := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"outputs": map[string]interface{}{
					"output1": "invalid",
				},
			},
		},
	}
	_, err = ExtractOutputsFromWorkflowCall(workflow6)
	if err == nil {
		t.Errorf("Expected error for invalid output definition type")
	}
}

// TestValidateCompositeAction tests validation of composite actions
func TestValidateCompositeAction(t *testing.T) {
	// Test valid composite action
	action := &ActionFile{
		Name:        "Composite Action",
		Description: "A composite action",
		Runs: RunsConfig{
			Using: "composite",
			Steps: []Step{
				{
					Name: "Step 1",
					Run:  "echo 'Hello'",
				},
			},
		},
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) > 0 {
		t.Errorf("Expected no validation errors for valid composite action, got %d: %v", len(errors), errors)
	}

	// Test composite action without steps
	action.Runs.Steps = []Step{}
	errors = validator.Validate(action)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error for composite action without steps, got %d", len(errors))
	} else if errors[0].Field != "runs.steps" {
		t.Errorf("Expected validation error for runs.steps, got %s", errors[0].Field)
	}
}

// TestValidateUnsupportedActionType tests validation of unsupported action types
func TestValidateUnsupportedActionType(t *testing.T) {
	action := &ActionFile{
		Name:        "Unsupported Action",
		Description: "An action with unsupported type",
		Runs: RunsConfig{
			Using: "unsupported-type",
		},
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error for unsupported action type, got %d", len(errors))
	} else if errors[0].Field != "runs.using" {
		t.Errorf("Expected validation error for runs.using, got %s", errors[0].Field)
	}
}

// TestValidateNode20Action tests validation of Node.js 20 actions
func TestValidateNode20Action(t *testing.T) {
	// Test valid Node.js 20 action
	action := &ActionFile{
		Name:        "Node.js 20 Action",
		Description: "A Node.js 20 action",
		Runs: RunsConfig{
			Using: "node20",
			Main:  "dist/index.js",
		},
	}

	validator := NewValidator()
	errors := validator.Validate(action)

	if len(errors) > 0 {
		t.Errorf("Expected no validation errors for valid Node.js 20 action, got %d: %v", len(errors), errors)
	}

	// Test Node.js 20 action without main
	action.Runs.Main = ""
	errors = validator.Validate(action)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error for Node.js 20 action without main, got %d", len(errors))
	} else if errors[0].Field != "runs.main" {
		t.Errorf("Expected validation error for runs.main, got %s", errors[0].Field)
	}
}

// TestParseDirWithSubdirectories tests ParseDir with subdirectories
func TestParseDirWithSubdirectories(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create YAML files in different directories
	yamlContent := `name: Test Action
description: A test action
runs:
  using: composite
  steps:
    - name: Test step
      run: echo "test"
`

	// File in root directory
	err = os.WriteFile(filepath.Join(tempDir, "action.yml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// File in subdirectory
	err = os.WriteFile(filepath.Join(subDir, "nested.yaml"), []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create nested test file: %v", err)
	}

	// Non-YAML file (should be ignored)
	err = os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("This is not YAML"), 0644)
	if err != nil {
		t.Fatalf("Failed to create non-YAML file: %v", err)
	}

	// Parse the directory
	actions, err := ParseDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to parse directory: %v", err)
	}

	// Should have parsed 2 YAML files
	if len(actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(actions))
	}

	// Check that both files were parsed correctly
	foundRoot := false
	foundNested := false
	for path, action := range actions {
		if path == "action.yml" {
			foundRoot = true
		} else if path == filepath.Join("subdir", "nested.yaml") {
			foundNested = true
		}
		if action.Name != "Test Action" {
			t.Errorf("Expected action name to be 'Test Action', got '%s'", action.Name)
		}
	}

	if !foundRoot {
		t.Errorf("Expected to find action.yml in results")
	}
	if !foundNested {
		t.Errorf("Expected to find nested.yaml in results")
	}
}

// TestParseDirWithInvalidYAML tests ParseDir with invalid YAML files
func TestParseDirWithInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()

	// Create an invalid YAML file
	invalidYAML := `name: Test Action
description: A test action
runs:
  using: composite
  steps:
    - name: Test step
      run: echo "test"
    invalid_yaml_here: [unclosed bracket
`

	err := os.WriteFile(filepath.Join(tempDir, "invalid.yml"), []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid YAML file: %v", err)
	}

	// ParseDir should return an error for invalid YAML
	_, err = ParseDir(tempDir)
	if err == nil {
		t.Errorf("Expected error when parsing directory with invalid YAML")
	}
}

// Benchmark tests for performance measurement

// BenchmarkParseFile benchmarks the ParseFile function
func BenchmarkParseFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseFile("testdata/action.yml")
		if err != nil {
			b.Fatalf("Failed to parse file: %v", err)
		}
	}
}

// BenchmarkParseWorkflow benchmarks parsing workflow files
func BenchmarkParseWorkflow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParseFile("testdata/workflow.yml")
		if err != nil {
			b.Fatalf("Failed to parse workflow: %v", err)
		}
	}
}

// BenchmarkValidateAction benchmarks the validation function
func BenchmarkValidateAction(b *testing.B) {
	action, err := ParseFile("testdata/action.yml")
	if err != nil {
		b.Fatalf("Failed to parse action for benchmark: %v", err)
	}

	validator := NewValidator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		validator.Validate(action)
	}
}

// BenchmarkMapOfStringInterface benchmarks the map conversion utility
func BenchmarkMapOfStringInterface(b *testing.B) {
	input := map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MapOfStringInterface(input)
		if err != nil {
			b.Fatalf("Failed to convert map: %v", err)
		}
	}
}

// BenchmarkIsReusableWorkflow benchmarks the IsReusableWorkflow function
func BenchmarkIsReusableWorkflow(b *testing.B) {
	workflow := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{},
			"push":          nil,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsReusableWorkflow(workflow)
	}
}

// BenchmarkExtractInputsFromWorkflowCall benchmarks the input extraction function
func BenchmarkExtractInputsFromWorkflowCall(b *testing.B) {
	workflow := &ActionFile{
		On: map[string]interface{}{
			"workflow_call": map[string]interface{}{
				"inputs": map[string]interface{}{
					"input1": map[string]interface{}{
						"description": "First input",
						"required":    true,
						"default":     "default1",
					},
					"input2": map[string]interface{}{
						"description": "Second input",
						"required":    false,
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractInputsFromWorkflowCall(workflow)
		if err != nil {
			b.Fatalf("Failed to extract inputs: %v", err)
		}
	}
}
