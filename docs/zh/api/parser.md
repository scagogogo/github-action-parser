# 解析函数

本页面记录了 GitHub Action Parser 库提供的核心解析函数。

## ParseFile

解析指定路径的 GitHub Action YAML 文件。

```go
func ParseFile(path string) (*ActionFile, error)
```

### 参数

- **path** (`string`): 要解析的 YAML 文件路径

### 返回值

- `*ActionFile`: 解析后的 action/workflow 结构
- `error`: 解析失败时的错误

### 描述

`ParseFile` 从文件系统打开并解析 YAML 文件。它支持 action 文件（`action.yml`, `action.yaml`）和 workflow 文件（`.github/workflows/*.yml`）。

### 使用示例

```go
// 解析 action 文件
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatalf("解析 action 失败: %v", err)
}

fmt.Printf("Action: %s\n", action.Name)
fmt.Printf("描述: %s\n", action.Description)

// 解析 workflow 文件
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
if err != nil {
    log.Fatalf("解析 workflow 失败: %v", err)
}

fmt.Printf("Workflow: %s\n", workflow.Name)
fmt.Printf("作业数: %d\n", len(workflow.Jobs))
```

### 错误处理

函数在以下情况下返回错误：
- 文件不存在或无法打开
- 文件包含无效的 YAML 语法
- 由于权限问题无法读取文件

```go
action, err := parser.ParseFile("nonexistent.yml")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("文件不存在")
    } else {
        fmt.Printf("解析错误: %v\n", err)
    }
}
```

## Parse

从 io.Reader 解析 GitHub Action YAML。

```go
func Parse(r io.Reader) (*ActionFile, error)
```

### 参数

- **r** (`io.Reader`): 包含 YAML 数据的 Reader

### 返回值

- `*ActionFile`: 解析后的 action/workflow 结构
- `error`: 解析失败时的错误

### 描述

`Parse` 从任何 `io.Reader` 读取 YAML 数据并解析为 `ActionFile` 结构。这对于从各种来源（如 HTTP 响应、嵌入文件或内存数据）解析 YAML 内容很有用。

### 使用示例

```go
// 从字符串解析
yamlContent := `
name: My Action
description: A sample action
runs:
  using: composite
  steps:
    - name: Hello
      run: echo "Hello World"
`

action, err := parser.Parse(strings.NewReader(yamlContent))
if err != nil {
    log.Fatalf("解析失败: %v", err)
}

fmt.Printf("Action: %s\n", action.Name)

// 从 HTTP 响应解析
resp, err := http.Get("https://example.com/action.yml")
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

action, err = parser.Parse(resp.Body)
if err != nil {
    log.Fatalf("解析远程 action 失败: %v", err)
}
```

### 错误处理

函数在以下情况下返回错误：
- Reader 读取数据时返回错误
- YAML 内容格式错误
- 数据无法反序列化为 ActionFile 结构

```go
// 处理不同类型的错误
action, err := parser.Parse(reader)
if err != nil {
    if strings.Contains(err.Error(), "unmarshal") {
        fmt.Println("无效的 YAML 结构")
    } else if strings.Contains(err.Error(), "read") {
        fmt.Println("读取数据失败")
    } else {
        fmt.Printf("解析错误: %v\n", err)
    }
}
```

## ParseDir

递归解析目录中的所有 GitHub Action YAML 文件。

```go
func ParseDir(dir string) (map[string]*ActionFile, error)
```

### 参数

- **dir** (`string`): 要扫描 YAML 文件的目录路径

### 返回值

- `map[string]*ActionFile`: 相对文件路径到解析结构的映射
- `error`: 解析失败时的错误

### 描述

`ParseDir` 递归遍历目录并解析所有 YAML 文件（`.yml` 和 `.yaml` 扩展名）。它返回一个映射，其中键是相对文件路径，值是解析的 `ActionFile` 结构。

### 使用示例

```go
// 解析 .github/workflows 中的所有 workflow
workflows, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatalf("解析 workflows 失败: %v", err)
}

for path, workflow := range workflows {
    fmt.Printf("文件: %s\n", path)
    fmt.Printf("  名称: %s\n", workflow.Name)
    fmt.Printf("  作业数: %d\n", len(workflow.Jobs))
    
    // 检查是否为可重用工作流
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("  类型: 可重用工作流\n")
    }
}

// 解析仓库中的所有 action 文件
actions, err := parser.ParseDir(".")
if err != nil {
    log.Fatalf("解析仓库失败: %v", err)
}

fmt.Printf("找到 %d 个 YAML 文件\n", len(actions))
```

### 文件过滤

`ParseDir` 根据扩展名自动过滤文件：
- 包含：`.yml`, `.yaml`
- 排除：所有其他文件类型、目录

### 错误处理

函数在以下情况下返回错误：
- 目录不存在或无法访问
- 读取目录或文件时权限被拒绝
- 任何单个文件解析失败（停止处理）

```go
actions, err := parser.ParseDir("workflows")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("目录不存在")
    } else if strings.Contains(err.Error(), "permission denied") {
        fmt.Println("权限被拒绝")
    } else {
        fmt.Printf("解析错误: %v\n", err)
    }
}
```

### 性能考虑

- 函数按顺序处理文件
- 包含许多文件的大目录可能需要时间处理
- 内存使用量随文件数量和大小而扩展
- 如需要，可考虑使用 goroutines 进行并行处理

```go
// 示例：处理大目录
start := time.Now()
actions, err := parser.ParseDir("large-repo")
if err != nil {
    log.Fatal(err)
}
duration := time.Since(start)

fmt.Printf("在 %v 内解析了 %d 个文件\n", duration, len(actions))
```

## 最佳实践

### 文件路径处理

始终使用适当的文件路径处理以实现跨平台兼容性：

```go
import "path/filepath"

// 好：使用 filepath.Join 实现跨平台路径
actionPath := filepath.Join("actions", "my-action", "action.yml")
action, err := parser.ParseFile(actionPath)

// 好：尽可能使用相对路径
workflow, err := parser.ParseFile(".github/workflows/ci.yml")
```

### 错误处理模式

实现全面的错误处理：

```go
func parseActionSafely(path string) (*parser.ActionFile, error) {
    action, err := parser.ParseFile(path)
    if err != nil {
        return nil, fmt.Errorf("解析 %s 失败: %w", path, err)
    }

    // 额外验证
    if action.Name == "" {
        return nil, fmt.Errorf("%s 中的 action 没有名称", path)
    }

    return action, nil
}
```

### 批量处理

处理多个文件时，考虑错误处理策略：

```go
func parseAllActions(dir string) (map[string]*parser.ActionFile, []error) {
    var errors []error
    results := make(map[string]*parser.ActionFile)

    actions, err := parser.ParseDir(dir)
    if err != nil {
        return nil, []error{err}
    }

    for path, action := range actions {
        // 每个文件的额外验证
        validator := parser.NewValidator()
        if validationErrors := validator.Validate(action); len(validationErrors) > 0 {
            for _, ve := range validationErrors {
                errors = append(errors, fmt.Errorf("%s: %s - %s", path, ve.Field, ve.Message))
            }
        } else {
            results[path] = action
        }
    }

    return results, errors
}
```

### 内存管理

对于大规模处理，考虑内存使用：

```go
// 逐个处理文件以减少内存使用
func processLargeRepository(dir string) error {
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && (strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml")) {
            action, err := parser.ParseFile(path)
            if err != nil {
                return fmt.Errorf("解析 %s 失败: %w", path, err)
            }

            // 立即处理 action，不要存储在内存中
            processAction(path, action)
        }

        return nil
    })
}
```
