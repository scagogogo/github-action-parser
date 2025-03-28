# GitHub Action Parser

一个用于解析、验证和处理GitHub Action YAML文件的Go库。

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/github-action-parser.svg)](https://pkg.go.dev/github.com/scagogogo/github-action-parser) [![Go CI](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml)

## 功能特点

- 解析GitHub Action YAML文件（`action.yml`/`action.yaml`）
- 解析GitHub Workflow文件（`.github/workflows/*.yml`）
- 根据GitHub的规范验证actions和workflows
- 支持复合型（composite）动作、Docker动作和JavaScript动作
- 提取元数据、输入、输出、作业和步骤信息
- 检测和处理可重用工作流（reusable workflows）
- 提供类型转换和数据处理的辅助工具函数
- 支持批量解析目录中的所有Action和Workflow文件

## 安装

```bash
go get github.com/scagogogo/github-action-parser
```

## 使用示例

安装后，导入解析器包：

```go
import "github.com/scagogogo/github-action-parser/pkg/parser"
```

### 示例程序

查看[示例目录](./examples)获取完整的示例应用程序：

1. [基本解析](./examples/01_basic_parsing/parse_action.go) - 解析Action文件并显示其结构
2. [工作流解析](./examples/02_workflow_parsing/parse_workflow.go) - 解析Workflow文件并分析其组成部分
3. [验证工具](./examples/03_validation/validate_action.go) - 验证Action/Workflow文件并提供修复建议
4. [可重用工作流分析](./examples/04_reusable_workflow/analyze_reusable_workflow.go) - 分析可重用工作流的结构和参数
5. [实用工具函数](./examples/05_utility_functions/utils_example.go) - 展示各种实用工具函数的使用方法

构建并运行示例：

```bash
# 构建基本解析示例
go build -o parse-action ./examples/01_basic_parsing/parse_action.go

# 使用测试文件运行
./parse-action pkg/parser/testdata/action.yml
```

### 解析Action文件

```go
package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 解析action.yml文件
	action, err := parser.ParseFile("path/to/action.yml")
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// 访问action元数据
	fmt.Printf("Action Name: %s\n", action.Name)
	fmt.Printf("Description: %s\n", action.Description)
	
	// 访问输入参数
	for name, input := range action.Inputs {
		fmt.Printf("Input %s: %s (Required: %t)\n", name, input.Description, input.Required)
	}
	
	// 访问输出参数
	for name, output := range action.Outputs {
		fmt.Printf("Output %s: %s\n", name, output.Description)
	}
	
	// 对于composite动作，访问步骤
	if action.Runs.Using == "composite" {
		for i, step := range action.Runs.Steps {
			fmt.Printf("Step %d: %s\n", i+1, step.Name)
		}
	}
}
```

### 解析Workflow文件

```go
package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 解析workflow.yml文件
	workflow, err := parser.ParseFile("path/to/.github/workflows/ci.yml")
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}
    
    // 访问工作流作业
    for jobId, job := range workflow.Jobs {
        fmt.Printf("Job: %s (%s)\n", jobId, job.Name)
        
        // 访问作业步骤
        for i, step := range job.Steps {
            if step.Name != "" {
                fmt.Printf("  Step %d: %s\n", i+1, step.Name)
            } else if step.Run != "" {
                fmt.Printf("  Step %d: Run command\n", i+1)
            } else if step.Uses != "" {
                fmt.Printf("  Step %d: Uses %s\n", i+1, step.Uses)
            }
        }
    }
    
    // 检查是否是可重用工作流
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("This is a reusable workflow")
        
        // 提取工作流调用中定义的输入
        inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
        for name, input := range inputs {
            fmt.Printf("Workflow input %s: %s (Required: %t)\n", 
                name, input.Description, input.Required)
        }
    }
}
```

### 验证Action或Workflow

```go
package main

import (
	"fmt"
	"os"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 解析任何GitHub Action/Workflow文件
	action, err := parser.ParseFile("path/to/file.yml")
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}
	
	// 创建验证器并验证
	validator := parser.NewValidator()
	errors := validator.Validate(action)
	
	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for _, err := range errors {
			fmt.Printf("- %s: %s\n", err.Field, err.Message)
		}
	} else {
		fmt.Println("The file is valid!")
	}
}
```

### 解析仓库中的所有Workflow文件

```go
package main

import (
	"fmt"

	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 解析目录中的所有workflow文件
	workflows, err := parser.ParseDir(".github/workflows")
	if err != nil {
		fmt.Printf("Error parsing workflows: %v\n", err)
		return
	}
	
	for path, workflow := range workflows {
		fmt.Printf("Workflow: %s\n", path)
		fmt.Printf("  Name: %s\n", workflow.Name)
		fmt.Printf("  Jobs: %d\n", len(workflow.Jobs))
	}
}
```

### 使用实用工具函数

```go
package main

import (
	"fmt"
	"github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
	// 使用MapOfStringInterface将interface{}映射转换为string->interface{}映射
	rawMap := map[interface{}]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	
	strMap, err := parser.MapOfStringInterface(rawMap)
	if err != nil {
		fmt.Printf("Error converting map: %v\n", err)
		return
	}
	
	// 使用MapOfStringString将映射转换为纯字符串映射
	// 这在处理某些需要字符串键值对的场景很有用
	strOnlyMap, err := parser.MapOfStringString(strMap)
	if err == nil {
		for k, v := range strOnlyMap {
			fmt.Printf("%s: %s\n", k, v)
		}
	}
	
	// 检查是否为可重用工作流
	workflow, _ := parser.ParseFile("path/to/workflow.yml")
	if parser.IsReusableWorkflow(workflow) {
		// 提取工作流输入
		inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
		// 提取工作流输出
		outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow) 
		
		fmt.Printf("Inputs: %d, Outputs: %d\n", len(inputs), len(outputs))
	}
}
```

## 支持的GitHub Action功能

- Action元数据（名称、描述、作者）
- 带验证要求的输入参数
- 带描述和值的输出参数
- Docker容器动作
- JavaScript动作
- 复合型动作
- 工作流作业定义
- 工作流触发器（事件）
- 可重用工作流
- 作业和步骤依赖关系
- 可重用工作流的密钥处理

## 许可证

本项目采用MIT许可证 - 详情请参阅LICENSE文件。 