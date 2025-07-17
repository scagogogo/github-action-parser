# 示例

本节提供了在各种场景下使用 GitHub Action Parser 库的全面示例。

## 概述

示例按功能和复杂性组织：

1. **[基本解析](/zh/examples/basic-parsing)** - 简单的 action 文件解析和数据访问
2. **[工作流解析](/zh/examples/workflow-parsing)** - 处理 GitHub workflow 文件
3. **[验证功能](/zh/examples/validation)** - 验证 actions 和 workflows
4. **[可重用工作流](/zh/examples/reusable-workflows)** - 处理可重用工作流
5. **[工具函数](/zh/examples/utilities)** - 高级工具函数和类型转换

## 快速开始示例

### 解析 Action 文件

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    action, err := parser.ParseFile("action.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Action: %s\n", action.Name)
    fmt.Printf("描述: %s\n", action.Description)
}
```

### 验证 Workflow

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    workflow, err := parser.ParseFile(".github/workflows/ci.yml")
    if err != nil {
        log.Fatal(err)
    }
    
    validator := parser.NewValidator()
    errors := validator.Validate(workflow)
    
    if len(errors) == 0 {
        fmt.Println("Workflow 有效！")
    } else {
        fmt.Printf("发现 %d 个验证错误\n", len(errors))
    }
}
```

### 检查可重用工作流

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatal(err)
    }
    
    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            fmt.Printf("可重用工作流: %s\n", path)
        }
    }
}
```

## 示例仓库结构

示例假设以下仓库结构：

```
my-project/
├── action.yml                    # 主 action 文件
├── .github/
│   └── workflows/
│       ├── ci.yml                # CI 工作流
│       ├── release.yml           # 发布工作流
│       └── reusable-build.yml    # 可重用工作流
├── actions/
│   ├── setup/
│   │   └── action.yml            # 自定义设置 action
│   └── deploy/
│       └── action.yml            # 自定义部署 action
└── main.go                       # 你的应用程序
```

## 示例 YAML 文件

### 示例 Action 文件 (action.yml)

```yaml
name: My Custom Action
description: 演示用的示例 action
author: Your Name

inputs:
  environment:
    description: 目标环境
    required: true
    default: development
  version:
    description: 要部署的版本
    required: false

outputs:
  deployment-url:
    description: 部署应用程序的 URL
    value: ${{ steps.deploy.outputs.url }}

runs:
  using: composite
  steps:
    - name: 设置
      run: echo "设置环境 ${{ inputs.environment }}"
      shell: bash
    - name: 部署
      id: deploy
      run: |
        echo "部署版本 ${{ inputs.version }}"
        echo "url=https://app.example.com" >> $GITHUB_OUTPUT
      shell: bash

branding:
  icon: rocket
  color: blue
```

### 示例 Workflow 文件 (.github/workflows/ci.yml)

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  NODE_VERSION: 18

jobs:
  test:
    name: 运行测试
    runs-on: ubuntu-latest
    
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        
      - name: 设置 Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: npm
          
      - name: 安装依赖
        run: npm ci
        
      - name: 运行测试
        run: npm test
        
      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        if: success()

  build:
    name: 构建应用程序
    runs-on: ubuntu-latest
    needs: test
    
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        
      - name: 构建
        run: |
          echo "构建应用程序..."
          mkdir -p dist
          echo "构建成功" > dist/app.txt
          
      - name: 上传构建产物
        uses: actions/upload-artifact@v3
        with:
          name: build-artifacts
          path: dist/
```

### 示例可重用工作流 (.github/workflows/reusable-build.yml)

```yaml
name: 可重用构建工作流

on:
  workflow_call:
    inputs:
      environment:
        description: 目标环境
        required: true
        type: string
      node-version:
        description: 要使用的 Node.js 版本
        required: false
        type: string
        default: '18'
    outputs:
      build-version:
        description: 构建应用程序的版本
        value: ${{ jobs.build.outputs.version }}
    secrets:
      NPM_TOKEN:
        description: NPM 认证令牌
        required: true

jobs:
  build:
    name: 为 ${{ inputs.environment }} 构建
    runs-on: ubuntu-latest
    
    outputs:
      version: ${{ steps.version.outputs.version }}
    
    steps:
      - name: 检出代码
        uses: actions/checkout@v4
        
      - name: 设置 Node.js
        uses: actions/setup-node@v4
        with:
          node-version: ${{ inputs.node-version }}
          registry-url: https://registry.npmjs.org/
          
      - name: 获取版本
        id: version
        run: |
          VERSION=$(date +%Y%m%d)-$(git rev-parse --short HEAD)
          echo "version=$VERSION" >> $GITHUB_OUTPUT
          
      - name: 安装依赖
        run: npm ci
        env:
          NODE_AUTH_TOKEN: ${{ secrets.NPM_TOKEN }}
          
      - name: 为环境构建
        run: |
          echo "为 ${{ inputs.environment }} 构建"
          npm run build:${{ inputs.environment }}
```

## 运行示例

每个示例都可以独立运行。确保你有：

1. 安装了 Go 1.20 或更高版本
2. 安装了 GitHub Action Parser 库：
   ```bash
   go get github.com/scagogogo/github-action-parser
   ```

3. 项目目录中有示例 YAML 文件

然后导航到每个示例部分获取详细说明和代码。
