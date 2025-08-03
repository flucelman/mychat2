<template>
            <div class="model">
            <div class="model-name-text">
                <ChatGPT class="model-icon" />
                ChatGPT
                <DropdownIcon class="dropdown-icon" />
            </div>
        </div>
</template>
    
<script setup>
import ChatGPT from '@/assets/icons/modelcons/ChatGPT.svg'
import DropdownIcon from '@/assets/icons/下拉.svg'
</script>
    
<style scoped>
    .model {
    font-size: 20px;
    font-weight: 600;
    margin-left: 10px;
    cursor: pointer;

    .model-name-text {
        background-color: var(--secondary-background);
        padding: 5px 10px;
        border-radius: 10px;
        color: var(--text-color);
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
}

.model-icon {
    width: 20px;
    height: 20px;
    color: var(--icon-color);
}
</style>
