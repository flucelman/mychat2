<template>
    <div class="model">
        <div class="model-name-text" @click="showModelList = !showModelList">
            <img :src="API.backend_url + '/assets/icons/modelLogo/' + chatConfigStore.modelList.find(item => item.name == chatConfigStore.AIConfig.model)?.logo"
                class="model-icon" />
            {{ chatConfigStore.AIConfig.model }}
            <DropdownIcon class="dropdown-icon" />
        </div>
        <transition name="dropdown">
            <div class="model-list" v-show="showModelList">
                <div class="model-list-item" v-for="model in chatConfigStore.modelList" :key="model.name" @click="chatConfigStore.AIConfig.model = model.name,showModelList = false">
                    <img :src="API.backend_url + '/assets/icons/modelLogo/' + model.logo" class="model-icon" />
                    {{ model.name }}
                </div>
            </div>
        </transition>
    </div>
</template>

<script setup>
import DropdownIcon from '@/assets/icons/下拉.svg'
import { useChatConfigStore } from '@/stores/chat_config'
const chatConfigStore = useChatConfigStore()
import { API } from '@/router/api'
import { ref } from 'vue'
const showModelList = ref(false)

</script>

<style scoped lang="scss">
.model {
    font-size: 18px;
    font-weight: 600;
    margin-left: 10px;
    cursor: pointer;
    position: relative;


    .model-name-text {
        background-color: var(--secondary-background);
        padding: 5px 10px;
        border-radius: 10px;
        color: var(--primary-text);
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 10px;
    }

    .model-icon {
        width: 25px;
        height: 25px;
        color: var(--icon-color);
    }

    .dropdown-icon {
        width: 30px;
        height: 30px;
        color: var(--icon-color);
    }

    .model-list {
        width: 100%;
        position: absolute;
        top: 110%;
        left: 0;
        background-color: var(--background-color);
        font-size: 16px;
        box-shadow: var(--shadow-color);
        border-radius: 10px;
        overflow: hidden;
        transform-origin: top center;
        will-change: transform, opacity, max-height;
        z-index: 10;
        border: 1px solid var(--border-color);
    }

    .model-list-item{
        display: flex;
        justify-content: start;
        align-items: center;
        gap: 10px;
        padding: 10px;
    }
    .model-list-item:hover{
        background-color: var(--secondary-background);
    }
}

/* Dropdown animations */
.dropdown-enter-active,
.dropdown-leave-active {
    transition: opacity 0.18s ease, transform 0.2s cubic-bezier(.2,.8,.2,1), max-height 0.22s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
    opacity: 0;
    transform: translateY(-6px) scale(0.98);
    max-height: 0;
}
.dropdown-enter-to,
.dropdown-leave-from {
    opacity: 1;
    transform: translateY(0) scale(1);
    max-height: 60vh;
}

.model-icon {
    width: 20px;
    height: 20px;
}
</style>
