-- login.lua
request = function()
    -- 替换为实际的登录数据
    local form_data = "username=admin&password=123"
    return wrk.format("POST", "/login", { ["Content-Type"] = "multipart/form-data" }, form_data)
end