n = 200000
window = 10

data = []
for i in range(n):
    data.append((i * 13) % 1000)

sum_val = 0
for i in range(window):
    sum_val += data[i]

out = 0
for i in range(window, n):
    sum_val += data[i]
    sum_val -= data[i-window]
    out += sum_val // window

print(out)
