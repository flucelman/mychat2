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
  }


  // 暴露方法给父组件调用
  defineExpose({ openFileManager })
  </script>