# 验证 API

验证 API 提供了根据 GitHub 规范验证 GitHub Action 和 Workflow 文件的工具。

## Validator

用于验证 ActionFile 结构的主要验证器结构。

```go
type Validator struct {
    errors []ValidationError
}
```

### 描述

`Validator` 根据 GitHub 的要求和规范检查 ActionFile 结构。它验证 GitHub Actions 和 Workflows，为发现的任何违规提供详细的错误信息。

## NewValidator

创建新的 Validator 实例。

```go
func NewValidator() *Validator
```

### 返回值

- `*Validator`: 准备使用的新验证器实例

### 使用示例

```go
validator := parser.NewValidator()
```

## Validate

根据 GitHub 的要求验证 ActionFile。

```go
func (v *Validator) Validate(action *ActionFile) []ValidationError
```

### 参数

- **action** (`*ActionFile`): 要验证的 action 或 workflow

### 返回值

- `[]ValidationError`: 发现的验证错误切片（如果有效则为空）

### 描述

`Validate` 对 ActionFile 结构执行全面验证，检查：

- **Action 验证**：对于具有 `runs` 配置的文件
  - 必需字段（name、description）
  - 运行时特定要求（Node.js 的 main 脚本、Docker 的 image、composite 的 steps）
  - 支持的运行时类型

- **Workflow 验证**：对于具有 `jobs` 配置的文件
  - 必需的触发事件（`on` 字段）
  - 作业要求（runs-on 或 uses）
  - 步骤验证（需要 uses 或 run）

### 使用示例

```go
// 基本验证
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

validator := parser.NewValidator()
errors := validator.Validate(action)

if len(errors) > 0 {
    fmt.Println("发现验证错误:")
    for _, err := range errors {
        fmt.Printf("- %s: %s\n", err.Field, err.Message)
    }
} else {
    fmt.Println("Action 有效!")
}

// 验证多个文件
files := []string{"action.yml", "workflow.yml"}
for _, file := range files {
    action, err := parser.ParseFile(file)
    if err != nil {
        fmt.Printf("解析 %s 失败: %v\n", file, err)
        continue
    }
    
    errors := validator.Validate(action)
    if len(errors) > 0 {
        fmt.Printf("%s 有 %d 个验证错误:\n", file, len(errors))
        for _, err := range errors {
            fmt.Printf("  - %s: %s\n", err.Field, err.Message)
        }
    } else {
        fmt.Printf("%s 有效\n", file)
    }
}
```

## IsValid

检查验证器在验证后是否没有错误。

```go
func (v *Validator) IsValid() bool
```

### 返回值

- `bool`: 如果没有验证错误则为 true，否则为 false

### 使用示例

```go
validator := parser.NewValidator()
validator.Validate(action)

if validator.IsValid() {
    fmt.Println("没有验证错误")
} else {
    fmt.Println("发现验证错误")
}
```

## ValidationError

表示带有字段和消息信息的验证错误。

```go
type ValidationError struct {
    Field   string
    Message string
}
```

### 字段说明

- **Field** (`string`): 发生错误的字段路径
- **Message** (`string`): 人类可读的错误消息

### 使用示例

```go
errors := validator.Validate(action)
for _, err := range errors {
    fmt.Printf("字段: %s\n", err.Field)
    fmt.Printf("错误: %s\n", err.Message)
    fmt.Println("---")
}
```

## 验证规则

### Action 验证规则

#### 必需字段
- `name`: Action 必须有名称
- `description`: Action 必须有描述
- `runs.using`: Action 必须指定运行时

#### 运行时特定规则

**Node.js Actions** (`using: "node16"` 或 `using: "node20"`):
- `runs.main`: 必须指定主入口点脚本

**Docker Actions** (`using: "docker"`):
- `runs.image`: 必须指定 Docker 镜像

**Composite Actions** (`using: "composite"`):
- `runs.steps`: 必须至少有一个步骤

#### 支持的运行时
- `node16`: Node.js 16 运行时
- `node20`: Node.js 20 运行时  
- `docker`: Docker 容器运行时
- `composite`: 复合 action 运行时

### Workflow 验证规则

#### 必需字段
- `on`: Workflow 必须至少有一个触发事件
- `jobs`: Workflow 必须至少有一个作业

#### 作业规则
- 每个作业必须指定 `runs-on` 或 `uses`
- 如果定义了 `steps`，必须包含至少一个步骤

#### 步骤规则
- 每个步骤必须有 `uses` 或 `run`

### 示例验证场景

#### 有效的 Action 示例

```yaml
# 有效的 Node.js Action
name: My Node Action
description: A Node.js action
runs:
  using: node20
  main: index.js

# 有效的 Docker Action  
name: My Docker Action
description: A Docker action
runs:
  using: docker
  image: Dockerfile

# 有效的 Composite Action
name: My Composite Action
description: A composite action
runs:
  using: composite
  steps:
    - name: Hello
      run: echo "Hello"
```

#### 有效的 Workflow 示例

```yaml
# 有效的 Workflow
name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: npm test

# 有效的可重用 Workflow
name: Reusable Workflow
on:
  workflow_call:
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Building"
```
