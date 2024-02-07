local key = KEYS[1]
local cntKey =key.. ":cnt"
--  key 和 ":cnt" 连接
-- 你准备存储的验证码
local  val = ARGV[1]

--  tonumber what   redis.call("command",object)
local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
-- key  存在   但是没有过期时间
 return -2       --  给go语言返回的-2

elseif ttl == -2  or ttl < 600 then
--     可以发验证码
   redis.call("set",key,val)
   --600 秒后过期
   redis.call("expire",key,600)
--    只能验证三次
-- 测试的时候避免流程中断要改改
   redis.call("set",cntKey,3)
--    这里已经设置了过期时间十分钟
   redis.call("expire",cntKey,600)
   return 0
else
--    发送太频繁
   return -1
end