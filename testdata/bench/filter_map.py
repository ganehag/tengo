records = []
n = 100000

for i in range(n):
    records.append({
        "id": i,
        "value": (i * 19) % 500,
        "active": i % 3 == 0,
    })

result = []

for r in records:
    if not r["active"]:
        continue
    if r["value"] > 100:
        result.append({
            "id": r["id"],
            "score": r["value"] * 2,
        })

out = len(result)
print(out)
