-- 如果使用 Lua 5.1
package.cpath = package.cpath .. "; /usr/lib/x86_64-linux-gnu/lua/5.1/?.so"
-- 配置参数
local config = {
    base_url = "http://localhost:9090",  -- 修改为实际接口地址
    username = "u1",               -- 测试用户名
    password = "1",           -- 测试密码
    login_path = "/login",               -- 登录接口路径
    uplink_path = "/uplink",             -- 报价接口路径
    min_price = 10.0,
    max_price = 1000.0,
    min_quantity = 1,
    max_quantity = 100
}

-- 全局token变量
local token = nil

-- 确保已经安装了 cjson 库
local cjson = require("cjson")

-- 生成随机字符串作为边界
function generateBoundary()
    local chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    local length = 16
    local boundary = ""
    for i = 1, length do
        boundary = boundary .. string.sub(chars, math.random(1, #chars), math.random(1, #chars))
    end
    return "WebKitFormBoundary" .. boundary
end

-- 生成随机报价数据
function generateRandomQuote()
    local boundary = generateBoundary()
    local price = math.random(config.min_price * 100, config.max_price * 100) / 100
    local quantity = math.random(config.min_quantity, config.max_quantity)
    local isSeller = math.random(0, 1)  -- 0或1表示是否为卖家

    local formData = ""
    formData = formData .. "--" .. boundary .. "\r\n"
    formData = formData .. "Content-Disposition: form-data; name=\"arg1\"\r\n\r\n" .. price .. "\r\n"
    formData = formData .. "--" .. boundary .. "\r\n"
    formData = formData .. "Content-Disposition: form-data; name=\"arg2\"\r\n\r\n" .. quantity .. "\r\n"
    formData = formData .. "--" .. boundary .. "\r\n"
    formData = formData .. "Content-Disposition: form-data; name=\"arg3\"\r\n\r\n" .. isSeller .. "\r\n"
    formData = formData .. "--" .. boundary .. "--\r\n"

    wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary
    return formData
end

-- 登录获取token
function login()
    print("正在获取token...")
    local boundary = generateBoundary()
    -- 构建登录请求
    local login_data = "--" .. boundary .. "\r\n"
    login_data = login_data .. "Content-Disposition: form-data; name=\"username\"\r\n\r\n" .. config.username .. "\r\n"
    login_data = login_data .. "--" .. boundary .. "\r\n"
    login_data = login_data .. "Content-Disposition: form-data; name=\"password\"\r\n\r\n" .. config.password .. "\r\n"
    login_data = login_data .. "--" .. boundary .. "--\r\n"

    -- 发送登录请求
    local resp = wrk.format("POST", config.login_path, {
        ["Content-Type"] = "multipart/form-data; boundary=" .. boundary
    }, login_data)
    
    -- 执行请求并获取响应
    local _, status, headers, body = wrk.http(resp)
    
    -- 解析token (假设返回JSON格式: {"jwt": "your_token_here"})
    if status == 200 then
        local success, result = pcall(cjson.decode, body)
        if success and result.jwt then
            token = result.jwt
            print("成功获取token:", token)
            return true
        end
    end
    
    print("获取token失败，状态码:", status)
    print("响应内容:", body)
    return false
end

-- 初始化测试
function setup(thread)
    -- 只在第一个线程中获取token（避免并发问题）
    if thread.id == 1 then
        if not login() then
            error("无法获取token，测试终止")
        end
    end
    
    -- 等待token获取完成
    local wait_time = 1
    local start_time = os.time()
    while os.time() - start_time < wait_time do
        -- 空循环等待
    end
    
    -- 设置线程特定的计数器
    thread:set("request_count", 0)
    thread:set("success_count", 0)
    thread:set("error_count", 0)
    
    --print(string.format("Thread %d started with token: %s", thread.id, token))
end

-- 每个请求前执行
function init(args)
    -- 生成随机报价数据
    wrk.body = generateRandomQuote()
    
    -- 设置请求头，包含token
    wrk.headers["Content-Type"] = "multipart/form-data"
    if token then
        wrk.headers["Authorization"] = "Bearer " .. token
    else
        print("Token is nil, skipping request")
        return nil -- 可以选择跳过请求
    end
    wrk.headers["Content-Length"] = #wrk.body
end

-- 处理响应
function response(status, headers, body)
    -- 获取当前线程计数器
    local request_count = tonumber(wrk.thread:get("request_count")) + 1
    local success_count = tonumber(wrk.thread:get("success_count"))
    local error_count = tonumber(wrk.thread:get("error_count"))
    
    -- 更新请求计数
    wrk.thread:set("request_count", request_count)
    
    -- 验证响应状态码
    if status == 200 then
        -- 检查响应内容是否包含成功标志
        if string.find(body, "success") then
            success_count = success_count + 1
        else
            error_count = error_count + 1
            print(string.format("Unexpected response: %s", body))
        end
    else
        error_count = error_count + 1
        print(string.format("Request failed with status: %d", status))
    end
    
    -- 更新成功/失败计数
    wrk.thread:set("success_count", success_count)
    wrk.thread:set("error_count", error_count)
    
    return status, headers, body
end

-- 测试结束后执行
function done(summary, latency, requests)
    -- 计算总请求数和成功率
    local total_requests = 0
    local total_success = 0
    local total_errors = 0
    
    for i = 1, wrk.threads do
        local thread = wrk.threads[i]
        total_requests = total_requests + tonumber(thread:get("request_count"))
        total_success = total_success + tonumber(thread:get("success_count"))
        total_errors = total_errors + tonumber(thread:get("error_count"))
    end
    
    -- 计算成功率
    local success_rate = (total_success / total_requests) * 100
    
    -- 输出详细统计信息
    print("\n=== Test Summary ===")
    print(string.format("Total Requests: %d", total_requests))
    print(string.format("Successful Requests: %d", total_success))
    print(string.format("Failed Requests: %d", total_errors))
    print(string.format("Success Rate: %.2f%%", success_rate))
    print(string.format("Total Time: %.2fms", summary.duration))
    print(string.format("Requests/sec: %.2f", summary.requests/summary.duration*1000))
    print(string.format("Latency (Avg): %.2fms", latency.avg))
    print(string.format("Latency (Max): %.2fms", latency.max))
    print(string.format("Latency (Stdev): %.2fms", latency.stdev))
    
    -- 输出百分位延迟
    print("\n=== Latency Distribution ===")
    print(string.format("  50%%: %.2fms", latency:percentile(50)))
    print(string.format("  75%%: %.2fms", latency:percentile(75)))
    print(string.format("  90%%: %.2fms", latency:percentile(90)))
    print(string.format("  99%%: %.2fms", latency:percentile(99)))
end