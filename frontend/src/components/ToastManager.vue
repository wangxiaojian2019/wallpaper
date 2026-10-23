<template>
  <teleport to="body">
    <div class="toast-container">
      <transition-group name="toast-anim">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="toast"
          :class="`toast-${toast.type}`"
          @click="remove(toast.id)"
        >
          <span class="toast-icon">{{ icons[toast.type] }}</span>
          <span class="toast-msg">{{ toast.message }}</span>
        </div>
      </transition-group>
    </div>
  </teleport>
</template>

<script setup>
import { useToast } from '@/composables/useToast'

const { toasts, remove } = useToast()

const icons = {
  success: '✓',
  error: '✕',
  info: 'ℹ',
  warning: '⚠'
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 80px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 18px;
  border-radius: var(--r-md);
  color: white;
  font-size: 0.88rem;
  font-weight: 500;
  box-shadow: var(--shadow-md);
  min-width: 220px;
  max-width: 360px;
  cursor: pointer;
  pointer-events: all;
  backdrop-filter: blur(4px);
}

.toast-icon {
  font-size: 0.9rem;
  font-weight: 700;
  flex-shrink: 0;
}

.toast-msg {
  flex: 1;
  line-height: 1.4;
}

.toast-success { background: rgba(52,199,89,0.95); }
.toast-error { background: rgba(255,59,48,0.95); }
.toast-info { background: rgba(0,113,227,0.95); }
.toast-warning { background: rgba(255,159,10,0.95); }

/* Animations */
.toast-anim-enter-active {
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.toast-anim-leave-active {
  transition: all 0.2s ease;
}
.toast-anim-enter-from {
  opacity: 0;
  transform: translateX(20px) scale(0.9);
}
.toast-anim-leave-to {
  opacity: 0;
  transform: translateX(20px) scale(0.9);
}
.toast-anim-move {
  transition: transform 0.3s ease;
}
</style>
