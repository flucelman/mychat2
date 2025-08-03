<template>
    <div class="eyes-container">
        <EyesIcon v-if="!chatConfigStore.openEyes" class="eyes-icon" @click="handleClick" />
        <CloseEyesIcon v-else class="eyes-icon" @click="handleClick" />
        <div class="controls" v-show="!chatConfigStore.openEyes">
            <ChangeTheme />
            <ChangeLangs />
        </div>
    </div>
</template>

<script setup>
import EyesIcon from '@/assets/icons/睁眼.svg'
import CloseEyesIcon from '@/assets/icons/闭眼.svg'
import ChangeLangs from '@/components/reuse/change_langs.vue'
import ChangeTheme from '@/components/reuse/change_theme.vue'
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()

const handleClick = () => {
    chatConfigStore.openEyes = !chatConfigStore.openEyes
}
</script>

<style scoped>
.eyes-container {
    position: relative;
}

.eyes-icon {
    width: 35px;
    height: 35px;
    cursor: pointer;
    color: var(--icon-color);
    transition: color 0.2s ease;
}

.controls {
    display: flex;
    gap: 20px;
    align-items: center;
    position: absolute;
    top: 50px;
    right: 0;
    z-index: 1000;
    flex-direction: column;
}
</style>
