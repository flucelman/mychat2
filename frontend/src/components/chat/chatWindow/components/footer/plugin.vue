<template>
    <div class="plugin-container">
       <div class="plugin-item">
        <DeepThinkingIcon class="plugin-item-icon" />
        {{ $t('message.deepThinking') }}
       </div>
       <div class="plugin-item">
        <OnlineSearchIcon class="plugin-item-icon" />
        {{ $t('message.onlineSearch') }}
       </div>
    </div>
</template>

<script setup>
import DeepThinkingIcon from '@/assets/icons/深度思考.svg'
import OnlineSearchIcon from '@/assets/icons/联网搜索.svg'
</script>

<style scoped>
.plugin-container{
    display: flex;
    gap: 10px;
}
.plugin-item{
    padding: 10px;
    border-radius: 10px;
    border: 1px solid var(--border-color);
    cursor: pointer;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 5px;
}
.plugin-item-icon{
    width: 20px;
    height: 20px;
}
</style>