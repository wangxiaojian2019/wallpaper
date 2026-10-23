import { ref } from 'vue'

const toasts = ref([])
let id = 0

function addToast(message, type = 'info', duration = 3000) {
  const toastId = ++id
  toasts.value.push({ id: toastId, message, type })
  setTimeout(() => {
    removeToast(toastId)
  }, duration)
}

function removeToast(toastId) {
  const idx = toasts.value.findIndex(t => t.id === toastId)
  if (idx !== -1) toasts.value.splice(idx, 1)
}

export function useToast() {
  return {
    toasts,
    toast: addToast,
    success: (msg) => addToast(msg, 'success'),
    error: (msg) => addToast(msg, 'error'),
    info: (msg) => addToast(msg, 'info'),
    warning: (msg) => addToast(msg, 'warning'),
    remove: removeToast
  }
}
