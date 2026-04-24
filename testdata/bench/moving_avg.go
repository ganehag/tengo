package main

import "fmt"

func main() {
	n := 200000
	window := 10

	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = (i * 13) % 1000
	}

	sum := 0
	for i := 0; i < window; i++ {
		sum += data[i]
	}

	out := 0
	for i := window; i < n; i++ {
		sum += data[i]
		sum -= data[i-window]
		out += sum / window
	}
	fmt.Println(out)
}
