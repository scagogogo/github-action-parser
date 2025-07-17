# 工具函数

此示例演示如何使用 GitHub Action Parser 库提供的工具函数进行类型转换和数据处理。

## 映射转换工具

### MapOfStringInterface

将各种映射类型转换为 `map[string]interface{}`：

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 带有 interface{} 键的示例映射
    interfaceMap := map[interface{}]interface{}{
        "key1": "value1",
        "key2": 42,
        "key3": true,
        "key4": []string{"a", "b", "c"},
    }
    
    // 转换为 map[string]interface{}
    stringMap, err := parser.MapOfStringInterface(interfaceMap)
    if err != nil {
        log.Fatalf("转换映射失败: %v", err)
    }
    
    fmt.Println("转换后的映射:")
    for key, value := range stringMap {
        fmt.Printf("  %s: %v (类型: %T)\n", key, value, value)
    }
}
```

### MapOfStringString

将各种映射类型转换为 `map[string]string`：

```go
// 带有 interface{} 值的示例映射
interfaceValueMap := map[string]interface{}{
    "key1": "value1",
    "key2": "value2",
    "key3": "value3",
}

// 转换为 map[string]string
stringStringMap, err := parser.MapOfStringString(interfaceValueMap)
if err != nil {
    log.Fatalf("转换映射失败: %v", err)
}

fmt.Println("\n转换后的字符串映射:")
for key, value := range stringStringMap {
    fmt.Printf("  %s: %s\n", key, value)
}

// 这将失败，因为不是所有值都是字符串
mixedMap := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
}

_, err = parser.MapOfStringString(mixedMap)
if err != nil {
    fmt.Printf("\n预期错误: %v\n", err)
}
```

## 处理可重用工作流

### IsReusableWorkflow

检查工作流是否可重用：

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析工作流文件
    workflow, err := parser.ParseFile(".github/workflows/ci.yml")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    // 检查是否为可重用工作流
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("这是一个可重用工作流")
    } else {
        fmt.Println("这是一个常规工作流")
    }
}
```

### ExtractInputsFromWorkflowCall

从可重用工作流提取输入参数：

```go
// 检查工作流是否可重用并提取输入
if parser.IsReusableWorkflow(workflow) {
    inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
    if err != nil {
        log.Fatalf("提取输入失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个输入:\n", len(inputs))
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

从可重用工作流提取输出参数：

```go
// 从可重用工作流提取输出
if parser.IsReusableWorkflow(workflow) {
    outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
    if err != nil {
        log.Fatalf("提取输出失败: %v", err)
    }
    
    fmt.Printf("找到 %d 个输出:\n", len(outputs))
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

使用 `StringOrStringSlice` 类型：

```go
package main

import (
    "fmt"
    "log"
    "gopkg.in/yaml.v3"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 带有字符串或字符串数组字段的示例 YAML
    yamlData := `
singleString: value1
stringArray:
  - item1
  - item2
  - item3
`
    
    // 解析 YAML
    var data struct {
        SingleString parser.StringOrStringSlice `yaml:"singleString"`
        StringArray  parser.StringOrStringSlice `yaml:"stringArray"`
    }
    
    err := yaml.Unmarshal([]byte(yamlData), &data)
    if err != nil {
        log.Fatalf("解析 YAML 失败: %v", err)
    }
    
    // 访问单个字符串
    fmt.Printf("SingleString.Value: %s\n", data.SingleString.Value)
    fmt.Printf("SingleString.Values: %v\n", data.SingleString.Values)
    
    // 访问字符串数组
    fmt.Printf("StringArray.Value: %s\n", data.StringArray.Value)
    fmt.Printf("StringArray.Values: %v\n", data.StringArray.Values)
    
    // 检查是否包含某个值
    if data.StringArray.Contains("item2") {
        fmt.Println("StringArray 包含 'item2'")
    }
    
    // 字符串表示
    fmt.Printf("StringArray 作为字符串: %s\n", data.StringArray.String())
}
```

## 错误处理模式

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 安全映射转换与错误处理的示例
    interfaceMap := map[interface{}]interface{}{
        "key1": "value1",
        "key2": 42,
        "key3": true,
    }
    
    // 尝试转换为字符串-字符串映射（将失败）
    stringMap, err := safeMapConversion(interfaceMap)
    if err != nil {
        fmt.Printf("警告: %v\n", err)
        fmt.Printf("回退映射: %v\n", stringMap)
    } else {
        fmt.Printf("转换后的映射: %v\n", stringMap)
    }
}

// 带回退的安全转换
func safeMapConversion(input interface{}) (map[string]string, error) {
    // 首先尝试直接转换
    result, err := parser.MapOfStringString(input)
    if err == nil {
        return result, nil
    }
    
    // 如果失败，尝试先转换为 map[string]interface{}
    interfaceMap, err := parser.MapOfStringInterface(input)
    if err != nil {
        return nil, fmt.Errorf("转换映射失败: %w", err)
    }
    
    // 然后手动将值转换为字符串
    result = make(map[string]string)
    var conversionErrors []string
    
    for key, value := range interfaceMap {
        switch v := value.(type) {
        case string:
            result[key] = v
        case int, int64, float64, bool:
            result[key] = fmt.Sprintf("%v", v)
        default:
            conversionErrors = append(conversionErrors, 
                fmt.Sprintf("无法转换 %s: %v (%T) 为字符串", key, v, v))
        }
    }
    
    if len(conversionErrors) > 0 {
        return result, fmt.Errorf("部分转换，有 %d 个错误", len(conversionErrors))
    }
    
    return result, nil
}
```

## 使用工具的批量处理

```go
package main

import (
    "fmt"
    "log"
    "github.com/scagogogo/github-action-parser/pkg/parser"
)

func main() {
    // 解析所有工作流
    workflows, err := parser.ParseDir(".github/workflows")
    if err != nil {
        log.Fatalf("解析工作流失败: %v", err)
    }
    
    // 处理每个工作流
    for path, workflow := range workflows {
        fmt.Printf("=== %s ===\n", path)
        
        // 处理工作流触发器
        processTriggers(workflow.On)
        
        // 处理环境变量
        processEnvVars(workflow.Env)
        
        // 处理作业
        for jobID, job := range workflow.Jobs {
            fmt.Printf("作业: %s\n", jobID)
            
            // 处理作业环境变量
            processEnvVars(job.Env)
            
            // 处理步骤
            for _, step := range job.Steps {
                if step.With != nil {
                    processStepInputs(step.With)
                }
            }
        }
        
        fmt.Println()
    }
}

func processTriggers(on interface{}) {
    fmt.Println("触发器:")
    
    switch v := on.(type) {
    case string:
        fmt.Printf("  %s\n", v)
    case []interface{}:
        for _, trigger := range v {
            fmt.Printf("  %v\n", trigger)
        }
    case map[string]interface{}:
        stringMap, err := parser.MapOfStringInterface(v)
        if err != nil {
            fmt.Printf("  转换触发器错误: %v\n", err)
            return
        }
        
        for event, config := range stringMap {
            fmt.Printf("  %s: %v\n", event, config)
        }
    default:
        fmt.Printf("  未知触发器类型: %T\n", on)
    }
}

func processEnvVars(env interface{}) {
    if env == nil {
        return
    }
    
    fmt.Println("环境变量:")
    
    envVars, err := parser.MapOfStringString(env)
    if err != nil {
        fmt.Printf("  转换环境变量错误: %v\n", err)
        return
    }
    
    for key, value := range envVars {
        fmt.Printf("  %s: %s\n", key, value)
    }
}

func processStepInputs(with interface{}) {
    fmt.Println("  步骤输入:")
    
    inputs, err := parser.MapOfStringInterface(with)
    if err != nil {
        fmt.Printf("    转换输入错误: %v\n", err)
        return
    }
    
    for key, value := range inputs {
        fmt.Printf("    %s: %v\n", key, value)
    }
}
```

## 完整工具示例

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
    
    // 解析工作流文件
    workflow, err := parser.ParseFile(workflowFile)
    if err != nil {
        log.Fatalf("解析 %s 失败: %v", workflowFile, err)
    }
    
    // 使用工具函数进行全面的工作流分析
    fmt.Printf("=== 分析 %s ===\n\n", workflowFile)
    
    // 基本信息
    fmt.Printf("名称: %s\n", workflow.Name)
    
    // 使用类型转换处理触发器
    fmt.Println("\n=== 触发器 ===")
    processTriggers(workflow.On)
    
    // 处理环境变量
    if workflow.Env != nil {
        fmt.Println("\n=== 环境变量 ===")
        envVars, err := parser.MapOfStringString(workflow.Env)
        if err != nil {
            fmt.Printf("错误: %v\n", err)
        } else {
            for key, value := range envVars {
                fmt.Printf("%s: %s\n", key, value)
            }
        }
    }
    
    // 检查是否可重用
    if parser.IsReusableWorkflow(workflow) {
        fmt.Println("\n=== 可重用工作流 ===")
        
        // 提取输入
        inputs, err := parser.ExtractInputsFromWorkflowCall(workflow)
        if err != nil {
            fmt.Printf("提取输入失败: %v\n", err)
        } else {
            fmt.Printf("输入: %d\n", len(inputs))
            for name, input := range inputs {
                fmt.Printf("  %s: %s (必需: %t)\n", 
                    name, input.Description, input.Required)
            }
        }
        
        // 提取输出
        outputs, err := parser.ExtractOutputsFromWorkflowCall(workflow)
        if err != nil {
            fmt.Printf("提取输出失败: %v\n", err)
        } else {
            fmt.Printf("输出: %d\n", len(outputs))
            for name, output := range outputs {
                fmt.Printf("  %s: %s\n", name, output.Description)
            }
        }
    }
    
    // 处理作业
    fmt.Printf("\n=== 作业 (%d) ===\n", len(workflow.Jobs))
    for jobID, job := range workflow.Jobs {
        fmt.Printf("\n作业: %s\n", jobID)
        
        // 处理 runs-on
        if job.RunsOn != nil {
            fmt.Printf("运行在: ")
            switch runsOn := job.RunsOn.(type) {
            case string:
                fmt.Printf("%s\n", runsOn)
            case []interface{}:
                fmt.Printf("%v\n", runsOn)
            default:
                fmt.Printf("%v (类型: %T)\n", runsOn, runsOn)
            }
        }
        
        // 处理步骤
        fmt.Printf("步骤: %d\n", len(job.Steps))
    }
}

func processTriggers(on interface{}) {
    switch v := on.(type) {
    case string:
        fmt.Printf("单个事件: %s\n", v)
    case []interface{}:
        fmt.Println("多个事件:")
        for _, event := range v {
            fmt.Printf("  - %v\n", event)
        }
    case map[interface{}]interface{}:
        fmt.Println("复杂事件:")
        events, err := parser.MapOfStringInterface(v)
        if err != nil {
            fmt.Printf("转换事件错误: %v\n", err)
            return
        }
        
        for event, config := range events {
            fmt.Printf("  %s: %v\n", event, config)
        }
    case map[string]interface{}:
        fmt.Println("复杂事件:")
        for event, config := range v {
            fmt.Printf("  %s: %v\n", event, config)
        }
    default:
        fmt.Printf("未知触发器类型: %T\n", on)
    }
}
```

## 下一步

- 查看详细文档的 [API 参考](/zh/api/utilities)
- 了解 [验证功能](/zh/examples/validation)
- 详细探索 [可重用工作流](/zh/examples/reusable-workflows)
