<template>
    <div class="drawer-header">
        <div class="drawer-header-title">{{ $t('message.chatHistory') }}</div>
        <DeleteIcon class="delete-icon" @click="deleteAllHistory" />
    </div>
</template>

<script setup>
import DeleteIcon from '@/assets/icons/删除.svg'
import { useChatConfigStore } from '@/stores/chat_config'
import { ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const chatConfigStore = useChatConfigStore()
const deleteAllHistory = () => {
    ElMessageBox.confirm(t('message.deleteAllHistory'), t('message.confirmTitle'), {
        confirmButtonText: t('message.confirm'),
        cancelButtonText: t('message.cancel'),
        type: 'warning',
    }).then(() => {
        chatConfigStore.deleteAllHistory()
    })
}
</script>

<style scoped>
.drawer-header {
    padding: 20px;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.drawer-header-title {
    color: var(--text-color);
    font-size: 18px;
    font-weight: 600;
}

.delete-icon {
    width: 20px;
    height: 20px;
    color: var(--text-color);
    cursor: pointer;
}
</style>