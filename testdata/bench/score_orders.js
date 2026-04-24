function score(order) {
    var total = order.total;
    var items = order.items;
    var country = order.country;
    var priority = order.priority;

    var s = 0;

    if (total > 1000) {
        s += 50;
    } else if (total > 250) {
        s += 20;
    } else {
        s += 5;
    }

    if (items > 10) {
        s += items * 2;
    } else {
        s += items;
    }

    if (country === "SE") {
        s += 7;
    } else if (country === "US") {
        s += 5;
    } else {
        s += 3;
    }

    if (priority) {
        s *= 2;
    }

    return s;
}

var orders = [];
var n = 100000;

for (var i = 0; i < n; i++) {
    orders.push({
        total: (i * 37) % 1500,
        items: (i % 17) + 1,
        country: i % 3 === 0 ? "SE" : (i % 3 === 1 ? "US" : "DE"),
        priority: i % 11 === 0,
    });
}

var out = 0;
for (var j = 0; j < orders.length; j++) {
    out += score(orders[j]);
}

console.log(out);
