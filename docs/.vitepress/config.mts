import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: '/cryo/',
  title: "Cryo 文档",
  description: "轻量级Golang聊天机器人框架",
  lastUpdated: true,
  themeConfig: {
    logo: '/img/ice.png',

    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '主页', link: '/' },
      { text: '指南', link: '/guides/intro' },
      { text: 'API 文档', link: 'https://pkg.go.dev/github.com/machinacanis/cryo' }
    ],

    sidebar: [
      {
        text: 'Cryo 速冻指南',
        items: [
          { text: '引言', link: '/guides/intro' },
          { text: '快速开始', link: '/guides/quick-start' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/machinacanis/cryo' }
    ]
  }
})
