-- 拿到key
-- 判断验证次数
-- 如果可以，就放行
-- 如果key的时间不够了，也不行
local key = KEYS[1]
local cntKey = key .. ":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

if cnt == nil or cnt <= 0 then
    return -1
end

if code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("desc", cntKey)
    return -2
end