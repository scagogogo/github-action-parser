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
		fmt.Println("Usage: validate_action <path_to_action_or_workflow_yml>")
		fmt.Println("Example: validate_action ../../pkg/parser/testdata/action.yml")
		os.Exit(1)
	}

	// 解析文件
	filePath := os.Args[1]
	action, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// 创建验证器并验证
	validator := parser.NewValidator()
	errors := validator.Validate(action)

	// 显示基本信息
	fmt.Println("==== 文件信息 ====")
	fmt.Printf("文件名: %s\n", filePath)
	fmt.Printf("名称: %s\n", action.Name)

	// 判断文件类型
	fileType := "未知类型"
	if action.Runs.Using != "" {
		fileType = "GitHub Action"
		fmt.Printf("Action 类型: %s\n", action.Runs.Using)
	} else if len(action.Jobs) > 0 {
		fileType = "GitHub Workflow"
		if parser.IsReusableWorkflow(action) {
			fileType = "可重用 GitHub Workflow"
		}
	}
	fmt.Printf("文件类型: %s\n", fileType)

	// 显示验证结果
	fmt.Println("\n==== 验证结果 ====")
	if len(errors) == 0 {
		fmt.Println("✓ 验证通过，没有发现问题")
	} else {
		fmt.Printf("✗ 发现 %d 个问题:\n", len(errors))
		for i, err := range errors {
			fmt.Printf("%d. 字段: %s\n   错误: %s\n", i+1, err.Field, err.Message)
		}
	}

	// 如果有错误，介绍如何修复
	if len(errors) > 0 {
		fmt.Println("\n==== 修复建议 ====")
		for _, err := range errors {
			suggestFix(err, action)
		}
	}
}

// 根据错误类型提供修复建议
func suggestFix(err parser.ValidationError, action *parser.ActionFile) {
	switch err.Field {
	case "name":
		fmt.Println("- 问题: 缺少名称")
		fmt.Println("  建议: 添加一个有意义的名称，例如:")
		fmt.Println("    name: 'My GitHub Action'")

	case "description":
		fmt.Println("- 问题: 缺少描述")
		fmt.Println("  建议: 添加一个清晰的描述，例如:")
		fmt.Println("    description: '这个 Action 的用途是...'")

	case "runs.using":
		fmt.Println("- 问题: 未指定运行环境")
		fmt.Println("  建议: 指定一个有效的运行环境，可选值包括 'composite', 'node16', 'node20', 'docker'")
		fmt.Println("    runs:")
		fmt.Println("      using: 'composite'")

	case "runs.main":
		fmt.Println("- 问题: JavaScript Action 缺少主入口点")
		fmt.Println("  建议: 添加一个 main 入口点文件，例如:")
		fmt.Println("    runs:")
		fmt.Println("      using: 'node16'")
		fmt.Println("      main: 'dist/index.js'")

	case "runs.image":
		fmt.Println("- 问题: Docker Action 缺少镜像")
		fmt.Println("  建议: 指定一个 Docker 镜像，例如:")
		fmt.Println("    runs:")
		fmt.Println("      using: 'docker'")
		fmt.Println("      image: 'Dockerfile'")
		fmt.Println("  或者:")
		fmt.Println("    runs:")
		fmt.Println("      using: 'docker'")
		fmt.Println("      image: 'docker://alpine:latest'")

	case "runs.steps":
		fmt.Println("- 问题: Composite Action 缺少步骤")
		fmt.Println("  建议: 至少添加一个步骤，例如:")
		fmt.Println("    runs:")
		fmt.Println("      using: 'composite'")
		fmt.Println("      steps:")
		fmt.Println("        - name: '执行任务'")
		fmt.Println("          run: 'echo \"Hello, World!\"'")
		fmt.Println("          shell: 'bash'")

	case "on":
		fmt.Println("- 问题: Workflow 缺少触发器")
		fmt.Println("  建议: 添加至少一个触发事件，例如:")
		fmt.Println("    on:")
		fmt.Println("      push:")
		fmt.Println("        branches: [ main ]")
		fmt.Println("      pull_request:")
		fmt.Println("        branches: [ main ]")

	case "jobs":
		fmt.Println("- 问题: Workflow 缺少作业")
		fmt.Println("  建议: 添加至少一个作业，例如:")
		fmt.Println("    jobs:")
		fmt.Println("      build:")
		fmt.Println("        runs-on: ubuntu-latest")
		fmt.Println("        steps:")
		fmt.Println("          - uses: actions/checkout@v3")
		fmt.Println("          - name: '运行构建'")
		fmt.Println("            run: 'echo \"Building...\"'")

	default:
		if err.Field[:4] == "jobs" {
			parts := err.Field
			if strings.Contains(parts, "steps") {
				fmt.Printf("- 问题: 作业的步骤配置有误 (%s)\n", err.Field)
				fmt.Println("  建议: 确保每个步骤都有 uses 或 run 字段")
				fmt.Println("    steps:")
				fmt.Println("      - name: '示例步骤'")
				fmt.Println("        run: 'echo \"Hello\"'")
				fmt.Println("    或:")
				fmt.Println("      - name: '示例步骤'")
				fmt.Println("        uses: 'actions/checkout@v3'")
			} else {
				fmt.Printf("- 问题: 作业配置有误 (%s)\n", err.Field)
				fmt.Println("  建议: 确保每个作业都有 runs-on 或 uses 字段")
				fmt.Println("    jobs:")
				fmt.Println("      example-job:")
				fmt.Println("        runs-on: ubuntu-latest")
				fmt.Println("        # ... 其他配置")
			}
		} else {
			fmt.Printf("- 问题: %s - %s\n", err.Field, err.Message)
			fmt.Println("  建议: 查阅 GitHub Actions 文档以了解正确的配置格式")
		}
	}
}
