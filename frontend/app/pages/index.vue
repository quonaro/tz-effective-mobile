<script setup lang="ts">
import { ref, computed } from 'vue'
import { useApi } from '~/composables/useApi'
import type { Subscription, SubscriptionListResponse } from '~/types'

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

const { $api } = useApi()
const subscriptions = ref<Subscription[]>([])
const total = ref(0)
const loading = ref(false)

// Pagination & Search
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const debouncedSearch = ref('')

const form = ref({
  service_name: '',
  price: 0,
  user_id: '',
  start_date: '',
  end_date: '',
})

// Computed for month inputs (YYYY-MM format for input type="month")
const startDateInput = computed({
  get: () => toInputMonth(form.value.start_date),
  set: (val) => { form.value.start_date = fromInputMonth(val) }
})

const endDateInput = computed({
  get: () => toInputMonth(form.value.end_date),
  set: (val) => { form.value.end_date = fromInputMonth(val) }
})

const editingId = ref<string | null>(null)

async function loadSubscriptions() {
  loading.value = true
  try {
    const offset = (currentPage.value - 1) * pageSize.value
    const query: Record<string, string | number> = {
      limit: pageSize.value,
      offset,
    }
    if (debouncedSearch.value) {
      query.service_name = debouncedSearch.value
    }

    const res = (await $api('/subscriptions', {
      query,
    })) as SubscriptionListResponse | null
    if (res) {
      subscriptions.value = res.data || []
      total.value = res.total || 0
    } else {
      subscriptions.value = []
      total.value = 0
    }
  } catch (error) {
    subscriptions.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Debounce search
let searchTimeout: ReturnType<typeof setTimeout> | null = null
watch(searchQuery, (newVal) => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    debouncedSearch.value = newVal
    currentPage.value = 1
    loadSubscriptions()
  }, 300)
})

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    loadSubscriptions()
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadSubscriptions()
  }
}

function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    loadSubscriptions()
  }
}

async function saveSubscription() {
  if (editingId.value) {
    // PUT update - без user_id
    const payload: any = {
      service_name: form.value.service_name,
      price: Number(form.value.price),
      start_date: form.value.start_date,
    }
    if (form.value.end_date) {
      payload.end_date = form.value.end_date
    }
    await $api(`/subscriptions/${editingId.value}`, {
      method: 'PUT',
      body: payload,
    })
  } else {
    // POST create - с user_id
    const payload: any = {
      service_name: form.value.service_name,
      price: Number(form.value.price),
      user_id: form.value.user_id,
      start_date: form.value.start_date,
    }
    if (form.value.end_date) {
      payload.end_date = form.value.end_date
    }
    await $api('/subscriptions', {
      method: 'POST',
      body: payload,
    })
  }

  resetForm()
  await loadSubscriptions()
}

async function deleteSubscription(id: string) {
  await $api(`/subscriptions/${id}`, { method: 'DELETE' })
  await loadSubscriptions()
}

function editSubscription(sub: Subscription) {
  editingId.value = sub.id
  form.value = {
    service_name: sub.service_name,
    price: sub.price,
    user_id: sub.user_id,
    start_date: formatDate(sub.start_date),
    end_date: sub.end_date ? formatDate(sub.end_date) : '',
  }
}

function resetForm() {
  editingId.value = null
  form.value = {
    service_name: '',
    price: 0,
    user_id: '',
    start_date: '',
    end_date: '',
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  // Handle ISO date format (2025-07-01T00:00:00Z or 2025-07-01)
  const isoMatch = dateStr.match(/^(\d{4})-(\d{2})/)
  if (isoMatch) {
    const [, year, month] = isoMatch
    return `${month}-${year}`
  }
  // If already in MM-YYYY format, return as-is
  if (/^\d{2}-\d{4}$/.test(dateStr)) {
    return dateStr
  }
  return dateStr
}

onMounted(loadSubscriptions)
</script>

<template>
  <div class="space-y-6">
    <div class="rounded-lg border border-border bg-card p-6 animate-fade-in">
      <h2 class="mb-4 text-lg font-medium text-brand">
        {{ editingId ? 'Edit Subscription' : 'New Subscription' }}
      </h2>
      <form class="grid gap-4 md:grid-cols-2" @submit.prevent="saveSubscription">
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">Service Name</label>
          <input v-model="form.service_name" type="text" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">Price (RUB)</label>
          <input v-model.number="form.price" type="number" required min="0" class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">User ID (UUID)</label>
          <div class="flex gap-2">
            <input v-model="form.user_id" type="text" required class="flex-1 rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <button
              type="button"
              class="rounded-md border border-border px-3 py-2 text-sm text-muted-foreground hover:bg-brand-dim hover:text-brand transition-colors whitespace-nowrap"
              @click="form.user_id = generateUUID()"
            >
              Generate
            </button>
          </div>
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">Start Date</label>
          <ClientOnly>
            <input v-model="startDateInput" type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <template #fallback>
              <input type="month" required class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            </template>
          </ClientOnly>
        </div>
        <div>
          <label class="mb-1 block text-base font-medium text-foreground">End Date (optional)</label>
          <ClientOnly>
            <input v-model="endDateInput" type="month" class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            <template #fallback>
              <input type="month" class="w-full rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand" />
            </template>
          </ClientOnly>
        </div>
        <div class="flex items-end gap-2 md:col-span-2">
          <button type="submit" class="rounded-md bg-brand px-4 py-2 text-sm font-medium text-[var(--brand-foreground)] hover:bg-brand/90 transition-colors">
            {{ editingId ? 'Update' : 'Create' }}
          </button>
          <button v-if="editingId" type="button" class="rounded-md border border-border px-4 py-2 text-sm text-muted-foreground hover:bg-brand-dim hover:text-brand transition-colors" @click="resetForm">
            Cancel
          </button>
        </div>
      </form>
    </div>

    <div class="rounded-lg border border-border bg-card animate-fade-in">
      <div class="border-b border-border px-6 py-4">
        <div class="flex items-center justify-between gap-4 flex-wrap">
          <div>
            <h2 class="text-lg font-medium text-brand">Subscriptions</h2>
            <p class="text-base text-muted-foreground">Total: {{ total }}</p>
          </div>
          <div class="flex items-center gap-2">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by service name..."
              class="w-64 rounded-md border border-border bg-muted px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:border-brand focus:outline-none focus:ring-1 focus:ring-brand"
            />
            <button
              v-if="searchQuery"
              type="button"
              class="rounded-md border border-border px-3 py-2 text-sm text-muted-foreground hover:bg-brand-dim hover:text-brand transition-colors"
              @click="searchQuery = ''; debouncedSearch = ''; currentPage = 1; loadSubscriptions()"
            >
              Clear
            </button>
          </div>
        </div>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead class="border-b border-border bg-muted">
            <tr>
              <th class="px-4 py-3 text-left font-medium text-foreground">Service</th>
              <th class="px-4 py-3 text-left font-medium text-foreground">Price</th>
              <th class="px-4 py-3 text-left font-medium text-foreground">User ID</th>
              <th class="px-4 py-3 text-left font-medium text-foreground">Start</th>
              <th class="px-4 py-3 text-left font-medium text-foreground">End</th>
              <th class="px-4 py-3 text-right font-medium text-foreground">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="sub in subscriptions" :key="sub.id" class="border-b border-border last:border-0 hover:bg-brand-dim transition-colors">
              <td class="px-4 py-3 text-foreground">{{ sub.service_name }}</td>
              <td class="px-4 py-3 text-foreground">{{ sub.price }} ₽</td>
              <td class="px-4 py-3 font-mono text-xs text-muted-foreground">{{ sub.user_id }}</td>
              <td class="px-4 py-3 text-muted-foreground">{{ formatDate(sub.start_date) }}</td>
              <td class="px-4 py-3 text-muted-foreground">{{ sub.end_date ? formatDate(sub.end_date) : '—' }}</td>
              <td class="px-4 py-3 text-right">
                <button class="mr-2 text-sm text-brand hover:underline transition-colors" @click="editSubscription(sub)">Edit</button>
                <button class="text-sm text-red-400 hover:underline transition-colors" @click="deleteSubscription(sub.id)">Delete</button>
              </td>
            </tr>
            <tr v-if="subscriptions.length === 0">
              <td colspan="6" class="px-4 py-8 text-center text-muted-foreground text-base">
                No subscriptions yet
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="border-t border-border px-6 py-4">
        <div class="flex items-center justify-between">
          <p class="text-sm text-muted-foreground">
            Page {{ currentPage }} of {{ totalPages }}
          </p>
          <div class="flex items-center gap-2">
            <button
              :disabled="currentPage === 1 || loading"
              class="rounded-md border border-border px-3 py-1 text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-brand-dim hover:text-brand transition-colors"
              @click="prevPage"
            >
              ← Prev
            </button>
            <div class="flex items-center gap-1">
              <button
                v-for="page in Math.min(5, totalPages)"
                :key="page"
                :class="[
                  'rounded-md px-3 py-1 text-sm transition-colors',
                  page === currentPage
                    ? 'bg-brand text-[var(--brand-foreground)]'
                    : 'border border-border hover:bg-brand-dim hover:text-brand'
                ]"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
              <span v-if="totalPages > 5" class="text-muted-foreground">...</span>
              <button
                v-if="totalPages > 5"
                :class="[
                  'rounded-md px-3 py-1 text-sm transition-colors',
                  currentPage > 5
                    ? 'bg-brand text-[var(--brand-foreground)]'
                    : 'border border-border hover:bg-brand-dim hover:text-brand'
                ]"
                @click="goToPage(totalPages)"
              >
                {{ totalPages }}
              </button>
            </div>
            <button
              :disabled="currentPage === totalPages || loading"
              class="rounded-md border border-border px-3 py-1 text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-brand-dim hover:text-brand transition-colors"
              @click="nextPage"
            >
              Next →
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
