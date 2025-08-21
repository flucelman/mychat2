<template>
    <div class="input-container">
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
        <Textarea1 />
        <div class="input-container-inner">
            <Plugin />
            <div class="input-container-inner-right">
                <UploadFiles :key="chatConfigStore.chatId" />
                <Send />
            </div>
        </div>
    </div>
</template>

<script setup>
import Plugin from './footer/plugin.vue'
import Send from './footer/send.vue'
import UploadFiles from './footer/uploadFiles.vue'
import showFilesPlaceholder from './footer/showFilesPlaceholder.vue'
import showFiles from './footer/showFiles.vue'
import Textarea1 from './footer/textarea.vue'
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()

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