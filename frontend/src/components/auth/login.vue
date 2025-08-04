<template>
    <div class="login-container">
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.email') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="email" :placeholder="$t('message.emailPlaceholder')" :disabled="loading" />
                </div>
            </div>
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.password') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="password" type="password" :placeholder="$t('message.passwordPlaceholder')" :disabled="loading" @keyup.enter="login" />
                </div>
            </div>
            <div class="input-item">
                <el-button type="primary" @click="login" :loading="loading">
                    {{ loading ? $t('message.loginInProgress') : $t('message.login') }}
                </el-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useGlobalSettingStore } from '@/stores/global_setting'
import { http } from '@/utils/http/client'
import { API } from '@/config/api'

const { t } = useI18n()
const globalSettingStore = useGlobalSettingStore()
const email = ref('')
const password = ref('')
const loading = ref(false)

const login = async () => {
    // 验证输入
    if (!email.value.trim() || !password.value.trim()) {
        ElMessage.error(t('message.emptyEmailOrPassword'))
        return
    }
    
    loading.value = true
    
    try {
        const response = await http.post(API.backend_url + '/api/auth/login', {
            email: email.value.trim(), 
            password: password.value
        })
        
        // HTTP客户端返回的数据在response.data中
        const data = response.data
        
        if (data.success) {
            // 登录成功
            globalSettingStore.isLogin = true
            globalSettingStore.userToken = data.token
            ElMessage.success(data.message || t('message.loginSuccess'))
        } else {
            // 登录失败，显示后端返回的具体错误信息
            ElMessage.error(data.message || t('message.loginFailed'))
        }
    } catch (err) {
        console.error('Login error:', err)
        
        // 根据错误类型显示不同的错误信息
        if (err.message.includes('404')) {
            ElMessage.error(t('message.userNotFound'))
        } else if (err.message.includes('401')) {
            ElMessage.error(t('message.invalidPassword'))
        } else if (err.message.includes('500')) {
            ElMessage.error(t('message.serverError'))
        } else if (err.name === 'TypeError' || err.message.includes('Failed to fetch')) {
            ElMessage.error(t('message.networkError'))
        } else {
            ElMessage.error(t('message.loginFailed'))
        }
    } finally {
        loading.value = false
    }
}
</script>

<style scoped>
.input-container{
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin: 20px;
}
.input-item{
    display: flex;
    align-items: center;
}
</style>