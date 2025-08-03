<template>
    <div class="upload-files-container">
        <uploadIcon/>
    </div>
</template>

<script setup>
import uploadIcon from '@/assets/icons/上传.svg'
</script>

<style scoped>
.upload-files-container{
    display: flex;
    justify-content: center;
    align-items: center;
    width: 25px;
    height: 25px;
    cursor: pointer;
    transition: transform 0.2s ease;
}
.upload-files-container:hover{
    transform: scale(1.1);
}

</style>