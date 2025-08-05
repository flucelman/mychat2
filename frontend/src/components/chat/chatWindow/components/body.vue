<template>
    <div class="body-container">
        <div class="message-container" v-for="message in chatConfigStore.baseMessageHistory" :key="message.id">
            <div class="message-item" :class="message.role">
                <div class="message-content">{{ message.content }}</div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
</script>

<style scoped>
.body-container {
    flex: 1;
    width: 100%;
    height: 100%;
    overflow-y: auto; /* 允许垂直滚动 */
    overflow-x: hidden; /* 隐藏水平滚动条 */
    padding: 10px;
    box-sizing: border-box;
}
.message-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
}
.message-item {
    display: flex;
    flex-direction: column;
}
.message-item.user {
    align-self: flex-end;
}
.message-item.assistant {
    align-self: flex-start;
}
.message-content {
    padding: 10px;
    border-radius: 10px;
}
</style>