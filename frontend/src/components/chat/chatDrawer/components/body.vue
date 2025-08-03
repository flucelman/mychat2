<template>
    <div class="drawer-body">
        <div class="drawer-body-item">
            <div>第1条记录</div>
            <ThreePointsIcon class="three-points-icon" />
        </div>
        <div class="drawer-body-item">
            <div>第2条记录</div>
            <ThreePointsIcon class="three-points-icon" />
        </div>
        <div class="drawer-body-item">
            <div>第3条记录</div>
            <ThreePointsIcon class="three-points-icon" />
        </div>
        <div class="drawer-body-item">
            <div>第4条记录</div>
            <ThreePointsIcon class="three-points-icon" />
        </div>
        
    </div>
</template>

<script setup>
import ThreePointsIcon from '@/assets/icons/三个点.svg'
</script>

<style scoped>
.drawer-body {
    overflow-y: auto;
    padding: 10px;
}

.drawer-body-item {
    padding: 10px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    border-radius: 10px;
}
.drawer-body-item:hover {
    background-color: var(--tertiary-background);
    .three-points-icon{
        display: block;
    }
}
.three-points-icon {
    width: 20px;
    height: 20px;
    color: var(--text-color);
    cursor: pointer;
    display: none;
}
</style>