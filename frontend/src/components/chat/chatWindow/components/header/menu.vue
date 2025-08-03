<template>
    <MenuIcon @click="handleMenuClick" class="menu-icon" />
</template>

<script setup>
import MenuIcon from '@/assets/icons/菜单.svg'
import { useChatConfigStore } from '@/stores/chat_config'

const chatConfigStore = useChatConfigStore()

const handleMenuClick = () => {
    chatConfigStore.showDrawer = !chatConfigStore.showDrawer
}
</script>

<style scoped>
.menu-icon {
    width: 35px;
    height: 35px;
    cursor: pointer;
    color: var(--icon-color);
    transition: color 0.2s ease;
}
</style>