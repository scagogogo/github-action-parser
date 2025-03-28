package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: analyze_reusable_workflow <path_to_workflow_yml>")
		fmt.Println("Example: analyze_reusable_workflow ../../pkg/parser/testdata/reusable-workflow.yml")
		os.Exit(1)
	}

	// 解析工作流文件
	filePath := os.Args[1]
	workflow, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing workflow file: %v\n", err)
		os.Exit(1)
	}

	// 检查是否是可重用工作流
	if !parser.IsReusableWorkflow(workflow) {
		fmt.Println("❌ 这个文件不是可重用工作流！")
		fmt.Println("可重用工作流必须包含 'workflow_call' 触发器。")
		os.Exit(1)
	}

	// 显示基本信息
	fmt.Println("==== 可重用工作流信息 ====")
	fmt.Printf("名称: %s\n", workflow.Name)
	if workflow.Description != "" {
		fmt.Printf("描述: %s\n", workflow.Description)
	}

	// 分析和显示输入参数
	fmt.Println("\n==== 输入参数 ====")
	inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
	if err != nil {
		fmt.Printf("无法提取输入参数: %v\n", err)
	} else if len(inputs) == 0 {
		fmt.Println("没有定义输入参数")
	} else {
		requiredCount := 0
		optionalCount := 0
		withDefaultCount := 0

		for name, input := range inputs {
			required := "可选"
			if input.Required {
				required = "必填"
				requiredCount++
			} else {
				optionalCount++
			}

			hasDefault := "无默认值"
			if input.Default != "" {
				hasDefault = fmt.Sprintf("默认值: %s", input.Default)
				withDefaultCount++
			}

			fmt.Printf("- %s (%s, %s)\n", name, required, hasDefault)
			fmt.Printf("  描述: %s\n", input.Description)
		}

		fmt.Printf("\n输入参数统计：共 %d 个，其中必填 %d 个，可选 %d 个，有默认值 %d 个\n",
			len(inputs), requiredCount, optionalCount, withDefaultCount)
	}

	// 分析和显示密钥
	fmt.Println("\n==== 密钥 ====")
	secrets := extractSecretsFromWorkflowCall(workflow)
	if len(secrets) == 0 {
		fmt.Println("没有定义密钥")
	} else {
		for name, secret := range secrets {
			required := "可选"
			if secret.Required {
				required = "必填"
			}

			fmt.Printf("- %s (%s)\n", name, required)
			if secret.Description != "" {
				fmt.Printf("  描述: %s\n", secret.Description)
			}
		}
	}

	// 分析和显示输出
	fmt.Println("\n==== 输出参数 ====")
	outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
	if err != nil {
		fmt.Printf("无法提取输出参数: %v\n", err)
	} else if len(outputs) == 0 {
		fmt.Println("没有定义输出参数")
	} else {
		for name, output := range outputs {
			fmt.Printf("- %s\n", name)
			fmt.Printf("  描述: %s\n", output.Description)
			fmt.Printf("  值来源: %s\n", output.Value)
		}
	}

	// 分析作业数量和类型
	fmt.Println("\n==== 作业分析 ====")
	if len(workflow.Jobs) == 0 {
		fmt.Println("没有定义作业")
	} else {
		fmt.Printf("包含 %d 个作业:\n", len(workflow.Jobs))

		// 分析每个作业
		for jobID, job := range workflow.Jobs {
			fmt.Printf("\n作业: %s\n", jobID)
			if job.Name != "" {
				fmt.Printf("名称: %s\n", job.Name)
			}

			// 统计步骤信息
			if len(job.Steps) > 0 {
				checkoutSteps := 0
				setupSteps := 0
				buildSteps := 0
				deploySteps := 0
				testSteps := 0
				otherSteps := 0

				for _, step := range job.Steps {
					if step.Uses != "" && strings.Contains(step.Uses, "checkout") {
						checkoutSteps++
					} else if step.Name != "" && (strings.Contains(strings.ToLower(step.Name), "setup") ||
						strings.Contains(strings.ToLower(step.Name), "install") ||
						strings.Contains(strings.ToLower(step.Name), "config")) {
						setupSteps++
					} else if step.Name != "" && (strings.Contains(strings.ToLower(step.Name), "build") ||
						strings.Contains(strings.ToLower(step.Name), "compile")) {
						buildSteps++
					} else if step.Name != "" && (strings.Contains(strings.ToLower(step.Name), "deploy") ||
						strings.Contains(strings.ToLower(step.Name), "publish") ||
						strings.Contains(strings.ToLower(step.Name), "release")) {
						deploySteps++
					} else if step.Name != "" && (strings.Contains(strings.ToLower(step.Name), "test") ||
						strings.Contains(strings.ToLower(step.Name), "lint") ||
						strings.Contains(strings.ToLower(step.Name), "check")) {
						testSteps++
					} else {
						otherSteps++
					}
				}

				fmt.Printf("步骤数量: %d\n", len(job.Steps))
				fmt.Printf("- 检出代码步骤: %d\n", checkoutSteps)
				fmt.Printf("- 环境设置步骤: %d\n", setupSteps)
				fmt.Printf("- 构建步骤: %d\n", buildSteps)
				fmt.Printf("- 部署步骤: %d\n", deploySteps)
				fmt.Printf("- 测试/检查步骤: %d\n", testSteps)
				fmt.Printf("- 其他步骤: %d\n", otherSteps)
			}

			// 检查作业是否提供输出
			if job.Outputs != nil && len(job.Outputs) > 0 {
				fmt.Printf("\n作业输出:\n")
				for outName, outValue := range job.Outputs {
					fmt.Printf("- %s: %s\n", outName, outValue)
				}
			}
		}
	}

	// 提供用法建议
	fmt.Println("\n==== 使用建议 ====")
	fmt.Println("在其他工作流中引用此可重用工作流的示例:")
	fmt.Println("```yaml")
	fmt.Println("jobs:")
	fmt.Println("  call-workflow:")
	fmt.Println("    uses: ./.github/workflows/" + filepath.Base(filePath)) // 使用filepath包获取文件名
	fmt.Println("    with:")

	// 添加必填输入
	hasRequiredInputs := false
	for name, input := range inputs {
		if input.Required {
			if !hasRequiredInputs {
				hasRequiredInputs = true
			}
			fmt.Printf("      %s: # 这里填写你的值\n", name)
		}
	}

	// 如果有密钥，添加密钥部分
	hasRequiredSecrets := false
	for name, secret := range secrets {
		if secret.Required {
			if !hasRequiredSecrets {
				fmt.Println("    secrets:")
				hasRequiredSecrets = true
			}
			fmt.Printf("      %s: ${{ secrets.%s }}\n", name, name)
		}
	}

	fmt.Println("```")
}

// Secret 表示工作流可以使用的密钥
type Secret struct {
	Description string
	Required    bool
}

// extractSecretsFromWorkflowCall 提取工作流调用中定义的密钥
func extractSecretsFromWorkflowCall(action *parser.ActionFile) map[string]Secret {
	secrets := make(map[string]Secret)

	switch on := action.On.(type) {
	case map[string]interface{}:
		workflowCall, ok := on["workflow_call"]
		if !ok {
			return secrets
		}

		workflowCallMap, err := parser.MapOfStringInterface(workflowCall)
		if err != nil {
			return secrets
		}

		secretsRaw, ok := workflowCallMap["secrets"]
		if !ok {
			return secrets
		}

		secretsMap, err := parser.MapOfStringInterface(secretsRaw)
		if err != nil {
			return secrets
		}

		for name, def := range secretsMap {
			secretDef, err := parser.MapOfStringInterface(def)
			if err != nil {
				continue
			}

			secret := Secret{}
			if desc, ok := secretDef["description"].(string); ok {
				secret.Description = desc
			}
			if required, ok := secretDef["required"].(bool); ok {
				secret.Required = required
			}

			secrets[name] = secret
		}
	}

	return secrets
}
