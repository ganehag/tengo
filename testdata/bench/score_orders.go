package main

import "fmt"

type order struct {
	total    int
	items    int
	country  string
	priority bool
}

func score(o order) int {
	s := 0
	if o.total > 1000 {
		s += 50
	} else if o.total > 250 {
		s += 20
	} else {
		s += 5
	}
	if o.items > 10 {
		s += o.items * 2
	} else {
		s += o.items
	}
	if o.country == "SE" {
		s += 7
	} else if o.country == "US" {
		s += 5
	} else {
		s += 3
	}
	if o.priority {
		s *= 2
	}
	return s
}

func main() {
	n := 100000
	orders := make([]order, n)
	for i := 0; i < n; i++ {
		var country string
		if i%3 == 0 {
			country = "SE"
		} else if i%3 == 1 {
			country = "US"
		} else {
			country = "DE"
		}
		orders[i] = order{
			total:    (i * 37) % 1500,
			items:    (i % 17) + 1,
			country:  country,
			priority: i%11 == 0,
		}
	}
	out := 0
	for _, o := range orders {
		out += score(o)
	}
	fmt.Println(out)
}
