<template>
    <div class="register-container">
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.username') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="username" placeholder="请输入用户名" />
                </div>
            </div>
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.email') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="email" placeholder="请输入邮箱" />
                </div>
            </div>
        </div>
        <div class="input-container">
            <div class="input-item">
                <div class="input-item-label">
                    {{ $t('message.password') }}:
                </div>
                <div class="input-item-input">
                    <el-input v-model="password" placeholder="请输入密码" />
                </div>
            </div>
            <div class="input-item">
                <el-button type="primary" @click="register">{{ $t('message.register') }}</el-button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { http } from '@/utils/http/client';  // 使用新的 HTTP 客户端
import { API } from '@/config/api';
import { useGlobalSettingStore } from '@/stores/global_setting'
const globalSettingStore = useGlobalSettingStore()
const { t } = useI18n()
const username = ref('')
const email = ref('')
const password = ref('')

const register = () => {
    http.post(API.backend_url + '/api/auth/register', {
        username: username.value,
        email: email.value,
        password: password.value
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            globalSettingStore.isLogin = true
            globalSettingStore.userToken = data.token
            ElMessage.success(t('message.registerSuccess'))
        } else {
            ElMessage.error(data.message)
        }
    })
    .catch(error => {
        ElMessage.error(error)
    })
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
}
</style>