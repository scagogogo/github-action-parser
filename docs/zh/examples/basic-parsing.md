# 基本解析

此示例演示了 GitHub Action Parser 库的基本解析功能。

## 解析 Action 文件

最基本的操作是解析 GitHub Action 文件：

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 action.yml 文件
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatalf("解析 action 失败: %v", err)
    }

    // 显示基本信息
    fmt.Printf("Action 名称: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
    fmt.Printf("作者: %s\n", action.Author)
}
```

## 访问输入参数

```go
// 访问输入参数
fmt.Println("\n输入:")
for name, input := range action.Inputs {
    fmt.Printf("  %s:\n", name)
    fmt.Printf("    描述: %s\n", input.Description)
    fmt.Printf("    必需: %t\n", input.Required)
    if input.Default != "" {
        fmt.Printf("    默认值: %s\n", input.Default)
    }
    if input.Deprecated {
        fmt.Printf("    已弃用: true\n")
    }
}
```

## 访问输出参数

```go
// 访问输出参数
fmt.Println("\n输出:")
for name, output := range action.Outputs {
    fmt.Printf("  %s:\n", name)
    fmt.Printf("    描述: %s\n", output.Description)
    if output.Value != "" {
        fmt.Printf("    值: %s\n", output.Value)
    }
}
```

## 检查 Action 类型

```go
// 检查这是什么类型的 action
fmt.Printf("\nAction 类型: %s\n", action.Runs.Using)

switch action.Runs.Using {
case "composite":
    fmt.Printf("复合 action，包含 %d 个步骤\n", len(action.Runs.Steps))
case "docker":
    fmt.Printf("Docker action，使用镜像: %s\n", action.Runs.Image)
case "node16", "node20":
    fmt.Printf("JavaScript action，主文件: %s\n", action.Runs.Main)
}
```

## 访问品牌信息

```go
// 访问品牌信息
if action.Branding.Icon != "" || action.Branding.Color != "" {
    fmt.Println("\n品牌:")
    if action.Branding.Icon != "" {
        fmt.Printf("  图标: %s\n", action.Branding.Icon)
    }
    if action.Branding.Color != "" {
        fmt.Printf("  颜色: %s\n", action.Branding.Color)
    }
}
```

## 完整示例

这是一个演示所有基本解析功能的完整示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("用法: go run main.go <action.yml>")
    }
    
    actionFile := os.Args[1]
    
    // 解析 action 文件
    action, err := parser.ParseFile(actionFile)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", actionFile, err)
    }
    
    // 显示全面信息
    displayActionInfo(action)
}

func displayActionInfo(action *parser.ActionFile) {
    fmt.Printf("=== Action 信息 ===\n")
    fmt.Printf("名称: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
    
    if action.Author != "" {
        fmt.Printf("作者: %s\n", action.Author)
    }
    
    // 显示输入
    if len(action.Inputs) > 0 {
        fmt.Printf("\n=== 输入 (%d) ===\n", len(action.Inputs))
        for name, input := range action.Inputs {
            fmt.Printf("• %s", name)
            if input.Required {
                fmt.Printf(" (必需)")
            }
            fmt.Printf("\n  %s\n", input.Description)
            if input.Default != "" {
                fmt.Printf("  默认值: %s\n", input.Default)
            }
        }
    }
    
    // 显示输出
    if len(action.Outputs) > 0 {
        fmt.Printf("\n=== 输出 (%d) ===\n", len(action.Outputs))
        for name, output := range action.Outputs {
            fmt.Printf("• %s\n", name)
            fmt.Printf("  %s\n", output.Description)
            if output.Value != "" {
                fmt.Printf("  值: %s\n", output.Value)
            }
        }
    }
    
    // 显示运行时信息
    fmt.Printf("\n=== 运行时 ===\n")
    fmt.Printf("使用: %s\n", action.Runs.Using)
    
    switch action.Runs.Using {
    case "composite":
        fmt.Printf("步骤: %d\n", len(action.Runs.Steps))
        for i, step := range action.Runs.Steps {
            fmt.Printf("  %d. %s\n", i+1, step.Name)
        }
    case "docker":
        fmt.Printf("镜像: %s\n", action.Runs.Image)
        if action.Runs.Entrypoint != "" {
            fmt.Printf("入口点: %s\n", action.Runs.Entrypoint)
        }
    case "node16", "node20":
        fmt.Printf("主文件: %s\n", action.Runs.Main)
        if action.Runs.Pre != "" {
            fmt.Printf("预执行: %s\n", action.Runs.Pre)
        }
        if action.Runs.Post != "" {
            fmt.Printf("后执行: %s\n", action.Runs.Post)
        }
    }
    
    // 显示品牌
    if action.Branding.Icon != "" || action.Branding.Color != "" {
        fmt.Printf("\n=== 品牌 ===\n")
        if action.Branding.Icon != "" {
            fmt.Printf("图标: %s\n", action.Branding.Icon)
        }
        if action.Branding.Color != "" {
            fmt.Printf("颜色: %s\n", action.Branding.Color)
        }
    }
}
```

## 错误处理

解析文件时始终适当处理错误：

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    // 检查特定错误类型
    if os.IsNotExist(err) {
        log.Fatal("Action 文件不存在")
    } else if strings.Contains(err.Error(), "yaml") {
        log.Fatal("Action 文件中的 YAML 语法无效")
    } else {
        log.Fatalf("解析 action 失败: %v", err)
    }
}
```

## 下一步

- 了解 [工作流解析](/zh/examples/workflow-parsing)
- 探索 [验证](/zh/examples/validation) 功能
- 查看 [工具函数](/zh/examples/utilities)
