-- 具体业务
local key = KEYS[1]
-- 是阅读数，点赞数还是收藏数
local cntKey = ARGV[1]

-- 要增加或减少的数量
local delta = tonumber(ARGV[2])
if not delta then
    -- 如果转换失败，返回-1表示输入参数错误
    return -1
end

-- 检查键是否存在
local exist = redis.call("EXISTS", key)
if exist == 1 then
    -- 如果存在，根据cntKey增加或减少delta
    local status, err = pcall(function()
        return redis.call("HINCRBY", key, cntKey, delta)
    end)

    if not status then
        -- 处理Redis调用中的错误
        return -2
    end

    -- 返回1表示操作成功
    return 1
else
    -- 返回0表示键不存在，操作失败
    return 0
end
