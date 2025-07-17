# 类型参考

本页面记录了 GitHub Action Parser 库使用的所有数据结构。

## ActionFile

表示 GitHub Action 或 Workflow 文件的主要结构。

```go
type ActionFile struct {
    Name        string                 `yaml:"name,omitempty"`
    Description string                 `yaml:"description,omitempty"`
    Author      string                 `yaml:"author,omitempty"`
    Inputs      map[string]Input       `yaml:"inputs,omitempty"`
    Outputs     map[string]Output      `yaml:"outputs,omitempty"`
    Runs        RunsConfig             `yaml:"runs,omitempty"`
    Branding    Branding               `yaml:"branding,omitempty"`
    On          interface{}            `yaml:"on,omitempty"`
    Jobs        map[string]Job         `yaml:"jobs,omitempty"`
    Env         map[string]string      `yaml:"env,omitempty"`
    Defaults    map[string]interface{} `yaml:"defaults,omitempty"`
    Permissions interface{}            `yaml:"permissions,omitempty"`
}
```

### 字段说明

- **Name** (`string`): action 或 workflow 的名称
- **Description** (`string`): action 或 workflow 的功能描述
- **Author** (`string`): action 的作者
- **Inputs** (`map[string]Input`): action 的输入参数
- **Outputs** (`map[string]Output`): action 的输出值
- **Runs** (`RunsConfig`): action 的运行配置
- **Branding** (`Branding`): action 的品牌信息
- **On** (`interface{}`): workflow 的触发事件
- **Jobs** (`map[string]Job`): workflow 中定义的作业
- **Env** (`map[string]string`): 环境变量
- **Defaults** (`map[string]interface{}`): 默认设置
- **Permissions** (`interface{}`): 权限设置

### 使用示例

```go
action, err := parser.ParseFile("action.yml")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Action: %s\n", action.Name)
fmt.Printf("描述: %s\n", action.Description)
fmt.Printf("作者: %s\n", action.Author)
```

## Input

表示 action 的输入参数。

```go
type Input struct {
    Description string `yaml:"description,omitempty"`
    Required    bool   `yaml:"required,omitempty"`
    Default     string `yaml:"default,omitempty"`
    Deprecated  bool   `yaml:"deprecated,omitempty"`
}
```

### 字段说明

- **Description** (`string`): 输入参数的描述
- **Required** (`bool`): 是否为必需参数
- **Default** (`string`): 未提供时的默认值
- **Deprecated** (`bool`): 是否已弃用

### 使用示例

```go
for name, input := range action.Inputs {
    fmt.Printf("输入: %s\n", name)
    fmt.Printf("  描述: %s\n", input.Description)
    fmt.Printf("  必需: %t\n", input.Required)
    if input.Default != "" {
        fmt.Printf("  默认值: %s\n", input.Default)
    }
}
```

## Output

表示 action 的输出值。

```go
type Output struct {
    Description string `yaml:"description,omitempty"`
    Value       string `yaml:"value,omitempty"`
}
```

### 字段说明

- **Description** (`string`): 输出的描述
- **Value** (`string`): 输出的值表达式

### 使用示例

```go
for name, output := range action.Outputs {
    fmt.Printf("输出: %s\n", name)
    fmt.Printf("  描述: %s\n", output.Description)
    fmt.Printf("  值: %s\n", output.Value)
}
```

## RunsConfig

action 的执行配置。

```go
type RunsConfig struct {
    Using      string            `yaml:"using,omitempty"`
    Main       string            `yaml:"main,omitempty"`
    Pre        string            `yaml:"pre,omitempty"`
    Post       string            `yaml:"post,omitempty"`
    Image      string            `yaml:"image,omitempty"`
    Entrypoint string            `yaml:"entrypoint,omitempty"`
    Args       []string          `yaml:"args,omitempty"`
    Env        map[string]string `yaml:"env,omitempty"`
    Steps      []Step            `yaml:"steps,omitempty"`
}
```

### 字段说明

- **Using** (`string`): 使用的运行时（如 "node20", "docker", "composite"）
- **Main** (`string`): JavaScript actions 的主入口点
- **Pre** (`string`): JavaScript actions 的预执行脚本
- **Post** (`string`): JavaScript actions 的后执行脚本
- **Image** (`string`): Docker actions 的 Docker 镜像
- **Entrypoint** (`string`): Docker 入口点
- **Args** (`[]string`): Docker actions 的参数
- **Env** (`map[string]string`): 环境变量
- **Steps** (`[]Step`): 复合 actions 的步骤

### 使用示例

```go
runs := action.Runs
fmt.Printf("运行时: %s\n", runs.Using)

switch runs.Using {
case "node20", "node16":
    fmt.Printf("主脚本: %s\n", runs.Main)
case "docker":
    fmt.Printf("Docker 镜像: %s\n", runs.Image)
case "composite":
    fmt.Printf("步骤数: %d\n", len(runs.Steps))
}
```

## Job

表示 GitHub workflow 中的作业。

```go
type Job struct {
    Name        string                 `yaml:"name,omitempty"`
    RunsOn      interface{}            `yaml:"runs-on,omitempty"`
    Needs       interface{}            `yaml:"needs,omitempty"`
    If          string                 `yaml:"if,omitempty"`
    Steps       []Step                 `yaml:"steps,omitempty"`
    Env         map[string]string      `yaml:"env,omitempty"`
    Defaults    map[string]interface{} `yaml:"defaults,omitempty"`
    Outputs     map[string]string      `yaml:"outputs,omitempty"`
    TimeoutMin  int                    `yaml:"timeout-minutes,omitempty"`
    Strategy    interface{}            `yaml:"strategy,omitempty"`
    ContinueOn  interface{}            `yaml:"continue-on-error,omitempty"`
    Container   interface{}            `yaml:"container,omitempty"`
    Services    map[string]interface{} `yaml:"services,omitempty"`
    Uses        string                 `yaml:"uses,omitempty"`
    With        map[string]interface{} `yaml:"with,omitempty"`
    Secrets     interface{}            `yaml:"secrets,omitempty"`
    Permissions interface{}            `yaml:"permissions,omitempty"`
}
```

### 字段说明

- **Name** (`string`): 作业的显示名称
- **RunsOn** (`interface{}`): 运行环境（字符串或数组）
- **Needs** (`interface{}`): 此作业前必须完成的作业
- **If** (`string`): 作业执行的条件表达式
- **Steps** (`[]Step`): 作业中要执行的步骤
- **Env** (`map[string]string`): 环境变量
- **Defaults** (`map[string]interface{}`): 默认设置
- **Outputs** (`map[string]string`): 作业输出
- **TimeoutMin** (`int`): 超时时间（分钟）
- **Strategy** (`interface{}`): 矩阵策略配置
- **ContinueOn** (`interface{}`): 出错时继续设置
- **Container** (`interface{}`): 容器配置
- **Services** (`map[string]interface{}`): 服务容器
- **Uses** (`string`): 可重用工作流引用
- **With** (`map[string]interface{}`): 可重用工作流的输入
- **Secrets** (`interface{}`): 可重用工作流的密钥
- **Permissions** (`interface{}`): 权限设置

### 使用示例

```go
for jobID, job := range workflow.Jobs {
    fmt.Printf("作业: %s\n", jobID)
    if job.Name != "" {
        fmt.Printf("  名称: %s\n", job.Name)
    }
    fmt.Printf("  步骤数: %d\n", len(job.Steps))

    if job.Uses != "" {
        fmt.Printf("  使用: %s\n", job.Uses)
    }
}
```

## Step

表示工作流作业中的单个步骤。

```go
type Step struct {
    ID         string                 `yaml:"id,omitempty"`
    If         string                 `yaml:"if,omitempty"`
    Name       string                 `yaml:"name,omitempty"`
    Uses       string                 `yaml:"uses,omitempty"`
    Run        string                 `yaml:"run,omitempty"`
    Shell      string                 `yaml:"shell,omitempty"`
    With       map[string]interface{} `yaml:"with,omitempty"`
    Env        map[string]string      `yaml:"env,omitempty"`
    ContinueOn interface{}            `yaml:"continue-on-error,omitempty"`
    TimeoutMin int                    `yaml:"timeout-minutes,omitempty"`
    WorkingDir string                 `yaml:"working-directory,omitempty"`
}
```

### 字段说明

- **ID** (`string`): 步骤的唯一标识符
- **If** (`string`): 步骤执行的条件表达式
- **Name** (`string`): 步骤的显示名称
- **Uses** (`string`): 要使用的 action
- **Run** (`string`): 要运行的命令
- **Shell** (`string`): 运行命令使用的 shell
- **With** (`map[string]interface{}`): action 的输入参数
- **Env** (`map[string]string`): 环境变量
- **ContinueOn** (`interface{}`): 出错时继续设置
- **TimeoutMin** (`int`): 超时时间（分钟）
- **WorkingDir** (`string`): 工作目录

### 使用示例

```go
for i, step := range job.Steps {
    fmt.Printf("步骤 %d:\n", i+1)
    if step.Name != "" {
        fmt.Printf("  名称: %s\n", step.Name)
    }
    if step.Uses != "" {
        fmt.Printf("  使用: %s\n", step.Uses)
    }
    if step.Run != "" {
        fmt.Printf("  运行: %s\n", step.Run)
    }
}
```
