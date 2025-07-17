# API 参考

GitHub Action Parser 库提供了一套全面的类型和函数，用于解析、验证和处理 GitHub Action 和 Workflow YAML 文件。

## 包概述

```go
import "github.com/scagogogo/github-action-parser/pkg/parser"
```

`parser` 包包含了处理 GitHub Actions 和 Workflows 所需的所有功能：

- **核心类型**：表示 GitHub Action 和 Workflow 组件的数据结构
- **解析函数**：解析 YAML 文件和目录的函数
- **验证功能**：根据 GitHub 规范验证解析文件的工具
- **工具函数**：类型转换和数据处理的辅助函数

## 快速参考

### 主要函数

| 函数 | 描述 |
|------|------|
| [`ParseFile(path string)`](/zh/api/parser#parsefile) | 解析单个 YAML 文件 |
| [`Parse(r io.Reader)`](/zh/api/parser#parse) | 从 io.Reader 解析 |
| [`ParseDir(dir string)`](/zh/api/parser#parsedir) | 解析目录中的所有 YAML 文件 |
| [`NewValidator()`](/zh/api/validation#newvalidator) | 创建新的验证器实例 |

### 核心类型

| 类型 | 描述 |
|------|------|
| [`ActionFile`](/zh/api/types#actionfile) | 表示 action 或 workflow 的主要结构 |
| [`Input`](/zh/api/types#input) | 输入参数定义 |
| [`Output`](/zh/api/types#output) | 输出参数定义 |
| [`Job`](/zh/api/types#job) | 工作流作业定义 |
| [`Step`](/zh/api/types#step) | 作业中的单个步骤 |
| [`RunsConfig`](/zh/api/types#runsconfig) | Action 执行配置 |

### 验证类型

| 类型 | 描述 |
|------|------|
| [`Validator`](/zh/api/validation#validator) | GitHub Action 规范验证器 |
| [`ValidationError`](/zh/api/validation#validationerror) | 验证错误信息 |

### 工具类型

| 类型 | 描述 |
|------|------|
| [`StringOrStringSlice`](/zh/api/utilities#stringorstringslice) | YAML 中灵活的字符串/数组类型 |

## 错误处理

所有解析函数都会返回错误，提供解析失败的详细信息：

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    // 处理解析错误
    fmt.Printf("解析失败: %v\n", err)
    return
}
```

验证错误以 `ValidationError` 结构体切片的形式返回：

```go
validator := parser.NewValidator()
errors := validator.Validate(action)
for _, err := range errors {
    fmt.Printf("字段 %s: %s\n", err.Field, err.Message)
}
```

## 类型安全

该库为所有 GitHub Action 和 Workflow 结构提供完整的类型安全。所有字段都根据 GitHub Actions 规范进行了适当的类型定义，对可选字段适当使用指针，对灵活数据类型使用接口。

## 性能

解析器针对性能进行了优化，可以高效处理：
- 大型 action 和 workflow 文件
- 多个文件的批量处理
- 递归目录解析
- 大型仓库的内存高效处理

## 下一步

- [类型参考](/zh/api/types) - 所有数据结构的详细文档
- [解析函数](/zh/api/parser) - 完整的解析 API 文档
- [验证功能](/zh/api/validation) - 验证功能和错误处理
- [工具函数](/zh/api/utilities) - 辅助函数和工具
