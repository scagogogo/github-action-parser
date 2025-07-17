# 工作流解析

此示例演示如何解析 GitHub Workflow 文件并提取作业信息、步骤和触发器。

## 解析 Workflow 文件

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 workflow 文件
    workflow, err := parser.ParseFile(".github/workflows/ci.yml")
    if err != nil {
        log.Fatalf("解析 workflow 失败: %v", err)
    }

    fmt.Printf("Workflow: %s\n", workflow.Name)
    fmt.Printf("作业数: %d\n", len(workflow.Jobs))
}
```

## 访问 Workflow 触发器

```go
// 访问 workflow 触发器（on 事件）
fmt.Println("\n触发器:")
switch on := workflow.On.(type) {
case string:
    fmt.Printf("  单个触发器: %s\n", on)
case []interface{}:
    fmt.Printf("  多个触发器:\n")
    for _, trigger := range on {
        fmt.Printf("    - %v\n", trigger)
    }
case map[string]interface{}:
    fmt.Printf("  复杂触发器:\n")
    for event, config := range on {
        fmt.Printf("    - %s: %v\n", event, config)
    }
}
```

## 分析作业

```go
// 分析每个作业
fmt.Println("\n作业:")
for jobID, job := range workflow.Jobs {
    fmt.Printf("  %s:\n", jobID)
    if job.Name != "" {
        fmt.Printf("    名称: %s\n", job.Name)
    }
    
    // 检查 runs-on
    switch runsOn := job.RunsOn.(type) {
    case string:
        fmt.Printf("    运行在: %s\n", runsOn)
    case []interface{}:
        fmt.Printf("    运行在: %v\n", runsOn)
    }
    
    // 检查依赖关系
    if job.Needs != nil {
        switch needs := job.Needs.(type) {
        case string:
            fmt.Printf("    需要: %s\n", needs)
        case []interface{}:
            fmt.Printf("    需要: %v\n", needs)
        }
    }
    
    fmt.Printf("    步骤数: %d\n", len(job.Steps))
}
```

## 分析作业步骤

```go
// 详细步骤分析
for jobID, job := range workflow.Jobs {
    fmt.Printf("\n=== 作业: %s ===\n", jobID)
    
    for i, step := range job.Steps {
        fmt.Printf("步骤 %d:\n", i+1)
        
        if step.Name != "" {
            fmt.Printf("  名称: %s\n", step.Name)
        }
        
        if step.ID != "" {
            fmt.Printf("  ID: %s\n", step.ID)
        }
        
        if step.Uses != "" {
            fmt.Printf("  使用: %s\n", step.Uses)
            
            // 显示步骤输入
            if len(step.With) > 0 {
                fmt.Printf("  输入:\n")
                for key, value := range step.With {
                    fmt.Printf("    %s: %v\n", key, value)
                }
            }
        }
        
        if step.Run != "" {
            fmt.Printf("  运行: %s\n", step.Run)
            if step.Shell != "" {
                fmt.Printf("  Shell: %s\n", step.Shell)
            }
        }
        
        if step.If != "" {
            fmt.Printf("  条件: %s\n", step.If)
        }
        
        // 显示步骤环境变量
        if len(step.Env) > 0 {
            fmt.Printf("  环境变量:\n")
            for key, value := range step.Env {
                fmt.Printf("    %s: %s\n", key, value)
            }
        }
        
        fmt.Println()
    }
}
```

## 检查可重用工作流

```go
// 检查是否有作业使用可重用工作流
fmt.Println("\n可重用工作流作业:")
for jobID, job := range workflow.Jobs {
    if job.Uses != "" {
        fmt.Printf("  %s 使用: %s\n", jobID, job.Uses)
        
        // 显示传递给可重用工作流的输入
        if len(job.With) > 0 {
            fmt.Printf("    输入:\n")
            for key, value := range job.With {
                fmt.Printf("      %s: %v\n", key, value)
            }
        }
        
        // 显示传递给可重用工作流的密钥
        if job.Secrets != nil {
            fmt.Printf("    密钥: %v\n", job.Secrets)
        }
    }
}
```

## 解析多个工作流

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 .github/workflows 目录中的所有工作流
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("解析 workflows 失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个工作流文件\n\n", len(workflows))
    
    for path, workflow := range workflows {
        analyzeWorkflow(path, workflow)
    }
}

func analyzeWorkflow(path string, workflow *parser.ActionFile) {
    fmt.Printf("=== %s ===\n", filepath.Base(path))
    
    if workflow.Name != "" {
        fmt.Printf("名称: %s\n", workflow.Name)
    }
    
    // 计算不同类型的触发器
    triggerCount := countTriggers(workflow.On)
    fmt.Printf("触发器: %d\n", triggerCount)
    
    // 分析作业
    fmt.Printf("作业: %d\n", len(workflow.Jobs))
    
    totalSteps := 0
    reusableJobs := 0
    
    for _, job := range workflow.Jobs {
        totalSteps += len(job.Steps)
        if job.Uses != "" {
            reusableJobs++
        }
    }
    
    fmt.Printf("总步骤数: %d\n", totalSteps)
    if reusableJobs > 0 {
        fmt.Printf("可重用工作流作业: %d\n", reusableJobs)
    }
    
    // 检查这是否是可重用工作流
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("类型: 可重用工作流\n")
        
        inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
        outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow)
        
        fmt.Printf("输入: %d\n", len(inputs))
        fmt.Printf("输出: %d\n", len(outputs))
    }
    
    fmt.Println()
}

func countTriggers(on interface{}) int {
    switch triggers := on.(type) {
    case string:
        return 1
    case []interface{}:
        return len(triggers)
    case map[string]interface{}:
        return len(triggers)
    default:
        return 0
    }
}
```

## 环境变量和密钥

```go
// 访问工作流级别的环境变量
if len(workflow.Env) > 0 {
    fmt.Println("\n工作流环境变量:")
    for key, value := range workflow.Env {
        fmt.Printf("  %s: %s\n", key, value)
    }
}

// 访问作业级别的环境变量
for jobID, job := range workflow.Jobs {
    if len(job.Env) > 0 {
        fmt.Printf("\n作业 %s 环境变量:\n", jobID)
        for key, value := range job.Env {
            fmt.Printf("  %s: %s\n", key, value)
        }
    }
}
```

## 完整的工作流分析示例

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
        log.Fatal("用法: go run main.go <workflow.yml>")
    }
    
    workflowFile := os.Args[1]
    
    workflow, err := parser.ParseFile(workflowFile)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", workflowFile, err)
    }
    
    analyzeCompleteWorkflow(workflow)
}

func analyzeCompleteWorkflow(workflow *parser.ActionFile) {
    fmt.Printf("=== 工作流分析 ===\n")
    
    if workflow.Name != "" {
        fmt.Printf("名称: %s\n", workflow.Name)
    }
    
    // 分析触发器
    fmt.Printf("\n触发器:\n")
    analyzeTriggers(workflow.On)
    
    // 分析作业
    fmt.Printf("\n作业 (%d):\n", len(workflow.Jobs))
    for jobID, job := range workflow.Jobs {
        analyzeJob(jobID, job)
    }
    
    // 检查是否可重用
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("\n=== 可重用工作流 ===\n")
        analyzeReusableWorkflow(workflow)
    }
}

func analyzeTriggers(on interface{}) {
    switch triggers := on.(type) {
    case string:
        fmt.Printf("  - %s\n", triggers)
    case []interface{}:
        for _, trigger := range triggers {
            fmt.Printf("  - %v\n", trigger)
        }
    case map[string]interface{}:
        for event, config := range triggers {
            fmt.Printf("  - %s:\n", event)
            if configMap, ok := config.(map[string]interface{}); ok {
                for key, value := range configMap {
                    fmt.Printf("      %s: %v\n", key, value)
                }
            }
        }
    }
}

func analyzeJob(jobID string, job parser.Job) {
    fmt.Printf("\n  %s:\n", jobID)
    
    if job.Name != "" {
        fmt.Printf("    名称: %s\n", job.Name)
    }
    
    if job.Uses != "" {
        fmt.Printf("    使用: %s (可重用工作流)\n", job.Uses)
    } else {
        fmt.Printf("    步骤: %d\n", len(job.Steps))
        
        // 显示前几个步骤
        for i, step := range job.Steps {
            if i >= 3 { // 限制为前 3 个步骤
                fmt.Printf("    ... 还有 %d 个步骤\n", len(job.Steps)-3)
                break
            }
            
            stepDesc := fmt.Sprintf("步骤 %d", i+1)
            if step.Name != "" {
                stepDesc = step.Name
            } else if step.Uses != "" {
                stepDesc = fmt.Sprintf("使用 %s", step.Uses)
            } else if step.Run != "" {
                stepDesc = "运行命令"
            }
            
            fmt.Printf("      %d. %s\n", i+1, stepDesc)
        }
    }
}

func analyzeReusableWorkflow(workflow *parser.ActionFile) {
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err == nil && len(inputs) > 0 {
        fmt.Printf("输入 (%d):\n", len(inputs))
        for name, input := range inputs {
            required := ""
            if input.Required {
                required = " (必需)"
            }
            fmt.Printf("  - %s%s: %s\n", name, required, input.Description)
        }
    }
    
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err == nil && len(outputs) > 0 {
        fmt.Printf("输出 (%d):\n", len(outputs))
        for name, output := range outputs {
            fmt.Printf("  - %s: %s\n", name, output.Description)
        }
    }
}
```

## 下一步

- 了解工作流的 [验证](/zh/examples/validation)
- 详细探索 [可重用工作流](/zh/examples/reusable-workflows)
- 查看高级处理的 [工具函数](/zh/examples/utilities)
