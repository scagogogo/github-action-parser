# Types Reference

This page documents all the data structures used by the GitHub Action Parser library.

## ActionFile

The main structure representing a GitHub Action or Workflow file.

```go
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
```

### Fields

- **Name** (`string`): The name of the action or workflow
- **Description** (`string`): A description of what the action or workflow does
- **Author** (`string`): The author of the action
- **Inputs** (`map[string]Input`): Input parameters for the action
- **Outputs** (`map[string]Output`): Output values from the action
- **Runs** (`RunsConfig`): Configuration for how the action runs
- **Branding** (`Branding`): Branding information for the action
- **On** (`interface{}`): Trigger events for workflows
- **Jobs** (`map[string]Job`): Jobs defined in a workflow
- **Env** (`map[string]string`): Environment variables
- **Defaults** (`map[string]interface{}`): Default settings
- **Permissions** (`interface{}`): Permission settings

### Usage Example

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Action: %s\n", action.Name)
fmt.Printf("Description: %s\n", action.Description)
fmt.Printf("Author: %s\n", action.Author)
```

## Input

Represents an input parameter for an action.

```go
type Input struct {
    Description string `yaml:"description,omitempty"`
    Required    bool   `yaml:"required,omitempty"`
    Default     string `yaml:"default,omitempty"`
    Deprecated  bool   `yaml:"deprecated,omitempty"`
}
```

### Fields

- **Description** (`string`): Description of the input parameter
- **Required** (`bool`): Whether the input is required
- **Default** (`string`): Default value if not provided
- **Deprecated** (`bool`): Whether the input is deprecated

### Usage Example

```go
for name, input := range action.Inputs {
    fmt.Printf("Input: %s\n", name)
    fmt.Printf("  Description: %s\n", input.Description)
    fmt.Printf("  Required: %t\n", input.Required)
    if input.Default != "" {
        fmt.Printf("  Default: %s\n", input.Default)
    }
}
```

## Output

Represents an output value from an action.

```go
type Output struct {
    Description string `yaml:"description,omitempty"`
    Value       string `yaml:"value,omitempty"`
}
```

### Fields

- **Description** (`string`): Description of the output
- **Value** (`string`): The value expression for the output

### Usage Example

```go
for name, output := range action.Outputs {
    fmt.Printf("Output: %s\n", name)
    fmt.Printf("  Description: %s\n", output.Description)
    fmt.Printf("  Value: %s\n", output.Value)
}
```

## RunsConfig

Configuration for how an action executes.

```go
type RunsConfig struct {
    Using      string            `yaml:"using,omitempty"`
    Main       string            `yaml:"main,omitempty"`
    Pre        string            `yaml:"pre,omitempty"`
    Post       string            `yaml:"post,omitempty"`
    Image      string            `yaml:"image,omitempty"`
    Entrypoint string            `yaml:"entrypoint,omitempty"`
    Args       []string          `yaml:"args,omitempty"`
    Env        map[string]string `yaml:"env,omitempty"`
    Steps      []Step            `yaml:"steps,omitempty"`
}
```

### Fields

- **Using** (`string`): The runtime to use (e.g., "node20", "docker", "composite")
- **Main** (`string`): Main entry point for JavaScript actions
- **Pre** (`string`): Pre-execution script for JavaScript actions
- **Post** (`string`): Post-execution script for JavaScript actions
- **Image** (`string`): Docker image for Docker actions
- **Entrypoint** (`string`): Docker entrypoint
- **Args** (`[]string`): Arguments for Docker actions
- **Env** (`map[string]string`): Environment variables
- **Steps** (`[]Step`): Steps for composite actions

### Usage Example

```go
runs := action.Runs
fmt.Printf("Runtime: %s\n", runs.Using)

switch runs.Using {
case "node20", "node16":
    fmt.Printf("Main script: %s\n", runs.Main)
case "docker":
    fmt.Printf("Docker image: %s\n", runs.Image)
case "composite":
    fmt.Printf("Steps: %d\n", len(runs.Steps))
}
```

## Job

Represents a job in a GitHub workflow.

```go
type Job struct {
    Name        string                 `yaml:"name,omitempty"`
    RunsOn      interface{}            `yaml:"runs-on,omitempty"`
    Needs       interface{}            `yaml:"needs,omitempty"`
    If          string                 `yaml:"if,omitempty"`
    Steps       []Step                 `yaml:"steps,omitempty"`
    Env         map[string]string      `yaml:"env,omitempty"`
    Defaults    map[string]interface{} `yaml:"defaults,omitempty"`
    Outputs     map[string]string      `yaml:"outputs,omitempty"`
    TimeoutMin  int                    `yaml:"timeout-minutes,omitempty"`
    Strategy    interface{}            `yaml:"strategy,omitempty"`
    ContinueOn  interface{}            `yaml:"continue-on-error,omitempty"`
    Container   interface{}            `yaml:"container,omitempty"`
    Services    map[string]interface{} `yaml:"services,omitempty"`
    Uses        string                 `yaml:"uses,omitempty"`
    With        map[string]interface{} `yaml:"with,omitempty"`
    Secrets     interface{}            `yaml:"secrets,omitempty"`
    Permissions interface{}            `yaml:"permissions,omitempty"`
}
```

### Fields

- **Name** (`string`): Display name for the job
- **RunsOn** (`interface{}`): Runner environment (string or array)
- **Needs** (`interface{}`): Jobs that must complete before this job
- **If** (`string`): Conditional expression for job execution
- **Steps** (`[]Step`): Steps to execute in the job
- **Env** (`map[string]string`): Environment variables
- **Defaults** (`map[string]interface{}`): Default settings
- **Outputs** (`map[string]string`): Job outputs
- **TimeoutMin** (`int`): Timeout in minutes
- **Strategy** (`interface{}`): Matrix strategy configuration
- **ContinueOn** (`interface{}`): Continue on error setting
- **Container** (`interface{}`): Container configuration
- **Services** (`map[string]interface{}`): Service containers
- **Uses** (`string`): Reusable workflow reference
- **With** (`map[string]interface{}`): Inputs for reusable workflows
- **Secrets** (`interface{}`): Secrets for reusable workflows
- **Permissions** (`interface{}`): Permission settings

### Usage Example

```go
for jobID, job := range workflow.Jobs {
    fmt.Printf("Job: %s\n", jobID)
    if job.Name != "" {
        fmt.Printf("  Name: %s\n", job.Name)
    }
    fmt.Printf("  Steps: %d\n", len(job.Steps))

    if job.Uses != "" {
        fmt.Printf("  Uses: %s\n", job.Uses)
    }
}
```

## Step

Represents a single step in a workflow job.

```go
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
```

### Fields

- **ID** (`string`): Unique identifier for the step
- **If** (`string`): Conditional expression for step execution
- **Name** (`string`): Display name for the step
- **Uses** (`string`): Action to use
- **Run** (`string`): Command to run
- **Shell** (`string`): Shell to use for run commands
- **With** (`map[string]interface{}`): Input parameters for actions
- **Env** (`map[string]string`): Environment variables
- **ContinueOn** (`interface{}`): Continue on error setting
- **TimeoutMin** (`int`): Timeout in minutes
- **WorkingDir** (`string`): Working directory

### Usage Example

```go
for i, step := range job.Steps {
    fmt.Printf("Step %d:\n", i+1)
    if step.Name != "" {
        fmt.Printf("  Name: %s\n", step.Name)
    }
    if step.Uses != "" {
        fmt.Printf("  Uses: %s\n", step.Uses)
    }
    if step.Run != "" {
        fmt.Printf("  Run: %s\n", step.Run)
    }
}
```

## Branding

Branding information for GitHub Actions.

```go
type Branding struct {
    Icon  string `yaml:"icon,omitempty"`
    Color string `yaml:"color,omitempty"`
}
```

### Fields

- **Icon** (`string`): Icon name from Feather icons
- **Color** (`string`): Background color (white, yellow, blue, green, orange, red, purple, gray-dark)

### Usage Example

```go
if action.Branding.Icon != "" {
    fmt.Printf("Icon: %s\n", action.Branding.Icon)
}
if action.Branding.Color != "" {
    fmt.Printf("Color: %s\n", action.Branding.Color)
}
```

## StringOrStringSlice

A utility type that can represent either a string or a slice of strings, commonly used in YAML files.

```go
type StringOrStringSlice struct {
    Value  string
    Values []string
}
```

### Fields

- **Value** (`string`): Single string value (first item if slice)
- **Values** (`[]string`): All values as a slice

### Methods

#### UnmarshalYAML

```go
func (s *StringOrStringSlice) UnmarshalYAML(unmarshal func(interface{}) error) error
```

Implements YAML unmarshaling to handle both string and string slice inputs.

#### Contains

```go
func (s *StringOrStringSlice) Contains(value string) bool
```

Checks if the given value is contained in the string or string slice.

#### String

```go
func (s *StringOrStringSlice) String() string
```

Returns a string representation. For single values, returns the value. For multiple values, returns a comma-separated list.

### Usage Example

```go
// This type is typically used internally by the parser
// but can be useful for custom processing

var trigger StringOrStringSlice
// Can unmarshal from: "push" or ["push", "pull_request"]

if trigger.Contains("push") {
    fmt.Println("Triggered by push events")
}

fmt.Printf("Triggers: %s\n", trigger.String())
```

## Type Conversion Notes

### Interface{} Fields

Several fields use `interface{}` to accommodate the flexible nature of YAML:

- **On**: Can be a string, array, or complex object with event configurations
- **RunsOn**: Can be a string (single runner) or array (multiple runners)
- **Needs**: Can be a string (single dependency) or array (multiple dependencies)
- **Permissions**: Can be a string ("read-all", "write-all") or object with specific permissions

### Working with Interface{} Fields

Use type assertions or the utility functions to work with these fields:

```go
// Example: Working with the 'on' field
switch on := workflow.On.(type) {
case string:
    fmt.Printf("Single trigger: %s\n", on)
case []interface{}:
    fmt.Printf("Multiple triggers: %v\n", on)
case map[string]interface{}:
    fmt.Printf("Complex trigger configuration\n")
}
```

The library provides utility functions like `MapOfStringInterface` and `MapOfStringString` to help with type conversions. See the [Utilities API](/api/utilities) for more details.
