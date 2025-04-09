# 引言 <Badge type="warning" text="beta" />

::: tip
当前的 Cryo 指南适配于 [β v0.1.3](https://github.com/machinacanis/cryo/releases/tag/v0.1.3) 版本
:::

::: warning 
文档还没有写完！在写了在写了！
:::

::: warning
🚧 Cryo 尚处于开发阶段，API的变动可能会很频繁

因此指南有可能还没有跟进API的变动，请以 [API Reference](https://pkg.go.dev/github.com/machinacanis/cryo) 为准
:::

**🧊Cryo** 是一个基于 LagrangeGo 的**轻量级聊天机器人框架**，通过一个额外的事件总线来处理多个 LagrangeGo 实例的消息事件，以此提供更好的多Bot连接支持。

和 [Nonebot](https://nonebot.dev/docs/) / [Koishi](https://koishi.chat/zh-CN/) / [ZeroBot](https://github.com/wdvxdr1123/ZeroBot) 这类常见的 Bot 框架不同， Cryo 不使用外部的协议实现，而是**内嵌了协议实现**（~~何尝不是某种意义上的返璞归真~~）。

尽管这基本牺牲了跨平台的灵活性，但是也更好的利用了 Golang 的单文件编译特色，可以真正的做到编译完随便跑，一个二进制文件就能实现 Bot 的全部功能。

如果你需要的是一个可以适配不同社交软件平台的 Bot 框架，那么 Cryo 可能并不适合你。

但是如果你在找的是一个**仅面向 NTQQ 平台**的，简单好用功能丰富的 Bot 框架，那么 Cryo 应该会是一个不错的选择。

## 你在找这些吗？

- [Nonebot](https://nonebot.dev/docs/) | 跨平台 Python 异步聊天机器人框架，它是 Cryo 许多设计的灵感来源（~~抄袭目标~~）
- [Koishi](https://koishi.chat/zh-CN/) | 跨平台 JavaScript / Typescript 聊天机器人框架，功能强大
- [Lagrange.Core](https://lagrangedev.github.io/Lagrange.Doc/Lagrange.Core/) | NTQQ 协议实现兼 C# 聊天机器人框架
- [Lagrange.Onebot](https://lagrangedev.github.io/Lagrange.Doc/Lagrange.OneBot/) | Lagrange.Core 的 Onebot V11 协议实现
- [go-cqhttp](https://github.com/LagrangeDev/go-cqhttp) | 用 LagrangeGo 成功打赢复活赛的 Onebot V11 协议实现
- [ZeroBot](https://github.com/wdvxdr1123/ZeroBot) | 基于 Onebot 协议的 Golang 机器人开发框架
- [Gensokyo](https://github.com/Hoshinonyaruko/Gensokyo) | 基于 QQ 官方机器人 API 的 Onebot V11 协议实现

## 其他有的没的

- Cryo 就是希腊语里的 **“冰(κρύος / krýos)”**，很冷吧（笑
- 也许后续会对 Cryo 进行解耦，让它不强依赖于 LagrangeGo，不过这还是个大饼

