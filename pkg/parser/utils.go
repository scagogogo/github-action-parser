package parser

import (
	"fmt"
	"strings"
)

// StringOrStringSlice represents a field that can be either a string or a slice of strings
// in GitHub Actions YAML. For example, 'on' can be either "push" or ["push", "pull_request"]
type StringOrStringSlice struct {
	Value  string
	Values []string
}

// UnmarshalYAML implements the yaml.Unmarshaler interface
func (s *StringOrStringSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to unmarshal as a string
	var str string
	if err := unmarshal(&str); err == nil {
		s.Value = str
		s.Values = []string{str}
		return nil
	}

	// Try to unmarshal as a slice of strings
	var slice []string
	if err := unmarshal(&slice); err == nil {
		s.Values = slice
		if len(slice) > 0 {
			s.Value = slice[0]
		}
		return nil
	}

	return fmt.Errorf("must be a string or a slice of strings")
}

// Contains checks if a string is in the StringOrStringSlice
func (s *StringOrStringSlice) Contains(value string) bool {
	for _, v := range s.Values {
		if v == value {
			return true
		}
	}
	return false
}

// String returns a string representation of the StringOrStringSlice
func (s *StringOrStringSlice) String() string {
	if len(s.Values) == 1 {
		return s.Value
	}
	return strings.Join(s.Values, ", ")
}

// MapOfStringInterface converts a YAML map to map[string]interface{}
func MapOfStringInterface(v interface{}) (map[string]interface{}, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		return value, nil
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range value {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("map key must be a string")
			}
			result[key] = v
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to map[string]interface{}", v)
	}
}

// MapOfStringString converts a YAML map to map[string]string
func MapOfStringString(v interface{}) (map[string]string, error) {
	switch value := v.(type) {
	case map[string]string:
		return value, nil
	case map[string]interface{}:
		result := make(map[string]string)
		for k, v := range value {
			str, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("value for key %q must be a string", k)
			}
			result[k] = str
		}
		return result, nil
	case map[interface{}]interface{}:
		result := make(map[string]string)
		for k, v := range value {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("map key must be a string")
			}
			str, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("value for key %q must be a string", key)
			}
			result[key] = str
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to map[string]string", v)
	}
}

// IsReusableWorkflow checks if a workflow is intended to be called by other workflows
func IsReusableWorkflow(action *ActionFile) bool {
	// Check if workflow has workflow_call event
	switch t := action.On.(type) {
	case map[string]interface{}:
		for event := range t {
			if event == "workflow_call" {
				return true
			}
		}
	case map[interface{}]interface{}:
		for event := range t {
			if eventStr, ok := event.(string); ok && eventStr == "workflow_call" {
				return true
			}
		}
	}
	return false
}

// ExtractInputsFromWorkflowCall extracts input definitions from a reusable workflow
func ExtractInputsFromWorkflowCall(action *ActionFile) (map[string]Input, error) {
	inputs := make(map[string]Input)

	switch on := action.On.(type) {
	case map[string]interface{}:
		workflowCall, ok := on["workflow_call"]
		if !ok {
			return nil, nil
		}

		workflowCallMap, err := MapOfStringInterface(workflowCall)
		if err != nil {
			return nil, err
		}

		inputsRaw, ok := workflowCallMap["inputs"]
		if !ok {
			return nil, nil
		}

		inputsMap, err := MapOfStringInterface(inputsRaw)
		if err != nil {
			return nil, err
		}

		for name, def := range inputsMap {
			inputDef, err := MapOfStringInterface(def)
			if err != nil {
				return nil, err
			}

			input := Input{}
			if desc, ok := inputDef["description"].(string); ok {
				input.Description = desc
			}
			if required, ok := inputDef["required"].(bool); ok {
				input.Required = required
			}
			if defaultVal, ok := inputDef["default"].(string); ok {
				input.Default = defaultVal
			}

			inputs[name] = input
		}
	}

	return inputs, nil
}

// ExtractOutputsFromWorkflowCall extracts output definitions from a reusable workflow
func ExtractOutputsFromWorkflowCall(action *ActionFile) (map[string]Output, error) {
	outputs := make(map[string]Output)

	switch on := action.On.(type) {
	case map[string]interface{}:
		workflowCall, ok := on["workflow_call"]
		if !ok {
			return nil, nil
		}

		workflowCallMap, err := MapOfStringInterface(workflowCall)
		if err != nil {
			return nil, err
		}

		outputsRaw, ok := workflowCallMap["outputs"]
		if !ok {
			return nil, nil
		}

		outputsMap, err := MapOfStringInterface(outputsRaw)
		if err != nil {
			return nil, err
		}

		for name, def := range outputsMap {
			outputDef, err := MapOfStringInterface(def)
			if err != nil {
				return nil, err
			}

			output := Output{}
			if desc, ok := outputDef["description"].(string); ok {
				output.Description = desc
			}
			if value, ok := outputDef["value"].(string); ok {
				output.Value = value
			}

			outputs[name] = output
		}
	}

	return outputs, nil
}
