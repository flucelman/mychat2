<template>
    <div class="login-container">
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.email') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="email" placeholder="请输入邮箱" />
                </div>
            </div>
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.password') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="password" placeholder="请输入密码" />
                </div>
            </div>
            <div class="input-item">
                <el-button type="primary" @click="login">{{ $t('message.login') }}</el-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useGlobalSettingStore } from '@/stores/global_setting'
import { http } from '@/utils/http/client';  // 使用新的 HTTP 客户端
import { API } from '@/config/api';

const { t } = useI18n()
const globalSettingStore = useGlobalSettingStore()
const email = ref('')
const password = ref('')

const login = async () => {
    if (email.value === '' || password.value === '') {
        ElMessage.error(t('message.loginFailed'))
        return
    }
    
    try {
        const response = await http.post(API.backend_url + '/api/auth/login', {
            email: email.value, 
            password: password.value
        })
        
        // 解析JSON响应
        const data = await response.json()
        
        if (response.ok && data.success) {
            // 根据实际响应结构：token直接在data对象中
            globalSettingStore.isLogin = true
            globalSettingStore.userToken = data.token  // 直接访问data.token
            ElMessage.success(t('message.loginSuccess'))
        } else {
            ElMessage.error(t('message.loginFailed'))
        }
    } catch (err) {
        ElMessage.error(err)
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