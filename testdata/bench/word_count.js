var text = "";
var n = 20000;

for (var i = 0; i < n; i++) {
    text += "lorem ipsum dolor sit amet consectetur adipiscing elit ";
}

var words = text.split(" ");
var counts = {};

for (var i = 0; i < words.length; i++) {
    var w = words[i];
    if (w === "") {
        continue;
    }
    if (counts[w] === undefined) {
        counts[w] = 1;
    } else {
        counts[w] += 1;
    }
}

var out = 0;
for (var k in counts) {
    out += counts[k];
}

console.log(out);
