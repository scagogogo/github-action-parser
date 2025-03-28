package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: utils_example <command> [args...]")
		fmt.Println("\nCommands:")
		fmt.Println("  1. parse_dir <directory_path> - 解析目录中的所有 Action/Workflow 文件")
		fmt.Println("  2. check_reusable <workflow_file> - 检查工作流是否可重用")
		fmt.Println("  3. extract_inputs <workflow_file> - 提取工作流的输入参数")
		fmt.Println("  4. extract_outputs <workflow_file> - 提取工作流的输出参数")
		fmt.Println("  5. convert_map <map_json_string> - 转换字符串或接口映射")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "parse_dir":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing directory path")
			fmt.Println("Usage: utils_example parse_dir <directory_path>")
			os.Exit(1)
		}
		parseDirectory(os.Args[2])

	case "check_reusable":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing workflow file path")
			fmt.Println("Usage: utils_example check_reusable <workflow_file>")
			os.Exit(1)
		}
		checkReusableWorkflow(os.Args[2])

	case "extract_inputs":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing workflow file path")
			fmt.Println("Usage: utils_example extract_inputs <workflow_file>")
			os.Exit(1)
		}
		extractInputs(os.Args[2])

	case "extract_outputs":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing workflow file path")
			fmt.Println("Usage: utils_example extract_outputs <workflow_file>")
			os.Exit(1)
		}
		extractOutputs(os.Args[2])

	case "convert_map":
		if len(os.Args) < 3 {
			fmt.Println("Error: Missing JSON map string")
			fmt.Println("Usage: utils_example convert_map '{\"key\":\"value\"}'")
			os.Exit(1)
		}
		convertMap(os.Args[2])

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Run 'utils_example' without arguments to see available commands")
		os.Exit(1)
	}
}

// parseDirectory 解析目录中的所有 Action/Workflow 文件
func parseDirectory(dirPath string) {
	fmt.Printf("解析目录: %s\n\n", dirPath)

	actions, err := parser.ParseDir(dirPath)
	if err != nil {
		fmt.Printf("Error: 解析目录失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个 Action/Workflow 文件:\n\n", len(actions))

	for path, action := range actions {
		fmt.Printf("文件: %s\n", path)
		fmt.Printf("  名称: %s\n", action.Name)

		// 检测文件类型
		fileType := "未知类型"
		if action.Runs.Using != "" {
			fileType = fmt.Sprintf("GitHub Action (%s)", action.Runs.Using)
		} else if len(action.Jobs) > 0 {
			fileType = "GitHub Workflow"
			if parser.IsReusableWorkflow(action) {
				fileType = "可重用 GitHub Workflow"
			}
		}
		fmt.Printf("  类型: %s\n", fileType)

		// 显示关键统计信息
		if len(action.Inputs) > 0 {
			fmt.Printf("  输入参数: %d 个\n", len(action.Inputs))
		}
		if len(action.Outputs) > 0 {
			fmt.Printf("  输出参数: %d 个\n", len(action.Outputs))
		}
		if len(action.Jobs) > 0 {
			fmt.Printf("  作业数量: %d 个\n", len(action.Jobs))

			totalSteps := 0
			for _, job := range action.Jobs {
				totalSteps += len(job.Steps)
			}
			fmt.Printf("  总步骤数: %d 个\n", totalSteps)
		}
		fmt.Println()
	}
}

// checkReusableWorkflow 检查工作流是否可重用
func checkReusableWorkflow(filePath string) {
	fmt.Printf("检查工作流是否可重用: %s\n\n", filePath)

	workflow, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error: 解析文件失败: %v\n", err)
		os.Exit(1)
	}

	isReusable := parser.IsReusableWorkflow(workflow)

	if isReusable {
		fmt.Println("✅ 这是一个可重用的工作流")
		fmt.Println("\n可重用工作流特点:")
		fmt.Println("- 包含 'workflow_call' 触发器")
		fmt.Println("- 可以被其他工作流通过 'uses' 关键字调用")
		fmt.Println("- 可以定义输入参数、密钥和输出参数")

		// 提取并显示调用信息
		var callTriggerMap map[string]interface{}

		switch on := workflow.On.(type) {
		case map[string]interface{}:
			if workflowCall, ok := on["workflow_call"]; ok {
				if wcMap, err := parser.MapOfStringInterface(workflowCall); err == nil {
					callTriggerMap = wcMap
				}
			}
		}

		if callTriggerMap != nil {
			fmt.Println("\n工作流调用配置:")
			for key, value := range callTriggerMap {
				fmt.Printf("- %s: %v\n", key, simplifyValue(value))
			}
		}
	} else {
		fmt.Println("❌ 这不是一个可重用的工作流")
		fmt.Println("\n要使工作流可重用，需要添加 'workflow_call' 触发器:")
		fmt.Println("```yaml")
		fmt.Println("on:")
		fmt.Println("  workflow_call:")
		fmt.Println("    inputs:")
		fmt.Println("      example-input:")
		fmt.Println("        description: '示例输入参数'")
		fmt.Println("        required: false")
		fmt.Println("        default: 'default-value'")
		fmt.Println("    secrets:")
		fmt.Println("      example-secret:")
		fmt.Println("        description: '示例密钥'")
		fmt.Println("        required: false")
		fmt.Println("    outputs:")
		fmt.Println("      example-output:")
		fmt.Println("        description: '示例输出参数'")
		fmt.Println("        value: ${{ jobs.example-job.outputs.result }}")
		fmt.Println("```")
	}
}

// extractInputs 提取工作流的输入参数
func extractInputs(filePath string) {
	fmt.Printf("提取工作流输入参数: %s\n\n", filePath)

	workflow, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error: 解析文件失败: %v\n", err)
		os.Exit(1)
	}

	if !parser.IsReusableWorkflow(workflow) {
		fmt.Println("❌ 这不是一个可重用的工作流，没有输入参数")
		os.Exit(1)
	}

	inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
	if err != nil {
		fmt.Printf("Error: 提取输入参数失败: %v\n", err)
		os.Exit(1)
	}

	if len(inputs) == 0 {
		fmt.Println("此工作流没有定义输入参数")
	} else {
		fmt.Printf("找到 %d 个输入参数:\n\n", len(inputs))

		// 表格头部
		fmt.Printf("%-20s %-10s %-20s %s\n", "名称", "是否必填", "默认值", "描述")
		fmt.Printf("%-20s %-10s %-20s %s\n", strings.Repeat("-", 20), strings.Repeat("-", 10), strings.Repeat("-", 20), strings.Repeat("-", 30))

		for name, input := range inputs {
			required := "否"
			if input.Required {
				required = "是"
			}

			defaultValue := "-"
			if input.Default != "" {
				defaultValue = input.Default
			}

			description := input.Description
			if len(description) > 30 {
				description = description[:27] + "..."
			}

			fmt.Printf("%-20s %-10s %-20s %s\n", name, required, defaultValue, description)
		}

		// 按 JSON 格式输出完整信息
		fmt.Println("\nJSON 格式的输入参数:")
		jsonData, _ := json.MarshalIndent(inputs, "", "  ")
		fmt.Println(string(jsonData))
	}
}

// extractOutputs 提取工作流的输出参数
func extractOutputs(filePath string) {
	fmt.Printf("提取工作流输出参数: %s\n\n", filePath)

	workflow, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error: 解析文件失败: %v\n", err)
		os.Exit(1)
	}

	if !parser.IsReusableWorkflow(workflow) {
		fmt.Println("❌ 这不是一个可重用的工作流，没有输出参数")
		os.Exit(1)
	}

	outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
	if err != nil {
		fmt.Printf("Error: 提取输出参数失败: %v\n", err)
		os.Exit(1)
	}

	if len(outputs) == 0 {
		fmt.Println("此工作流没有定义输出参数")
	} else {
		fmt.Printf("找到 %d 个输出参数:\n\n", len(outputs))

		for name, output := range outputs {
			fmt.Printf("名称: %s\n", name)
			fmt.Printf("  描述: %s\n", output.Description)
			fmt.Printf("  值来源: %s\n", output.Value)
			fmt.Println()
		}

		// 分析输出参数来源
		fmt.Println("输出参数来源分析:")
		for name, output := range outputs {
			jobName := extractJobNameFromOutputValue(output.Value)
			fmt.Printf("- %s: 来自作业 %s\n", name, jobName)
		}
	}
}

// convertMap 演示如何使用 MapOfStringInterface 和 MapOfStringString 函数
func convertMap(jsonStr string) {
	fmt.Println("演示地图转换工具函数:")
	fmt.Printf("输入 JSON: %s\n\n", jsonStr)

	// 解析 JSON 到 map[string]interface{}
	var inputMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &inputMap)
	if err != nil {
		fmt.Printf("Error: 无法解析 JSON: %v\n", err)
		os.Exit(1)
	}

	// 使用 MapOfStringInterface
	fmt.Println("使用 MapOfStringInterface:")
	result1, err := parser.MapOfStringInterface(inputMap)
	if err != nil {
		fmt.Printf("Error: MapOfStringInterface 失败: %v\n", err)
	} else {
		printMap("结果", result1)
	}

	// 尝试使用 MapOfStringString
	fmt.Println("\n使用 MapOfStringString:")
	result2, err := parser.MapOfStringString(inputMap)
	if err != nil {
		fmt.Printf("注意: MapOfStringString 失败: %v\n", err)
		fmt.Println("这是预期的，如果输入包含非字符串值")

		// 创建只有字符串值的映射
		stringMap := make(map[string]interface{})
		for k, v := range inputMap {
			if str, ok := v.(string); ok {
				stringMap[k] = str
			} else {
				stringMap[k] = fmt.Sprintf("%v", v)
			}
		}

		fmt.Println("\n使用只包含字符串值的映射:")
		result3, err := parser.MapOfStringString(stringMap)
		if err != nil {
			fmt.Printf("Error: MapOfStringString 仍然失败: %v\n", err)
		} else {
			printMap("结果", result3)
		}
	} else {
		printMap("结果", result2)
	}

	// 演示使用 map[interface{}]interface{} 类型
	fmt.Println("\n使用 map[interface{}]interface{} 类型:")

	// 创建一个 map[interface{}]interface{}
	mixedMap := make(map[interface{}]interface{})
	for k, v := range inputMap {
		mixedMap[k] = v
	}

	// 使用 MapOfStringInterface
	fmt.Println("使用 MapOfStringInterface:")
	result4, err := parser.MapOfStringInterface(mixedMap)
	if err != nil {
		fmt.Printf("Error: MapOfStringInterface 失败: %v\n", err)
	} else {
		printMap("结果", result4)
	}
}

// 辅助函数

func simplifyValue(v interface{}) string {
	switch val := v.(type) {
	case map[string]interface{}:
		return fmt.Sprintf("对象 (包含 %d 个键)", len(val))
	case map[interface{}]interface{}:
		return fmt.Sprintf("混合对象 (包含 %d 个键)", len(val))
	case []interface{}:
		return fmt.Sprintf("数组 (包含 %d 个项)", len(val))
	case string:
		if len(val) > 30 {
			return val[:27] + "..."
		}
		return val
	default:
		return fmt.Sprintf("%v", v)
	}
}

func printMap(label string, m interface{}) {
	switch typedMap := m.(type) {
	case map[string]interface{}:
		fmt.Printf("%s:\n", label)
		for k, v := range typedMap {
			fmt.Printf("  %s: %v\n", k, simplifyValue(v))
		}
	case map[string]string:
		fmt.Printf("%s:\n", label)
		for k, v := range typedMap {
			fmt.Printf("  %s: %s\n", k, v)
		}
	default:
		fmt.Printf("%s: %v\n", label, m)
	}
}

func extractJobNameFromOutputValue(value string) string {
	// 尝试从 ${{ jobs.X.outputs.Y }} 格式中提取作业名称
	if strings.Contains(value, "jobs.") && strings.Contains(value, ".outputs.") {
		// 提取 jobs.X.outputs 部分
		start := strings.Index(value, "jobs.") + 5
		end := strings.Index(value[start:], ".") + start
		if start < end {
			return value[start:end]
		}
	}
	return "未知"
}
