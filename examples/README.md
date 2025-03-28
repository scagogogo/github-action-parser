# GitHub Action Parser - 示例用法

本目录包含展示如何使用GitHub Action Parser库的多个示例。

## 示例列表

1. [基本解析](./01_basic_parsing/parse_action.go) - 解析Action文件并展示其结构
2. [工作流解析](./02_workflow_parsing/parse_workflow.go) - 解析Workflow文件并分析其组成部分
3. [验证工具](./03_validation/validate_action.go) - 验证Action/Workflow文件并提供修复建议
4. [可重用工作流分析](./04_reusable_workflow/analyze_reusable_workflow.go) - 分析可重用工作流的结构和参数
5. [实用工具函数](./05_utility_functions/utils_example.go) - 展示各种实用工具函数的使用方法

## 01 - 基本解析

此示例演示如何解析GitHub Action文件并显示其内容，包括：

- 基本元数据（名称、描述、作者）
- 品牌相关设置（图标、颜色）
- 输入参数详情（名称、描述、是否必填）
- 输出参数详情（名称、描述）
- 执行配置（运行类型、步骤）

### 运行方法

```bash
# 从项目根目录
go build -o parse-action ./example/01_basic_parsing/parse_action.go

# 使用测试文件运行
./parse-action pkg/parser/testdata/action.yml
```

## 02 - 工作流解析

此示例演示如何解析GitHub Workflow文件并分析其内容，包括：

- 工作流基本信息（名称、描述）
- 触发器配置
- 全局环境变量
- 作业配置和依赖关系
- 步骤详情

### 运行方法

```bash
# 从项目根目录
go build -o parse-workflow ./example/02_workflow_parsing/parse_workflow.go

# 使用测试文件运行
./parse-workflow pkg/parser/testdata/workflow.yml
```

## 03 - 验证工具

此示例演示如何验证GitHub Action/Workflow文件，包括：

- 识别文件类型
- 检查符合GitHub规范
- 显示验证错误
- 提供修复建议

### 运行方法

```bash
# 从项目根目录
go build -o validate-action ./example/03_validation/validate_action.go

# 使用测试文件运行
./validate-action pkg/parser/testdata/action.yml
```

## 04 - 可重用工作流分析

此示例演示如何分析可重用工作流，包括：

- 检查工作流是否可重用
- 分析工作流的输入参数
- 分析工作流的密钥配置
- 分析工作流的输出参数
- 分析工作流的作业和步骤
- 提供使用建议

### 运行方法

```bash
# 从项目根目录
go build -o analyze-reusable ./example/04_reusable_workflow/analyze_reusable_workflow.go

# 使用测试文件运行
./analyze-reusable pkg/parser/testdata/reusable-workflow.yml
```

## 05 - 实用工具函数

此示例演示了库中提供的多种实用函数，包括：

- 解析目录中的所有Action/Workflow文件
- 检查工作流是否可重用
- 提取可重用工作流的输入/输出参数
- 使用地图转换工具函数

### 运行方法

```bash
# 从项目根目录
go build -o utils-example ./example/05_utility_functions/utils_example.go

# 运行不同命令
./utils-example parse_dir pkg/parser/testdata/
./utils-example check_reusable pkg/parser/testdata/reusable-workflow.yml
./utils-example extract_inputs pkg/parser/testdata/reusable-workflow.yml
./utils-example extract_outputs pkg/parser/testdata/reusable-workflow.yml
./utils-example convert_map '{"key1":"value1","key2":"value2"}'
```

## 示例输出

### 解析Action文件

```
==== 基本信息 ====
文件名: action.yml
Action 名称: Example GitHub Action
描述: An example GitHub Action for testing the parser
作者: GitHub

==== 品牌设置 ====
图标: code
颜色: blue

==== 输入参数 ====
- file-path (必填, 默认值: 无): Path to the file to process
- output-format (可选, 默认值: json): Format of the output (json, yaml, or text)
- verbose (可选, 默认值: false): Enable verbose output

==== 输出参数 ====
- result: The result of the action
- status: The status of the operation

==== 执行配置 ====
运行类型: node16
主脚本: dist/index.js
预执行脚本: dist/setup.js
后执行脚本: dist/cleanup.js
```

### 解析可重用工作流

```
==== 可重用工作流信息 ====
名称: Reusable Workflow Example
描述: This is an example of a reusable workflow

==== 输入参数 ====
- version (可选, 默认值: latest)
  描述: The version to use
- environment (必填, 无默认值)
  描述: The environment to deploy to

输入参数统计：共 2 个，其中必填 1 个，可选 1 个，有默认值 1 个

==== 密钥 ====
- token (必填)
  描述: GitHub token for authentication

==== 输出参数 ====
- result
  描述: The result of the deployment
  值来源: ${{ jobs.deploy.outputs.result }}
```

### 使用工具函数

```
使用地图转换工具函数:
输入 JSON: {"key1":"value1","key2":123}

使用 MapOfStringInterface:
结果:
  key1: value1
  key2: 123

使用 MapOfStringString:
注意: MapOfStringString 失败: value for key 'key2' is not a string

使用只包含字符串值的映射:
结果:
  key1: value1
  key2: 123
``` 