<template>
    <div class="show-files-container">
        <img class="show-files-img" :src="fileType === 'pdf' ? pdfIcon : fileType === 'ppt' ? pptIcon : fileType === 'excel' ? excelIcon : fileType === 'image' ? props.url : fileIcon" alt="file" />
        <div class="show-files-info">
            <span>{{ fileName }}</span>
            <span>{{ fileSuffix }}</span>
        </div>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import pdfIcon from '@/assets/icons/pdf.svg?url'
import pptIcon from '@/assets/icons/ppt.svg?url'
import excelIcon from '@/assets/icons/excel.svg?url'
import fileIcon from '@/assets/icons/file.svg?url'

const props = defineProps({
    url: {
        type: String,
        required: true
    }
})

const fileName = computed(() => {
    const cleaned = props.url.split('?')[0].split('#')[0]
    return cleaned.split('/').pop().split('.').shift()
})

const fileSuffix = computed(() => {
    const cleaned = props.url.split('?')[0].split('#')[0]
    return cleaned.split('.').pop()?.toLowerCase()
})

const fileType = computed(() => {
    const cleaned = props.url.split('?')[0].split('#')[0]
    const suffix = cleaned.split('.').pop()?.toLowerCase()
    if (suffix === 'jpg' || suffix === 'png' || suffix === 'jpeg' || suffix === 'gif' || suffix === 'bmp' || suffix === 'webp') {
        return 'image'
    } else if (suffix === 'pdf') {
        return 'pdf'
    } else if (suffix === 'pptx' || suffix === 'ppt') {
        return 'ppt'
    } else if (suffix === 'xlsx' || suffix === 'xls') {
        return 'excel'
    } else {
        console.log('后缀',suffix)
        return 'file'
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
    gap: 2px;
}
</style>
