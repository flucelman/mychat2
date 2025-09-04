<template>
  <!-- 遮罩层 -->
  <div v-if="!globalSettingStore.isLogin" class="auth-overlay">
    <Auth class="auth-modal"/>
  </div>
  <router-view></router-view>
</template>

<script setup>
import '@/assets/css/theme.scss'
import { watch, onMounted, onUnmounted } from 'vue'
import { useGlobalSettingStore } from '@/stores/global_setting'
import Auth from '@/views/auth.vue'

const globalSettingStore = useGlobalSettingStore()

// 应用主题到HTML元素
const applyTheme = (theme) => {
  document.documentElement.setAttribute('theme', theme)
}


// 判断移动端和pc端 - 更准确的检测方法
const checkIsMobile = () => {
  // 方法1: 检测用户代理字符串
  const userAgent = navigator.userAgent.toLowerCase()
  const mobileKeywords = [
    'android', 'iphone', 'ipad', 'ipod', 'blackberry', 'windows phone',
    'mobile', 'webos', 'opera mini', 'iemobile', 'kindle'
  ]
  const isMobileUserAgent = mobileKeywords.some(keyword => userAgent.includes(keyword))
  
  // 方法2: 检测触摸支持
  const hasTouchSupport = 'ontouchstart' in window || navigator.maxTouchPoints > 0
  
  // 方法3: 检测屏幕尺寸和像素密度
  const screenWidth = window.innerWidth
  const screenHeight = window.innerHeight
  const isSmallScreen = screenWidth <= 768 || screenHeight <= 1024
  
  // 方法4: 检测设备方向支持
  const hasOrientationSupport = 'orientation' in window || 'onorientationchange' in window
  
  // 方法5: 检测CSS媒体查询
  const isPointerCoarse = window.matchMedia('(pointer: coarse)').matches
  const isHoverNone = window.matchMedia('(hover: none)').matches
  
  // 综合判断 - 多个条件满足则认为是移动端
  const mobileIndicators = [
    isMobileUserAgent,
    hasTouchSupport && isSmallScreen,
    isPointerCoarse,
    isHoverNone,
    hasOrientationSupport && isSmallScreen
  ]
  
  const mobileScore = mobileIndicators.filter(Boolean).length
  
  // 如果满足2个或以上条件，则认为是移动端
  globalSettingStore.isMobile = mobileScore >= 2
}

// 初始化主题
onMounted(() => {
  applyTheme(globalSettingStore.theme)
  checkIsMobile()
  
  // 监听窗口尺寸变化
  const handleResize = () => {
    checkIsMobile()
  }
  
  // 监听设备方向变化
  const handleOrientationChange = () => {
    // 延迟执行，确保方向变化完成后再检测
    setTimeout(() => {
      checkIsMobile()
    }, 100)
  }
  
  // 添加事件监听器
  window.addEventListener('resize', handleResize)
  window.addEventListener('orientationchange', handleOrientationChange)
  
  // 使用onUnmounted清理事件监听器
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    window.removeEventListener('orientationchange', handleOrientationChange)
  })
})

// 监听主题变化
watch(
  () => globalSettingStore.theme,
  (newTheme) => {
    applyTheme(newTheme)
  }
)
</script>

<style>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  /* 禁用移动端点击高亮效果 */
  -webkit-tap-highlight-color: transparent;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

html, body {
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  background-color: var(--background-color);
  color: var(--text-color);
}

#app {
  width: 100%;
  height: 100vh;
  height: 100dvh;
}

.el-drawer__body{
    padding: 0 !important;
}

/* 遮罩层样式 */
.auth-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px); /* 添加模糊效果 */
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

/* 弹窗内容样式 */
.auth-modal {
  background: var(--background-color);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  max-width: 400px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  animation: modalFadeIn 0.3s ease-out;
}

/* 弹窗动画 */
@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: scale(0.9) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.el-message {
    z-index: 99999 !important;
    position: fixed !important;
}
</style>