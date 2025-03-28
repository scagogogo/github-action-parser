package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_workflow <path_to_workflow_yml>")
		fmt.Println("Example: parse_workflow ../../pkg/parser/testdata/workflow.yml")
		os.Exit(1)
	}

	// 解析工作流文件
	filePath := os.Args[1]
	workflow, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing workflow file: %v\n", err)
		os.Exit(1)
	}

	// 显示基本信息
	fmt.Println("==== 工作流基本信息 ====")
	fmt.Printf("工作流名称: %s\n", workflow.Name)
	if workflow.Description != "" {
		fmt.Printf("描述: %s\n", workflow.Description)
	}

	// 显示触发器
	fmt.Println("\n==== 触发器 ====")
	displayTriggers(workflow.On)

	// 显示环境变量
	if len(workflow.Env) > 0 {
		fmt.Println("\n==== 全局环境变量 ====")
		for name, value := range workflow.Env {
			fmt.Printf("%s: %s\n", name, value)
		}
	}

	// 显示作业信息
	if len(workflow.Jobs) > 0 {
		fmt.Println("\n==== 作业 ====")
		for jobID, job := range workflow.Jobs {
			fmt.Printf("\n## %s", jobID)
			if job.Name != "" {
				fmt.Printf(" (%s)", job.Name)
			}
			fmt.Println()

			// 显示 runs-on
			if job.RunsOn != nil {
				fmt.Printf("运行环境: %v\n", formatRunsOn(job.RunsOn))
			}

			// 显示 needs
			if job.Needs != nil {
				fmt.Printf("依赖作业: %v\n", formatNeeds(job.Needs))
			}

			// 显示 if 条件
			if job.If != "" {
				fmt.Printf("条件: %s\n", job.If)
			}

			// 显示环境
			if job.Env != nil && len(job.Env) > 0 {
				fmt.Println("环境变量:")
				for name, value := range job.Env {
					fmt.Printf("  %s: %s\n", name, value)
				}
			}

			// 显示步骤
			if len(job.Steps) > 0 {
				fmt.Println("步骤:")
				for i, step := range job.Steps {
					stepDesc := fmt.Sprintf("  %d. ", i+1)
					if step.Name != "" {
						stepDesc += step.Name
					} else if step.Run != "" {
						runPart := step.Run
						if len(runPart) > 30 {
							runPart = runPart[:27] + "..."
						}
						stepDesc += fmt.Sprintf("运行: %s", runPart)
					} else if step.Uses != "" {
						stepDesc += fmt.Sprintf("使用: %s", step.Uses)
					} else {
						stepDesc += "<未命名步骤>"
					}
					fmt.Println(stepDesc)

					// 显示步骤的详细信息
					if step.If != "" {
						fmt.Printf("    条件: %s\n", step.If)
					}
					if step.ContinueOn != nil {
						fmt.Printf("    错误继续: %v\n", step.ContinueOn)
					}
					if step.WorkingDir != "" {
						fmt.Printf("    工作目录: %s\n", step.WorkingDir)
					}
				}
			}
		}
	}

	// 检查是否是可重用工作流
	if parser.IsReusableWorkflow(workflow) {
		fmt.Println("\n==== 可重用工作流信息 ====")
		fmt.Println("这是一个可重用的工作流")

		// 提取工作流输入参数
		inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
		if err != nil {
			fmt.Printf("提取输入参数时出错: %v\n", err)
		} else if len(inputs) > 0 {
			fmt.Println("\n工作流输入参数:")
			for name, input := range inputs {
				required := "可选"
				if input.Required {
					required = "必填"
				}

				defaultValue := "无"
				if input.Default != "" {
					defaultValue = input.Default
				}

				fmt.Printf("- %s (%s, 默认值: %s): %s\n",
					name, required, defaultValue, input.Description)
			}
		}

		// 提取工作流输出参数
		outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
		if err != nil {
			fmt.Printf("提取输出参数时出错: %v\n", err)
		} else if len(outputs) > 0 {
			fmt.Println("\n工作流输出参数:")
			for name, output := range outputs {
				fmt.Printf("- %s: %s\n", name, output.Description)
				if output.Value != "" {
					fmt.Printf("  值: %s\n", output.Value)
				}
			}
		}
	}
}

// 显示触发器信息
func displayTriggers(on interface{}) {
	switch v := on.(type) {
	case nil:
		fmt.Println("未定义触发器")
	case string:
		fmt.Printf("触发事件: %s\n", v)
	case []string:
		fmt.Printf("触发事件: %s\n", strings.Join(v, ", "))
	case map[string]interface{}:
		fmt.Println("触发事件:")
		for event, config := range v {
			if config == nil {
				fmt.Printf("- %s\n", event)
			} else {
				fmt.Printf("- %s (带配置)\n", event)
			}
		}
	case map[interface{}]interface{}:
		fmt.Println("触发事件:")
		for event, config := range v {
			eventStr, ok := event.(string)
			if !ok {
				eventStr = fmt.Sprintf("%v", event)
			}
			if config == nil {
				fmt.Printf("- %s\n", eventStr)
			} else {
				fmt.Printf("- %s (带配置)\n", eventStr)
			}
		}
	default:
		fmt.Printf("其他类型的触发器: %T\n", on)
	}
}

// 格式化 runs-on 字段
func formatRunsOn(runsOn interface{}) string {
	switch v := runsOn.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	case []interface{}:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
		return strings.Join(strs, ", ")
	default:
		return fmt.Sprintf("%v", runsOn)
	}
}

// 格式化 needs 字段
func formatNeeds(needs interface{}) string {
	switch v := needs.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	case []interface{}:
		strs := make([]string, len(v))
		for i, item := range v {
			strs[i] = fmt.Sprintf("%v", item)
		}
		return strings.Join(strs, ", ")
	default:
		return fmt.Sprintf("%v", needs)
	}
}
