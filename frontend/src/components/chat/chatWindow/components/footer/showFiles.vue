<template>
    <div class="show-files-container">
        <div class="close-icon" @click="chatConfigStore.deleteFile(props.info.file_id,props.info.file_url)">
            <closeIcon class="close-icon-svg"/>
        </div>
        <img class="show-files-img" :src="fileType === 'pdf' ? pdfIcon : fileType === 'ppt' ? pptIcon : fileType === 'excel' ? excelIcon : fileType === 'image' ? props.info.file_url : fileIcon" alt="file" />
        <div class="show-files-info">
            <span>{{ displayFileName }}</span>
            <div class="show-files-detail">
                <span>{{ props.info.suffix }}</span>
                <span>|</span>
                <span>{{ formattedFileSize }}</span>
            </div>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
import pdfIcon from '@/assets/icons/pdf.svg?url'
import pptIcon from '@/assets/icons/ppt.svg?url'
import excelIcon from '@/assets/icons/excel.svg?url'
import fileIcon from '@/assets/icons/file.svg?url'
import closeIcon from '@/assets/icons/叉.svg'

const props = defineProps({
    info: {
        type: Object,
        required: true
    }
})

const chatConfigStore = useChatConfigStore()


const fileType = computed(() => {
    const suffix = props.info.suffix.toLowerCase()
    if (suffix === '.jpg' || suffix === '.png' || suffix === '.jpeg' || suffix === '.gif' || suffix === '.bmp' || suffix === '.webp') {
        return 'image'
    } else if (suffix === '.pdf') {
        return 'pdf'
    } else if (suffix === '.pptx' || suffix === '.ppt') {
        return 'ppt'
    } else if (suffix === '.xlsx' || suffix === '.xls') {
        return 'excel'
    } else {
        console.log('后缀',suffix)
        return 'file'
    }
})

// 限制文件名最多显示10个字符
const displayFileName = computed(() => {
    const name = props.info.name
    if (name.length > 10) {
        return name.substring(0, 10) + '...'
    }
    return name
})

// 格式化文件大小
const formattedFileSize = computed(() => {
    const size = props.info.size
    if (size < 1024) {
        return size + 'B'
    } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(1) + 'KB'
    } else if (size < 1024 * 1024 * 1024) {
        return (size / (1024 * 1024)).toFixed(1) + 'MB'
    } else {
        return (size / (1024 * 1024 * 1024)).toFixed(1) + 'GB'
    }
})
</script>

<style scoped>
.show-files-container {
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

.show-files-img {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 4px;
}
.show-files-info {
    display: flex;
    flex-direction: column;
    gap: 5px;
}
.show-files-detail {
    display: flex;
    gap: 5px;
    justify-content: flex-start;
    align-items: center;
    font-size: 12px;
    color: var(--secondary-text);

}
.close-icon {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 18px;
    height: 18px;
    cursor: pointer;
    background-color: var(--tertiary-background);
    border-radius: 50%;
    position: absolute;
    right: -8px;
    top: -8px;
    display: none;
}
.close-icon-svg {
    width: 10px;
    height: 10px;
    color: var(--secondary-text);
}
.show-files-container:hover .close-icon {
    display: flex;
}
</style>
