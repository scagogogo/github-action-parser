import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'GitHub Action Parser',
  description: 'A Go library for parsing, validating and processing GitHub Action YAML files',
  
  // Base URL for GitHub Pages
  base: '/github-action-parser/',
  
  // Language configuration
  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'GitHub Action Parser',
      description: 'A Go library for parsing, validating and processing GitHub Action YAML files',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'API Reference', link: '/api/' },
          { text: 'Examples', link: '/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/github-action-parser' }
        ],
        sidebar: {
          '/api/': [
            {
              text: 'API Reference',
              items: [
                { text: 'Overview', link: '/api/' },
                { text: 'Types', link: '/api/types' },
                { text: 'Parser Functions', link: '/api/parser' },
                { text: 'Validation', link: '/api/validation' },
                { text: 'Utilities', link: '/api/utilities' }
              ]
            }
          ],
          '/examples/': [
            {
              text: 'Examples',
              items: [
                { text: 'Overview', link: '/examples/' },
                { text: 'Basic Parsing', link: '/examples/basic-parsing' },
                { text: 'Workflow Parsing', link: '/examples/workflow-parsing' },
                { text: 'Validation', link: '/examples/validation' },
                { text: 'Reusable Workflows', link: '/examples/reusable-workflows' },
                { text: 'Utility Functions', link: '/examples/utilities' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/github-action-parser' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024 GitHub Action Parser'
        }
      }
    },
    zh: {
      label: '简体中文',
      lang: 'zh-CN',
      title: 'GitHub Action Parser',
      description: '用于解析、验证和处理GitHub Action YAML文件的Go库',
      themeConfig: {
        nav: [
          { text: '首页', link: '/zh/' },
          { text: '快速开始', link: '/zh/getting-started' },
          { text: 'API 参考', link: '/zh/api/' },
          { text: '示例', link: '/zh/examples/' },
          { text: 'GitHub', link: 'https://github.com/scagogogo/github-action-parser' }
        ],
        sidebar: {
          '/zh/api/': [
            {
              text: 'API 参考',
              items: [
                { text: '概述', link: '/zh/api/' },
                { text: '类型定义', link: '/zh/api/types' },
                { text: '解析函数', link: '/zh/api/parser' },
                { text: '验证功能', link: '/zh/api/validation' },
                { text: '工具函数', link: '/zh/api/utilities' }
              ]
            }
          ],
          '/zh/examples/': [
            {
              text: '示例',
              items: [
                { text: '概述', link: '/zh/examples/' },
                { text: '基本解析', link: '/zh/examples/basic-parsing' },
                { text: '工作流解析', link: '/zh/examples/workflow-parsing' },
                { text: '验证功能', link: '/zh/examples/validation' },
                { text: '可重用工作流', link: '/zh/examples/reusable-workflows' },
                { text: '工具函数', link: '/zh/examples/utilities' }
              ]
            }
          ]
        },
        socialLinks: [
          { icon: 'github', link: 'https://github.com/scagogogo/github-action-parser' }
        ],
        footer: {
          message: '基于 MIT 许可证发布。',
          copyright: 'Copyright © 2024 GitHub Action Parser'
        }
      }
    }
  },
  
  themeConfig: {
    search: {
      provider: 'local'
    }
  }
})
