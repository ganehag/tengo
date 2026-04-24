local n = 200000
local window = 10

local data = {}
for i = 0, n - 1 do
    table.insert(data, (i * 13) % 1000)
end

local sum = 0
for i = 1, window do
    sum = sum + data[i]
end

local out = 0
for i = window + 1, n do
    sum = sum + data[i]
    sum = sum - data[i - window]
    out = out + math.floor(sum / window)
end

print(out)
