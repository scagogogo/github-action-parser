package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ActionFile represents the structure of a GitHub Action YAML file
type ActionFile struct {
	Name        string                 `yaml:"name,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Author      string                 `yaml:"author,omitempty"`
	Inputs      map[string]Input       `yaml:"inputs,omitempty"`
	Outputs     map[string]Output      `yaml:"outputs,omitempty"`
	Runs        RunsConfig             `yaml:"runs,omitempty"`
	Branding    Branding               `yaml:"branding,omitempty"`
	On          interface{}            `yaml:"on,omitempty"`
	Jobs        map[string]Job         `yaml:"jobs,omitempty"`
	Env         map[string]string      `yaml:"env,omitempty"`
	Defaults    map[string]interface{} `yaml:"defaults,omitempty"`
	Permissions interface{}            `yaml:"permissions,omitempty"`
}

// Input represents an input parameter for the action
type Input struct {
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
	Default     string `yaml:"default,omitempty"`
	Deprecated  bool   `yaml:"deprecated,omitempty"`
}

// Output represents an output value from the action
type Output struct {
	Description string `yaml:"description,omitempty"`
	Value       string `yaml:"value,omitempty"`
}

// RunsConfig defines how the action is executed
type RunsConfig struct {
	Using      string                 `yaml:"using,omitempty"`
	Main       string                 `yaml:"main,omitempty"`
	Pre        string                 `yaml:"pre,omitempty"`
	PreIf      string                 `yaml:"pre-if,omitempty"`
	Post       string                 `yaml:"post,omitempty"`
	PostIf     string                 `yaml:"post-if,omitempty"`
	Steps      []Step                 `yaml:"steps,omitempty"`
	Image      string                 `yaml:"image,omitempty"`
	Entrypoint string                 `yaml:"entrypoint,omitempty"`
	Args       []string               `yaml:"args,omitempty"`
	Env        map[string]string      `yaml:"env,omitempty"`
	Shell      string                 `yaml:"shell,omitempty"`
	Command    string                 `yaml:"command,omitempty"`
	With       map[string]interface{} `yaml:"with,omitempty"`
}

// Step represents a single step in a workflow job
type Step struct {
	ID         string                 `yaml:"id,omitempty"`
	If         string                 `yaml:"if,omitempty"`
	Name       string                 `yaml:"name,omitempty"`
	Uses       string                 `yaml:"uses,omitempty"`
	Run        string                 `yaml:"run,omitempty"`
	Shell      string                 `yaml:"shell,omitempty"`
	With       map[string]interface{} `yaml:"with,omitempty"`
	Env        map[string]string      `yaml:"env,omitempty"`
	ContinueOn interface{}            `yaml:"continue-on-error,omitempty"`
	TimeoutMin int                    `yaml:"timeout-minutes,omitempty"`
	WorkingDir string                 `yaml:"working-directory,omitempty"`
}

// Job represents a workflow job
type Job struct {
	Name           string                 `yaml:"name,omitempty"`
	Needs          interface{}            `yaml:"needs,omitempty"`
	RunsOn         interface{}            `yaml:"runs-on,omitempty"`
	Container      interface{}            `yaml:"container,omitempty"`
	Services       map[string]interface{} `yaml:"services,omitempty"`
	Outputs        map[string]string      `yaml:"outputs,omitempty"`
	Env            map[string]string      `yaml:"env,omitempty"`
	Defaults       map[string]interface{} `yaml:"defaults,omitempty"`
	If             string                 `yaml:"if,omitempty"`
	Steps          []Step                 `yaml:"steps,omitempty"`
	TimeoutMin     int                    `yaml:"timeout-minutes,omitempty"`
	Strategy       map[string]interface{} `yaml:"strategy,omitempty"`
	ContinueOn     interface{}            `yaml:"continue-on-error,omitempty"`
	Permissions    interface{}            `yaml:"permissions,omitempty"`
	ConcurrencyKey string                 `yaml:"concurrency,omitempty"`
	Uses           string                 `yaml:"uses,omitempty"`
	With           map[string]interface{} `yaml:"with,omitempty"`
	Secrets        interface{}            `yaml:"secrets,omitempty"`
}

// Branding defines the visual branding of the action
type Branding struct {
	Icon  string `yaml:"icon,omitempty"`
	Color string `yaml:"color,omitempty"`
}

// ParseFile parses a GitHub Action YAML file at the specified path
func ParseFile(path string) (*ActionFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return Parse(file)
}

// Parse parses a GitHub Action YAML from an io.Reader
func Parse(r io.Reader) (*ActionFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	var action ActionFile
	if err := yaml.Unmarshal(data, &action); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &action, nil
}

// ParseDir parses all GitHub Action YAML files in a directory recursively
func ParseDir(dir string) (map[string]*ActionFile, error) {
	result := make(map[string]*ActionFile)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process YAML files
		ext := filepath.Ext(path)
		if ext != ".yml" && ext != ".yaml" {
			return nil
		}

		action, err := ParseFile(path)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		relativePath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		result[relativePath] = action
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return result, nil
}
