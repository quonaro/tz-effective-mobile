<script setup lang="ts">
import type { Toast } from '~/composables/useToast'

interface Props {
  toast: Toast
}

const props = defineProps<Props>()
const emit = defineEmits<{
  dismiss: [id: string]
}>()

const variantClasses = computed(() => {
  switch (props.toast.variant) {
    case 'destructive':
      return 'border-red-500 bg-red-500/10 text-red-500'
    default:
      return 'border-border bg-card text-foreground'
  }
})

const icon = computed(() => {
  switch (props.toast.variant) {
    case 'destructive':
      return '✕'
    default:
      return '✓'
  }
})
</script>

<template>
  <div
    class="pointer-events-auto relative flex w-full items-center justify-between gap-4 overflow-hidden rounded-md border p-4 pr-6 shadow-lg transition-all animate-in slide-in-from-bottom-5 fade-in duration-300"
    :class="variantClasses"
    role="alert"
  >
    <div class="flex items-start gap-3">
      <span class="mt-0.5 text-lg">{{ icon }}</span>
      <div class="grid gap-1">
        <div v-if="toast.title" class="text-sm font-semibold">
          {{ toast.title }}
        </div>
        <div v-if="toast.description" class="text-sm opacity-90">
          {{ toast.description }}
        </div>
      </div>
    </div>
    <button
      class="absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-0 transition-opacity hover:text-foreground focus:opacity-100 focus:outline-none group-hover:opacity-100"
      @click="emit('dismiss', toast.id)"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M18 6 6 18" /><path d="m6 6 12 12" />
      </svg>
    </button>
  </div>
</template>
