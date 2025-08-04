<template>
    <div class="register-container">
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.username') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="username" :placeholder="$t('message.usernamePlaceholder')" :disabled="loading" />
                </div>
            </div>
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.email') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="email" :placeholder="$t('message.emailPlaceholder')" :disabled="loading" />
                </div>
            </div>
        </div>
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.password') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="password" type="password" :placeholder="$t('message.passwordPlaceholder')" :disabled="loading" />
                </div>
            </div>
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.confirmPassword') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="confirmPassword" type="password" :placeholder="$t('message.confirmPasswordPlaceholder')" :disabled="loading" @keyup.enter="register" />
                </div>
            </div>
            <div class="input-item">
                <el-button type="primary" @click="register" :loading="loading">
                    {{ loading ? $t('message.registerInProgress') : $t('message.register') }}
                </el-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { http } from '@/utils/http/client'
import { API } from '@/config/api'
import { useGlobalSettingStore } from '@/stores/global_setting'

const globalSettingStore = useGlobalSettingStore()
const { t } = useI18n()
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

// 邮箱格式验证
const isValidEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
}

// 输入验证
const validateInputs = () => {
    if (!username.value.trim()) {
        ElMessage.error(t('message.usernameRequired'))
        return false
    }
    
    if (!email.value.trim()) {
        ElMessage.error(t('message.emailRequired'))
        return false
    }
    
    if (!isValidEmail(email.value.trim())) {
        ElMessage.error(t('message.invalidEmailFormat'))
        return false
    }
    
    if (!password.value) {
        ElMessage.error(t('message.passwordRequired'))
        return false
    }
    
    if (password.value.length < 8) {
        ElMessage.error(t('message.passwordTooShort'))
        return false
    }
    
    if (password.value !== confirmPassword.value) {
        ElMessage.error(t('message.passwordsNotMatch'))
        return false
    }
    
    return true
}

const register = async () => {
    // 验证输入
    if (!validateInputs()) {
        return
    }
    
    loading.value = true
    
    try {
        const response = await http.post(API.backend_url + '/api/auth/register', {
            username: username.value.trim(),
            email: email.value.trim(),
            password: password.value
        })
        
        // HTTP客户端返回的数据在response.data中
        const data = response.data
        
        if (data.success) {
            // 注册成功，自动登录
            globalSettingStore.isLogin = true
            globalSettingStore.userToken = data.token
            ElMessage.success(data.message || t('message.registerSuccess'))
        } else {
            // 注册失败，显示后端返回的具体错误信息
            ElMessage.error(data.message || t('message.registerFailed'))
        }
    } catch (err) {
        console.error('Register error:', err)
        
        // 根据错误类型显示不同的错误信息
        if (err.message.includes('400')) {
            // 可能是用户已存在或其他客户端错误
            ElMessage.error(t('message.userAlreadyExists'))
        } else if (err.message.includes('500')) {
            ElMessage.error(t('message.serverError'))
        } else if (err.name === 'TypeError' || err.message.includes('Failed to fetch')) {
            ElMessage.error(t('message.networkError'))
        } else {
            ElMessage.error(t('message.registerFailed'))
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
.input-item-label {
    min-width: 100px;
    text-align: right;
    margin-right: 10px;
}
.input-item-input {
    flex: 1;
}
</style>