<template>
    <div class="selector-container">
        <div class="selector-item" :class="{ active: globalSettingStore.theme === 'light' }">
            <div class="item-content" @click="globalSettingStore.theme = 'light'">
                <span class="item-icon">☀️</span>
                <span class="item-text">{{ $t('message.light') }}</span>
            </div>
        </div>
        <div class="selector-item" :class="{ active: globalSettingStore.theme === 'dark' }">
            <div class="item-content" @click="globalSettingStore.theme = 'dark'">
                <span class="item-icon">🌙</span>
                <span class="item-text">{{ $t('message.dark') }}</span>
            </div>
        </div>
    </div>
</template>

<script setup>
import { useGlobalSettingStore } from '@/stores/global_setting'

const globalSettingStore = useGlobalSettingStore()
</script>

<style scoped>
.selector-container {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 12px;
    min-width: 140px;
    background: var(--background-color, #ffffff);
    border-radius: 16px;
    box-shadow: 
        0 8px 32px rgba(0, 0, 0, 0.08),
        0 4px 16px rgba(0, 0, 0, 0.04),
        0 0 0 1px rgba(0, 0, 0, 0.04);
    backdrop-filter: blur(16px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
}

.selector-container::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, 
        rgba(255, 255, 255, 0.8) 0%, 
        rgba(255, 255, 255, 0.4) 100%);
    pointer-events: none;
    opacity: 0.6;
    z-index: 0;
}

.selector-container:hover {
    transform: translateY(-2px);
    box-shadow: 
        0 12px 40px rgba(0, 0, 0, 0.12),
        0 6px 20px rgba(0, 0, 0, 0.08),
        0 0 0 1px rgba(0, 0, 0, 0.06);
}

.selector-item {
    position: relative;
    z-index: 1;
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    background: transparent;
}

.selector-item::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, 
        rgba(99, 102, 241, 0.1) 0%, 
        rgba(168, 85, 247, 0.1) 100%);
    opacity: 0;
    transition: all 0.25s ease;
    border-radius: 12px;
}

.selector-item:hover::before {
    opacity: 1;
}

.selector-item:hover {
    transform: scale(1.02);
    background: rgba(255, 255, 255, 0.5);
}

.selector-item.active {
    background: linear-gradient(135deg, 
        #6366f1 0%, 
        #a855f7 100%);
    box-shadow: 
        0 4px 16px rgba(99, 102, 241, 0.4),
        0 2px 8px rgba(168, 85, 247, 0.3);
    transform: scale(1.02);
}

.selector-item.active::before {
    display: none;
}

.item-content {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 12px 16px;
    position: relative;
    z-index: 1;
}

.item-icon {
    font-size: 16px;
    transition: all 0.25s ease;
}

.item-text {
    font-size: 14px;
    font-weight: 500;
    color: var(--text-color, #374151);
    transition: all 0.25s ease;
    letter-spacing: 0.3px;
}

.selector-item.active .item-text {
    color: white;
    font-weight: 600;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.selector-item.active .item-icon {
    filter: brightness(1.2);
    transform: scale(1.1);
}

.selector-item:active {
    transform: scale(0.98);
}

/* 响应式设计 */
@media (max-width: 768px) {
    .selector-container {
        min-width: 120px;
        padding: 10px;
        gap: 4px;
    }
    
    .item-content {
        padding: 10px 14px;
        gap: 6px;
    }
    
    .item-text {
        font-size: 13px;
    }
    
    .item-icon {
        font-size: 14px;
    }
}

/* 深色模式适配 */
@media (prefers-color-scheme: dark) {
    .selector-container {
        background: var(--background-color, #1f2937);
        box-shadow: 
            0 8px 32px rgba(0, 0, 0, 0.4),
            0 4px 16px rgba(0, 0, 0, 0.2),
            0 0 0 1px rgba(255, 255, 255, 0.1);
    }
    
    .selector-container::before {
        background: linear-gradient(135deg, 
            rgba(255, 255, 255, 0.1) 0%, 
            rgba(255, 255, 255, 0.05) 100%);
    }
    
    .selector-item:hover {
        background: rgba(255, 255, 255, 0.1);
    }
    
    .item-text {
        color: var(--text-color, #e5e7eb);
    }
}
</style>    