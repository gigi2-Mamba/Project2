-- 滑动窗口限流。  达到阈值要限流。   窗口大小  窗口起始时间戳。
-- KEY 作为限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口起始时间戳
local min = now - window

redis.call('ZREMRANGEBYSCORE',key,'-inf',min)
-- infinite:无穷的
local cnt = redis.call('ZCOUNT',key,'-inf','+inf')

if cnt >= threshold then
--   限流
    return "true"
else
--   怎么要把score 和 member 设置为now    为了唯一性。
    redis.call('ZADD',key,now,now)
    -- 在go里传入就做了毫秒处理就为了适应这个命令还是追求实时性，两者皆有。实时在这里约等于唯一。保持的zset键不重复。
    redis.call('PEXPIRE',key,window)
    return "false"
end


