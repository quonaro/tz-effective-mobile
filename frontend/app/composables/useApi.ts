import { useToast } from './useToast'

export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string
  const { error } = useToast()

  const $api = $fetch.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
    async onResponseError({ response }) {
      const status = response.status
      const data = response._data

      let title = 'Error'
      let description = 'Something went wrong'

      if (status === 400) {
        title = 'Bad Request'
        description = data?.detail || 'Invalid request data'
      } else if (status === 404) {
        title = 'Not Found'
        description = data?.detail || 'Resource not found'
      } else if (status === 500) {
        title = 'Server Error'
        description = data?.errors?.[0]?.message || data?.detail || 'Internal server error'
      } else if (status >= 400 && status < 500) {
        title = 'Request Failed'
        description = data?.detail || `HTTP ${status}`
      }

      error(title, description)
    },
  })

  return {
    $api,
    baseURL,
  }
}
