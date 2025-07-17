# 可重用工作流

此示例演示如何处理可重用工作流，包括检测、输入/输出提取和分析。

## 检测可重用工作流

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析工作流文件
    workflow, err := parser.ParseFile(".github/workflows/reusable.yml")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    // 检查是否为可重用工作流
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("✅ 这是一个可重用工作流")
        analyzeReusableWorkflow(workflow)
    } else {
        fmt.Println("❌ 这不是可重用工作流")
    }
}

func analyzeReusableWorkflow(workflow *parser.ActionFile) {
    fmt.Printf("名称: %s\n", workflow.Name)
    
    // 提取输入
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        log.Printf("提取输入失败: %v", err)
    } else {
        fmt.Printf("输入: %d\n", len(inputs))
    }
    
    // 提取输出
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        log.Printf("提取输出失败: %v", err)
    } else {
        fmt.Printf("输出: %d\n", len(outputs))
    }
}
```

## 提取并显示输入

```go
// 提取并显示详细的输入信息
inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("提取输入失败: %v", err)
}

if len(inputs) > 0 {
    fmt.Printf("\n=== 输入 (%d) ===\n", len(inputs))
    for name, input := range inputs {
        fmt.Printf("• %s", name)
        
        if input.Required {
            fmt.Printf(" (必需)")
        } else {
            fmt.Printf(" (可选)")
        }
        
        fmt.Printf("\n  描述: %s\n", input.Description)
        
        if input.Default != "" {
            fmt.Printf("  默认值: %s\n", input.Default)
        }
        
        if input.Deprecated {
            fmt.Printf("  ⚠️  已弃用\n")
        }
        
        fmt.Println()
    }
} else {
    fmt.Println("未定义输入")
}
```

## 提取并显示输出

```go
// 提取并显示详细的输出信息
outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("提取输出失败: %v", err)
}

if len(outputs) > 0 {
    fmt.Printf("\n=== 输出 (%d) ===\n", len(outputs))
    for name, output := range outputs {
        fmt.Printf("• %s\n", name)
        fmt.Printf("  描述: %s\n", output.Description)
        
        if output.Value != "" {
            fmt.Printf("  值: %s\n", output.Value)
        }
        
        fmt.Println()
    }
} else {
    fmt.Println("未定义输出")
}
```

## 查找所有可重用工作流

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析所有工作流
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    fmt.Printf("扫描 %d 个工作流文件...\n\n", len(workflows))
    
    reusableCount := 0
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableCount++
            analyzeReusableWorkflowDetailed(filepath.Base(path), workflow)
        }
    }
    
    fmt.Printf("\n=== 摘要 ===\n")
    fmt.Printf("总工作流: %d\n", len(workflows))
    fmt.Printf("可重用工作流: %d\n", reusableCount)
    fmt.Printf("常规工作流: %d\n", len(workflows)-reusableCount)
}

func analyzeReusableWorkflowDetailed(filename string, workflow *parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filename)
    
    if workflow.Name != "" {
        fmt.Printf("名称: %s\n", workflow.Name)
    }
    
    // 分析输入
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ 提取输入失败: %v\n", err)
    } else {
        fmt.Printf("输入: %d", len(inputs))
        if len(inputs) > 0 {
            requiredCount := 0
            for _, input := range inputs {
                if input.Required {
                    requiredCount++
                }
            }
            fmt.Printf(" (%d 必需, %d 可选)", requiredCount, len(inputs)-requiredCount)
        }
        fmt.Println()
    }
    
    // 分析输出
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ 提取输出失败: %v\n", err)
    } else {
        fmt.Printf("输出: %d\n", len(outputs))
    }
    
    // 分析作业
    fmt.Printf("作业: %d\n", len(workflow.Jobs))
    
    fmt.Println()
}
```

## 验证可重用工作流使用

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析所有工作流以查找可重用工作流使用
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    // 查找可重用工作流及其调用者
    reusableWorkflows := make(map[string]*parser.ActionFile)
    callerWorkflows := make(map[string]*parser.ActionFile)
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableWorkflows[path] = workflow
        } else {
            // 检查此工作流是否调用任何可重用工作流
            for _, job := range workflow.Jobs {
                if job.Uses != "" {
                    callerWorkflows[path] = workflow
                    break
                }
            }
        }
    }
    
    fmt.Printf("找到 %d 个可重用工作流和 %d 个调用者工作流\n\n", 
        len(reusableWorkflows), len(callerWorkflows))
    
    // 分析使用情况
    for callerPath, callerWorkflow := range callerWorkflows {
        analyzeReusableWorkflowUsage(callerPath, callerWorkflow, reusableWorkflows)
    }
}

func analyzeReusableWorkflowUsage(callerPath string, caller *parser.ActionFile, reusableWorkflows map[string]*parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filepath.Base(callerPath))
    
    for jobID, job := range caller.Jobs {
        if job.Uses == "" {
            continue
        }
        
        fmt.Printf("作业 '%s' 使用: %s\n", jobID, job.Uses)
        
        // 检查是否为本地可重用工作流
        if strings.HasPrefix(job.Uses, "./") {
            localPath := strings.TrimPrefix(job.Uses, "./")
            if reusableWorkflow, exists := reusableWorkflows[localPath]; exists {
                validateReusableWorkflowCall(jobID, job, reusableWorkflow)
            } else {
                fmt.Printf("  ⚠️  未找到本地可重用工作流: %s\n", localPath)
            }
        }
        
        // 显示传递给可重用工作流的输入
        if len(job.With) > 0 {
            fmt.Printf("  传递的输入:\n")
            for key, value := range job.With {
                fmt.Printf("    %s: %v\n", key, value)
            }
        }
        
        // 显示传递给可重用工作流的密钥
        if job.Secrets != nil {
            fmt.Printf("  密钥: %v\n", job.Secrets)
        }
        
        fmt.Println()
    }
}

func validateReusableWorkflowCall(jobID string, job parser.Job, reusableWorkflow *parser.ActionFile) {
    // 从可重用工作流提取预期输入
    expectedInputs, err := parser.ExtractInputsFromWorkflowCall(reusableWorkflow)
    if err != nil {
        fmt.Printf("  ❌ 提取预期输入失败: %v\n", err)
        return
    }
    
    // 检查是否提供了所有必需输入
    providedInputs := make(map[string]bool)
    for key := range job.With {
        providedInputs[key] = true
    }
    
    missingRequired := []string{}
    extraInputs := []string{}
    
    // 检查缺失的必需输入
    for name, input := range expectedInputs {
        if input.Required && !providedInputs[name] {
            missingRequired = append(missingRequired, name)
        }
    }
    
    // 检查额外输入
    for name := range providedInputs {
        if _, exists := expectedInputs[name]; !exists {
            extraInputs = append(extraInputs, name)
        }
    }
    
    // 报告验证结果
    if len(missingRequired) == 0 && len(extraInputs) == 0 {
        fmt.Printf("  ✅ 输入验证通过\n")
    } else {
        if len(missingRequired) > 0 {
            fmt.Printf("  ❌ 缺少必需输入: %v\n", missingRequired)
        }
        if len(extraInputs) > 0 {
            fmt.Printf("  ⚠️  额外输入（未定义）: %v\n", extraInputs)
        }
    }
}
```

## 生成可重用工作流文档

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 查找所有可重用工作流
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    reusableWorkflows := make(map[string]*parser.ActionFile)
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableWorkflows[path] = workflow
        }
    }
    
    if len(reusableWorkflows) == 0 {
        fmt.Println("未找到可重用工作流")
        return
    }
    
    // 生成文档
    generateReusableWorkflowDocs(reusableWorkflows)
}

func generateReusableWorkflowDocs(workflows map[string]*parser.ActionFile) {
    fmt.Println("# 可重用工作流文档\n")
    
    // 按名称排序工作流
    var paths []string
    for path := range workflows {
        paths = append(paths, path)
    }
    sort.Strings(paths)
    
    for _, path := range paths {
        workflow := workflows[path]
        generateWorkflowDoc(filepath.Base(path), workflow)
    }
}

func generateWorkflowDoc(filename string, workflow *parser.ActionFile) {
    fmt.Printf("## %s\n\n", strings.TrimSuffix(filename, filepath.Ext(filename)))
    
    if workflow.Name != "" {
        fmt.Printf("**名称:** %s\n\n", workflow.Name)
    }
    
    // 提取输入
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ 提取输入失败: %v\n\n", err)
    } else if len(inputs) > 0 {
        fmt.Printf("### 输入\n\n")
        fmt.Printf("| 名称 | 描述 | 必需 | 默认值 |\n")
        fmt.Printf("|------|------|------|--------|\n")
        
        // 按名称排序输入
        var inputNames []string
        for name := range inputs {
            inputNames = append(inputNames, name)
        }
        sort.Strings(inputNames)
        
        for _, name := range inputNames {
            input := inputs[name]
            required := "否"
            if input.Required {
                required = "是"
            }
            
            defaultValue := input.Default
            if defaultValue == "" {
                defaultValue = "-"
            }
            
            fmt.Printf("| `%s` | %s | %s | `%s` |\n", 
                name, input.Description, required, defaultValue)
        }
        fmt.Println()
    }
    
    // 提取输出
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        fmt.Printf("❌ 提取输出失败: %v\n\n", err)
    } else if len(outputs) > 0 {
        fmt.Printf("### 输出\n\n")
        fmt.Printf("| 名称 | 描述 |\n")
        fmt.Printf("|------|------|\n")
        
        // 按名称排序输出
        var outputNames []string
        for name := range outputs {
            outputNames = append(outputNames, name)
        }
        sort.Strings(outputNames)
        
        for _, name := range outputNames {
            output := outputs[name]
            fmt.Printf("| `%s` | %s |\n", name, output.Description)
        }
        fmt.Println()
    }
    
    // 使用示例
    fmt.Printf("### 使用示例\n\n")
    fmt.Printf("```yaml\n")
    fmt.Printf("jobs:\n")
    fmt.Printf("  call-reusable-workflow:\n")
    fmt.Printf("    uses: ./.github/workflows/%s\n", filename)
    
    if len(inputs) > 0 {
        fmt.Printf("    with:\n")
        for name, input := range inputs {
            if input.Required {
                fmt.Printf("      %s: # 必需 - %s\n", name, input.Description)
            }
        }
    }
    
    fmt.Printf("```\n\n")
    
    fmt.Println("---\n")
}
```

## 下一步

- 了解高级处理的 [工具函数](/zh/examples/utilities)
- 查看可重用工作流函数的 [API 参考](/zh/api/utilities)
- 探索全面工作流验证的 [验证功能](/zh/examples/validation)
