import { createI18n } from 'vue-i18n'
import zhCN from './langs/zh-CN'
import en from './langs/en'

const messages = {
    zhCN,
    en
}

const i18n = createI18n({
    legacy: false,
    locale: 'zhCN', // 设置默认语言
    messages
})
export default i18n