var n = 200000;
var w = 10;

var data = [];
for (var i = 0; i < n; i++) {
    data.push((i * 13) % 1000);
}

var sum = 0;
for (var i = 0; i < w; i++) {
    sum += data[i];
}

var out = 0;
for (var i = w; i < n; i++) {
    sum += data[i];
    sum -= data[i - w];
    out += Math.floor(sum / w);
}

console.log(out);
