<script setup lang="ts">
const colorMode = useColorMode()
const config = useRuntimeConfig()
const buildId = config.public.buildId as string

function toggleTheme() {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'
}
</script>

<template>
  <div class="min-h-screen flex flex-col bg-background text-foreground">
    <header class="border-b border-border bg-card px-6 py-4">
      <div class="mx-auto flex w-full max-w-screen-xl items-center justify-between">
        <div class="flex items-center gap-3">
          <h1 class="text-xl font-semibold">Subscriptions Manager</h1>
          <span class="px-2 py-0.5 text-xs rounded-full bg-brand/20 text-brand border border-brand/30 font-mono" :title="'Build: ' + buildId">{{ buildId }}</span>
        </div>
        <div class="flex items-center gap-4">
          <NuxtLink
            to="/"
            class="text-base text-muted-foreground hover:text-foreground"
            active-class="text-brand font-medium"
          >
            Subscriptions
          </NuxtLink>
          <NuxtLink
            to="/total"
            class="text-base text-muted-foreground hover:text-foreground"
            active-class="text-brand font-medium"
          >
            Total Cost
          </NuxtLink>
          <NuxtLink
            to="/tz"
            class="text-base text-muted-foreground hover:text-foreground"
            active-class="text-brand font-medium"
          >
            Spec
          </NuxtLink>
          <button
            type="button"
            class="p-2 rounded-md hover:bg-muted transition-colors"
            @click="toggleTheme"
            :title="colorMode.value === 'dark' ? 'Switch to light' : 'Switch to dark'"
          >
            <span v-if="colorMode.value === 'dark'">☀️</span>
            <span v-else>🌙</span>
          </button>
        </div>
      </div>
    </header>
    <main class="mx-auto w-full max-w-screen-xl p-6 flex-grow">
      <slot />
    </main>
    <footer class="border-t border-border bg-card px-6 py-6 mt-auto">
      <div class="mx-auto w-full max-w-screen-xl text-center">
        <p class="text-base text-muted-foreground mb-2">
          Subscriptions Manager — application for managing online subscriptions
        </p>
        <a
          href="https://github.com/quonaro/tz-effective-mobile"
          target="_blank"
          rel="noopener noreferrer"
          class="text-base text-brand hover:underline transition-colors"
        >
          GitHub: quonaro/tz-effective-mobile
        </a>
      </div>
    </footer>
  </div>
</template>
