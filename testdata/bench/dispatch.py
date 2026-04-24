handlers = {
    "add": lambda x, y: x + y,
    "mul": lambda x, y: x * y,
    "sub": lambda x, y: x - y,
}

events = []
n = 200000

for i in range(n):
    op = "add" if i % 3 == 0 else ("mul" if i % 3 == 1 else "sub")
    events.append({
        "op": op,
        "a": i,
        "b": i % 7,
    })

out = 0
for e in events:
    fn = handlers[e["op"]]
    out += fn(e["a"], e["b"])

print(out)
