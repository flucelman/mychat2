<template>
    <textarea 
        ref="textareaRef"
        v-model="chatConfigStore.userMessage" 
        @input="adjustHeight"
        class="input-textarea" 
        :placeholder="$t('message.input_placeholder')"
        rows="1"
    ></textarea>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
const textareaRef = ref(null)

// 自动调整textarea高度
const adjustHeight = async () => {
    await nextTick()
    if (textareaRef.value) {
        // 重置高度以获取正确的scrollHeight
        textareaRef.value.style.height = 'auto'
        // 设置新高度
        const scrollHeight = textareaRef.value.scrollHeight
        const maxHeight = parseFloat(getComputedStyle(textareaRef.value).maxHeight)
        
        if (scrollHeight < maxHeight) {
            textareaRef.value.style.height = scrollHeight + 'px'
        } else {
            textareaRef.value.style.height = maxHeight + 'px'
        }
    }
}

// 组件挂载时初始化高度
onMounted(() => {
    adjustHeight()
})
</script>

<style scoped>
.input-textarea{
    width: 100%;
    min-height: 40px; /* 最小高度 */
    max-height: 120px; /* 最大约5行高度 */
    height: auto;
    background-color: var(--background-color);
    border-radius: 10px;
    border: 0px solid var(--border-color);
    padding: 10px;
    align-self: flex-start;
    outline: none;
    resize: none; /* 禁用手动调整，使用自动调整 */
    font-size: 16px;
    font-family: inherit;
    color: var(--text-color);
    line-height: 1.5;
    overflow-y: auto;
    transition: height 0.1s ease;
    box-sizing: border-box;
}
</style>