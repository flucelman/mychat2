<template>
    <div class="settings-container">
        <!-- 模型 -->
        <div class="settings-item">
            <div class="settings-item-title">{{ $t('message.model') }}</div>
            <div class="settings-item-content">
                <el-select v-model="chatConfigStore.AIConfig.model" placeholder="请选择模型" style="width: 100%;">
                    <el-option v-for="model in chatConfigStore.modelList" :label="model.name" :value="model.name">
                        <div class="model-option">
                            <span class="model-show">{{ model.name }}</span>
                            <div class="model-abilities">
                                <div class="ability-icons">
                                    <span 
                                        v-for="ability in parseAbilities(model.ability)" 
                                        :key="ability"
                                        class="ability-tag"
                                        :class="`ability-${ability}`"
                                    >
                                        <span class="ability-icon">{{ getAbilityIcon(ability) }}</span>
                                    </span>
                                </div>
                            </div>
                            <span class="model-show model-price">{{ model.price }}</span>
                        </div>
                    </el-option>
                </el-select>
            </div>
        </div>
        <!-- 模型温度 -->
        <div class="settings-item">
            <div class="settings-item-title">
                {{ $t('message.temperature') }} 
                <span class="setting-value">{{ chatConfigStore.AIConfig.temperature }}</span>
            </div>
            <div class="settings-item-content">
                <el-slider 
                    v-model="chatConfigStore.AIConfig.temperature" 
                    :min="0" 
                    :max="2" 
                    :step="0.1" 
                />
            </div>
        </div>
        <!-- 模型最大token -->
        <div class="settings-item">
            <div class="settings-item-title">
                {{ $t('message.maxTokens') }} 
                <span class="setting-value">{{ chatConfigStore.AIConfig.max_tokens }}</span>
            </div>
            <div class="settings-item-content">
                <el-slider 
                    v-model="chatConfigStore.AIConfig.max_tokens" 
                    :min="100" 
                    :max="16000" 
                    :step="100"
                />
            </div>
        </div>
        <!-- 模型top_p -->
        <div class="settings-item">
            <div class="settings-item-title">
                {{ $t('message.topP') }} 
                <span class="setting-value">{{ chatConfigStore.AIConfig.top_p }}</span>
            </div>
            <div class="settings-item-content">
                <el-slider 
                    v-model="chatConfigStore.AIConfig.top_p" 
                    :min="0" 
                    :max="1" 
                    :step="0.05"
                />
            </div>
        </div>
        <!-- 模型频率惩罚 -->
        <div class="settings-item">
            <div class="settings-item-title">
                {{ $t('message.frequencyPenalty') }} 
                <span class="setting-value">{{ chatConfigStore.AIConfig.frequency_penalty }}</span>
            </div>
            <div class="settings-item-content">
                <el-slider 
                    v-model="chatConfigStore.AIConfig.frequency_penalty" 
                    :min="-2" 
                    :max="2" 
                    :step="0.1"
                />
            </div>
        </div>
        <!-- 语言 -->
        <div class="settings-item">
            <div class="settings-item-title">{{ $t('message.language') }}</div>
            <div class="settings-item-content">
                <ChangeLangs />
            </div>
        </div>
        <!-- 主题 -->
        <div class="settings-item">
            <div class="settings-item-title">{{ $t('message.theme') }}</div>
            <div class="settings-item-content">
                <ChangeTheme />
            </div>
        </div>
    </div>
</template>

<script setup>
import { useChatConfigStore } from '@/stores/chat_config'
import ChangeLangs from '@/components/reuse/change_langs.vue'
import ChangeTheme from '@/components/reuse/change_theme.vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const chatConfigStore = useChatConfigStore()

// 解析ability字符串，按逗号分割并过滤空字符串
const parseAbilities = (abilityString) => {
    if (!abilityString || abilityString.trim() === '') {
        return []
    }
    return abilityString.split(',').map(item => item.trim()).filter(item => item !== '')
}

// 获取能力对应的图标
const getAbilityIcon = (ability) => {
    const iconMap = {
        'image': '🖼️',
        'audio': '🎵',
        'video': '🎬',
        'text': '📝'
    }
    return iconMap[ability] || '⚡'
}

// 获取能力对应的显示文本
const getAbilityText = (ability) => {
    const textMap = {
        'image': t('message.image'),
        'audio': t('message.audio'), 
        'video': t('message.video'),
        'text': t('message.text')
    }
    return textMap[ability] || ability
}
</script>

<style scoped>
.settings-container {
    width: 100%;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    overflow-y: auto;
}

.settings-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.settings-item-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 4px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.setting-value {
    font-size: 14px;
    font-weight: 500;
    color: #409eff;
    background-color: #ecf5ff;
    padding: 2px 8px;
    border-radius: 12px;
    min-width: 40px;
    text-align: center;
}

.settings-item-content {
    width: 100%;
    max-width: 400px;
}

/* 调整子组件样式以适应设置页面 */
.settings-item-content :deep(.selector-container) {
    width: 100%;
    min-width: unset;
    margin: 0;
}

.settings-item-content :deep(.el-select) {
    width: 100%;
}

.settings-item-content :deep(.el-select .el-select__wrapper) {
    min-height: 42px;
}

.settings-item-content :deep(.el-option) {
    padding: 8px 12px;
    min-height: 44px;
}

.settings-item-content :deep(.el-option:hover) {
    background-color: #f5f7fa;
}

.settings-item-content :deep(.el-option.is-selected) {
    background-color: #ecf5ff;
    color: #409eff;
}

/* 模型选项布局 */
.model-option {
    display: flex;
    align-items: center;
    width: 100%;
    gap: 1px;
    padding: 4px 0;
    min-height: 36px;
}

.model-show {
    font-size: 14px;
    font-weight: 500;
    color: #333;
    flex-shrink: 0;
    min-width: 120px;
}

.model-price {
    color: #666;
    font-weight: 500;
    font-size: 13px;
    flex-shrink: 0;
    min-width: 40px;
    text-align: right;
}

/* 能力标签容器 */
.model-abilities {
    display: flex;
    gap: 8px;
    flex: 1;
    align-items: center;
    justify-content: flex-end;
}

.ability-icons {
    display: flex;
    gap: 6px;
    flex-wrap: wrap;
    align-items: center;
}

.ability-texts {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
    align-items: center;
}

/* 能力标签基础样式 */
.ability-tag {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    font-size: 11px;
    font-weight: 500;
    background-color: #e8f5e8;
    color: #2d5a2d;
    border: 1px solid #c8e6c9;
    transition: all 0.2s ease;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.ability-tag:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 不同能力类型的样式 */
.ability-image {
    background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
    color: #1565c0;
    border-color: #90caf9;
}

.ability-audio {
    background: linear-gradient(135deg, #fce4ec 0%, #f8bbd9 100%);
    color: #c2185b;
    border-color: #f48fb1;
}

.ability-video {
    background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
    color: #ef6c00;
    border-color: #ffb74d;
}

.ability-text {
    color: #7b1fa2;
    border-color: #ce93d8;
}

.ability-icon {
    font-size: 12px;
    line-height: 1;
}

.ability-text {
    font-size: 11px;
    font-weight: 500;
    line-height: 1;
    color: #666;
    white-space: nowrap;
}


/* 当没有能力标签时的样式 */
.model-abilities:empty {
    display: none;
}

/* 滑动条样式 */
.settings-item-content :deep(.el-slider) {
    margin: 8px 0;
}

.settings-item-content :deep(.el-slider__runway) {
    background-color: #e4e7ed;
    border-radius: 3px;
    height: 6px;
}

.settings-item-content :deep(.el-slider__bar) {
    background-color: #409eff;
    border-radius: 3px;
}

.settings-item-content :deep(.el-slider__button) {
    width: 16px;
    height: 16px;
    border: 2px solid #409eff;
    background-color: #fff;
}

.settings-item-content :deep(.el-slider__button:hover) {
    transform: scale(1.2);
}

.settings-item-content :deep(.el-slider__stop) {
    background-color: #c0c4cc;
    width: 4px;
    height: 4px;
    border-radius: 50%;
}
</style>
