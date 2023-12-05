local biz = KEYS[1]

local action = ARGV[1]

-- in case there's cancel
local delta = tonumber(ARGV[2])

local exists = redis.call("EXISTS", biz)
if exists == 1 then
    redis.call("HINCRBY", action, delta)
    return 1
else
    return 0
end