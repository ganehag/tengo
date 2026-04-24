local function fib(x, a, b)
    if x == 0 then return a
    elseif x == 1 then return b
    else return fib(x-1, b, a+b)
    end
end
print(fib(35, 0, 1))
