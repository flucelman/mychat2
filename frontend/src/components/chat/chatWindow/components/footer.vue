<template>
        <div class="input-container">
            <Textarea1/>
            <div class="input-container-inner">
                <Plugin/>
                <div class="input-container-inner-right">
                    <UploadFiles/>
                    <Send/>
                </div>
            </div>
        </div>
</template>

<script setup>
import Plugin from './footer/plugin.vue'
import Send from './footer/send.vue'
import UploadFiles from './footer/uploadFiles.vue'
import Textarea1 from './footer/textarea.vue'

</script>

<style scoped>
@media screen and (max-width: 1200px) {
    .input-container{
        width: 100%;
    }
}
@media screen and (min-width: 1200px) {
    .input-container{
        width: 50%;
        max-width: 800px;
    }
}

.input-container{
    display: flex;
    justify-content: center;
    align-items: flex-start;
    gap: 10px;
    flex-direction: column;
    border: 1px solid var(--border-color);
    border-radius: 12px;
    padding: 6px 16px 16px 16px;
    box-shadow: var(--shadow-color);
    transition: box-shadow 0.3s ease, transform 0.2s ease;
}

.input-container:hover{
    box-shadow: var(--shadow-color);
    transform: translateY(-2px);
}

.input-container:focus-within{
    box-shadow: var(--shadow-color);
    border-color: var(--border-color-hover);
    border-width: 2px;
}



.input::placeholder{
    color: var(#999);
    font-size: 16px;
    opacity: 0.7;
}
.input-container-inner{
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
}
.input-container-inner-right{
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}
</style>