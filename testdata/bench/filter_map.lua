local records = {}
local n = 100000

for i = 0, n - 1 do
    table.insert(records, {
        id = i,
        value = (i * 19) % 500,
        active = i % 3 == 0,
    })
end

local result = {}

for _, r in ipairs(records) do
    if r.active then
        if r.value > 100 then
            table.insert(result, {
                id = r.id,
                score = r.value * 2,
            })
        end
    end
end

local out = #result
print(out)
