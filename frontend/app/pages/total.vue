<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useApi } from '~/composables/useApi'
import type { TotalCostResponse, ServicesResponse } from '~/types'

const { $api } = useApi()

const filters = ref({
  user_id: '',
  service_name: '',
  start_month: '',
  end_month: '',
})

const result = ref<TotalCostResponse | null>(null)
const loading = ref(false)

// Service selector state
const services = ref<string[]>([])
const servicesLoading = ref(false)
const showDropdown = ref(false)
const searchQuery = ref('')
const selectedIndex = ref(-1)

const filteredServices = computed(() => {
  if (!searchQuery.value) return services.value
  const query = searchQuery.value.toLowerCase()
  return services.value.filter(s => s.toLowerCase().includes(query))
})

// Computed for month inputs (YYYY-MM format for input type="month")
const startMonthInput = computed({
  get: () => toInputMonth(filters.value.start_month),
  set: (val) => { filters.value.start_month = fromInputMonth(val) }
})

const endMonthInput = computed({
  get: () => toInputMonth(filters.value.end_month),
  set: (val) => { filters.value.end_month = fromInputMonth(val) }
})

async function loadServices() {
  servicesLoading.value = true
  try {
    const res = (await $api('/subscriptions/services')) as ServicesResponse | null
    if (res) {
      services.value = res.services || []
    }
  } catch (error) {
    services.value = []
  } finally {
    servicesLoading.value = false
  }
}

function selectService(name: string) {
  filters.value.service_name = name
  searchQuery.value = ''
  showDropdown.value = false
  selectedIndex.value = -1
}

function clearService() {
  filters.value.service_name = ''
  searchQuery.value = ''
}

function generateUUID(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID()
  }
  // Fallback for older browsers
  const array = new Uint8Array(16)
  crypto.getRandomValues(array)
  const val6 = array[6]
  const val8 = array[8]
  if (val6 === undefined || val8 === undefined) {
    return '00000000-0000-4000-8000-000000000000'
  }
  array[6] = (val6 & 0x0f) | 0x40
  array[8] = (val8 & 0x3f) | 0x80
  const hex = Array.from(array, b => b.toString(16).padStart(2, '0'))
  return `${hex[0]}${hex[1]}${hex[2]}${hex[3]}-${hex[4]}${hex[5]}-${hex[6]}${hex[7]}-${hex[8]}${hex[9]}-${hex[10]}${hex[11]}${hex[12]}${hex[13]}${hex[14]}${hex[15]}`
}

// Date format conversion helpers: MM-YYYY ↔ YYYY-MM
function toInputMonth(value: string): string {
  if (!value) return ''
  const match = value.match(/^(\d{2})-(\d{4})$/)
  if (!match) return ''
  return `${match[2]}-${match[1]}`
}

function fromInputMonth(value: string): string {
  if (!value) return ''
  const match = value.match(/^(\d{4})-(\d{2})$/)
  if (!match) return ''
  return `${match[2]}-${match[1]}`
}

function onKeydown(e: KeyboardEvent) {
  if (!showDropdown.value) return

  const items = filteredServices.value
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    selectedIndex.value = (selectedIndex.value + 1) % items.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selectedIndex.value = selectedIndex.value <= 0 ? items.length - 1 : selectedIndex.value - 1
  } else if (e.key === 'Enter' && selectedIndex.value >= 0) {
    e.preventDefault()
    const item = items[selectedIndex.value]
    if (item) {
      selectService(item)
    }
  } else if (e.key === 'Escape') {
    showDropdown.value = false
    selectedIndex.value = -1
  }
}

onMounted(loadServices)

async function calculateTotal() {
  loading.value = true
  try {
    const params: any = {
      start_month: filters.value.start_month,
      end_month: filters.value.end_month,
    }
    if (filters.value.user_id) {
      params.user_id = filters.value.user_id
    }
    if (filters.value.service_name) {
      params.service_name = filters.value.service_name
    }

    const res = (await $api('/subscriptions/total', {
      query: params,
    })) as TotalCostResponse | null
    result.value = res
  } catch (error) {
    result.value = null
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-lg border border-border bg-card p-6 animate-fade-in">
      <h2 class="mb-4 text-lg font-medium text-brand">Calculate Total Cost</h2>
      <form class="grid gap-4 md:grid-cols-2" @submit.prevent="calculateTotal">
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">Start Month</label>
          <ClientOnly>
            <input v-model="startMonthInput" type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <template #fallback>
              <input type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            </template>
          </ClientOnly>
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">End Month</label>
          <ClientOnly>
            <input v-model="endMonthInput" type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <template #fallback>
              <input type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            </template>
          </ClientOnly>
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">User ID (optional)</label>
          <div class="flex gap-2">
            <input v-model="filters.user_id" type="text" placeholder="UUID" class="flex-1 rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <button
              type="button"
              class="rounded-md border border-border px-3 py-2 text-sm text-muted-foreground hover:bg-brand-dim hover:text-brand transition-colors whitespace-nowrap"
              @click="filters.user_id = generateUUID()"
            >
              Generate
            </button>
          </div>
        </div>
        <div class="relative">
          <label class="mb-1 block text-base font-medium text-foreground">Service Name (optional)</label>
          <div class="relative">
            <input
              v-model="filters.service_name"
              type="text"
              placeholder="Netflix"
              class="w-full rounded-md border border-border bg-muted px-3 py-2 pr-8 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
              @focus="showDropdown = true"
              @keydown="onKeydown"
            />
            <button
              v-if="filters.service_name"
              type="button"
              class="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              @click="clearService"
            >
              ×
            </button>
          </div>
          <!-- Dropdown -->
          <div
            v-if="showDropdown"
            class="absolute z-10 mt-1 w-full rounded-md border border-border bg-card shadow-lg"
          >
            <!-- Search input in dropdown -->
            <div class="border-b border-border p-2">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search services..."
                class="w-full rounded border border-border bg-muted px-2 py-1 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none"
                @click.stop
              />
            </div>
            <!-- Service list -->
            <div class="max-h-48 overflow-auto py-1">
              <div v-if="servicesLoading" class="px-3 py-2 text-sm text-muted-foreground">
                Loading...
              </div>
              <div v-else-if="filteredServices.length === 0" class="px-3 py-2 text-sm text-muted-foreground">
                No services found
              </div>
              <button
                v-for="(service, index) in filteredServices"
                :key="service"
                type="button"
                class="block w-full px-3 py-2 text-left text-sm text-foreground hover:bg-brand-dim"
                :class="{ 'bg-brand-dim': index === selectedIndex }"
                @click="selectService(service)"
              >
                {{ service }}
              </button>
            </div>
          </div>
          <!-- Click outside to close -->
          <div
            v-if="showDropdown"
            class="fixed inset-0 z-0"
            @click="showDropdown = false"
          />
        </div>
        <div class="md:col-span-2">
          <button type="submit" class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-white hover:bg-brand/90 transition-colors">
            Calculate
          </button>
        </div>
      </form>
    </div>

    <div v-if="result" class="rounded-lg border border-border bg-card p-6 animate-fade-in">
      <h3 class="text-base text-muted-foreground">Total Cost</h3>
      <p class="mt-1 text-3xl font-bold text-brand">{{ result.total_cost }} ₽</p>
    </div>
  </div>
</template>
