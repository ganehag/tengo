local handlers = {
    add = function(x, y) return x + y end,
    mul = function(x, y) return x * y end,
    sub = function(x, y) return x - y end,
}

local events = {}
local n = 200000

for i = 0, n - 1 do
    local op = "sub"
    if i % 3 == 0 then
        op = "add"
    elseif i % 3 == 1 then
        op = "mul"
    end
    table.insert(events, {
        op = op,
        a = i,
        b = i % 7,
    })
end

local out = 0
for _, e in ipairs(events) do
    local fn = handlers[e.op]
    out = out + fn(e.a, e.b)
end

print(out)
