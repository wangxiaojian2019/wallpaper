import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(uni.getStorageSync('token') || '')
  const userInfo = ref(uni.getStorageSync('user_info') ? JSON.parse(uni.getStorageSync('user_info')) : null)

  const isLoggedIn = computed(() => !!token.value)
  const isMember = computed(() => {
    if (!userInfo.value) return false
    if (userInfo.value.role === 'admin') return true
    if (userInfo.value.membership_type === 'lifetime') return true
    if (!userInfo.value.is_member) return false
    if (!userInfo.value.membership_expire) return false
    return new Date(userInfo.value.membership_expire) > new Date()
  })

  const membershipType = computed(() => {
    if (!userInfo.value) return 'none'
    if (userInfo.value.role === 'admin') return 'lifetime'
    return userInfo.value.membership_type || 'none'
  })

  function setAuth(data) {
    token.value = data.token
    userInfo.value = data.user
    uni.setStorageSync('token', data.token)
    uni.setStorageSync('user_info', JSON.stringify(data.user))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    uni.removeStorageSync('token')
    uni.removeStorageSync('user_info')
  }

  async function login(data) {
    const res = await authAPI.login(data)
    const tokenVal = res.token || res.data?.token
    const userVal = res.user || res.data?.user
    if (tokenVal) {
      setAuth({ token: tokenVal, user: userVal })
    }
    return res
  }

  async function register(data) {
    const res = await authAPI.register(data)
    const tokenVal = res.token || res.data?.token
    const userVal = res.user || res.data?.user
    if (tokenVal) {
      setAuth({ token: tokenVal, user: userVal })
    }
    return res
  }

  async function fetchProfile() {
    if (!token.value) return null
    try {
      const res = await authAPI.me()
      const user = res.user || res.data?.user || res
      userInfo.value = user
      uni.setStorageSync('user_info', JSON.stringify(user))
      return user
    } catch (err) {
      console.warn('Failed to refresh user profile:', err)
      return null
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isMember,
    membershipType,
    setAuth,
    login,
    register,
    logout,
    fetchProfile
  }
})
