<template>
    <div class="drawer-container">
        <Header1 class="drawer-header"/>
        <Body1 class="drawer-body"/>
        <Footer1 class="drawer-footer"/>
    </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
import Header1 from './components/header.vue'
import Body1 from './components/body.vue'
import Footer1 from './components/footer.vue'
const chatConfigStore = useChatConfigStore()
onMounted(() => {
    chatConfigStore.getChatHistory()
})
</script>

<style scoped>
.drawer-container {
    width: 300px;
    height: 100vh;
    background-color: var(--secondary-background);
    border-right: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
}

.drawer-header {
    min-height: 80px;
}

.drawer-body {
    flex: 1;
    overflow-y: auto;
}

.drawer-footer {
    min-height: 80px;
}


</style>