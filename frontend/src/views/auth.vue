<template>
    <div class="auth-container">
            <div class="auth-title">
                <div class="auth-title-item" @click="activeTab = 'login'">{{ $t('message.login') }}</div>
                <div class="auth-title-item" @click="activeTab = 'register'">{{ $t('message.register') }}</div>
            </div>
        <div class="auth-content">
            <Login v-if="activeTab === 'login'" />
            <Register v-if="activeTab === 'register'" />
        </div>
    </div>
</template>

<script setup>
import Login from '@/components/auth/login.vue'
import Register from '@/components/auth/register.vue'
import { ref } from 'vue'

const activeTab = ref('login')

</script>

<style scoped>
.auth-container {
    display: flex;
    align-items: center;
    background-color: var(--el-bg-color);
    border-radius: 10px;
    box-shadow: 0 0 10px 0 rgba(0, 0, 0, 0.1);
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    flex-direction: column;
}
.auth-title{
    display: flex;
    justify-content: space-between;
    width: 100%;
    padding: 0 20px;
    font-size: 20px;
    font-weight: bold;
    color: var(--el-text-color-regular);
    margin-bottom: 20px;
}
.auth-title-item{
    cursor: pointer;
}
.auth-title-item:hover{
    color: var(--el-color-primary);
}
</style>