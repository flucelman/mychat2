<template>
  <div class="user-container">
    <el-card class="user-card" shadow="hover">
      <div class="user-header">
        <el-avatar class="user-avatar" :size="48">
          {{ (globalSettingStore.userInfo.username || globalSettingStore.userInfo.email || 'U').slice(0, 1).toUpperCase() }}
        </el-avatar>
        <div class="user-name-email">
          <div class="name">{{ globalSettingStore.userInfo.username || '匿名用户' }}</div>
          <div class="email">{{ globalSettingStore.userInfo.email || '—' }}</div>
        </div>
      </div>

      <el-divider class="user-divider"/>

      <el-descriptions class="user-desc" :column="1" :border="false">
        <el-descriptions-item label="邮箱">
          {{ globalSettingStore.userInfo.email || '—' }}
        </el-descriptions-item>
        <el-descriptions-item label="用户名">
          {{ globalSettingStore.userInfo.username || '—' }}
        </el-descriptions-item>
        <el-descriptions-item label="注册日期">
          {{ formatDate(globalSettingStore.userInfo.created_at) || '—' }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="actions">
        <el-button type="danger" plain round @click="globalSettingStore.logout()">
          退出登录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { useGlobalSettingStore } from '@/stores/global_setting'
const globalSettingStore = useGlobalSettingStore()


const formatDate = (value) => {
  if (!value) return ''
  const s = String(value)
  if (s.includes('T')) return s.split('T')[0]
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return s
  return d.toISOString().split('T')[0]
}

</script>

<style scoped>
.user-container {
  padding: 12px;
}

.user-card {
  border-radius: 12px;
  border: 1px solid #e0e0e0;
}

:deep(.el-card__body) {
  padding: 16px 16px 12px 16px;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  background-color: #007bff;
  color: #fff;
  font-weight: 700;
}

.user-name-email {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name-email .name {
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.user-name-email .email {
  font-size: 12px;
  color: #666;
}

.user-divider {
  margin: 12px 0;
}

.user-desc {
  margin-top: 4px;
}

.user-desc :deep(.el-descriptions__label) {
  width: 88px;
  color: #666;
}

.user-desc :deep(.el-descriptions__content) {
  color: #333;
}

.actions {
  display: flex;
  justify-content: center;
  margin-top: 8px;
}

</style>