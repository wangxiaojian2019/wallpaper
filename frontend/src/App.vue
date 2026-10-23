<template>
  <div id="app-root" :class="{ 'in-admin': isAdminRoute }">
    <AppNavbar v-if="!isAdminRoute" />
    <main class="page-content" :class="{ 'no-nav': isAdminRoute }">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
    <AppFooter v-if="!isAdminRoute" />
    <AppTabBar v-if="!isAdminRoute" />
    <ToastManager />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from '@/components/AppFooter.vue'
import AppTabBar from '@/components/AppTabBar.vue'
import ToastManager from '@/components/ToastManager.vue'

const route = useRoute()
const isAdminRoute = computed(() => {
  return route.path ? route.path.startsWith('/admin') : false
})
</script>

<style>
#app-root {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.page-content.no-nav {
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}

@media (max-width: 768px) {
  .page-content:not(.no-nav) {
    padding-bottom: calc(64px + env(safe-area-inset-bottom, 0px));
  }
}
</style>

