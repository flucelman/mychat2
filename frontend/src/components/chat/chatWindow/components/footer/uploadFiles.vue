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
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const chatConfigStore = useChatConfigStore()
const { t } = useI18n()
const uploader = ref(null)
const handleUploadClick = () => {
    if (uploader.value && typeof uploader.value.openFileManager === 'function') {
        uploader.value.openFileManager()
    }
}

// 发送文件到后端
const onFilesSelected = async (files) => {
    if (!files || files.length === 0) return
    const formData = new FormData()
    formData.append('chat_id', chatConfigStore.chatId)
    // 后端使用 form.File["files"], 需要逐个追加文件
    files.forEach((file) => {
        formData.append('files', file, file.name)
    })
    chatConfigStore.isUploading = true
    chatConfigStore.uploadingFiles = files.length

    try {
        const res = await http.post(API.backend_url + "/api/oss/uploadFile", formData)
        // 兼容不同返回结构：优先取 res.data.data，其次取 res.data
        const payload = (res && res.data && (res.data.data ?? res.data)) || {}
        const failedList = Array.isArray(payload.failed) ? payload.failed : []
        const successList = Array.isArray(payload.success) ? payload.success : []

        // 失败提示（从对象数组中提取文件名）
        if (failedList.length > 0) {
            const failedNames = failedList.map(item => (item && item.name) ? item.name : String(item)).filter(Boolean)
            ElMessage.error(failedNames.join(', ') + ' ' + t('message.uploadFailed'))
        }

        // 成功结果写入 fileUrl（支持 string 或对象 {url|path}）
        if (successList.length > 0) {
            chatConfigStore.filesInfo.push(...successList)
            console.log(chatConfigStore.filesInfo)
            ElMessage.success(successList.length + ' ' + t('message.uploadSuccess'))
        }

        // 若全部失败
        if (failedList.length > 0 && successList.length === 0) {
            // 已在上面弹出失败详情，这里可再提示汇总
            // 保持一次错误提示即可，无需重复
        }
    } catch (err) {
        console.error('Upload failed:', err)
        ElMessage.error(t('message.uploadFailed'))
    } finally {
        chatConfigStore.isUploading = false
        chatConfigStore.uploadingFiles = 0
    }
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