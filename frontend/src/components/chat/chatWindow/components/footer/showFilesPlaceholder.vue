<template>
    <div class="show-files-container">
        这是占位符
    </div>
</template>

<script setup>
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()

</script>

<style scoped>
.show-files-container{
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: 50px;
    height: 50px;
    background-color: #e23737;
}
</style>