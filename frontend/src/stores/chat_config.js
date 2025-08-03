import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatConfigStore = defineStore('chatConfig', () => {
    // 是否显示历史记录抽屉
    const showDrawer = ref(true)
    // 是否展示lang、theme设置
    const openEyes = ref(false)
    return {
        showDrawer,
        openEyes
    }
}, {
    persist: ['showDrawer']
})