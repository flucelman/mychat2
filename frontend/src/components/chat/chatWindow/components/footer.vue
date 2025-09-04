<template>
    <div class="input-container" id="input-container-id" v-show="showInputBox">
        <DropdownIcon class="dropdown-icon" @click="showInputBox=false" />
        <div class="files-container">
            <div v-if="chatConfigStore.filesInfo.length > 0" class="uploading-files-container">
                <div v-for="info in chatConfigStore.filesInfo" :key="info.file_id">
                    <showFiles :fileInfo="info" />
                </div>
            </div>
            <div v-if="chatConfigStore.isUploading" class="uploading-files-container">
                <div v-for="i in chatConfigStore.uploadingFiles" :key="i">
                    <showFilesPlaceholder />
                </div>
            </div>
        </div>
        <Textarea1 ref="textareaRef" @focus="handleFocus" @blur="handleBlur"/>
        <div class="input-container-inner">
            <Plugin />
            <div class="input-container-inner-right">
                <UploadFiles :key="chatConfigStore.chatId" />
                <Send :onAdjustHeight="handleAdjustHeight" />
            </div>
        </div>
    </div>
    <UpIcon class="up-icon" @click="showInputBox=true" v-if="!showInputBox" />
</template>

<script setup>
import { ref } from 'vue'
import Plugin from './footer/plugin.vue'
import Send from './footer/send.vue'
import UploadFiles from './footer/uploadFiles.vue'
import showFilesPlaceholder from './footer/showFilesPlaceholder.vue'
import showFiles from './footer/showFiles.vue'
import Textarea1 from './footer/textarea.vue'
import { useChatConfigStore } from '@/stores/chat_config'
import { useGlobalSettingStore } from '@/stores/global_setting'
import DropdownIcon from '@/assets/icons/下拉.svg'
import UpIcon from '@/assets/icons/上拉.svg'

const globalSettingStore = useGlobalSettingStore()
const chatConfigStore = useChatConfigStore()
const textareaRef = ref(null)
const showInputBox = ref(false)

// 处理调整高度的方法
const handleAdjustHeight = () => {
    if (textareaRef.value && textareaRef.value.adjustHeight) {
        textareaRef.value.adjustHeight()
    }
}

const handleFocus = () => {
    if (!globalSettingStore.isMobile) {
        return
    }
    const inputContainer = document.getElementById("input-container-id");
    inputContainer.style.position = "absolute";
    inputContainer.style.bottom = "50%";
    inputContainer.style.transform = "translateY(50%)";
    inputContainer.style.zIndex = "1000";
    inputContainer.style.backgroundColor = "var(--background-color)";
}

const handleBlur = () => {
    // 添加短暂延迟，确保点击事件能够正常触发
    if (!globalSettingStore.isMobile) {
        return
    }
    setTimeout(() => {
        const inputContainer = document.getElementById("input-container-id");
        inputContainer.style.position = "";
        inputContainer.style.bottom = "";
        inputContainer.style.transform = "";
        inputContainer.style.zIndex = "";
        inputContainer.style.backgroundColor = "";
    }, 100);
}

</script>

<style scoped>
@media screen and (max-width: 1200px) {
    .input-container {
        width: 100%;
    }
}

@media screen and (min-width: 1200px) {
    .input-container {
        width: 50%;
        max-width: 800px;
    }
}

.dropdown-icon {
        width: 30px;
        height: 30px;
        color: var(--icon-color-light);
        position: absolute;
        top: 6px;
        right: 6px;
        cursor: pointer;
}

.up-icon {
    width: 30px;
    height: 30px;
    color: var(--icon-color-light);
    cursor: pointer;
}

.input-container {
    display: flex;
    justify-content: center;
    align-items: flex-start;
    gap: 8px;
    flex-direction: column;
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 2px 16px 16px 16px;
    box-shadow: var(--shadow-color);
    transition: box-shadow 0.3s ease, transform 0.2s ease;
    position: relative;
}

.input-container:hover {
    box-shadow: var(--shadow-color);
    transform: translateY(-2px);
}

.input-container:focus-within {
    box-shadow: var(--shadow-color);
    border-color: var(--border-color-hover);
    border-width: 2px;
}



.input::placeholder {
    color: var(#999);
    font-size: 16px;
    opacity: 0.7;
}

.input-container-inner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
}

.input-container-inner-right {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}

.files-container {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 10px;
}

.uploading-files-container {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 10px;
}
</style>