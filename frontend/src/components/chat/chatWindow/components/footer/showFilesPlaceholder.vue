<template>
    <div class="show-files-container">
        <div class="upload-content">
            <div class="upload-icon">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M14.83 9.17L12 6.34L9.17 9.17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M12 6.34V17" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <path d="M20 15V19C20 19.5304 19.7893 20.0391 19.4142 20.4142C19.0391 20.7893 18.5304 21 18 21H6C5.46957 21 4.96086 20.7893 4.58579 20.4142C4.21071 20.0391 4 19.5304 4 19V15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
            </div>
            <div class="upload-text">{{ $t('message.uploading') }}</div>
        </div>
        
        <!-- 进度条 -->
        <div class="progress-container">
            <div class="progress-bar">
                <div class="progress-fill" :style="{ width: progress + '%' }"></div>
                <div class="progress-shimmer"></div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()
const progress = ref(0)

// 模拟进度条动画
onMounted(() => {
    const interval = setInterval(() => {
        progress.value += 2
        if (progress.value >= 100) {
            progress.value = 0
        }
    }, 100)
    
    // 清理定时器（可选，如果需要停止动画）
    // setTimeout(() => clearInterval(interval), 10000)
})
</script>

<style scoped>
.show-files-container {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    width: 120px;
    height: 50px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    padding: 8px;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    position: relative;
    overflow: hidden;
}

.show-files-container:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.show-files-container::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    animation: shimmer 2s infinite;
}

@keyframes shimmer {
    0% { left: -100%; }
    100% { left: 100%; }
}

.upload-content {
    display: flex;
    align-items: center;
    gap: 6px;
    color: white;
    flex: 1;
}

.upload-icon {
    width: 16px;
    height: 16px;
    opacity: 0.9;
}

.upload-icon svg {
    width: 100%;
    height: 100%;
}

.upload-text {
    font-size: 11px;
    font-weight: 500;
    letter-spacing: 0.3px;
}

.progress-container {
    width: 100%;
    height: 3px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 2px;
    overflow: hidden;
    position: relative;
}

.progress-bar {
    width: 100%;
    height: 100%;
    position: relative;
    border-radius: 2px;
    overflow: hidden;
}

.progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #ffffff, #f0f0f0);
    border-radius: 2px;
    transition: width 0.3s ease;
    position: relative;
}

.progress-shimmer {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, 
        transparent 0%, 
        rgba(255, 255, 255, 0.6) 50%, 
        transparent 100%);
    animation: progressShimmer 1.5s infinite;
}

@keyframes progressShimmer {
    0% { transform: translateX(-100%); }
    100% { transform: translateX(200%); }
}

/* 响应式调整 */
@media (max-width: 768px) {
    .upload-text {
        font-size: 10px;
    }
    
    .upload-icon {
        width: 14px;
        height: 14px;
    }
}

/* 暗色主题变体（可选） */
.show-files-container.dark {
    background: linear-gradient(135deg, #2d3748 0%, #4a5568 100%);
    box-shadow: 0 4px 15px rgba(45, 55, 72, 0.3);
}

.show-files-container.dark:hover {
    box-shadow: 0 8px 25px rgba(45, 55, 72, 0.4);
}
</style>