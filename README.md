# 🧊Cryo
![Go Badge](https://img.shields.io/badge/Go-1.24%2B-cyan?logo=go)
![GitHub Tag](https://img.shields.io/github/v/release/machinacanis/cryo)
[![goreportcard](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/machinacanis/cryo)
![GitHub License](https://img.shields.io/github/license/machinacanis/cryo)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](http://pkg.go.dev/github.com/machinacanis/cryo)

🚧开发中...

cryo 是一个轻量级聊天机器人开发框架，通过嵌入协议实现  [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)  来实现简单的部署和迁移。

## 特性

- 内嵌协议实现
- 不出意外的话可以单文件部署
- 为多Bot连接设计
- 消息去重 / 负载均衡
- 可启用的Web后台

## 安装

```bash
go get -u github.com/machinacanis/cryo
```

## 快速开始

[`script/example.go`](https://github.com/machinacanis/cryo/blob/main/script/example.go)是一个最小化的聊天机器人示例，展示了如何使用 cryo 框架登录账号并处理消息。

```go
// 尚处于开发阶段，API 可能会有变动
// 仅供参考
package main

import (
	"github.com/machinacanis/cryo"
	"github.com/machinacanis/cryo/log"
	"github.com/sirupsen/logrus"
)

func main() {
	bot := cryo.NewBot()
	bot.Init(cryo.Config{
		LogLevel:                     logrus.DebugLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	bot.OnType(cryo.PrivateMessageEventType, cryo.GroupMessageEventType).
		Handle(func(e *cryo.PrivateMessageEvent) {
			log.Info("响应事件 " + e.EventId)
			// 自定义逻辑
		}, cryo.AsyncMiddlewareType).
		Register()

	bot.AutoConnect()
	bot.Start()
}

```

查看 [文档](https://machinacanis.github.io/cryo/) 以查看完整的框架功能介绍及一个更全面的示例。

> [cryo-plugin-echo](https://github.com/machinacanis/cryo-echo-plugin) 是一个简单的 cryo 插件示例，展示了如何使用插件系统来更方便的组织代码。


## Thanks！

cryo 基于这些开源项目：

- [Lagrange.Core](https://github.com/LagrangeDev/Lagrange.Core) | NTQQ 协议实现
- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) | Lagrange.Core 的 Go 语言实现
- [LagrangeGo-Template](https://github.com/ExquisiteCore/LagrangeGo-Template) | LagrangeGo 的模板示例
- [logrus](https://github.com/sirupsen/logrus) | 优雅的 Go 日志库
- [freecache](https://github.com/coocood/freecache) | 高性能的内存缓存库
- [uuid](https://github.com/google/uuid) | UUID 生成器
- [go-qrcode](https://github.com/skip2/go-qrcode) | 二维码生成 / 解析工具

向这些项目的贡献者们致以最诚挚的感谢！