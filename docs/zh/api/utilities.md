# 工具 API

工具 API 提供了类型转换、数据处理和处理可重用工作流的辅助函数。

## 类型转换函数

### MapOfStringInterface

将 YAML 映射转换为 `map[string]interface{}`。

```go
func MapOfStringInterface(v interface{}) (map[string]interface{}, error)
```

#### 参数

- **v** (`interface{}`): 要转换的输入值

#### 返回值

- `map[string]interface{}`: 转换后的映射
- `error`: 转换失败时的错误

#### 描述

此函数处理常见的 YAML 反序列化场景，其中映射可能具有不同的键类型。它转换：
- `map[string]interface{}` → 按原样返回
- `map[interface{}]interface{}` → 将键转换为字符串
- `nil` → 返回 `nil`

#### 使用示例

```go
// 处理工作流触发器
switch on := workflow.On.(type) {
case map[interface{}]interface{}:
    // 转换为字符串键映射
    onMap, err := parser.MapOfStringInterface(on)
    if err != nil {
        log.Fatal(err)
    }
    
    for event, config := range onMap {
        fmt.Printf("触发器: %s\n", event)
    }
}

// 处理作业配置
if job.With != nil {
    withMap, err := parser.MapOfStringInterface(job.With)
    if err != nil {
        log.Fatal(err)
    }
    
    for key, value := range withMap {
        fmt.Printf("输入 %s: %v\n", key, value)
    }
}
```

#### 错误情况

- 如果输入包含无法转换的非字符串键，则返回错误
- 对于不支持的输入类型返回错误

### MapOfStringString

将 YAML 映射转换为 `map[string]string`。

```go
func MapOfStringString(v interface{}) (map[string]string, error)
```

#### 参数

- **v** (`interface{}`): 要转换的输入值

#### 返回值

- `map[string]string`: 转换后的字符串值映射
- `error`: 转换失败时的错误

#### 描述

将各种映射类型转换为仅字符串映射。处理：
- `map[string]string` → 按原样返回
- `map[string]interface{}` → 将值转换为字符串
- `map[interface{}]interface{}` → 将键和值转换为字符串
- `nil` → 返回 `nil`

#### 使用示例

```go
// 处理环境变量
if job.Env != nil {
    envMap, err := parser.MapOfStringString(job.Env)
    if err != nil {
        log.Fatal(err)
    }
    
    for key, value := range envMap {
        fmt.Printf("ENV %s=%s\n", key, value)
    }
}

// 转换步骤环境变量
if step.Env != nil {
    stepEnv, err := parser.MapOfStringString(step.Env)
    if err != nil {
        log.Printf("警告: 无法转换步骤环境变量: %v", err)
    } else {
        for k, v := range stepEnv {
            fmt.Printf("  %s: %s\n", k, v)
        }
    }
}
```

#### 错误情况

- 如果值无法转换为字符串，则返回错误
- 如果键无法转换为字符串，则返回错误
- 对于不支持的输入类型返回错误

## 可重用工作流函数

### IsReusableWorkflow

检查工作流是否旨在被其他工作流调用。

```go
func IsReusableWorkflow(action *ActionFile) bool
```

#### 参数

- **action** (`*ActionFile`): 要检查的工作流

#### 返回值

- `bool`: 如果工作流可重用则为 true，否则为 false

#### 描述

通过检查 `on` 字段中的 `workflow_call` 触发事件来确定工作流是否可重用。

#### 使用示例

```go
workflow, err := parser.ParseFile(".github/workflows/reusable.yml")
if err != nil {
    log.Fatal(err)
}

if parser.IsReusableWorkflow(workflow) {
    fmt.Println("这是一个可重用工作流")
    
    // 提取输入和输出
    inputs, _ := parser.ExtractInputsFromWorkflowCall(workflow)
    outputs, _ := parser.ExtractOutputsFromWorkflowCall(workflow)
    
    fmt.Printf("输入: %d\n", len(inputs))
    fmt.Printf("输出: %d\n", len(outputs))
} else {
    fmt.Println("这是一个常规工作流")
}

// 批量检查工作流
workflows, err := parser.ParseDir(".github/workflows")
if err != nil {
    log.Fatal(err)
}

reusableCount := 0
for path, workflow := range workflows {
    if parser.IsReusableWorkflow(workflow) {
        fmt.Printf("可重用: %s\n", path)
        reusableCount++
    }
}

fmt.Printf("找到 %d 个可重用工作流\n", reusableCount)
```

### ExtractInputsFromWorkflowCall

从可重用工作流中提取输入定义。

```go
func ExtractInputsFromWorkflowCall(action *ActionFile) (map[string]Input, error)
```

#### 参数

- **action** (`*ActionFile`): 可重用工作流

#### 返回值

- `map[string]Input`: 输入名称到 Input 定义的映射
- `error`: 提取失败时的错误

#### 描述

从 `workflow_call` 触发器配置中提取输入参数定义。如果工作流不可重用或未定义输入，则返回 `nil`。

#### 使用示例

```go
workflow, err := parser.ParseFile("reusable-workflow.yml")
if err != nil {
    log.Fatal(err)
}

if !parser.IsReusableWorkflow(workflow) {
    fmt.Println("不是可重用工作流")
    return
}

inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("提取输入失败: %v", err)
}

if len(inputs) == 0 {
    fmt.Println("未定义输入")
} else {
    fmt.Printf("工作流输入:\n")
    for name, input := range inputs {
        fmt.Printf("  %s:\n", name)
        fmt.Printf("    描述: %s\n", input.Description)
        fmt.Printf("    必需: %t\n", input.Required)
        if input.Default != "" {
            fmt.Printf("    默认值: %s\n", input.Default)
        }
    }
}
```

### ExtractOutputsFromWorkflowCall

从可重用工作流中提取输出定义。

```go
func ExtractOutputsFromWorkflowCall(action *ActionFile) (map[string]Output, error)
```

#### 参数

- **action** (`*ActionFile`): 可重用工作流

#### 返回值

- `map[string]Output`: 输出名称到 Output 定义的映射
- `error`: 提取失败时的错误

#### 描述

从 `workflow_call` 触发器配置中提取输出参数定义。如果工作流不可重用或未定义输出，则返回 `nil`。

#### 使用示例

```go
workflow, err := parser.ParseFile("reusable-workflow.yml")
if err != nil {
    log.Fatal(err)
}

outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
if err != nil {
    log.Fatalf("提取输出失败: %v", err)
}

if len(outputs) == 0 {
    fmt.Println("未定义输出")
} else {
    fmt.Printf("工作流输出:\n")
    for name, output := range outputs {
        fmt.Printf("  %s:\n", name)
        fmt.Printf("    描述: %s\n", output.Description)
        if output.Value != "" {
            fmt.Printf("    值: %s\n", output.Value)
        }
    }
}
```

## StringOrStringSlice 类型

用于处理可以是字符串或字符串数组的 YAML 字段的工具类型。

### 方法

#### Contains

```go
func (s *StringOrStringSlice) Contains(value string) bool
```

检查给定值是否包含在字符串或字符串切片中。

#### String

```go
func (s *StringOrStringSlice) String() string
```

返回字符串表示。对于单个值，返回该值。对于多个值，返回逗号分隔的列表。

#### 使用示例

```go
// 此类型在内部使用，但对自定义处理很有用
var triggers StringOrStringSlice

// 模拟 YAML 反序列化
yaml.Unmarshal([]byte(`["push", "pull_request"]`), &triggers)

if triggers.Contains("push") {
    fmt.Println("由 push 事件触发")
}

fmt.Printf("所有触发器: %s\n", triggers.String())
// 输出: 所有触发器: push, pull_request
```

## 错误处理模式

### 优雅的类型转换

```go
func safeMapConversion(v interface{}) map[string]string {
    result, err := parser.MapOfStringString(v)
    if err != nil {
        // 回退到字符串接口映射
        if interfaceMap, err2 := parser.MapOfStringInterface(v); err2 == nil {
            result = make(map[string]string)
            for k, val := range interfaceMap {
                result[k] = fmt.Sprintf("%v", val)
            }
        }
    }
    return result
}
```

### 带错误收集的批量处理

```go
func processReusableWorkflows(dir string) {
    workflows, err := parser.ParseDir(dir)
    if err != nil {
        log.Fatal(err)
    }

    var errors []error
    reusableWorkflows := make(map[string]*parser.ActionFile)

    for path, workflow := range workflows {
        if parser.IsReusableWorkflow(workflow) {
            reusableWorkflows[path] = workflow

            // 验证可以提取输入
            if _, err := parser.ExtractInputsFromWorkflowCall(workflow); err != nil {
                errors = append(errors, fmt.Errorf("%s: 提取输入失败: %w", path, err))
            }

            // 验证可以提取输出
            if _, err := parser.ExtractOutputsFromWorkflowCall(workflow); err != nil {
                errors = append(errors, fmt.Errorf("%s: 提取输出失败: %w", path, err))
            }
        }
    }

    fmt.Printf("找到 %d 个可重用工作流\n", len(reusableWorkflows))
    if len(errors) > 0 {
        fmt.Printf("遇到 %d 个错误:\n", len(errors))
        for _, err := range errors {
            fmt.Printf("  - %v\n", err)
        }
    }
}
```
