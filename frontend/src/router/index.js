import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue')
  },
  {
    path: '/wallpaper/:id',
    name: 'ProductDetail',
    component: () => import('@/views/ProductDetail.vue')
  },
  {
    path: '/product/:id',
    redirect: to => `/wallpaper/${to.params.id}`
  },
  {
    path: '/membership',
    name: 'Membership',
    component: () => import('@/views/Membership.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/payment/success',
    name: 'PaymentSuccess',
    component: () => import('@/views/PaymentSuccess.vue')
  },
  {
    path: '/payment/cancel',
    name: 'PaymentCancel',
    component: () => import('@/views/PaymentCancel.vue')
  },
  // Admin routes
  {
    path: '/admin',
    component: () => import('@/views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/AdminDashboard.vue')
      },
      {
        path: 'products',
        name: 'AdminProducts',
        component: () => import('@/views/admin/AdminProducts.vue')
      },
      {
        path: 'products/new',
        name: 'AdminProductNew',
        component: () => import('@/views/admin/AdminProductEdit.vue')
      },
      {
        path: 'products/:id/edit',
        name: 'AdminProductEdit',
        component: () => import('@/views/admin/AdminProductEdit.vue')
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/AdminOrders.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/AdminUsers.vue')
      },
      {
        path: 'visitors',
        name: 'AdminVisitors',
        component: () => import('@/views/admin/AdminVisitors.vue')
      },
      {
        path: 'stats',
        name: 'AdminStats',
        component: () => import('@/views/admin/AdminStats.vue')
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: () => import('@/views/admin/AdminSettings.vue')
      },
      {
        path: 'categories',
        name: 'AdminCategories',
        component: () => import('@/views/admin/AdminCategories.vue')
      },
      {
        path: 'tutorial',
        name: 'AdminTutorial',
        component: () => import('@/views/admin/AdminTutorial.vue')
      },
      {
        path: 'redeem-codes',
        name: 'AdminRedeemCodes',
        component: () => import('@/views/admin/AdminRedeemCodes.vue')
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0, behavior: 'smooth' }
  }
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  const isHostAdmin = typeof window !== 'undefined' && (
    window.location.hostname.includes('admin') || 
    window.location.port === '5174'
  )

  // On admin domain or admin port (5174), root path '/' should always direct to '/admin'
  if (to.path === '/' && isHostAdmin) {
    next(auth.isLoggedIn && auth.isAdmin ? { path: '/admin' } : { name: 'Login', query: { redirect: '/admin' } })
    return
  }

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (to.meta.requiresAdmin && !auth.isAdmin) {
    next(isHostAdmin ? { name: 'Login', query: { redirect: to.fullPath } } : { name: 'Home' })
  } else if (to.meta.guestOnly && auth.isLoggedIn) {
    next(auth.isAdmin || isHostAdmin ? { path: '/admin' } : { name: 'Home' })
  } else {
    next()
  }
})

export default router
