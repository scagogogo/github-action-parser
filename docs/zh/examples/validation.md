# 验证功能

此示例演示如何使用内置验证功能验证 GitHub Actions 和 Workflows。

## 基本验证

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 action 文件
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatalf("解析 action 失败: %v", err)
    }
    
    // 创建验证器并验证
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) == 0 {
        fmt.Println("✅ Action 有效!")
    } else {
        fmt.Printf("❌ 发现 %d 个验证错误:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    }
}
```

## 验证多个文件

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析当前目录中的所有 YAML 文件
    files, err := parser.ParseDir(".")
    if err != nil {
        log.Fatalf("解析目录失败: %v", err)
    }
    
    validator := parser.NewValidator()
    totalErrors := 0
    
    for path, action := range files {
        fmt.Printf("\n=== 验证 %s ===\n", filepath.Base(path))
        
        errors := validator.Validate(action)
        if len(errors) == 0 {
            fmt.Println("✅ 有效")
        } else {
            fmt.Printf("❌ %d 个错误:\n", len(errors))
            for _, err := range errors {
                fmt.Printf("  - %s: %s\n", err.Field, err.Message)
            }
            totalErrors += len(errors)
        }
    }
    
    fmt.Printf("\n=== 摘要 ===\n")
    fmt.Printf("检查的文件: %d\n", len(files))
    fmt.Printf("总错误数: %d\n", totalErrors)
}
```

## 详细报告验证

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
        log.Fatal("用法: go run main.go <file.yml>")
    }
    
    filename := os.Args[1]
    
    // 解析文件
    action, err := parser.ParseFile(filename)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", filename, err)
    }
    
    // 详细验证报告
    validateWithDetails(filename, action)
}

func validateWithDetails(filename string, action *parser.ActionFile) {
    fmt.Printf("=== 验证 %s ===\n\n", filename)
    
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) == 0 {
        fmt.Println("✅ 文件有效!")
        displayFileInfo(action)
        return
    }
    
    // 按类别分组错误
    fieldErrors := make(map[string][]parser.ValidationError)
    for _, err := range errors {
        fieldErrors[err.Field] = append(fieldErrors[err.Field], err)
    }
    
    fmt.Printf("❌ 发现 %d 个验证错误:\n\n", len(errors))
    
    // 按类别显示错误
    for field, errs := range fieldErrors {
        fmt.Printf("字段: %s\n", field)
        for _, err := range errs {
            fmt.Printf("  ❌ %s\n", err.Message)
            provideSuggestion(err)
        }
        fmt.Println()
    }
    
    // 提供一般建议
    fmt.Println("💡 一般建议:")
    provideGeneralSuggestions(action, errors)
}

func provideSuggestion(err parser.ValidationError) {
    suggestions := map[string]string{
        "name":        "为你的 action 添加描述性名称",
        "description": "添加清楚的描述，说明你的 action 的功能",
        "runs.using":  "指定支持的运行时: node16, node20, docker, 或 composite",
        "runs.main":   "对于 JavaScript actions，指定主入口点文件",
        "runs.image":  "对于 Docker actions，指定 Docker 镜像或 Dockerfile",
        "runs.steps":  "对于复合 actions，添加至少一个步骤",
        "on":          "为工作流指定至少一个触发事件",
        "jobs":        "为你的工作流添加至少一个作业",
    }
    
    if suggestion, exists := suggestions[err.Field]; exists {
        fmt.Printf("    💡 %s\n", suggestion)
    }
}

func provideGeneralSuggestions(action *parser.ActionFile, errors []parser.ValidationError) {
    // 根据 action 类型建议
    if action.Runs.Using != "" {
        switch action.Runs.Using {
        case "composite":
            fmt.Println("  - 对于复合 actions，确保每个步骤都有 'uses' 或 'run'")
            fmt.Println("  - 考虑为步骤添加名称以提高可读性")
        case "docker":
            fmt.Println("  - 对于 Docker actions，确保你的 Dockerfile 存在")
            fmt.Println("  - 如需要，考虑指定入口点")
        case "node16", "node20":
            fmt.Println("  - 对于 JavaScript actions，确保你的主文件存在")
            fmt.Println("  - 如需要，考虑添加 pre/post 脚本")
        }
    }
    
    // 根据错误模式建议
    hasRequiredFieldErrors := false
    for _, err := range errors {
        if err.Field == "name" || err.Field == "description" {
            hasRequiredFieldErrors = true
            break
        }
    }
    
    if hasRequiredFieldErrors {
        fmt.Println("  - 必需字段（name, description）对 GitHub Actions 至关重要")
        fmt.Println("  - 这些帮助用户理解你的 action 的功能")
    }
    
    // 工作流特定建议
    if len(action.Jobs) > 0 {
        fmt.Println("  - 对于工作流，确保每个作业都有 'runs-on' 或 'uses'")
        fmt.Println("  - 检查所有引用的 actions 是否存在且可访问")
    }
}

func displayFileInfo(action *parser.ActionFile) {
    fmt.Println("\n📋 文件信息:")
    
    if action.Name != "" {
        fmt.Printf("  名称: %s\n", action.Name)
    }
    
    if action.Description != "" {
        fmt.Printf("  描述: %s\n", action.Description)
    }
    
    if action.Runs.Using != "" {
        fmt.Printf("  类型: %s action\n", action.Runs.Using)
    }
    
    if len(action.Jobs) > 0 {
        fmt.Printf("  类型: 包含 %d 个作业的工作流\n", len(action.Jobs))
    }
    
    if len(action.Inputs) > 0 {
        fmt.Printf("  输入: %d\n", len(action.Inputs))
    }
    
    if len(action.Outputs) > 0 {
        fmt.Printf("  输出: %d\n", len(action.Outputs))
    }
}
```

## 批量验证摘要

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 验证 .github/workflows 中的所有工作流
    validateDirectory(".github/workflows", "工作流")
    
    // 验证当前目录中的 action 文件
    validateDirectory(".", "Actions")
}

func validateDirectory(dir, category string) {
    fmt.Printf("\n=== 验证 %s 中的 %s ===\n", dir, category)
    
    files, err := parser.ParseDir(dir)
    if err != nil {
        log.Printf("解析 %s 失败: %v", dir, err)
        return
    }
    
    if len(files) == 0 {
        fmt.Printf("在 %s 中未找到 YAML 文件\n", dir)
        return
    }
    
    validator := parser.NewValidator()
    
    validFiles := 0
    totalErrors := 0
    errorsByType := make(map[string]int)
    
    for path, action := range files {
        errors := validator.Validate(action)
        
        filename := filepath.Base(path)
        if len(errors) == 0 {
            fmt.Printf("✅ %s\n", filename)
            validFiles++
        } else {
            fmt.Printf("❌ %s (%d 个错误)\n", filename, len(errors))
            totalErrors += len(errors)
            
            // 统计错误类型
            for _, err := range errors {
                errorsByType[err.Field]++
            }
        }
    }
    
    // 显示摘要
    fmt.Printf("\n📊 %s 摘要:\n", category)
    fmt.Printf("  总文件数: %d\n", len(files))
    fmt.Printf("  有效文件: %d\n", validFiles)
    fmt.Printf("  有错误的文件: %d\n", len(files)-validFiles)
    fmt.Printf("  总错误数: %d\n", totalErrors)
    
    if len(errorsByType) > 0 {
        fmt.Printf("\n🔍 最常见的错误:\n")
        for field, count := range errorsByType {
            fmt.Printf("  %s: %d 次出现\n", field, count)
        }
    }
}
```

## 下一步

- 了解 [可重用工作流](/zh/examples/reusable-workflows)
- 探索高级处理的 [工具函数](/zh/examples/utilities)
- 查看详细验证文档的 [API 参考](/zh/api/validation)
