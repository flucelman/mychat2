<template>
    <div class="drawer-body">
        <div class="drawer-body-item" v-for="item in chatConfigStore.chatHistory" :key="item.id"
            @click="chatConfigStore.getChatMessage(item.chat_id)"
            :class="{ 'active': item.chat_id == chatConfigStore.chatId }">
            <div>{{ item.title }}</div>
            <ThreePointsIcon class="three-points-icon" />
        </div>

    </div>
</template>

<script setup>
import ThreePointsIcon from '@/assets/icons/三个点.svg'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
</script>

<style scoped>
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
}

.drawer-body-item.active {
    background-color: var(--tertiary-background);
}

.drawer-body-item:hover {
    background-color: var(--tertiary-background);

    .three-points-icon {
        display: block;
    }
}

.three-points-icon {
    width: 20px;
    height: 20px;
    color: var(--text-color);
    cursor: pointer;
    display: none;
}
</style>