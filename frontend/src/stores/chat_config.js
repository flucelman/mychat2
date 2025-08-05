import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { http } from '@/utils/http/client';
import { API } from '@/config/api'

export const useChatConfigStore = defineStore('chatConfig', () => {
    const { t } = useI18n()
    // 是否显示历史记录抽屉
    const showDrawer = ref(true)
    // 是否展示lang、theme设置
    const openEyes = ref(false)
    // ================ AIChat ================
    const chatId = ref('')
    const systemPrompt = ref('you are a helpful assistant')
    const AIConfig = reactive({
        model: 'gpt-4o-mini',
        temperature: 0.5,
        max_tokens: 1000,
        top_p: 1,
        frequency_penalty: 0,
    })
    // 基础消息历史（不包含system消息）
    const baseMessageHistory = reactive([])
    
    // 使用计算属性，自动包含system消息
    const sendMessageHistory = computed(() => [
        {
            role: 'system',
            content: systemPrompt.value
        },
        ...baseMessageHistory
    ])
    const userMessage = ref('')
    const assistantMessage = ref('')
    const fileUrl = ref([])
    
    /*
    发送用户消息
    {
        "chat_id": chatId,
        "AI_config": AIConfig,
        "chat_message": {
            "role": "user",
            "content": userMessage.value,
            "model": AIConfig.model
        },
        "file_url": fileUrl
    }
    */
    const sendUserMessage = () => {
        if (userMessage.value === '') {
            ElMessage.error(t('message.input_placeholder'))
            return
        }
        
        // 保存用户消息的值，避免被清空
        const currentUserMessage = userMessage.value
        // 创建用户消息的副本并推入历史记录
        baseMessageHistory.push({
            role: 'user',
            content: currentUserMessage  // 使用保存的值
        })
                
        
        // 创建要发送的消息对象
        const messageToSend = {
            chat_id: chatId.value,
            AI_config: AIConfig,
            message_history: sendMessageHistory.value,  // 使用.value获取计算属性的值
            file_url: fileUrl.value
        }
        

        // 清空输入框
        userMessage.value = ''
        
        // 发送用户消息
        http.post(API.backend_url + '/api/chat/add_user_message', messageToSend).then(res => {
            // 先打印响应数据，看看实际结构
            console.log("完整响应数据:", res)
            console.log("res.data:", res.data)
            
            // 根据实际数据结构调整
            let chat_id, AI_response
            
            chat_id = res.data.data.chat_id
            AI_response = res.data.data.AI_response
            
            // 保存助手消息
            chatId.value = chat_id
            assistantMessage.value = AI_response
            
            // 创建助手消息的副本并推入历史记录
            baseMessageHistory.push({
                role: 'assistant',
                content: AI_response
            })
            
        }).catch(err => {
            console.log("AI_response error: ", err)
            ElMessage.error(t('message.AI_response_error'))
        })
        console.log("messageHistory: ", sendMessageHistory)
    }
    
    return {
        showDrawer,
        openEyes,
        AIConfig,
        sendUserMessage,
        userMessage,
        fileUrl,
        baseMessageHistory
    }
}, {
    persist: ['showDrawer']
})