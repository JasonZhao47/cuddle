-- 设置一个6位code，让用户每次获取这个code作为验证
-- 这个key一定要有过期时间
-- 设置好后发给用户
-- 设置后的一分钟内不能重复发给用户
-- key: code:biz:phoneNumber
-- 检测key code的失效时间
-- 防止有人误操作

-- 获得keys
local key = KEYS[1]
-- 还能使用几次
-- ..是lua的字符串拼接
local cntKey = key .. ":cnt"
-- 获得value
local val = ARGV[1]

local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- 有key但没有时间，错误数据
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 设置后还没满1分钟
    return -1
end