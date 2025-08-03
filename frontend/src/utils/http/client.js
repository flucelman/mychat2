import { fetch as tauriFetch } from '@tauri-apps/plugin-http';

// 检查是否在 Tauri 环境中
const isTauri = () => {
    return window.__TAURI__ !== undefined;
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

// 封装的 HTTP 请求方法
export const http = {
    async get(url, options = {}) {
        return httpFetch(url, {
            method: 'GET',
            ...options
        });
    },

    async post(url, data, options = {}) {
        return httpFetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            body: JSON.stringify(data),
            ...options
        });
    },

    async put(url, data, options = {}) {
        return httpFetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            body: JSON.stringify(data),
            ...options
        });
    },

    async delete(url, options = {}) {
        return httpFetch(url, {
            method: 'DELETE',
            ...options
        });
    }
}
