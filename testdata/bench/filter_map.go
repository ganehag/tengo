package main

import "fmt"

type record struct {
	id     int
	value  int
	active bool
}

type scored struct {
	id    int
	score int
}

func main() {
	n := 100000
	records := make([]record, n)
	for i := 0; i < n; i++ {
		records[i] = record{
			id:     i,
			value:  (i * 19) % 500,
			active: i%3 == 0,
		}
	}

	var result []scored
	for _, r := range records {
		if !r.active {
			continue
		}
		if r.value > 100 {
			result = append(result, scored{id: r.id, score: r.value * 2})
		}
	}

	fmt.Println(len(result))
}
