export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  const $api = $fetch.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  return {
    $api,
    baseURL,
  }
}
