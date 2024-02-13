--一个增加计数的lua，用于点赞，阅读，收藏
--使用HINCRBY   如果存在就更新。 不存在key或者field就初始化并设置0值

-- 具体业务
local key = KEYS[1]

--具体的原子业务
--阅读数/收藏数/点赞数
local cntKey = ARGV[1]

local delta = tonumber(ARGV[2])

local exist = redis.call("EXISTS",key)


--
-- return redis.call("HINCRBY",key,cntKey,delta)
if exist ==1 then
    redis.call("HINCRBY",key,cntKey,delta)
    return 1
else
    return 0
end




