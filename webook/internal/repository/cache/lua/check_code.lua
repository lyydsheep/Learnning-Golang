-- 验证码的key
local key = KEYS[1]
-- 剩余验证次数的key
local keyCnt = key.."cnt"
-- 用户输入
local inputCode = ARGV[1]
-- 获取验证码和剩余验证次数
local val = redis.call("get", key)
local cnt = redis.call("get", keyCnt)

if val == nil then
    -- 系统错误
    return -3
elseif cnt <= 0 then
    -- 有人攻击
    return -2
elseif val == inputCode then
    -- 验证成功，那么这个验证码就失效了，直到被重置
    redis.call("set", keyCnt, -1)
    return 0
end
-- 手误，输入错了
redis.call("decr", keyCnt)
return -1