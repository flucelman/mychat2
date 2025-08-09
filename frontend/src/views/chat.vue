<template>

    <div class="chat-container">
        <!-- 历史记录抽屉 -->
        <el-drawer v-if="allowDrawer" v-model="chatConfigStore.showDrawer" style="background-color: var(--secondary-background); padding: 0;"  :with-header="false" :size="0" direction="ltr">
            <ChatDrawerMain />
        </el-drawer>
        <!-- 使用Transition包装非抽屉模式的ChatDrawerMain -->
        <Transition name="drawer-slide" appear>
            <ChatDrawerMain v-if="!allowDrawer && chatConfigStore.showDrawer" class="chat-drawer"/>
        </Transition>
        <!-- 聊天窗口 -->
        <ChatWindowMain />
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
import { useGlobalSettingStore } from '@/stores/global_setting'
import ChatDrawerMain from '@/components/chat/chatDrawer/chatDrawerMain.vue'
import ChatWindowMain from '@/components/chat/chatWindow/chatWindowMain.vue'
const chatConfigStore = useChatConfigStore()
const globalSettingStore = useGlobalSettingStore()
const allowDrawer = ref(false)
onMounted(() => {
    // 判断是否允许抽屉
    if (globalSettingStore.isMobile) {
        allowDrawer.value = true
    }else{
        allowDrawer.value = false
    }
    globalSettingStore.checkToken()
  globalSettingStore.getUserInfo()
  chatConfigStore.getModelList()
})

</script>

<style scoped>
.chat-container {
    width: 100%;
    height: 100vh;
    background-color: var(--background-color);
    color: var(--text-color);
    display: flex;
    flex-direction: row;
    overflow: hidden; /* 防止容器产生滚动条 */
}


/* 抽屉滑动动画 */
.drawer-slide-enter-active {
    transition: all 0.3s ease-in-out;
}

.drawer-slide-leave-active {
    transition: none;
}

.drawer-slide-enter-from {
    transform: translateX(-100%);
    opacity: 0;
}

.drawer-slide-leave-to {
    transform: translateX(-100%);
    opacity: 0;
}

.drawer-slide-enter-to,
.drawer-slide-leave-from {
    transform: translateX(0);
    opacity: 1;
}

.el-drawer__body{
    padding: 0;
}
</style>