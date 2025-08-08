<template>
    <div class="drawer-body">
        <div class="drawer-body-item" v-for="item in chatConfigStore.chatHistory" :key="item.chat_id"
            @click="getChatHistory(item.chat_id)"
            :class="{ 'active': item.chat_id == chatConfigStore.chatId }">
            <div>{{ item.title }}</div>
            <el-popover placement="right" :width="180">
                <template #reference>
                    <ThreePointsIcon class="three-points-icon" :chat_id="item.chat_id" />
                </template>
                <div class="delete-main" @click="chatConfigStore.deleteSingleHistory(item.chat_id)">
                    <DeleteIcon class="delete-icon" />
                    {{ $t('message.delete') }}
                </div>
            </el-popover>
        </div>

    </div>
</template>

<script setup>
import ThreePointsIcon from '@/assets/icons/三个点.svg'
import DeleteIcon from '@/assets/icons/删除.svg'
import { useChatConfigStore } from '@/stores/chat_config'
import { useGlobalSettingStore } from '@/stores/global_setting'
const chatConfigStore = useChatConfigStore()
const globalSettingStore = useGlobalSettingStore()
const getChatHistory = async (chat_id) => {
    chatConfigStore.getChatMessage(chat_id);
    if (globalSettingStore.isMobile) {
        chatConfigStore.showDrawer = false;
    }
}
</script>

<style scoped lang="scss">
.drawer-body {
    overflow-y: auto;
    padding: 10px;
}

.drawer-body-item {
    padding: 10px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    border-radius: 10px;
    position: relative;
}

.drawer-body-item.active {
    background-color: var(--tertiary-background);
}

.drawer-body-item:hover {
    background-color: var(--tertiary-background);

    .three-points-icon {
        visibility: visible;
    }
}

/* 当弹出框激活时，保持图标可见 */
.drawer-body-item:has(.el-popover__reference:hover) .three-points-icon,
.drawer-body-item .three-points-icon:hover {
    visibility: visible;
}

.three-points-icon {
    color: var(--text-color);
    cursor: pointer;
    visibility: hidden;
    width: 20px;
    height: 20px;
}
.delete-main {
    display: flex;
    justify-content: start;
    align-items: center;
    width: 100%;
    height: 40px;
    border-radius: 10px;
    cursor: pointer;
 
    font-size: 14px;
    font-weight: 500;
    transition: all 0.3s ease;
    padding: 10px;
}

.delete-main:hover {
    background-color: #f0f0f0;
}

.delete-icon {
    width: 20px;
    height: 20px;
}
</style>