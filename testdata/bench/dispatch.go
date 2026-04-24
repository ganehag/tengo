package main

import "fmt"

func main() {
	handlers := map[string]func(int, int) int{
		"add": func(x, y int) int { return x + y },
		"mul": func(x, y int) int { return x * y },
		"sub": func(x, y int) int { return x - y },
	}

	type event struct {
		op   string
		a, b int
	}

	n := 200000
	events := make([]event, n)
	for i := 0; i < n; i++ {
		var op string
		if i%3 == 0 {
			op = "add"
		} else if i%3 == 1 {
			op = "mul"
		} else {
			op = "sub"
		}
		events[i] = event{op: op, a: i, b: i % 7}
	}

	out := 0
	for _, e := range events {
		out += handlers[e.op](e.a, e.b)
	}
	fmt.Println(out)
}
