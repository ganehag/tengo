package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type benchEntry struct {
	name  string
	label string
}

type langEntry struct {
	name   string
	link   string
	binary string
	ext    string
	flags  []string
}

var benches = []benchEntry{
	{name: "fib", label: "fib(35)"},
	{name: "score_orders", label: "score orders"},
	{name: "word_count", label: "word count"},
	{name: "moving_avg", label: "moving avg"},
	{name: "filter_map", label: "filter/map"},
	{name: "dispatch", label: "dispatch"},
}

var languages = []langEntry{
	{name: "**Tengo**", link: "https://github.com/tengolang/tengo", binary: "tengo", ext: "tengo"},
	{name: "go-lua", link: "https://github.com/Shopify/go-lua", binary: "go-lua", ext: "lua"},
	{name: "GopherLua", link: "https://github.com/yuin/gopher-lua", binary: "glua", ext: "lua"},
	{name: "goja", link: "https://github.com/dop251/goja", binary: "goja", ext: "js"},
	{name: "starlark-go", link: "https://github.com/google/starlark-go", binary: "starlark", ext: "star", flags: []string{"-recursion"}},
	{name: "gpython", link: "https://github.com/go-python/gpython", binary: "gpython", ext: "py"},
	{name: "Yaegi", link: "https://github.com/traefik/yaegi", binary: "yaegi", ext: "yaegi", flags: []string{"run"}},
	{name: "otto", link: "https://github.com/robertkrimen/otto", binary: "otto", ext: "js"},
	{name: "Anko", link: "https://github.com/mattn/anko", binary: "anko", ext: "ank"},
	{name: "-"},
	{name: "Go", binary: "go", ext: "go", flags: []string{"run"}},
	{name: "Lua", binary: "lua", ext: "lua"},
	{name: "Python 2", binary: "python2.7", ext: "py"},
	{name: "Python 3", binary: "python3", ext: "py"},
}

func main() {
	// n := flag.Int("n", 35, "fibonacci input (detailed mode)")
	count := flag.Int("count", 3, "number of runs (minimum is reported)")
	timeout := flag.Duration("timeout", 60*time.Second, "per-run timeout")
	markdown := flag.Bool("markdown", false, "print comparison table in markdown")
	dir := flag.String("dir", "testdata/bench", "directory containing benchmark scripts")
	flag.Parse()

	if *markdown {
		printMarkdown(*dir, *count, *timeout)
		return
	}
}
func printMarkdown(dir string, count int, timeout time.Duration) {
	type benchResult struct {
		duration time.Duration
		ok       bool
	}

	type result struct {
		times []benchResult
		skip  bool
	}

	results := make([]result, len(languages))

	for i, lang := range languages {
		if lang.name == "-" {
			continue
		}

		if _, err := exec.LookPath(lang.binary); err != nil {
			results[i].skip = true
			continue
		}

		r := result{
			times: make([]benchResult, len(benches)),
		}

		for j, bench := range benches {
			path := scriptPath(lang, dir, bench.name)
			if _, err := os.Stat(path); err != nil {
				r.times[j] = benchResult{ok: false}
				continue
			}

			d, ok := measureExternal(buildCmd(lang, dir, bench.name), count, timeout)
			r.times[j] = benchResult{
				duration: d,
				ok:       ok,
			}
		}

		results[i] = r
	}

	fmt.Print("| |")
	for _, bench := range benches {
		fmt.Printf(" %s |", bench.label)
	}
	fmt.Println()

	fmt.Print("| :--- |")
	for range benches {
		fmt.Print(" ---: |")
	}
	fmt.Println()

	for i, lang := range languages {
		if lang.name == "-" {
			fmt.Print("| - |")
			for range benches {
				fmt.Print(" - |")
			}
			fmt.Println()
			continue
		}

		if results[i].skip {
			continue
		}

		nameCol := lang.name
		if lang.link != "" {
			nameCol = fmt.Sprintf("[%s](%s)", lang.name, lang.link)
		}

		fmt.Printf("| %s |", nameCol)
		for _, t := range results[i].times {
			if !t.ok {
				fmt.Print(" - |")
				continue
			}
			fmt.Printf(" `%sms` |", formatMs(t.duration))
		}
		fmt.Println()
	}
}

func scriptPath(lang langEntry, dir, base string) string {
	return filepath.Join(dir, base+"."+lang.ext)
}

func buildCmd(lang langEntry, dir, base string) []string {
	script := scriptPath(lang, dir, base)

	args := make([]string, 0, 2+len(lang.flags))
	args = append(args, lang.binary)
	args = append(args, lang.flags...)
	args = append(args, script)
	return args
}


func measureExternal(args []string, count int, timeout time.Duration) (time.Duration, bool) {
	min := time.Duration(math.MaxInt64)
	ok := false

	for i := 0; i < count; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)

		start := time.Now()
		err := cmd.Run()
		d := time.Since(start)
		cancel()

		if err != nil {
			continue
		}

		ok = true
		if d < min {
			min = d
		}
	}

	if !ok {
		return 0, false
	}
	return min, true
}

func formatMs(d time.Duration) string {
	if d == math.MaxInt64 {
		return "?"
	}
	n := (d.Nanoseconds() + 500_000) / 1_000_000
	if n == 0 {
		return "0"
	}
	neg := n < 0
	a := n
	if neg {
		a = -a
	}
	digits := int(math.Log10(float64(a))) + 1
	outLen := digits + (digits-1)/3
	if neg {
		outLen++
	}
	out := make([]rune, outLen)
	var i, j int
	for a > 0 {
		out[outLen-j-1] = '0' + rune(a%10)
		i++
		j++
		if i%3 == 0 && j < outLen {
			out[outLen-j-1] = ','
			j++
		}
		a /= 10
	}
	if neg {
		out[0] = '-'
	}
	return string(out)
}
