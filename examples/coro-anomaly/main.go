// Demonstrates a four-stage coroutine pipeline that detects anomalies in a
// stream of sensor readings using a sliding-window average.
//
//	go run ./examples/coro-anomaly
package main

import (
	"fmt"
	"os"

	"github.com/ganehag/tengo/v3"
	"github.com/ganehag/tengo/v3/stdlib"
)

const script = `
coro := import("coro")
fmt  := import("fmt")

// Stage 1: raw sensor readings — index 7 is a spike
source := func(yield) {
    data := [10, 11, 10, 12, 11, 10, 11, 45, 11, 10, 12, 11]
    for _, v in data {
        yield(v)
    }
}

// Stage 2: sliding-window average over the last n values
//   emits {raw, avg} maps; maintains its own rolling buffer
sliding_avg := func(yield, src, n) {
    win := []
    for {
        v, ok := src.resume()
        if !ok { break }
        win = append(win, v)
        if len(win) > n {
            win = win[1:]
        }
        sum := 0
        for _, w in win { sum += w }
        yield({raw: v, avg: sum / len(win)})
    }
}

// Stage 3: absolute deviation from the sliding average
deviation := func(yield, src) {
    for {
        rec, ok := src.resume()
        if !ok { break }
        d := rec.raw - rec.avg
        if d < 0 { d = -d }
        yield({raw: rec.raw, avg: rec.avg, dev: d})
    }
}

// Stage 4: suppress normal readings; only pass through anomalies
anomalies := func(yield, src, threshold) {
    for {
        rec, ok := src.resume()
        if !ok { break }
        if rec.dev > threshold {
            yield(rec)
        }
    }
}

src  := coro.new(source)
avg  := coro.new(sliding_avg, src, 4)
dev  := coro.new(deviation,  avg)
anom := coro.new(anomalies,  dev, 5)

for rec in anom {
    fmt.println("anomaly  raw:", rec.raw, " avg:", rec.avg, " dev:", rec.dev)
}
anom.close()
`

func main() {
	s := tengo.NewScript([]byte(script))
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))

	if _, err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
