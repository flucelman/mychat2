import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useGlobalSettingStore = defineStore('globalSetting', () => {
    // 语言
    const lang = ref('zhCN')
    
    // 主题
    const theme = ref('light')
    
    // 判断移动端和pc端
    const isMobile = ref(false)

    // 是否登录
    const isLogin = ref(false)
    const userToken = ref('')

    return {
        lang,
        theme,
        isMobile,
        isLogin,
        userToken
    }
}, {
    persist: ['lang', 'theme', 'isLogin', 'userToken']
})