<template>
    <div class="drawer-container">
        <Header class="drawer-header"/>
        <Body class="drawer-body"/>
        <Footer class="drawer-footer"/>
    </div>
</template>

<script setup>
import { useChatConfigStore } from '@/stores/chat_config'
import Header from './components/header.vue'
import Body from './components/body.vue'
import Footer from './components/footer.vue'
const chatConfigStore = useChatConfigStore()

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