<template>
    <div class="drawer-footer">
        <div class="drawer-footer-item new-chat" @click="chatConfigStore.newChat()">
            <NewChatIcon class="new-chat-icon" />
            {{ $t('message.newChat') }}
        </div>
        <div class="drawer-footer-item" @click="settingsDialogVisible = true">
            <SettingsIcon class="settings-icon" />
            {{ $t('message.settings') }}
        </div>
        <div class="drawer-footer-item" @click="globalSettingStore.logout()">
            <PersonalCenterIcon class="personal-center-icon" />
            {{ $t('message.personalCenter') }}
        </div>
        <!-- 对话框 -->
        <el-dialog v-model="settingsDialogVisible" :title="$t('message.settings')" width="400" align-center>
            <Settings />
        </el-dialog>
    </div>
</template>

<script setup>
import NewChatIcon from '@/assets/icons/新建聊天.svg'
import SettingsIcon from '@/assets/icons/设置.svg'
import PersonalCenterIcon from '@/assets/icons/个人中心.svg'
import { useGlobalSettingStore } from '@/stores/global_setting'
import { useChatConfigStore } from '@/stores/chat_config'
import { ref } from 'vue'
import Settings from './settings.vue'
const globalSettingStore = useGlobalSettingStore()
const chatConfigStore = useChatConfigStore()
const settingsDialogVisible = ref(false)
</script>

<style scoped>
.drawer-footer {
    padding: 20px;
    border-top: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    flex-direction: column;
}

.drawer-footer-item {
    color: var(--text-color);
    font-size: 18px;
    margin-top: 10px;
    background-color: var(--secondary-background);
    padding: 10px;
    border-radius: 10px;
    cursor: pointer;
    width: 85%;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;

}

.drawer-footer-item:hover {
    background-color: var(--tertiary-background);
}

.new-chat {
    background-color: var(--tertiary-background);
}

.new-chat-icon,
.personal-center-icon,
.settings-icon {
    width: 20px;
    height: 20px;
    color: var(--text-color);
}
</style>