package redisstore

import goredis "github.com/redis/go-redis/v9"

var enqueueScript = goredis.NewScript(`
if redis.call("SADD", KEYS[1], ARGV[1]) == 1 then
  redis.call("LPUSH", KEYS[2], ARGV[2])
  return 1
end
return 0
`)

var countScript = goredis.NewScript(`
return {redis.call("LLEN", KEYS[1]), redis.call("LLEN", KEYS[2])}
`)

var recoverScript = goredis.NewScript(`
local jobs = redis.call("LRANGE", KEYS[1], 0, -1)
for _, job in ipairs(jobs) do
  redis.call("LPUSH", KEYS[2], job)
end
redis.call("DEL", KEYS[1])
return #jobs
`)
