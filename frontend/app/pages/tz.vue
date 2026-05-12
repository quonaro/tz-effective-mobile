<script setup lang="ts">
import { ref, onMounted } from 'vue'

const content = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await fetch('/tz.md')
    if (!res.ok) {
      throw new Error('Failed to load TZ')
    }
    content.value = await res.text()
  } catch (e) {
    content.value = 'Failed to load technical specification'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="rounded-lg border border-border bg-card p-6 animate-fade-in">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-semibold text-foreground">
        Техническое задание (ТЗ)
      </h1>
      <NuxtLink
        to="/"
        class="rounded-md border border-border px-4 py-2 text-sm text-muted-foreground hover:bg-brand-dim hover:text-brand transition-colors"
      >
        ← Back
      </NuxtLink>
    </div>

    <div v-if="loading" class="py-8 text-center text-muted-foreground">
      Loading...
    </div>

    <div v-else class="prose prose-invert max-w-none">
      <pre class="whitespace-pre-wrap font-mono text-sm text-foreground bg-muted p-4 rounded-md">{{ content }}</pre>
    </div>
  </div>
</template>
