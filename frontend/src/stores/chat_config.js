import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { http, httpFetch } from '@/utils/http/client';
import { API } from '@/config/api'
import { useGlobalSettingStore } from '@/stores/global_setting'
import { v4 as uuidv4 } from 'uuid'

export const useChatConfigStore = defineStore('chatConfig', () => {
    const { t } = useI18n()
    // 是否显示历史记录抽屉
    const showDrawer = ref(true)
    // 是否展示lang、theme设置
    const openEyes = ref(false)
    // ================ AIChat ================
    const chatHistory = ref([])
    const chatId = ref('')
    const systemPrompt = ref('you are a helpful assistant')
    const AIConfig = reactive({
        model: 'gpt-4o-mini',
        temperature: 0.5,
        max_tokens: 4096,
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

    // 添加 AbortController 来管理请求取消
    let currentAbortController = null
    const assistantMessageId = ref(null)
    // 发送消息
    const sendUserMessage = async () => {
        if (userMessage.value.trim() === '') {
            ElMessage.error(t('message.input_placeholder'))
            return
        }
        isReceiving.value = true

        // 创建新的 AbortController
        currentAbortController = new AbortController()

        // 保存用户消息的值，避免被清空
        const currentUserMessage = userMessage.value
        // 创建用户消息的副本并推入历史记录
        baseMessageHistory.push({
            role: 'user',
            content: currentUserMessage,  // 使用保存的值
            model: AIConfig.model,
            message_id: uuidv4()
        })

        
        // 创建要发送的消息对象
        const messageToSend = {
            chat_id: chatId.value,
            AI_config: AIConfig,
            message_history: sendMessageHistory.value,  // 使用.value获取计算属性的值
            file_url: fileUrl.value
        }

        // 把当前聊天放到最前面
        const index = chatHistory.value.findIndex(item => item.chat_id == chatId.value);
        if (index > 0) { // 如果找到且不在第一位
            const item = chatHistory.value.splice(index, 1)[0]; // 删除并获取元素
            chatHistory.value.unshift(item); // 添加到最前面
        }
        // 清空输入框和之前的AI响应
        userMessage.value = ''
        instantAssistantMessage.value = ''
        
        try {
            // 获取全局设置store来获取token
            const globalSettingStore = useGlobalSettingStore()

            // 直接使用httpFetch发送请求，添加 AbortController 的 signal
            const response = await httpFetch(API.backend_url + '/api/chat/addChatMessage', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': globalSettingStore.userToken,
                },
                body: JSON.stringify(messageToSend),
                signal: currentAbortController.signal  // 添加取消信号
            })

            if (!response.ok) {
                throw new Error(`请求失败: ${response.status}`)
            }

            const reader = response.body.getReader()
            const decoder = new TextDecoder()


            while (true) {
                // 检查是否已被取消
                if (currentAbortController.signal.aborted) {
                    reader.cancel()
                    break
                }

                const { done, value } = await reader.read()

                if (done) {
                    // 将完整的AI响应添加到消息历史
                    baseMessageHistory.push({
                        role: 'assistant',
                        content: instantAssistantMessage.value,
                        model: AIConfig.model,
                        message_id: assistantMessageId.value
                    })
                    instantAssistantMessage.value = ''
                    break
                }
                const chunk = decoder.decode(value, { stream: true })

                if (chunk.startsWith('event:start')) {
                    const lines = chunk.split('\n')
                    const jsonStr = lines[1].replace('data:', ''); // 去掉 "data:" 前缀
                    const data = JSON.parse(jsonStr);
                    chatId.value = data.chat_id
                    assistantMessageId.value = data.assistant_message_id
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
                    const lines = chunk.split('\n')
                    const jsonStr = lines[1].replace('data:', ''); // 去掉 "data:" 前缀
                    const data = JSON.parse(jsonStr);

                    isReceiving.value = false
                }


            }
        } catch (error) {
            // 检查是否是因为取消导致的错误
            if (error.name === 'AbortError') {
                console.log("请求已被用户取消")
                ElMessage.info(t('message.request_cancelled'))
            } else {
                console.error("发送消息时出错:", error)
                ElMessage.error(t('message.AI_response_error'))
            }
        } finally {
            isReceiving.value = false
            currentAbortController = null
            checkUpdateChatHistory()
        }
    }

    // 检查是否需要更新聊天历史表
    const checkUpdateChatHistory = () => {
        // 检查是否已存在当前聊天记录
        const existingChat = chatHistory.value.find(item => item.chat_id == chatId.value)
        if (!existingChat) {
            chatHistory.value.unshift({
                chat_id: chatId.value,
                title: baseMessageHistory[0].content,
                update_time: new Date().toISOString()
            })
        }
    }

    // 取消连接
    const cancelConnection = () => {

        // 如果有正在进行的请求，取消它
        if (currentAbortController && !currentAbortController.signal.aborted) {
            currentAbortController.abort()
        }

        // 重置状态
        isReceiving.value = false
        currentAbortController = null

        if (instantAssistantMessage.value != '') {
            baseMessageHistory.push({
                role: 'assistant',
                content: instantAssistantMessage.value,
                model: AIConfig.model,
                message_id: assistantMessageId.value
            })
            instantAssistantMessage.value = ''
        }
        // 清理 AbortController
    }

    // 获取聊天历史
    const getChatHistory = async () => {
        const response = await http.get(API.backend_url + '/api/chat/getChatHistory')
        chatHistory.value = response.data.data
    }

    // 获取单个聊天记录
    const getChatMessage = async (chat_id) => {
        chatId.value = chat_id
        const response = await http.get(API.backend_url + `/api/chat/getChatMessage/${chat_id}`)
        // 清空数组并添加新数据，而不是直接赋值
        baseMessageHistory.splice(0, baseMessageHistory.length, ...response.data)
    }

    // 删除所有聊天记录
    const deleteAllHistory = async () => {
        await http.delete(API.backend_url + '/api/chat/deleteAllHistory')
        chatHistory.value = []
        baseMessageHistory.splice(0, baseMessageHistory.length)
    }

    // 删除单个聊天记录
    const deleteSingleHistory = async (chat_id) => {
        await http.delete(API.backend_url + `/api/chat/deleteSingleHistory/${chat_id}`)
        chatHistory.value = chatHistory.value.filter(item => item.chat_id != chat_id)
        if (chatId.value == chat_id) {
            baseMessageHistory.splice(0, baseMessageHistory.length)
        }
    }

    // 新建聊天
    const newChat = async () => {
        chatId.value = ''
        baseMessageHistory.splice(0, baseMessageHistory.length)
        ElMessage.success(t('message.newChatSuccess'))
    }

    return {
        showDrawer,
        openEyes,
        chatId,
        AIConfig,
        sendUserMessage,
        userMessage,
        instantAssistantMessage,
        fileUrl,
        baseMessageHistory,
        isReceiving,
        cancelConnection,
        getChatHistory,
        chatHistory,
        getChatMessage,
        deleteAllHistory,
        deleteSingleHistory,
        newChat
    }
}, {
    persist: {
        pick: ['showDrawer', 'openEyes']  // 只持久化 showDrawer 字段
    }
})