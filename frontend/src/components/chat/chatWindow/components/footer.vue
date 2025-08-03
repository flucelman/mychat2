<template>
        <div class="input-container">
            <textarea 
                ref="textareaRef"
                v-model="inputMessage" 
                @input="adjustHeight"
                class="input-textarea" 
                :placeholder="$t('message.input_placeholder')"
                rows="1"
            ></textarea>
            <div class="input-container-inner">
                <Plugin/>
                <div class="input-container-inner-right">
                    <UploadFiles/>
                    <Send/>
                </div>
            </div>
        </div>
</template>

<script setup>
import { ref, nextTick, onMounted } from 'vue'
import Plugin from './footer/plugin.vue'
import Send from './footer/send.vue'
import UploadFiles from './footer/uploadFiles.vue'

const inputMessage = ref('')
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
@media screen and (max-width: 1200px) {
    .input-container{
        width: 100%;
    }
}
@media screen and (min-width: 1200px) {
    .input-container{
        width: 50%;
        max-width: 800px;
    }
}

.input-container{
    display: flex;
    justify-content: center;
    align-items: flex-start;
    gap: 10px;
    flex-direction: column;
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 6px 16px 16px 16px;
    box-shadow: var(--shadow-color);
    transition: box-shadow 0.3s ease, transform 0.2s ease;
}

.input-container:hover{
    box-shadow: var(--shadow-color);
    transform: translateY(-2px);
}

.input-container:focus-within{
    box-shadow: var(--shadow-color);
    border-color: var(--border-color-hover);
    border-width: 2px;
}

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

.input::placeholder{
    color: var(#999);
    font-size: 16px;
    opacity: 0.7;
}
.input-container-inner{
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
}
.input-container-inner-right{
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}
</style>