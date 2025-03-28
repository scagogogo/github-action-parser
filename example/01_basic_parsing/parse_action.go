package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		fmt.Println("Usage: parse_action <path_to_action_yml>")
		fmt.Println("Example: parse_action ../../pkg/parser/testdata/action.yml")
		os.Exit(1)
	}

	// 解析 action.yml 文件
	filePath := os.Args[1]
	action, err := parser.ParseFile(filePath)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// 显示基本信息
	fmt.Println("==== 基本信息 ====")
	fmt.Printf("文件名: %s\n", filepath.Base(filePath))
	fmt.Printf("Action 名称: %s\n", action.Name)
	fmt.Printf("描述: %s\n", action.Description)
	if action.Author != "" {
		fmt.Printf("作者: %s\n", action.Author)
	}

	// 显示品牌相关信息
	if action.Branding.Icon != "" || action.Branding.Color != "" {
		fmt.Println("\n==== 品牌设置 ====")
		if action.Branding.Icon != "" {
			fmt.Printf("图标: %s\n", action.Branding.Icon)
		}
		if action.Branding.Color != "" {
			fmt.Printf("颜色: %s\n", action.Branding.Color)
		}
	}

	// 显示输入参数
	if len(action.Inputs) > 0 {
		fmt.Println("\n==== 输入参数 ====")
		for name, input := range action.Inputs {
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

	// 显示输出参数
	if len(action.Outputs) > 0 {
		fmt.Println("\n==== 输出参数 ====")
		for name, output := range action.Outputs {
			fmt.Printf("- %s: %s\n", name, output.Description)
			if output.Value != "" {
				fmt.Printf("  值: %s\n", output.Value)
			}
		}
	}

	// 显示执行配置
	fmt.Println("\n==== 执行配置 ====")
	fmt.Printf("运行类型: %s\n", action.Runs.Using)

	switch action.Runs.Using {
	case "composite":
		fmt.Printf("步骤数量: %d\n", len(action.Runs.Steps))
		if len(action.Runs.Steps) > 0 {
			fmt.Println("\n步骤列表:")
			for i, step := range action.Runs.Steps {
				fmt.Printf("%d. %s\n", i+1, step.Name)
			}
		}
	case "docker":
		fmt.Printf("Docker 镜像: %s\n", action.Runs.Image)
		if action.Runs.Entrypoint != "" {
			fmt.Printf("入口点: %s\n", action.Runs.Entrypoint)
		}
	case "node16", "node20":
		fmt.Printf("主脚本: %s\n", action.Runs.Main)
		if action.Runs.Pre != "" {
			fmt.Printf("预执行脚本: %s\n", action.Runs.Pre)
		}
		if action.Runs.Post != "" {
			fmt.Printf("后执行脚本: %s\n", action.Runs.Post)
		}
	}
}
