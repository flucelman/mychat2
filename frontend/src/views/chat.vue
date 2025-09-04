<template>

    <div class="chat-container">
        <!-- 历史记录抽屉 -->
        <el-drawer v-if="isMobileView" v-model="chatConfigStore.showDrawer" class="el-drawer-container"  :with-header="false" :size="0" direction="ltr">
            <ChatDrawerMain />
        </el-drawer>
        <!-- 使用Transition包装非抽屉模式的ChatDrawerMain -->
        <Transition name="drawer-slide" appear>
            <ChatDrawerMain v-if="!isMobileView && chatConfigStore.showDrawer" class="chat-drawer"/>
        </Transition>
        <!-- 聊天窗口 -->
        <ChatWindowMain />
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useChatConfigStore } from '@/stores/chat_config'
import { useGlobalSettingStore } from '@/stores/global_setting'
import ChatDrawerMain from '@/components/chat/chatDrawer/chatDrawerMain.vue'
import ChatWindowMain from '@/components/chat/chatWindow/chatWindowMain.vue'

const chatConfigStore = useChatConfigStore()
const globalSettingStore = useGlobalSettingStore()
const isMobileView = ref(false)

// 检查是否为移动端视图（屏幕宽度小于768px）
const checkMobileView = () => {
    isMobileView.value = window.innerWidth < 768
    
    // 如果是移动端视图，默认隐藏抽屉
    if (isMobileView.value) {
        chatConfigStore.showDrawer = false
    }
}

// 使用媒体查询监听器
let mediaQueryList = null

const handleMediaQueryChange = (e) => {
    isMobileView.value = e.matches
    
    // 当切换到移动端视图时，隐藏抽屉
    if (isMobileView.value) {
        chatConfigStore.showDrawer = false
    }
}

onMounted(() => {
    // 初始检查
    checkMobileView()
    
    // 创建媒体查询
    mediaQueryList = window.matchMedia('(max-width: 767px)')
    
    // 添加监听器
    mediaQueryList.addEventListener('change', handleMediaQueryChange)
    
    // 同时添加resize事件作为备用方案
    window.addEventListener('resize', checkMobileView)
    
    // 获取用户信息等其他初始化操作
    globalSettingStore.checkToken()
    globalSettingStore.getUserInfo()
    chatConfigStore.getModelList()
})

onUnmounted(() => {
    // 清理事件监听器
    if (mediaQueryList) {
        mediaQueryList.removeEventListener('change', handleMediaQueryChange)
    }
    window.removeEventListener('resize', checkMobileView)
})

</script>

<style scoped>
.chat-container {
    width: 100%;
    height: 100vh;
    height: 100dvh;
    background-color: var(--background-color);
    color: var(--text-color);
    display: flex;
    flex-direction: row;
    overflow: hidden; /* 防止容器产生滚动条 */
}

.el-drawer-container{
    background-color: var(--secondary-background); 
    padding: 100px;
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