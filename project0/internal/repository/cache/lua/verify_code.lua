-- 键，  键:cnt 验证次数  键为:服务-手机号
local key = KEYS[1]
--  牛逼语法  拼接构成新的key    一个local key 和 ..凭借 “字符”
local cntKey = key.. ":cnt"

local expectCode = ARGV[1]

local code = redis.call("get",key)
local cnt = tonumber(redis.call("get",cntKey))


--  验证验证码， 用到验证码次数统计的key
if cnt == nil or  cnt <= 0 then
--验证次数耗尽了
 return -1
end

if  code == expectCode  then
--    把验证码标记为0,表示不可用
   redis.call("set",cntKey,0)
   return 0
else
  redis.call("decr",cntKey)
--   不相等用户输错了
  return  -2
end
