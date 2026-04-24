var handlers = {
    "add": function(x, y) { return x + y; },
    "mul": function(x, y) { return x * y; },
    "sub": function(x, y) { return x - y; },
};

var events = [];
var n = 200000;

for (var i = 0; i < n; i++) {
    var op = i % 3 === 0 ? "add" : (i % 3 === 1 ? "mul" : "sub");
    events.push({
        "op": op,
        "a": i,
        "b": i % 7,
    });
}

var out = 0;
for (var i = 0; i < events.length; i++) {
    var e = events[i];
    var fn = handlers[e.op];
    out += fn(e.a, e.b);
}

console.log(out);
