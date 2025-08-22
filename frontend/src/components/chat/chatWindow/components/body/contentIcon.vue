<template>
    <div class="content-icon-container"
        :class="{ 'user-icons': message.role == 'user', 'assistant-icons': message.role == 'assistant' }">
        <deleteIcon class="content-icon" @click="deleteMessage" />
        <refreshIcon class="content-icon" v-if="message.role == 'assistant'" />
        <el-dropdown @command="handleCopyCommand" trigger="click" :teleported="false">
            <copyIcon class="content-icon" />
            <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item :command="{ type: 'plain', content: message.content }">
                        {{ $t('message.copyPlainText') }}
                    </el-dropdown-item>
                    <el-dropdown-item :command="{ type: 'markdown', content: message.content }">
                        {{ $t('message.copyMarkdown') }}
                    </el-dropdown-item>
                </el-dropdown-menu>
            </template>
        </el-dropdown>
    </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import copyIcon from '@/assets/icons/复制.svg'
import deleteIcon from '@/assets/icons/删除.svg'
import refreshIcon from '@/assets/icons/刷新.svg'
import { http } from '@/utils/http/client'
import { API } from '@/router/api'
import { useChatConfigStore } from '@/stores/chat_config'

const { t: $t } = useI18n()

const chatConfigStore = useChatConfigStore()

// 定义组件属性
const props = defineProps({
    message: {
        type: Object,
        required: true
    }
})

// 处理复制命令
const handleCopyCommand = (command) => {
    const { type, content } = command
    let textToCopy = content

    if (type === 'plain') {
        // 复制原始文本（去除markdown格式）
        textToCopy = convertMarkdownToPlainText(content)
    } else if (type === 'markdown') {
        // 复制markdown文本（保持原格式）
        textToCopy = content
    }

    copyToClipboard(textToCopy)
}

// 复制到剪贴板
const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text).then(() => {
        ElMessage.success($t('message.copySuccess'))
    }).catch(() => {
        ElMessage.error($t('message.copyFailed'))
    })
}

// 将markdown转换为纯文本
const convertMarkdownToPlainText = (markdown) => {
    return markdown
        // 移除代码块
        .replace(/```[\s\S]*?```/g, '')
        .replace(/`([^`]+)`/g, '$1')
        // 移除标题标记
        .replace(/^#{1,6}\s+/gm, '')
        // 移除粗体和斜体
        .replace(/\*\*([^*]+)\*\*/g, '$1')
        .replace(/\*([^*]+)\*/g, '$1')
        .replace(/__([^_]+)__/g, '$1')
        .replace(/_([^_]+)_/g, '$1')
        // 移除链接，保留链接文本
        .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
        // 移除图片
        .replace(/!\[([^\]]*)\]\([^)]+\)/g, '$1')
        // 移除列表标记
        .replace(/^[\s]*[-*+]\s+/gm, '')
        .replace(/^[\s]*\d+\.\s+/gm, '')
        // 移除引用标记
        .replace(/^>\s+/gm, '')
        // 移除水平线
        .replace(/^---+$/gm, '')
        // 清理多余的空行
        .replace(/\n\s*\n/g, '\n')
        .trim()
}

// 删除消息
const deleteMessage = () => {
    console.log('删除消息', props.message.message_id)
    ElMessageBox.confirm($t('message.deleteMessage'), $t('message.confirmTitle'), {
        confirmButtonText: $t('message.confirm'),
        cancelButtonText: $t('message.cancel'),
        type: 'warning',
    }).then(async () => {
        try {
            // 删除chatconfigStore中的消息
            chatConfigStore.deleteMessage(props.message.message_id)
            const res = await http.post(API.backend_url + '/api/chat/deleteSingleMessage', {
                message_id: props.message.message_id
            })
            if (res.data.success) {
                ElMessage.success($t('message.deleteSuccess'))
            } else {
                ElMessage.error($t('message.deleteFailed'))
            }
        } catch (err) {
            ElMessage.error($t('message.deleteFailed'))
        }
    })
}
</script>

<style scoped>
.content-icon-container {
    display: flex;
    gap: 10px;
    justify-content: flex-start;
    margin-top: 10px;
}

.content-icon-container.user-icons {
    justify-content: flex-end;
}

.content-icon-container.assistant-icons {
    opacity: 1;
}

.content-icon {
    width: 18px;
    height: 18px;
    color: var(--secondary-text);
    cursor: pointer;
}

.content-icon:hover {
    color: var(--text-color);
}

/* 复制下拉菜单样式 */
:deep(.el-dropdown-menu) {
    background-color: var(--secondary-background) !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
}

:deep(.el-dropdown-menu__item) {
    color: var(--text-color) !important;
    font-size: 14px !important;
    padding: 8px 16px !important;
}

:deep(.el-dropdown-menu__item:hover) {
    background-color: var(--tertiary-background) !important;
    color: var(--text-color) !important;
}
</style>
