import http from '@/utils/request'

export const productAPI = {
  list: (params) => http.get('/products', params),
  get: (id, params) => http.get(`/products/${id}`, params),
  getDownloads: (id, params) => http.get(`/products/${id}/download`, params),
  unlockWithCode: (id, code) => http.post(`/products/${id}/unlock-code`, { code })
}

export const authAPI = {
  login: (data) => http.post('/auth/login', data),
  register: (data) => http.post('/auth/register', data),
  me: () => http.get('/auth/me')
}

export const orderAPI = {
  create: (data) => http.post('/orders', data),
  lookup: (email) => http.post('/orders/lookup', { email }),
  my: () => http.get('/orders/my'),
  get: (id) => http.get(`/orders/${id}`)
}

export const paypalAPI = {
  createOrder: (data) => http.post('/payment/paypal/create', data),
  captureOrder: (orderId) => http.post(`/payment/paypal/capture/${orderId}`),
  createSubscription: (planType) => http.post('/payment/paypal/subscription', { plan_type: planType })
}

export const configAPI = {
  getPublic: () => http.get('/config/public')
}

export const categoryAPI = {
  list: () => http.get('/categories')
}

export const tutorialAPI = {
  get: () => http.get('/config/tutorial')
}
