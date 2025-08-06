import { defineStore } from 'pinia'
import { ref } from 'vue'
import { http } from '@/utils/http/client'
import { ElMessage } from 'element-plus'
import { API } from '@/config/api'

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
    const checkToken = async () => {
        try {
            const response = await http.post(API.backend_url + '/api/auth/checkToken')
            if (response.data.code != 200) {
                isLogin.value = false
                userToken.value = ''
            } else {
                isLogin.value = true
            }
        } catch (error) {
            console.error('Token验证失败:', error)
            isLogin.value = false
            userToken.value = ''
        }
    }

    return {
        lang,
        theme,
        isMobile,
        isLogin,
        userToken,
        checkToken
    }
}, {
    persist: {
        pick: ['lang', 'theme', 'isLogin', 'userToken']
    }
})