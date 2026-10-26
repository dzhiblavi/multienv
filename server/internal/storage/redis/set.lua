local path = '$'

for i = 3, #ARGV, 2 do
    local types = redis.call('JSON.TYPE', KEYS[1], path)
    if not types or not types[1] then
        local value = ARGV[2]
        for j = #ARGV - 1, i, -2 do
            value = '{' .. ARGV[j] .. ':' .. value .. '}'
        end

        redis.call('JSON.SET', KEYS[1], path, value)
        return 1
    end

    if types[1] ~= 'object' then
        return redis.error_reply('storage parent is not an object')
    end

    path = path .. '[' .. ARGV[i + 1] .. ']'
end

local result
if ARGV[1] == 'NX' then
    result = redis.call('JSON.SET', KEYS[1], path, ARGV[2], 'NX')
else
    result = redis.call('JSON.SET', KEYS[1], path, ARGV[2])
end

if result then
    return 1
end

return 0
