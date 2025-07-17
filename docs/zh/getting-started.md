# 快速开始

## 安装

使用 Go modules 安装 GitHub Action Parser 库：

```bash
go get github.com/scagogogo/github-action-parser
```

## 基本用法

在你的 Go 代码中导入 parser 包：

```go
import "github.com/scagogogo/github-action-parser/pkg/parser"
```

### 解析 Action 文件

```go
package main

import (
    "fmt"
    "os"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 action.yml 文件
    action, err := parser.ParseFile("path/to/action.yml")
    if err != nil {
        fmt.Printf("解析文件错误: %v\n", err)
        os.Exit(1)
    }

    // 访问 action 元数据
    fmt.Printf("Action 名称: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
    
    // 访问输入参数
    for name, input := range action.Inputs {
        fmt.Printf("输入 %s: %s (必填: %t)\n", 
            name, input.Description, input.Required)
    }
    
    // 访问输出参数
    for name, output := range action.Outputs {
        fmt.Printf("输出 %s: %s\n", name, output.Description)
    }
}
```

### 解析工作流文件

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 workflow.yml 文件
    workflow, err := parser.ParseFile("path/to/.github/workflows/ci.yml")
    if err != nil {
        fmt.Printf("解析文件错误: %v\n", err)
        return
    }
    
    // 访问工作流作业
    for jobId, job := range workflow.Jobs {
        fmt.Printf("作业: %s (%s)\n", jobId, job.Name)
        
        // 访问作业步骤
        for i, step := range job.Steps {
            if step.Name != "" {
                fmt.Printf("  步骤 %d: %s\n", i+1, step.Name)
            } else if step.Run != "" {
                fmt.Printf("  步骤 %d: 运行命令\n", i+1)
            } else if step.Uses != "" {
                fmt.Printf("  步骤 %d: 使用 %s\n", i+1, step.Uses)
            }
        }
    }
}
```

### 验证文件

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析任何 GitHub Action/Workflow 文件
    action, err := parser.ParseFile("path/to/file.yml")
    if err != nil {
        fmt.Printf("解析文件错误: %v\n", err)
        return
    }
    
    // 创建验证器并验证
    validator := parser.NewValidator()
    errors := validator.Validate(action)
    
    if len(errors) > 0 {
        fmt.Println("验证错误:")
        for _, err := range errors {
            fmt.Printf("- %s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Println("文件有效!")
    }
}
```

### 解析目录

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析目录中的所有工作流文件
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        fmt.Printf("解析工作流错误: %v\n", err)
        return
    }
    
    for path, workflow := range workflows {
        fmt.Printf("工作流: %s\n", path)
        fmt.Printf("  名称: %s\n", workflow.Name)
        fmt.Printf("  作业数: %d\n", len(workflow.Jobs))
    }
}
```

## 下一步

- 探索 [API 参考](/zh/api/) 获取详细文档
- 查看 [示例](/zh/examples/) 了解更多用例
- 了解 [验证功能](/zh/api/validation)
- 发现 [工具函数](/zh/api/utilities) 的高级用法
