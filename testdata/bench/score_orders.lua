local function score(order)
    local total = order.total
    local items = order.items
    local country = order.country
    local priority = order.priority

    local s = 0

    if total > 1000 then
        s = s + 50
    elseif total > 250 then
        s = s + 20
    else
        s = s + 5
    end

    if items > 10 then
        s = s + items * 2
    else
        s = s + items
    end

    if country == "SE" then
        s = s + 7
    elseif country == "US" then
        s = s + 5
    else
        s = s + 3
    end

    if priority then
        s = s * 2
    end

    return s
end

local orders = {}
local n = 100000

for i = 0, n - 1 do
    local country = "DE"
    if i % 3 == 0 then
        country = "SE"
    elseif i % 3 == 1 then
        country = "US"
    end
    table.insert(orders, {
        total = (i * 37) % 1500,
        items = (i % 17) + 1,
        country = country,
        priority = (i % 11 == 0)
    })
end

local out = 0
for _, order in ipairs(orders) do
    out = out + score(order)
end

print(out)
