import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const token = ref(localStorage.getItem('token') || '')

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isMember = computed(() => {
    if (!user.value?.membership_expire) return false
    return new Date(user.value.membership_expire) > new Date()
  })
  const membershipType = computed(() => user.value?.membership_type || 'none')

  function setAuth(userData, tokenStr) {
    user.value = userData
    token.value = tokenStr
    localStorage.setItem('user', JSON.stringify(userData))
    localStorage.setItem('token', tokenStr)
    axios.defaults.headers.common['Authorization'] = `Bearer ${tokenStr}`
  }

  function logout() {
    user.value = null
    token.value = ''
    localStorage.removeItem('user')
    localStorage.removeItem('token')
    delete axios.defaults.headers.common['Authorization']
  }

  function init() {
    if (token.value) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
    }
  }

  return { user, token, isLoggedIn, isAdmin, isMember, membershipType, setAuth, logout, init }
})

export const useCartStore = defineStore('cart', () => {
  const items = ref(JSON.parse(localStorage.getItem('cart') || '[]'))

  function addItem(product) {
    if (!items.value.find(i => i.id === product.id)) {
      items.value.push(product)
      localStorage.setItem('cart', JSON.stringify(items.value))
    }
  }

  function removeItem(id) {
    items.value = items.value.filter(i => i.id !== id)
    localStorage.setItem('cart', JSON.stringify(items.value))
  }

  function clearCart() {
    items.value = []
    localStorage.setItem('cart', JSON.stringify([]))
  }

  const total = computed(() => items.value.reduce((sum, i) => sum + Number(i.price), 0))
  const count = computed(() => items.value.length)

  return { items, total, count, addItem, removeItem, clearCart }
})
