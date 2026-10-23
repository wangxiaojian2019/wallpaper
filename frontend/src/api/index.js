import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 15000
})

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res.data,
  err => {
    const msg = err.response?.data?.message || err.message || 'Network error'
    return Promise.reject(new Error(msg))
  }
)

// Auth
export const authAPI = {
  login: (data) => api.post('/auth/login', data),
  register: (data) => api.post('/auth/register', data),
  me: () => api.get('/auth/me')
}

// Products
export const productAPI = {
  list: (params) => api.get('/products', { params }),
  get: (id) => api.get(`/products/${id}`),
  adminGet: (id) => api.get(`/admin/products/${id}`),
  create: (data) => api.post('/admin/products', data),
  batchAction: (action, ids) => api.post('/admin/products/batch', { action, ids }),
  update: (id, data) => api.put(`/admin/products/${id}`, data),
  togglePin: (id) => api.put(`/admin/products/${id}/pin`),
  updateSort: (id, sortOrder) => api.put(`/admin/products/${id}/sort-order`, { sort_order: sortOrder }),
  delete: (id) => api.delete(`/admin/products/${id}`),
  uploadImage: (formData) => api.post('/admin/upload/image', formData, { timeout: 120000 }),
  uploadInteractive: (formData) => api.post('/admin/upload/interactive', formData, { timeout: 120000 })
}

// Categories
export const categoryAPI = {
  listPublic: () => api.get('/categories'),
  adminList: () => api.get('/admin/categories'),
  create: (data) => api.post('/admin/categories', data),
  update: (id, data) => api.put(`/admin/categories/${id}`, data),
  delete: (id) => api.delete(`/admin/categories/${id}`)
}

// Universal Tutorial Template
export const tutorialAPI = {
  getPublic: () => api.get('/config/tutorial'),
  adminGet: () => api.get('/admin/tutorial'),
  adminUpdate: (data) => api.put('/admin/tutorial', data)
}

// Orders
export const orderAPI = {
  create: (data) => api.post('/orders', data),
  list: () => api.get('/orders/my'),
  get: (id) => api.get(`/orders/${id}`),
  lookup: (data) => api.post('/orders/lookup', data),
  adminList: (params) => api.get('/admin/orders', { params }),
  updateStatus: (id, status) => api.put(`/admin/orders/${id}/status`, { status })
}

// PayPal
export const paypalAPI = {
  createOrder: (data) => api.post('/payment/paypal/create', data),
  captureOrder: (orderId) => api.post(`/payment/paypal/capture/${orderId}`),
  createSubscription: (planType) => api.post('/payment/paypal/subscription', { plan_type: planType }),
  cancelSubscription: () => api.post('/payment/paypal/subscription/cancel')
}

// Users (Admin)
export const userAPI = {
  list: (params) => api.get('/admin/users', { params }),
  get: (id) => api.get(`/admin/users/${id}`),
  create: (data) => api.post('/admin/users', data),
  update: (id, data) => api.put(`/admin/users/${id}`, data),
  delete: (id) => api.delete(`/admin/users/${id}`),
  resetPassword: (id, password) => api.post(`/admin/users/${id}/reset-password`, { password }),
  forceLogout: (id) => api.post(`/admin/users/${id}/force-logout`),
  stats: () => api.get('/admin/stats')
}

// Visitors (Admin)
export const visitorAPI = {
  list: (params) => api.get('/admin/visitors', { params }),
  stats: () => api.get('/admin/stats/visits')
}

// Payment Channels
export const paymentChannelAPI = {
  listPublic: () => api.get('/payment/channels'),
  adminList: () => api.get('/admin/payment/channels'),
  create: (data) => api.post('/admin/payment/channels', data),
  update: (id, data) => api.put(`/admin/payment/channels/${id}`, data),
  delete: (id) => api.delete(`/admin/payment/channels/${id}`)
}

// Download (unlock)
export const downloadAPI = {
  getLinks: (productId) => api.get(`/products/${productId}/download`),
  unlockWithCode: (productId, data) => api.post(`/products/${productId}/unlock-code`, data)
}

// Config & Settings
export const configAPI = {
  getPublic: () => api.get('/config/public')
}

export const settingsAPI = {
  getPayment: () => api.get('/admin/settings/payment'),
  updatePayment: (data) => api.post('/admin/settings/payment', data),
  getAll: () => api.get('/admin/settings/all'),
  updateAll: (data) => api.post('/admin/settings/all', data)
}

// Redeem Codes (Tokens & Vouchers)
export const redeemAPI = {
  list: (params) => api.get('/admin/redeem-codes', { params }),
  create: (data) => api.post('/admin/redeem-codes', data),
  toggle: (id) => api.put(`/admin/redeem-codes/${id}`),
  delete: (id) => api.delete(`/admin/redeem-codes/${id}`),
  logs: (params) => api.get('/admin/redeem-logs', { params }),
  apply: (data) => api.post('/redeem/apply', data)
}

export default api
