---
layout: home

hero:
  name: "GitHub Action Parser"
  text: "GitHub Actions Go 库"
  tagline: "轻松解析、验证和处理 GitHub Action YAML 文件"
  image:
    src: /logo.svg
    alt: GitHub Action Parser
  actions:
    - theme: brand
      text: 快速开始
      link: /zh/getting-started
    - theme: alt
      text: API 参考
      link: /zh/api/
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/github-action-parser

features:
  - icon: 📄
    title: 解析 Action 文件
    details: 解析 GitHub Action YAML 文件（action.yml/action.yaml），提供完整的类型安全和验证。
  - icon: ⚙️
    title: 工作流支持
    details: 解析 GitHub 工作流文件（.github/workflows/*.yml）并提取作业定义、步骤和触发器。
  - icon: ✅
    title: 验证功能
    details: 根据 GitHub 规范验证 actions 和 workflows，提供详细的错误报告。
  - icon: 🔄
    title: 可重用工作流
    details: 检测和处理可重用工作流，支持输入/输出参数提取。
  - icon: 🛠️
    title: 工具函数
    details: 提供类型转换和数据处理工具函数，用于处理 YAML 数据结构。
  - icon: 📁
    title: 批量处理
    details: 通过单个函数调用递归解析目录中的所有 Action 和 Workflow 文件。
---

## 快速开始

安装库：

```bash
go get github.com/scagogogo/github-action-parser
```

解析 action 文件：

```go
package main

import (
    "fmt"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
}
```

## 支持的 GitHub Action 功能

- **Action 元数据**：名称、描述、作者信息
- **输入参数**：带验证要求和默认值
- **输出参数**：带描述和值
- **Docker Actions**：基于容器的 actions
- **JavaScript Actions**：Node.js 16/20 actions
- **复合 Actions**：多步骤复合 actions
- **工作流作业**：作业定义和依赖关系
- **工作流触发器**：基于事件的触发器
- **可重用工作流**：带参数的可调用工作流
- **密钥处理**：可重用工作流密钥处理

## 为什么选择 GitHub Action Parser？

- **类型安全**：为所有 GitHub Action 结构提供完整的 Go 类型定义
- **功能全面**：支持所有 GitHub Action 和 Workflow 功能
- **经过验证**：根据 GitHub 规范内置验证
- **测试充分**：98.9% 测试覆盖率，全面的测试套件
- **易于使用**：简单的 API，清晰的文档和示例
- **高性能**：针对高效解析大量文件进行优化
