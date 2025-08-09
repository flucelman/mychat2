<template>
    <div class="body-container">
        <div class="message-container" v-for="message in chatConfigStore.baseMessageHistory" :key="message.id">
            <div class="message-item" :class="message.role">
                <div class="message-model user" v-if="message.role == 'user'">{{ globalSettingStore.userInfo.username }}</div>
                <div class="message-model assistant" v-if="message.role == 'assistant'">
                    <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == message.model)?.logo" class="model-icon" />
                    <div class="model-name-text">
                        {{ message.model }}
                    </div>
                </div>
                <div class="message-content">{{ message.content }}</div>
            </div>
        </div>
        <div v-if="chatConfigStore.isReceiving == true" class="message-item assistant">
            <div class="message-model assistant" >
                <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == chatConfigStore.AIConfig.model)?.logo" class="model-icon" />
                <div class="model-name-text">
                    {{ chatConfigStore.AIConfig.model }}
                </div>
            </div>
            <div class="message-content">{{ chatConfigStore.instantAssistantMessage }}</div>
        </div>
    </div>
</template>
    
<script setup>
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
import { useGlobalSettingStore } from '@/stores/global_setting'
const globalSettingStore = useGlobalSettingStore()
import { API } from '@/router/api'
</script>

<style scoped>
.model-icon {
    width: 25px;
    height: 25px;
}
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
.message-model {
    font-size: 14px;
    color: var(--secondary-text);
    padding-left: 10px;
    display: flex;
    align-items: center;
    gap: 5px;
}
.message-model.assistant {
    color: var(--secondary-text);
}
.message-content {
    padding: 10px;
    border-radius: 10px;
}
</style>