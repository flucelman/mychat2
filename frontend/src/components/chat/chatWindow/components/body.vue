<template>
    <div class="body-container">
        <div class="message-container" v-for="(group, groupIndex) in groupedMessages" :key="groupIndex">
            <!-- 非文件消息 -->
            <div v-if="group.type == 'message'" class="message-item" :class="group.message.role">
                <div class="message-model user username"
                    v-if="group.message.role == 'user' && (groupIndex === 0 || groupedMessages[groupIndex - 1].type !== 'files')">
                    {{ globalSettingStore.userInfo.username }}</div>
                <div class="message-model assistant"
                    v-if="group.message.role == 'assistant'">
                    <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == group.message.model)?.logo"
                        class="model-icon" />
                    <div class="model-name-text" v-if="group.message.role != 'search'">
                        {{ group.message.model }}
                    </div>
                </div>
                <div class="message-content">
                    <MarkdownView :content="group.message.content" />
                    <div class="content-icons-wrapper">
                        <ContentIcon :message="group.message" />
                    </div>
                </div>
            </div>
            <!-- 文件消息组 -->
            <div v-else-if="group.type == 'files'" class="message-item file">
                <div class="message-model user username">{{ globalSettingStore.userInfo.username }}</div>
                <div class="files-container">
                    <div class="file-item" v-for="fileMessage in group.files" :key="fileMessage.id">
                        <img class="file-icon" :src="getFileIcon(fileMessage)" alt="file" />
                        <div class="file-info">
                            <span class="file-name">{{ getDisplayFileName(fileMessage.file_name) }}</span>
                            <div class="file-detail">
                                <span>{{ getFileType(fileMessage.file_name) }}</span>
                                <span>|</span>
                                <span>{{ formatFileSize(fileMessage.file_size) }}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <!-- 联网搜索消息组 -->
            <div v-else-if="group.type == 'search'" class="message-item search">
                <div class="message-model assistant">
                    <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == group.search.model)?.logo"
                        class="model-icon" />
                    <div class="model-name-text">
                        {{ group.search.model }}
                    </div>
                </div>
                <div class="search-result" @click="openSearchDrawer(group.search.content)">{{ $t('message.searchResult') }}</div>
            </div>
        </div>
        <!-- 流式展示助手消息 -->
        <div v-if="chatConfigStore.isReceiving == true" class="message-item assistant">

            <div class="message-model assistant">
                <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == chatConfigStore.AIConfig.model)?.logo"
                    class="model-icon" />
                <div class="model-name-text">
                    {{ chatConfigStore.AIConfig.model }}
                </div>
            </div>
            <div class="message-content">
                <div v-if="chatConfigStore.instantAssistantMessage == ''" class="loading-container">
                    <Vue3Lottie 
                        :animationData="AILoadingAnimation" 
                        :height="60" 
                        :width="60" 
                        :loop="true" 
                        :autoplay="true"
                    />
                </div>
                <MarkdownView v-else :content="chatConfigStore.instantAssistantMessage" />
            </div>
        </div>
        
        <!-- 统一的搜索结果抽屉 -->
        <el-drawer v-model="searchDrawer" :title="$t('message.onlineSearch')" :direction="'rtl' " class="search-drawer">
            <search-card :search="currentSearchContent" />
        </el-drawer>
    </div>
</template>

<script setup>
import { computed, ref, nextTick } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
import { useGlobalSettingStore } from '@/stores/global_setting'
const globalSettingStore = useGlobalSettingStore()
import { API } from '@/router/api'
import pdfIcon from '@/assets/icons/pdf.svg?url'
import pptIcon from '@/assets/icons/ppt.svg?url'
import excelIcon from '@/assets/icons/excel.svg?url'
import fileIcon from '@/assets/icons/file.svg?url'
import videoIcon from '@/assets/icons/video.svg?url'
import SearchCard from './body/searchCard.vue'
import MarkdownView from './body/markdownView.vue'
import ContentIcon from './body/contentIcon.vue'
import { Vue3Lottie } from 'vue3-lottie'
import AILoadingAnimation from '@/assets/lottie/AI_loading.json'

const searchDrawer = ref(false)
const currentSearchContent = ref('')

// 打开搜索抽屉并设置当前搜索内容
const openSearchDrawer = (content) => {
    // 先设置内容，再打开抽屉，确保内容更新
    currentSearchContent.value = content
    // 使用nextTick确保DOM更新后再打开抽屉
    nextTick(() => {
        searchDrawer.value = true
    })
}



// 获取文件类型
const getFileType = (fileName) => {
    const suffix = fileName.split('.').pop()
    if (suffix === 'jpg' || suffix === 'png' || suffix === 'jpeg' || suffix === 'gif' || suffix === 'bmp' || suffix === 'webp') {
        return 'image'
    } else if (suffix === 'pdf') {
        return 'pdf'
    } else if (suffix === 'pptx' || suffix === 'ppt') {
        return 'ppt'
    } else if (suffix === 'xlsx' || suffix === 'xls') {
        return 'excel'
    } else if (suffix === 'mp4' || suffix === 'avi' || suffix === 'mov' || suffix === 'wmv' || suffix === 'flv' || suffix === 'mkv') {
        return 'video'
    } else {
        return 'file'
    }
}

// 获取文件图标
const getFileIcon = (fileMessage) => {
    const type = getFileType(fileMessage.file_name)
    if (type === 'pdf') return pdfIcon
    if (type === 'ppt') return pptIcon
    if (type === 'excel') return excelIcon
    if (type === 'image') return fileMessage.file_url
    if (type === 'video') return videoIcon
    return fileIcon
}

// 限制文件名最多显示10个字符
const getDisplayFileName = (fileName) => {
    if (fileName.length > 10) {
        return fileName.substring(0, 10) + '...'
    }
    return fileName
}

// 格式化文件大小
const formatFileSize = (size) => {
    if (size < 1024) {
        return size + 'B'
    } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(1) + 'KB'
    } else if (size < 1024 * 1024 * 1024) {
        return (size / (1024 * 1024)).toFixed(1) + 'MB'
    } else {
        return (size / (1024 * 1024 * 1024)).toFixed(1) + 'GB'
    }
}

// 将消息分组，连续的文件消息合并为一组
const groupedMessages = computed(() => {
    const groups = []
    let currentFileGroup = null

    chatConfigStore.baseMessageHistory.forEach(message => {
        if (message.role === 'file') {
            if (!currentFileGroup) {
                currentFileGroup = {
                    type: 'files',
                    files: []
                }
                groups.push(currentFileGroup)
            }
            currentFileGroup.files.push(message)
        } else if (message.role === 'search') {
            currentFileGroup = null
            groups.push({
                type: 'search',
                search: message
            })
        }else if (message.role === 'system') {
            currentFileGroup = null
            groups.push({
                type: 'system',
                system: message
            })
        } else {
            currentFileGroup = null
            groups.push({
                type: 'message',
                message: message
            })
        }
    })

    return groups
})
</script>

<style >

.model-icon {
    width: 25px;
    height: 25px;
}

.body-container {
    flex: 1;
    width: 100%;
    height: 100%;
    overflow-y: auto;
    /* 允许垂直滚动 */
    overflow-x: hidden;
    /* 隐藏水平滚动条 */
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

.message-item.search {
    align-self: flex-start;
    margin-bottom: 10px;
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
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    align-items: flex-start;
}

/* 用户消息的图标hover效果（通过包装器控制） */
.message-item.user .content-icons-wrapper {
    opacity: 0;
}

.message-item.user:hover .content-icons-wrapper {
    opacity: 1;
}



.username {
    align-self: flex-end;
}

.message-item.file {
    align-self: flex-end;
}

.files-container {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: flex-end;
    margin-top: 10px;
}

.file-item {
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: var(--secondary-background);
    border-radius: 8px;
    border: 1px solid var(--border-color);
    padding: 5px;
    font-size: 14px;
    gap: 5px;
    position: relative;
}

.file-icon {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 4px;
}

.file-info {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.file-name {
    font-size: 14px;
    color: var(--text-color);
}

.file-detail {
    display: flex;
    gap: 5px;
    justify-content: flex-start;
    align-items: center;
    font-size: 12px;
    color: var(--secondary-text);
}

.search-drawer {
    background-color: var(--secondary-background);
    width: 350px !important;
}
@media (max-width: 768px) {
    .search-drawer {
        width: 80vw !important;
    }
}
@media (min-width: 1024px) {
    .body-container {
        width: 70%;
        margin: 0 auto;
    }
}
.search-result {
    color: var(--secondary-text);
    font-size: 14px;
    cursor: pointer;
    margin-top: 10px;
    padding-left: 40px;
}

/* 美化 body-container 滚动条 */
.body-container::-webkit-scrollbar {
  width: 8px;
}

.body-container::-webkit-scrollbar-track {
  background: transparent;
  border-radius: 4px;
}

.body-container::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 4px;
  transition: background-color 0.3s ease;
}

.body-container::-webkit-scrollbar-thumb:hover {
  background: var(--secondary-text);
}

/* Firefox 滚动条样式 */
.body-container {
  scrollbar-width: thin;
  scrollbar-color: var(--border-color) transparent;
}

/* 隐藏滚动条但保留滚动（Chrome / Edge / Safari）*/
.search-drawer .el-drawer__body::-webkit-scrollbar {
  width: 0;
  height: 0;
}

/* 隐藏 Firefox 滚动条 */
.search-drawer .el-drawer__body {
  scrollbar-width: none;
  -ms-overflow-style: none; /* IE 10+ */
}

/* 加载动画容器样式 */
.loading-container {
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  padding: 10px 0;
}



</style>