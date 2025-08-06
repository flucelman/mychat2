<template>
    <div class="send-container" 
        :class="{ 'container-has-user-message': chatConfigStore.userMessage.trim().length > 0 || chatConfigStore.isReceiving == true }">
        <upIcon v-if="chatConfigStore.isReceiving == false" class="send-icon" @click="chatConfigStore.sendUserMessage()"
            :class="{ 'send-icon-has-user-message': chatConfigStore.userMessage.trim().length > 0 }" />
        <squareIcon v-else style="width: 16px; height: 16px;" @click="chatConfigStore.cancelConnection()"
            :class="{ 'send-icon-has-user-message': chatConfigStore.userMessage.trim().length > 0 }" />
    </div>
</template>

<script setup>
import upIcon from '@/assets/icons/上.svg'
import squareIcon from '@/assets/icons/正方形.svg'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()

</script>

<style scoped>
.send-container {
    width: 32px;
    height: 32px;
    background-color: #cfcfcf;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    transition: transform 0.2s ease;
}

.send-container:hover {
    transform: scale(1.1);
}

.send-icon {
    width: 25px;
    height: 25px;
    color: #696969;
}

.container-has-user-message {
    background-color: #409EFF;
}

.send-icon-has-user-message {
    color: #ffffff;
}
</style>