export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',

  future: {
    compatibilityVersion: 4,
  },

  app: {
    head: {
      title: 'Subscriptions Manager',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      ],
    },
  },

  modules: [
    '@nuxtjs/color-mode',
    '@vueuse/nuxt',
    'shadcn-nuxt',
  ],

  shadcn: {
    prefix: '',
    componentDir: 'app/components/ui',
  },

  colorMode: {
    preference: 'system',
    fallback: 'light',
    classSuffix: '',
    storageKey: 'theme',
  },

  css: ['./app/assets/css/main.css'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
      buildId: process.env.NUXT_PUBLIC_BUILD_ID || 'dev',
    },
  },

  devServer: {
    port: 3000,
    host: '127.0.0.1',
  },

  typescript: {
    strict: true,
    typeCheck: false,
  },

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  vite: {
    optimizeDeps: {
      include: [
        '@tanstack/vue-query',
        'vee-validate',
        '@vee-validate/zod',
        'zod',
        'lucide-vue-next',
        'class-variance-authority',
        'reka-ui',
        'clsx',
        'tailwind-merge',
        'radix-vue',
        'shadcn-nuxt',
        'tailwindcss-animate',
      ],
    },
    server: {
      proxy: {
        '/api': {
          target: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
          changeOrigin: true,
          rewrite: (path: string) => path.replace(/^\/api/, ''),
        },
      },
    },
  },
})
