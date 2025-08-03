import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import ElementPlus from 'element-plus'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
// 导入中文语言包
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import i18n from './i18n'
import '@/assets/css/theme.scss'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

const app = createApp(App)

// 配置axios为全局属性，而不是插件
app.config.globalProperties.$axios = axios

app.use(pinia)
  .use(router)
  .use(ElementPlus, {
    locale: zhCn
  })
  .use(i18n)

// 在 Pinia 激活后设置语言
import { useGlobalSettingStore } from '@/stores/global_setting'
const globalSettingStore = useGlobalSettingStore()
i18n.global.locale.value = globalSettingStore.lang
  
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.mount('#app')