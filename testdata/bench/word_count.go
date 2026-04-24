package main

import (
	"fmt"
	"strings"
)

func main() {
	n := 20000
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("lorem ipsum dolor sit amet consectetur adipiscing elit ")
	}
	text := sb.String()

	words := strings.Split(text, " ")
	counts := make(map[string]int)
	for _, w := range words {
		if w == "" {
			continue
		}
		counts[w]++
	}

	out := 0
	for _, v := range counts {
		out += v
	}
	fmt.Println(out)
}
