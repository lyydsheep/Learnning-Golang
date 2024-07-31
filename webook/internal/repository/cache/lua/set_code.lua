-- key，比如：code:login:17812341234
local key = KEYS[1]
-- val，比如：123456
local val = ARGV[1]
-- 验证码可以验证的次数
local keyCnt = key.."cnt"
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    -- 没有expiration，说明设置错误
    return -2
elseif ttl == -2 or ttl < 540 then
    -- 没有这个key
    -- 或者val快要过期了，可以重发一次

    -- 设置k-v
    redis.call("set", key, val, "EX", 600)
    -- 设置可以验证的次数
    redis.call("set", keyCnt, 3, "EX", 600)
    return 0
else
    return -1
end