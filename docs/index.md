---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "“🧊Cryo”"
  text: "轻量级 \nGo 聊天机器人框架"
  tagline: Make go chatbot great again.
  actions:
    - theme: brand
      text: ⚡ 快速上手
      link: /guides/intro
    - theme: alt
      text: API 文档
      link: https://pkg.go.dev/github.com/machinacanis/cryo

features:
  - title: 基于 LagrangeGo
    details: 内嵌 QQNT 协议实现，无需额外连接 Onebot 协议实现，单文件运行So Easy
  - title: 自动并发
    details: 默认使用 goroutine 自动并发处理事件，你只管塞处理函数，并发交给运行时
  - title: 事件驱动
    details: 基于订阅 / 发布模型，高度可定制的三段式事件处理流程
  - title: 插件系统
    details: 内置了插件接口，你可以利用插件来编写、组织代码，也可以导入开源的社区插件
---

