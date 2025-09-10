<template>
  <div class="markdown-body" v-html="renderedContent" ref="containerRef"></div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-light.css'
import 'katex/dist/katex.min.css'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n() 

const props = defineProps({
  content: {
    type: String,
    required: true
  }
})

// 创建markdown-it实例并配置插件
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><button class="copy-btn" data-code="${encodeURIComponent(str)}">复制</button><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`;
      } catch (__) {}
    }
    return `<pre class="hljs"><button class="copy-btn" data-code="${encodeURIComponent(str)}">复制</button><code>${md.utils.escapeHtml(str)}</code></pre>`;
  }
})

// 配置链接在新窗口打开
const defaultRender = md.renderer.rules.link_open || function(tokens, idx, options, env, self) {
  return self.renderToken(tokens, idx, options);
};
md.renderer.rules.link_open = function (tokens, idx, options, env, self) {
  tokens[idx].attrPush(['target', '_blank']);
  return defaultRender(tokens, idx, options, env, self);
};

// 解析自定义标签的函数
const parseCustomTags = (content) => {
  if (!content) return '';
  
  // 处理 toDoList 标签
  content = content.replace(/<toDoList>(.*?)<\/toDoList>/gs, (match, innerContent) => {
    return `<div class="custom-card todo-card">
      <div class="card-header">
        <div class="card-icon todo-icon">📋</div>
        <div class="card-title">${t('message.todoList')}</div>
      </div>
      <div class="card-content">${innerContent.trim()}</div>
    </div>`;
  });
  
  // 处理 task 标签
  content = content.replace(/<task>(.*?)<\/task>/gs, (match, innerContent) => {
    return `<div class="custom-card task-card">
      <div class="card-header">
        <div class="card-icon task-icon">⚡</div>
        <div class="card-title">${t('message.task')}</div>
      </div>
      <div class="card-content">${innerContent.trim()}</div>
    </div>`;
  });
  
  // 处理 summary 标签（空标签，只显示卡片）
  content = content.replace(/<summary\s*\/?>(?:<\/summary>)?/g, () => {
    return `<div class="summary-compact">
      <div class="summary-compact-content">
        <div class="summary-compact-icon">📝</div>
        <span class="summary-compact-text">${t('message.summary')}</span>
      </div>
    </div>`;
  });
  
  return content;
};

// 计算渲染后的HTML
const renderedContent = computed(() => {
  const processedContent = parseCustomTags(props.content || '');
  return md.render(processedContent);
});

const containerRef = ref(null)

const handleClick = (event) => {
  const target = event.target;
  if (target && target.classList && target.classList.contains('copy-btn')) {
    const raw = target.getAttribute('data-code') || '';
    const code = decodeURIComponent(raw);
    navigator.clipboard.writeText(code)
      .then(() => ElMessage.success('复制成功'))
      .catch(() => ElMessage.error('复制失败'));
  }
}

onMounted(() => {
  if (containerRef.value) {
    containerRef.value.addEventListener('click', handleClick)
  }
})

onBeforeUnmount(() => {
  if (containerRef.value) {
    containerRef.value.removeEventListener('click', handleClick)
  }
})
</script>

<style lang="scss">
.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  line-height: 2;
  color: var(--text-color);
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
  user-select: text;
}

/* 确保内部元素也可被选中（例如代码块、链接等） */
.markdown-body, .markdown-body * {
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
  user-select: text;
}

/* 内联代码样式 */
.markdown-body code {
  font-family: 'JetBrains Mono', 'Fira Code', 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  border: 1px solid rgba(99, 102, 241, 0.15);
  border-radius: 6px;
  padding: 0.25em 0.5em;
  font-size: 0.875em;
  font-weight: 500;
  color: var(--text-color);
  transition: all 0.2s ease;
  word-break: break-all;
  white-space: pre-wrap;
}

.markdown-body code:hover {
  border-color: rgba(99, 102, 241, 0.25);
  transform: translateY(-1px);
}

/* 代码块容器样式 */
.markdown-body pre {
  background: linear-gradient(145deg, #f8fafc 0%, #f1f5f9 100%);
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 12px;
  padding: 20px;
  margin: 16px 0;
  overflow: auto;
  position: relative;
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.8);
  transition: all 0.3s ease;
  word-wrap: break-word;
  overflow-wrap: break-word;
  white-space: pre-wrap;
  max-width: 100%;
}

.markdown-body pre:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 8px 25px -5px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.9);
  border-color: rgba(99, 102, 241, 0.3);
}

/* 代码块内的代码样式 */
.markdown-body pre code {
  background-color: transparent;
  padding: 0;
  border: none;
  color: inherit;
  font-size: 0.875em;
  line-height: 1.5;
  word-break: break-all;
  white-space: pre-wrap;
  overflow-wrap: break-word;
}

/* highlight.js 样式重置 */
.hljs {
  padding: 0;
  background: transparent;
  border-radius: 0;
}

/* 复制按钮样式 */
.copy-btn {
  cursor: pointer;
  font-size: 12px;
  font-weight: 600;
  position: absolute;
  right: 12px;
  top: 12px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(248, 250, 252, 0.95) 100%);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 8px;
  padding: 8px 14px;
  color: #4f46e5;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  backdrop-filter: blur(8px);
  box-shadow: 
    0 2px 4px rgba(0, 0, 0, 0.1),
    0 1px 2px rgba(0, 0, 0, 0.06);
  z-index: 10;
  letter-spacing: 0.025em;
}

.copy-btn:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  border-color: #4f46e5;
  color: white;
  transform: translateY(-2px) scale(1.05);
  box-shadow: 
    0 8px 25px rgba(79, 70, 229, 0.25),
    0 4px 10px rgba(79, 70, 229, 0.15);
}

.copy-btn:active {
  transform: translateY(-1px) scale(1.02);
  box-shadow: 
    0 4px 15px rgba(79, 70, 229, 0.2),
    0 2px 5px rgba(79, 70, 229, 0.1);
}

/* 暗色主题适配 */
[theme="dark"] .markdown-body code {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(168, 85, 247, 0.15) 100%);
  border-color: var(--border-color-hover);
  color: var(--text-color);
}

[theme="dark"] .markdown-body pre {
  background: linear-gradient(145deg, var(--secondary-background) 0%, var(--background-color) 100%);
  border-color: var(--border-color);
  box-shadow: var(--shadow-color);
}

[theme="dark"] .copy-btn {
  background: linear-gradient(135deg, rgba(51, 51, 51, 0.95) 0%, rgba(41, 41, 41, 0.95) 100%);
  border-color: var(--border-color-hover);
  color: var(--text-color);
}

[theme="dark"] .copy-btn:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

/* 自定义卡片样式 */
.custom-card {
  background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
  border: 1px solid rgba(148, 163, 184, 0.2);
  border-radius: 16px;
  margin: 20px 0;
  overflow: hidden;
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.custom-card:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 10px 25px -3px rgba(0, 0, 0, 0.1),
    0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.05) 0%, rgba(168, 85, 247, 0.05) 100%);
  border-bottom: 1px solid rgba(148, 163, 184, 0.1);
}

.card-icon {
  font-size: 20px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.8) 0%, rgba(248, 250, 252, 0.8) 100%);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--secondary-text);
  letter-spacing: 0.025em;
}

.card-content {
  padding: 20px;
  color: var(--text-color);
  line-height: 1.6;
  font-size: 14px;
  // 加粗
  font-weight: 600;
  color: var(--secondary-text);
}

/* 待办清单卡片特殊样式 */
.todo-card {
  border-left: 4px solid #10b981;
}

.todo-card .card-header {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.05) 0%, rgba(5, 150, 105, 0.05) 100%);
}

.todo-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.1) 100%);
  color: #10b981;
}

/* 任务步骤卡片特殊样式 */
.task-card {
  border-left: 4px solid #f59e0b;
}

.task-card .card-header {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.05) 0%, rgba(217, 119, 6, 0.05) 100%);
}

.task-icon {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.1) 0%, rgba(217, 119, 6, 0.1) 100%);
  color: #f59e0b;
}

/* 总结紧凑样式 */
.summary-compact {
  display: block;
  margin-bottom: 16px;
  text-align: left;
}

.summary-compact-content {
  display: inline-flex;
  align-items: center;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.08) 0%, rgba(79, 70, 229, 0.08) 100%);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 10px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  color: #6366f1;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 4px rgba(99, 102, 241, 0.1);
}

.summary-compact-content:hover {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.12) 0%, rgba(79, 70, 229, 0.12) 100%);
  border-color: rgba(99, 102, 241, 0.3);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(99, 102, 241, 0.15);
}

.summary-compact-icon {
  font-size: 16px;
  margin-right: 8px;
  display: flex;
  align-items: center;
}

.summary-compact-text {
  color: var(--text-color);
  opacity: 0.8;
}

/* 自定义卡片暗色主题适配 */
[theme="dark"] .custom-card {
  background: linear-gradient(145deg, var(--secondary-background) 0%, var(--background-color) 100%);
  border-color: var(--border-color);
}

[theme="dark"] .card-header {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(168, 85, 247, 0.1) 100%);
  border-bottom-color: var(--border-color);
}

[theme="dark"] .card-icon {
  background: linear-gradient(135deg, var(--tertiary-background) 0%, var(--secondary-background) 100%);
}

[theme="dark"] .todo-card .card-header {
  background: linear-gradient(135deg, rgba(147, 224, 198, 0.1) 0%, rgba(122, 239, 202, 0.1) 100%);
}

[theme="dark"] .task-card .card-header {
  background: linear-gradient(135deg, rgba(252, 220, 164, 0.1) 0%, rgba(245, 199, 146, 0.1) 100%);
}

[theme="dark"] .summary-compact-content {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.12) 0%, rgba(79, 70, 229, 0.12) 100%);
  border-color: var(--border-color-hover);
  color: var(--text-color);
}

[theme="dark"] .summary-compact-content:hover {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.18) 0%, rgba(79, 70, 229, 0.18) 100%);
  border-color: var(--border-color-hover);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .markdown-body pre {
    padding: 16px;
    margin: 12px 0;
    border-radius: 8px;
  }
  
  .copy-btn {
    right: 8px;
    top: 8px;
    padding: 6px 10px;
    font-size: 11px;
  }
  
  .custom-card {
    margin: 16px 0;
    border-radius: 12px;
  }
  
  .card-header {
    padding: 12px 16px;
  }
  
  .card-icon {
    width: 28px;
    height: 28px;
    font-size: 18px;
    margin-right: 10px;
  }
  
  .card-title {
    font-size: 15px;
  }
  
  .card-content {
    padding: 16px;
    font-size: 13px;
  }
  
  .summary-compact {
    margin: 12px 0;
  }
  
  .summary-compact-content {
    padding: 6px 12px;
    font-size: 13px;
    border-radius: 16px;
  }
  
  .summary-compact-icon {
    font-size: 14px;
    margin-right: 6px;
  }
}
</style>