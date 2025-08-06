import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { http, httpFetch } from '@/utils/http/client';
import { API } from '@/config/api'
import { useGlobalSettingStore } from '@/stores/global_setting'

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
    const instantAssistantMessage = ref('')
    const fileUrl = ref([])
    const isReceiving = ref(false)

    const sendUserMessage = async () => {
        if (userMessage.value.trim() === '') {
            ElMessage.error(t('message.input_placeholder'))
            return
        }
        isReceiving.value = true
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


        // 清空输入框和之前的AI响应
        userMessage.value = ''
        instantAssistantMessage.value = ''

        try {
            // 获取全局设置store来获取token
            const globalSettingStore = useGlobalSettingStore()
            
            // 直接使用httpFetch发送请求，不使用http.post以避免响应体被提前消耗
            const response = await httpFetch(API.backend_url + '/api/chat/add_user_message', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': globalSettingStore.userToken,
                },
                body: JSON.stringify(messageToSend)
            })

            if (!response.ok) {
                throw new Error(`请求失败: ${response.status}`)
            }

            const reader = response.body.getReader()
            const decoder = new TextDecoder()

            console.log("开始读取流式数据")

            while (true) {
                const { done, value } = await reader.read()
                if (done) {
                    console.log("流式数据读取完成")
                    // 将完整的AI响应添加到消息历史
                    baseMessageHistory.push({
                        role: 'assistant',
                        content: instantAssistantMessage.value
                    })
                    instantAssistantMessage.value = ''
                    break
                }
                
                const chunk = decoder.decode(value, { stream: true })
                
                if (chunk.startsWith('event:start')) {
                    console.log("event:start")
                    const lines = chunk.split('\n')
                    const jsonStr = lines[1].replace('data:', ''); // 去掉 "data:" 前缀
                    const data = JSON.parse(jsonStr);
                    chatId.value = data.chat_id
                    console.log("已保存chatId:", chatId.value)
                }
                else if (chunk.startsWith('event:content')) {
                    const lines = chunk.split('\n')
                    for (const line of lines) {
                        if (line.startsWith('data:')) {
                            const jsonStr = line.replace('data:', ''); // 去掉 "data:" 前缀
                            const data = JSON.parse(jsonStr);
                            instantAssistantMessage.value += data.content
                        }
                    }
                }
                else if (chunk.startsWith('event:end')) {
                    console.log("event:end")
                    const lines = chunk.split('\n')
                    const jsonStr = lines[1].replace('data:', ''); // 去掉 "data:" 前缀
                    const data = JSON.parse(jsonStr);
                    console.log("AI回复结束")
                    isReceiving.value = false
                }
            }
        } catch (error) {
            console.error("发送消息时出错:", error)
            ElMessage.error(t('message.AI_response_error'))
        }
    }
    return {
        showDrawer,
        openEyes,
        AIConfig,
        sendUserMessage,
        userMessage,
        instantAssistantMessage,
        fileUrl,
        baseMessageHistory,
        isReceiving
    }
}, {
    persist: {
        pick: ['showDrawer']  // 只持久化 showDrawer 字段
    }
})