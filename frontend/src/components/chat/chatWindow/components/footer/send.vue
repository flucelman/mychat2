<template>
    <div class="send-container" @click="handleMessage">
        <sendIcon class="send-icon"/>
    </div>
</template>

<script setup>
import sendIcon from '@/assets/icons/上.svg'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
const handleMessage = () => {
    chatConfigStore.sendUserMessage()
}
</script>

<style scoped>
.send-container{
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
.send-container:hover{
    transform: scale(1.1);
}
.send-icon{
    width: 25px;
    height: 25px;
    color: #696969;
}
</style>