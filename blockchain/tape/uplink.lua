-- 设置随机数种子，基于当前时间确保每次测试生成不同的随机值
math.randomseed(os.time())

-- 生成唯一的 offerId（时间戳+随机数）
local timestamp = os.time()
local random_id = math.random(1000, 9999)
local offer_id = string.format("OFFER-%d-%d", timestamp, random_id)

-- 生成随机价格 (0-1)
local price = math.random() 
price = math.floor(price * 100) / 100  -- 保留两位小数

-- 生成随机电量 (10-200 单位)
local quantity = math.random(0, 200)

-- 随机生成用户角色 (true 表示卖方，false 表示买方)
local is_seller = (math.random(1, 2) == 1)

-- 构建 multipart/form-data 请求体
local boundary = "---------------------------1234567890"
local body = "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"userId\"\r\n\r\n1929768823477506048\r\n" ..
             "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"offerId\"\r\n\r\n" .. offer_id .. "\r\n" ..
             "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"price\"\r\n\r\n" .. price .. "\r\n" ..
             "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"quantity\"\r\n\r\n" .. quantity .. "\r\n" ..
             "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"isSeller\"\r\n\r\n" .. tostring(is_seller) .. "\r\n" ..
             "--" .. boundary .. "\r\n" ..
             "Content-Disposition: form-data; name=\"adminId\"\r\n\r\n1917140260131704832\r\n" ..
             "--" .. boundary .. "--"

-- 设置请求参数
wrk.method = "POST"
wrk.path = "/uplink"
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary
wrk.body = body