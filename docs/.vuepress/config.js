module.exports = {
  theme: 'cosmos',
  title: 'Gridiron Documentation',
  locales: {
    '/': {
      lang: 'en-US',
    },
  },
  markdown: {
    extendMarkdown: (md) => {
      md.use(require('markdown-it-katex'))
    },
  },
  head: [
    [
      'link',
      {
        rel: 'stylesheet',
        href:
          'https://cdnjs.cloudflare.com/ajax/libs/KaTeX/0.5.1/katex.min.css',
      },
    ],
    [
      'link',
      {
        rel: 'stylesheet',
        href:
          'https://cdn.jsdelivr.net/github-markdown-css/2.2.1/github-markdown.css',
      },
    ],
    [
      'link',
      {rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon32.png'}],
    [
      'link',
      {rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon16.png'}],
    ['link', {rel: 'manifest', href: '/site.webmanifest'}],
    ['meta', {name: 'msapplication-TileColor', content: '#2e3148'}],
    ['meta', {name: 'theme-color', content: '#ffffff'}],
    ['link', {rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg'}],
  ],
  base: process.env.VUEPRESS_BASE || '/',
  plugins: [
    'vuepress-plugin-element-tabs',
  ],
  patterns: [
    '**/*.md',
    '**/*.vue',
    '!ethermint/**',
    '!cosmos-sdk/**',
    '!ibc-go/**'],
  themeConfig: {
    repo: 'gridiron-zone/gridiron',
    docsRepo: 'gridiron-zone/gridiron',
    docsBranch: 'main',
    docsDir: 'docs',
    editLinks: true,
    custom: true,
    logo: {
      src: '/gridiron-black.svg',
    },
    topbar: {
      banner: false,
    },
    sidebar: {
      auto: false,
      nav: [
        {
          title: 'Reference',
          children: [
            {
              title: 'Introduction',
              directory: true,
              path: '/intro',
            },
            {
              title: 'Basics',
              directory: true,
              path: '/basics',
            },
          ],
        },
        {
          title: 'APIs',
          children: [
            {
              title: 'JSON-RPC',
              directory: true,
              path: '/api/json-rpc',
            },
            {
              title: 'Protobuf Reference',
              directory: false,
              path: '/api/proto-docs',
            },
          ],
        },
        {
          title: 'Specifications',
          children: [
            {
              title: 'Modules',
              directory: true,
              path: '/modules',
            }],
        },
      ],
    },
    gutter: {
      title: 'Help & Support',
      chat: {
        title: 'Developer Chat',
        text: 'Chat with Gridiron developers on Discord.',
        url: 'https://discord.gg/gridiron',
        bg: 'linear-gradient(103.75deg, #1B1E36 0%, #22253F 100%)',
      },
      forum: {
        title: 'Gridiron Developer Forum',
        text: 'Join the Gridiron Developer Forum to learn more.',
        url: 'https://forum.cosmos.network/c/gridiron', // TODO: replace with commonwealth link
        bg: 'linear-gradient(221.79deg, #3D6B99 -1.08%, #336699 95.88%)',
        logo: 'ethereum-white',
      },
      github: {
        title: 'Found an Issue?',
        text: 'Help us improve this page by suggesting edits on GitHub.',
        bg: '#F8F9FC',
      },
    },
    footer: {
      logo: '/gridiron-black.svg',
      textLink: {
        text: 'gridiron.zone',
        url: 'https://gridiron.zone',
      },
      services: [
        {
          service: 'github',
          url: 'https://github.com/gridiron-zone/gridiron',
        },
        {
          service: 'medium',
          url: 'https://blog.gridiron.zone',
        },
        {
          service: 'twitter',
          url: 'https://twitter.com/GridironZone',
        },
        {
          service: 'telegram',
          url: 'https://t.me/GridironZone',
        },
        {
          service: 'discord',
          url: 'https://discord.gg/gridiron',
        },
      ],
      smallprint: 'This website is maintained by Gridiron Dev Team',
    },
  },
}
