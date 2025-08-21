<template>
    <div class="search-card">
        <div class="results-container">
            <div v-for="item in searchResult" :key="item.name" class="result-item" @click="openSearchResult(item)">
                <div class="result-header">
                    <div class="site-info">
                        <img v-if="item.siteIcon" :src="item.siteIcon" class="site-icon" alt="网站图标" />
                        <span class="site-name">{{ item.siteName }}</span>
                    </div>
                    <div class="date-info">{{ item.dateLastCrawled.split('T')[0] }}</div>
                </div>
                <h3 class="result-title">{{ item.name }}</h3>
                <p class="result-snippet">{{ item.snippet }}</p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch } from 'vue';
const props = defineProps({
    search: {
        type: String,
        required: true
    }
})

const searchResult = ref([])

// 初始解析搜索结果
const parseSearchResult = () => {
    try {
        searchResult.value = JSON.parse(props.search)
    } catch (error) {
        console.error('解析搜索结果失败:', error)
        searchResult.value = []
    }
}

// 初始化
parseSearchResult()

// 监听search属性变化，实时更新搜索结果
watch(() => props.search, (newValue) => {
    parseSearchResult()
}, { immediate: true })

const openSearchResult = (item) => {
    window.open(item.url, '_blank')
}

</script>

<style scoped>
.search-card {
    width: 100%;
    padding: 10px;
}

.results-container {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    justify-content: flex-start;
    margin: 0 auto;
}

.result-item {
    background-color: var(--background-color);
    border-radius: 8px;
    padding: 16px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    transition: all 0.3s ease;
    flex-grow: 1;
}

.result-item:hover {
    box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
}

.result-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
}

.site-info {
    display: flex;
    align-items: center;
    gap: 8px;
}

.site-icon {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    object-fit: cover;
}

.site-icon-placeholder {
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background-color: #e0e0e0;
}

.site-name {
    font-size: 12px;
}

.date-info {
    font-size: 12px;
}

.result-title {
    font-size: 16px;
    font-weight: 600;
    margin: 8px 0;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
}

.result-snippet {
    font-size: 14px;
    line-height: 1.5;
    margin: 0;
    overflow: hidden;
    text-overflow: ellipsis;
}

@media (max-width: 600px) {
    .result-item {
        width: 100%;
    }
}
</style>