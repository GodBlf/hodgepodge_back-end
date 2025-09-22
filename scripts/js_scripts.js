(function poll() {
    const url = 'http://10.242.194.32:8080/r';

    async function getData() {
        try {
            const response = await fetch(url, {
                method: 'GET',
                // 如果需要跨域携带凭证，可以添加 credentials
                // credentials: 'include'
            });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const text = await response.text();
            console.log(`[${new Date().toLocaleTimeString()}] 返回数据:`, text);
        } catch (err) {
            console.error('请求失败:', err);
        }
    }

    // 先立即执行一次
    getData();
    // 每隔 5 秒执行一次
    setInterval(getData, 1000);
})();
