<template>
    <div>
      <!-- 隐藏的文件输入框 -->
      <input 
        ref="fileInput" 
        type="file" 
        style="display: none" 
        @change="handleFileSelect"
        :accept="acceptString"
        multiple
      />
      
    </div>
  </template>
  
  <script setup>
  import { ref, computed } from 'vue'
  
  const props = defineProps({
    allowedExtensions: {
      type: Array,
      default: () => [
        '.jpg', '.jpeg', '.png', '.gif', '.webp',
        '.mp4', '.mov', '.avi', '.mkv',
        '.mp3', '.wav', '.aac', '.flac',
        '.pdf', '.doc', '.docx', '.xls', '.xlsx', '.ppt', '.pptx',
        '.txt', '.md', '.html'
      ]
    }
  })
  
  const emit = defineEmits(['files-selected', 'files-rejected'])
  
  const fileInput = ref(null)
  const selectedFiles = ref([])
  
  const acceptString = computed(() => props.allowedExtensions.join(','))
  
  // 打开文件管理器
  const openFileManager = () => {
    if (fileInput.value) fileInput.value.click()
  }
  
  // 处理文件选择
  const handleFileSelect = (event) => {
    const rawFiles = Array.from(event.target.files || [])
    const allowed = []
    const rejected = []

    const lowerAllowed = props.allowedExtensions.map(e => e.toLowerCase())

    for (const file of rawFiles) {
      const idx = file.name.lastIndexOf('.')
      const ext = idx !== -1 ? file.name.slice(idx).toLowerCase() : ''
      if (lowerAllowed.includes(ext)) {
        allowed.push(file)
      } else {
        rejected.push(file.name)
      }
    }

    selectedFiles.value = allowed
    emit('files-selected', allowed)
    if (rejected.length) {
      emit('files-rejected', rejected)
      // 简单提示（父组件可通过事件自定义更友好的提示）
      window.alert('以下文件类型不被允许：\n' + rejected.join('\n'))
    }

    // 关键：重置 input 的值，确保再次选择相同文件也能触发 change
    event.target.value = ''
    if (fileInput.value) fileInput.value.value = ''
  }


  // 暴露方法给父组件调用
  defineExpose({ openFileManager })
  </script>