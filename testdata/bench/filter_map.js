var records = [];
var n = 100000;

for (var i = 0; i < n; i++) {
    records.push({
        id: i,
        value: (i * 19) % 500,
        active: i % 3 === 0,
    });
}

var result = [];

for (var i = 0; i < records.length; i++) {
    var r = records[i];
    if (!r.active) {
        continue;
    }
    if (r.value > 100) {
        result.push({
            id: r.id,
            score: r.value * 2,
        });
    }
}

var out = result.length;
console.log(out);
