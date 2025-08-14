<template>
    <div>
      <!-- 隐藏的文件输入框 -->
      <input 
        ref="fileInput" 
        type="file" 
        style="display: none" 
        @change="handleFileSelect"
        multiple
      />
      
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  
  const fileInput = ref(null)
  const selectedFiles = ref([])
  
  // 暴露内容给父组件调用
  const emit = defineEmits(['files-selected'])
  
  // 打开文件管理器
  const openFileManager = () => {
    fileInput.value.click()
  }
  
  // 处理文件选择
  const handleFileSelect = (event) => {
    const files = Array.from(event.target.files)
    selectedFiles.value = files
    emit('files-selected', files)
    // 关键：重置 input 的值，确保再次选择相同文件也能触发 change
    // event.target.value = ''
    // if (fileInput.value) fileInput.value.value = ''
  }


  // 暴露方法给父组件调用
  defineExpose({ openFileManager })
  </script>