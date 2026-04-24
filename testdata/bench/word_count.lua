local text = ""
local n = 20000

for i = 0, n - 1 do
    text = text .. "lorem ipsum dolor sit amet consectetur adipiscing elit "
end

local counts = {}

-- Plain-search split (go-lua doesn't support %S+ pattern in string.find)
local pos = 1
local len = #text
while pos <= len do
    local sp = string.find(text, " ", pos, true)
    if not sp then sp = len + 1 end
    if sp > pos then
        local w = string.sub(text, pos, sp - 1)
        counts[w] = (counts[w] or 0) + 1
    end
    pos = sp + 1
end

local out = 0
for k, v in pairs(counts) do
    out = out + v
end

print(out)
