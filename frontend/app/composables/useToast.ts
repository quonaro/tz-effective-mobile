export interface Toast {
  id: string
  title: string
  description?: string
  variant?: 'default' | 'destructive'
  duration?: number
}

let toastIdCounter = 0

export function useToast() {
  const toasts = useState<Toast[]>('toasts', () => [])

  function generateId(): string {
    return `toast-${Date.now()}-${++toastIdCounter}`
  }

  function add(toast: Omit<Toast, 'id'>) {
    const id = generateId()
    const newToast: Toast = {
      id,
      duration: 5000,
      ...toast,
    }
    toasts.value.push(newToast)

    // Auto remove after duration
    setTimeout(() => {
      remove(id)
    }, newToast.duration)

    return id
  }

  function remove(id: string) {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  function success(title: string, description?: string) {
    return add({ title, description, variant: 'default' })
  }

  function error(title: string, description?: string) {
    return add({ title, description, variant: 'destructive' })
  }

  return {
    toasts: readonly(toasts),
    add,
    remove,
    success,
    error,
  }
}
