wrk.method = "POST"
wrk.path = "/getOfferHistory"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
-- 请将 YOUR_USER_ID 替换为实际的用户 ID
wrk.body = "userId=1929768823477506048"