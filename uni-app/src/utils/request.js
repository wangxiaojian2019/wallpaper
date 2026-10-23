// Multi-terminal unified request client (Web, H5, WeChat Mini Program)
const BASE_URL = (function () {
  // #ifdef H5
  if (typeof window !== 'undefined' && window.location) {
    return import.meta.env?.VITE_API_URL || (window.location.origin + '/api')
  }
  return '/api'
  // #endif
  // #ifndef H5
  return import.meta.env?.VITE_API_URL || 'https://wallpaper.icangoanyway.com/api'
  // #endif
})()

export function request(options = {}) {
  const token = uni.getStorageSync('token') || ''
  const header = {
    'Content-Type': 'application/json',
    ...(options.header || {})
  }

  if (token) {
    header['Authorization'] = `Bearer ${token}`
  }

  let url = options.url.startsWith('http') ? options.url : `${BASE_URL}${options.url.startsWith('/') ? options.url : '/' + options.url}`
  const method = (options.method || 'GET').toUpperCase()
  if (method === 'GET' && options.data && typeof options.data === 'object') {
    const queryParts = []
    Object.keys(options.data).forEach(key => {
      const val = options.data[key]
      if (val !== undefined && val !== null && val !== '') {
        queryParts.push(`${encodeURIComponent(key)}=${encodeURIComponent(val)}`)
      }
    })
    if (queryParts.length > 0) {
      url += (url.includes('?') ? '&' : '?') + queryParts.join('&')
    }
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url,
      method,
      data: options.data || {},
      header,
      timeout: options.timeout || 15000,
      success: (res) => {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data)
        } else if (res.statusCode === 401) {
          uni.removeStorageSync('token')
          uni.removeStorageSync('user_info')
          const err = new Error(res.data?.message || '登录已过期，请重新登录')
          err.statusCode = 401
          err.data = res.data
          reject(err)
        } else {
          const errMsg = res.data?.message || res.data?.error || `请求失败 (${res.statusCode})`
          const err = new Error(errMsg)
          err.statusCode = res.statusCode
          err.data = res.data
          reject(err)
        }
      },
      fail: (err) => {
        const msg = err.errMsg?.includes('timeout') ? '网络连接超时，请重试' : '无法连接服务器，请检查网络'
        const error = new Error(msg)
        error.original = err
        reject(error)
      }
    })
  })
}

export default {
  get: (url, data, options = {}) => request({ ...options, url, data, method: 'GET' }),
  post: (url, data, options = {}) => request({ ...options, url, data, method: 'POST' }),
  put: (url, data, options = {}) => request({ ...options, url, data, method: 'PUT' }),
  delete: (url, data, options = {}) => request({ ...options, url, data, method: 'DELETE' }),
  BASE_URL
}
