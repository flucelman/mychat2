import { fetch as tauriFetch } from '@tauri-apps/plugin-http';
import { useGlobalSettingStore } from '@/stores/global_setting'

// 检查是否在 Tauri 环境中
const isTauri = () => {
    return window.__TAURI__ !== undefined;
}

// 获取 store 的辅助函数
const getGlobalSettingStore = () => {
    return useGlobalSettingStore();
}

// 统一的 fetch 方法
export const httpFetch = async (url, options = {}) => {
    if (isTauri()) {
        // 在 Tauri 环境中使用 Tauri 的 fetch
        return tauriFetch(url, options);
    } else {
        // 在浏览器环境中使用原生 fetch
        return window.fetch(url, options);
    }
}

// 处理响应的辅助函数
const handleResponse = async (response) => {
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
        const jsonData = await response.json();
        // 返回一个包含 data 属性的对象，保持与原有代码的兼容性
        return { data: jsonData };
    }
    
    return { data: await response.text() };
}

// 封装的 HTTP 请求方法
export const http = {
    async get(url, headers = {}, params = {}, options = {}) {
        const globalSettingStore = getGlobalSettingStore();
        const response = await httpFetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': globalSettingStore.userToken,
                ...headers
            },
            params: params,
            ...options
        });
        return handleResponse(response);
    },

    async post(url, data, headers = {}, options = {}) {
        const globalSettingStore = getGlobalSettingStore();
        const response = await httpFetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': globalSettingStore.userToken,
                ...headers
            },
            body: JSON.stringify(data),
            ...options
        });
        return handleResponse(response);
    },

    async put(url, data, headers = {}, options = {}) {
        const globalSettingStore = getGlobalSettingStore();
        const response = await httpFetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': globalSettingStore.userToken,
                ...headers
            },
            body: JSON.stringify(data),
            ...options
        });
        return handleResponse(response);
    },

    async delete(url, headers = {}, options = {}) {
        const globalSettingStore = getGlobalSettingStore();
        const response = await httpFetch(url, {
            method: 'DELETE',
            headers: {
                'Authorization': globalSettingStore.userToken,
                ...headers
            },
            ...options
        });
        return handleResponse(response);
    }
}
