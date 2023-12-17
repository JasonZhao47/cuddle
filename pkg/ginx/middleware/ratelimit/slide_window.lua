--    1. 记录每个请求来的timestamp，每个客户端对应一个key
--    2. 只保留window内的所有访问数量
--    3. 当前时间，判断如果window内访问数量，超限了就返回1，否则返回0
local key = KEY[1]
--
local window = tonumber(ARGV[1])
local threshold = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local windowStart = now - window

-- 只保留window内的所有访问数量
redis.call('ZREMRANGEBYSCORE', key, '-inf', windowStart)
local cnt = redis.call('ZCOUNT', key, windowStart, '+inf')

if cnt >= threshold then
    return "true"
else
    redis.call("ZADD", key, now, now)
    redis.call("PEXPIRE", key, window)
    return "false"
end