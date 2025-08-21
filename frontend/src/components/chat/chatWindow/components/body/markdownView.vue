<template>
  <div class="markdown-body" v-html="renderedContent" ref="containerRef"></div>
</template>

<script setup>
import { computed, onMounted, onBeforeUnmount, ref } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-light.css'
import 'katex/dist/katex.min.css'
import { ElMessage } from 'element-plus';

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

// 计算渲染后的HTML
const renderedContent = computed(() => {
  return md.render(props.content || '');
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
  color: var(--primary-text);
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
  color: var(--primary-text);
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
@media (prefers-color-scheme: dark) {
  .markdown-body code {
    background: linear-gradient(135deg, rgba(99, 102, 241, 0.15) 0%, rgba(168, 85, 247, 0.15) 100%);
    border-color: rgba(99, 102, 241, 0.3);
    color: #a5b4fc;
  }
  
  .markdown-body pre {
    background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
    border-color: rgba(148, 163, 184, 0.3);
    box-shadow: 
      0 4px 6px -1px rgba(0, 0, 0, 0.3),
      0 2px 4px -1px rgba(0, 0, 0, 0.2),
      inset 0 1px 0 rgba(255, 255, 255, 0.1);
  }
  
  .copy-btn {
    background: linear-gradient(135deg, rgba(30, 41, 59, 0.95) 0%, rgba(15, 23, 42, 0.95) 100%);
    border-color: rgba(99, 102, 241, 0.4);
    color: #a5b4fc;
  }
  
  .copy-btn:hover {
    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
    color: white;
  }
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
}
</style>