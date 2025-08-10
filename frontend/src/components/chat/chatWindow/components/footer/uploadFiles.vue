<template>
    <div class="upload-files-container" @click="handleUploadClick">
        <uploadIcon/>
        <upload_files ref="uploader" @files-selected="onFilesSelected"/>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import uploadIcon from '@/assets/icons/上传.svg'
import upload_files from '@/components/reuse/upload_files.vue'
import { http } from '@/utils/http/client'
import {API} from '@/router/api'
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()

const uploader = ref(null)
const handleUploadClick = () => {
    if (uploader.value && typeof uploader.value.openFileManager === 'function') {
        uploader.value.openFileManager()
    }
}

const onFilesSelected = (files) => {
    if (!files || files.length === 0) return
    const formData = new FormData()
    // 后端使用 form.File["files"], 需要逐个追加文件
    files.forEach((file) => {
        formData.append('files', file, file.name)
    })
    chatConfigStore.isUploading = true
    chatConfigStore.uploadingFiles = files.length

    http.post(API.backend_url + "/api/oss/uploadFile", formData)
        .then(res => {
            chatConfigStore.isUploading = false
            chatConfigStore.uploadingFiles = 0
            chatConfigStore.fileUrl = res.data.file_urls
        })
        .catch(err => {
            console.error('Upload failed:', err)
        })
}
</script>

<style scoped>
.upload-files-container{
    display: flex;
    justify-content: center;
    align-items: center;
    width: 25px;
    height: 25px;
    cursor: pointer;
    transition: transform 0.2s ease;
}
.upload-files-container:hover{
    transform: scale(1.1);
}

</style>