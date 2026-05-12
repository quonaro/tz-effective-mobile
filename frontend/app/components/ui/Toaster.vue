<script setup lang="ts">
import { useToast } from '~/composables/useToast'

const { toasts, remove } = useToast()
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 w-full max-w-sm">
      <Toast
        v-for="toast in toasts"
        :key="toast.id"
        :toast="toast"
        @dismiss="remove"
      />
    </div>
  </Teleport>
</template>
