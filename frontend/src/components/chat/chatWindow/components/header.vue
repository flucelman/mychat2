<template>
    <div class="header-container">
        <div class="left-items">
            <Menu1 />
            <Models />
        </div>
        <Eyes />
    </div>
</template>

<script setup>
import Menu1 from './header/menu.vue'
import Models from './header/models.vue'
import Eyes from './header/eyes.vue'
</script>

<style scoped>
.header-container {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.left-items {
    display: flex;
    align-items: center;
}
</style>