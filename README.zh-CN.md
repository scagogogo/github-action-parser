# GitHub Action Parser

[![Go Reference](https://pkg.go.dev/badge/github.com/scagogogo/github-action-parser.svg)](https://pkg.go.dev/github.com/scagogogo/github-action-parser) 
[![Go CI](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/github-action-parser/actions/workflows/ci.yml)
[![Documentation](https://github.com/scagogogo/github-action-parser/actions/workflows/docs.yml/badge.svg)](https://scagogogo.github.io/github-action-parser/)
[![Coverage](https://img.shields.io/badge/coverage-98.9%25-brightgreen)](https://github.com/scagogogo/github-action-parser)

一个用于解析、验证和处理 GitHub Action YAML 文件的 Go 库。

**📖 [文档](https://scagogogo.github.io/github-action-parser/zh/) | [English Documentation](https://scagogogo.github.io/github-action-parser/)**

## 功能特点

- 解析 GitHub Action YAML 文件（`action.yml`/`action.yaml`）
- 解析 GitHub Workflow 文件（`.github/workflows/*.yml`）
- 根据 GitHub 规范验证 actions 和 workflows
- 支持复合型（composite）、Docker 和 JavaScript actions
- 提取元数据、输入、输出、作业和步骤信息
- 检测和处理可重用工作流
- 类型转换和数据处理工具函数
- 批量解析目录中的所有 Action 和 Workflow 文件

## 安装

```bash
go get github.com/scagogogo/github-action-parser
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析 action 文件
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
    
    // 验证 action
    validator := parser.NewValidator()
    if errors := validator.Validate(action); len(errors) > 0 {
        fmt.Printf("验证错误: %v\n", errors)
    } else {
        fmt.Println("Action 有效!")
    }
}
```

## 文档

- **📖 [完整文档](https://scagogogo.github.io/github-action-parser/zh/)** - 完整的 API 参考和指南
- **🚀 [快速开始](https://scagogogo.github.io/github-action-parser/zh/getting-started)** - 快速入门指南
- **📚 [API 参考](https://scagogogo.github.io/github-action-parser/zh/api/)** - 详细的 API 文档
- **💡 [示例](https://scagogogo.github.io/github-action-parser/zh/examples/)** - 代码示例和用例

### English Documentation

- **📖 [Full Documentation](https://scagogogo.github.io/github-action-parser/)** - Complete API reference and guides
- **🚀 [Getting Started](https://scagogogo.github.io/github-action-parser/getting-started)** - Quick start guide
- **📚 [API Reference](https://scagogogo.github.io/github-action-parser/api/)** - Detailed API documentation
- **💡 [Examples](https://scagogogo.github.io/github-action-parser/examples/)** - Code examples and use cases

## 示例

### 解析 Action 文件

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Action: %s\n", action.Name)
for name, input := range action.Inputs {
    fmt.Printf("输入 %s: 必填=%t\n", name, input.Required)
}
```

### 解析 Workflow 文件

```go
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
if err != nil {
    log.Fatal(err)
}

for jobID, job := range workflow.Jobs {
    fmt.Printf("作业 %s 有 %d 个步骤\n", jobID, len(job.Steps))
}
```

### 验证文件

```go
validator := parser.NewValidator()
errors := validator.Validate(action)

if len(errors) > 0 {
    for _, err := range errors {
        fmt.Printf("%s 中的错误: %s\n", err.Field, err.Message)
    }
}
```

### 解析目录

```go
actions, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("找到 %d 个工作流文件\n", len(actions))
```

### 检查可重用工作流

```go
if parser.IsReusableWorkflow(workflow) {
    inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
    fmt.Printf("可重用工作流，有 %d 个输入\n", len(inputs))
}
```

## 支持的 GitHub Action 功能

- ✅ Action 元数据（名称、描述、作者）
- ✅ 带验证要求的输入参数
- ✅ 带描述和值的输出参数
- ✅ Docker 容器 actions
- ✅ JavaScript actions（Node.js 16/20）
- ✅ 复合 actions
- ✅ 工作流作业定义
- ✅ 工作流触发器（事件）
- ✅ 可重用工作流
- ✅ 作业和步骤依赖关系
- ✅ 可重用工作流的密钥处理

## 测试

该库具有全面的测试覆盖率（98.9%），包括：

- 所有函数的单元测试
- 使用真实 GitHub Action 文件的集成测试
- GitHub 规范的验证测试
- 性能基准测试

```bash
go test ./pkg/parser/
go test -bench=. ./pkg/parser/
```

## 贡献

欢迎贡献！请随时提交 Pull Request。对于重大更改，请先开启 issue 讨论您想要更改的内容。

## 许可证

该项目基于 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 链接

- **文档**: https://scagogogo.github.io/github-action-parser/zh/
- **Go 包**: https://pkg.go.dev/github.com/scagogogo/github-action-parser
- **GitHub 仓库**: https://github.com/scagogogo/github-action-parser
- **问题反馈**: https://github.com/scagogogo/github-action-parser/issues
