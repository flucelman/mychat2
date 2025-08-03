<template>
    <div class="chat-window-container">
        <HeaderMain/>
        <BodyMain />
        <FooterMain/>
    </div>
</template>

<script setup>
import HeaderMain from './components/header.vue'
import BodyMain from './components/body.vue'
import FooterMain from './components/footer.vue'
</script>

<style scoped>
.chat-window-container {
    width: 100%;
    height: 100vh; /* 减去padding的高度 */
    background-color: var(--background-color);
    padding: 10px;
    /* 移动端顶部安全区域适配 */
    padding-top: env(safe-area-inset-top);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    overflow: hidden; /* 防止内容溢出 */
}

</style>