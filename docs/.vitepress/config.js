import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Bolt",
  description: "A fast, single-binary developer Swiss Army knife CLI.",
  base: '/bolt/',
  cleanUrls: true,
  themeConfig: {
    logo: '⚡',
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/installation' },
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Installation', link: '/guide/installation' },
          { text: 'Usage', link: '/guide/usage' }
        ]
      },
      {
        text: 'Commands',
        items: [
          { text: '🔑 ID Generators', link: '/guide/commands/id-generators' },
          { text: '📦 Data & Encoding', link: '/guide/commands/data-encoding' },
          { text: '🔒 Security', link: '/guide/commands/security' },
          { text: '⏰ Time & Date', link: '/guide/commands/time-date' },
          { text: '🌐 HTTP & Network', link: '/guide/commands/http-network' },
          { text: '✏️ Text & Strings', link: '/guide/commands/text-strings' },
          { text: '📁 Files & System', link: '/guide/commands/files-system' },
          { text: '🖼️ Image Processing', link: '/guide/commands/image-processing' },
          { text: '🧰 Utilities', link: '/guide/commands/utilities' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/AkshayS96/bolt' },
      { icon: 'twitter', link: 'https://x.com/__akshaysolanki' }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present Akshay Solanki'
    }
  }
})
