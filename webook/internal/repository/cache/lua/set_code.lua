-- code:biz:phone
local key = KEYS[1]
local cntKey = key..":cnt"
-- 为什么这里是argv？
local val =ARGV[1]

local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
-- key exists, but no expiration
    return -2
elseif ttl==-2 or ttl<540 then
-- can send sms
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- send too many
    return -1 

end 