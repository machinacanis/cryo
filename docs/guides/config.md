# 配置项

::: tip
这篇文档适用于 Cryo [β v0.1.7](https://github.com/machinacanis/cryo/releases/tag/v0.1.7) 版本
:::

Cryo 通过一个 [`Config`](https://pkg.go.dev/github.com/machinacanis/cryo#Config) 结构体来向 Bot 传递配置项，就像这样：

```go
func main() {
    // ...
    config := cryo.Config{
        EnableMessagePrintMiddleware: true,
    }
    bot.Init(logger, config) // 初始化 Bot 时传入配置项
    // ...
}

```

你只需要在初始化配置项时传入你需要修改的配置， [`Bot.Init()`](https://pkg.go.dev/github.com/machinacanis/cryo#Bot.Init) 方法会自动处理并给没有传入的配置项设置默认值。

每个 Bot 实例都可以拥有自己的配置项，配置项只会在 Bot 启动时被读取用于初始化，之后的修改不会影响 Bot 的运行，且 Bot 实例会在连接客户端时自动将配置项传递给客户端。

## 可用的配置项

| 配置项                            | 类型         | 默认值                 | 简介                                                                                                                |
|--------------------------------|------------|---------------------|-------------------------------------------------------------------------------------------------------------------|
| `SignServers`                  | `[]string` | `DefaultSignServer` | NTQQ 签名服务器地址列表                                                                                                    |
| `EnablePluginAutoLoad`         | `bool`     | `true`              | 是否启用插件自动加载，启用时通过 [`Bot.AddPlugin()`](https://pkg.go.dev/github.com/machinacanis/cryo#Bot.AddPlugin) 方法添加的插件会被自动启用 |
| `EnableClientAutoSave`         | `bool`     | `true`              | 是否启用客户端信息自动保存                                                                                                     |
| `EnablePrintLogo`              | `bool`     | `true`              | 是否在启用时向终端打印 ASCII 字符画 Logo                                                                                        |
| `EnableConnectPrintMiddleware` | `bool`     | `true`              | 是否启用内置的连接事件打印中间件                                                                                                  |
| `EnableMessagePrintMiddleware` | `bool`     | `true`              | 是否启用内置的消息打印中间件                                                                                                    |
| `EnableEventDebugMiddleware`   | `bool`     | `false`             | 是否启用内置的事件调试中间件                                                                                                    |
| `EnableCronScheduler`          | `bool`     | `false`             | 是否启用内置的gocron定时任务调度器                                                                                              |                                                                                                                   |

同时使用多个 Logger 实例高频率的进行 Log 是有些影响性能表现的，如果你的 Bot 需要处理特别大量的消息事件，建议在生产环境中关闭终端输出的日志，仅将日志输出到 `.log` 或 `.json` 文件中。