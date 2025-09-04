
// 滚动条到底
const scrollToBottom = (element) => {
    if (!element) return
    // 使用 requestAnimationFrame 确保在下一帧渲染时滚动
    requestAnimationFrame(() => {
        element.scrollTop = element.scrollHeight
    })
}

// 滚动条到顶
const scrollToTop = (element) => {
    if (!element) return
    element.scrollTop = 0
}

export const useScrollBar = () => {
    return {
        scrollToBottom,
        scrollToTop
    }
}
