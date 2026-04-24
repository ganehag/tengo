package stdlib_test

import (
	"testing"

	"github.com/tengolang/tengo/v3"
	"github.com/tengolang/tengo/v3/require"
	"github.com/tengolang/tengo/v3/stdlib"
)

// runCoroVar compiles and runs src with the coro module, returns the named variable.
func runCoroVar(t *testing.T, src, name string) interface{} {
	t.Helper()
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	err = compiled.Run()
	require.NoError(t, err)
	return compiled.Get(name).Value()
}

func TestCoroBasicGenerator(t *testing.T) {
	// Simple generator: yields 1, 2, 3 then finishes.
	src := `
coro := import("coro")
gen := func(yield) {
    yield(1)
    yield(2)
    yield(3)
}
co := coro.new(gen)
v1, ok1 := co.resume()
v2, ok2 := co.resume()
v3, ok3 := co.resume()
v4, ok4 := co.resume()
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())

	require.Equal(t, int64(1), compiled.Get("v1").Int64())
	require.Equal(t, true, compiled.Get("ok1").Bool())
	require.Equal(t, int64(2), compiled.Get("v2").Int64())
	require.Equal(t, true, compiled.Get("ok2").Bool())
	require.Equal(t, int64(3), compiled.Get("v3").Int64())
	require.Equal(t, true, compiled.Get("ok3").Bool())
	require.Equal(t, false, compiled.Get("ok4").Bool()) // coroutine dead
}

func TestCoroForIn(t *testing.T) {
	// Coroutine used as an iterable in for-in.
	src := `
coro := import("coro")
gen := func(yield) {
    for i := 1; i <= 5; i++ {
        yield(i)
    }
}
sum := 0
for v in coro.new(gen) {
    sum += v
}
`
	sum := runCoroVar(t, src, "sum")
	require.Equal(t, int64(15), sum)
}

func TestCoroWithArgs(t *testing.T) {
	// Coroutine function with user arguments.
	src := `
coro := import("coro")
counter := func(yield, start, step) {
    i := start
    for {
        yield(i)
        i += step
        if i > 10 { break }
    }
}
co := coro.new(counter, 2, 3)
sum := 0
for v in co {
    sum += v
}
`
	sum := runCoroVar(t, src, "sum")
	// yields: 2, 5, 8 → sum = 15
	require.Equal(t, int64(15), sum)
}

func TestCoroStatus(t *testing.T) {
	src := `
coro := import("coro")
gen := func(yield) { yield(1) }
co := coro.new(gen)
s1 := co.status
_, _ := co.resume()
s2 := co.status
_, _ := co.resume()
s3 := co.status
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())

	require.Equal(t, "suspended", compiled.Get("s1").String())
	require.Equal(t, "suspended", compiled.Get("s2").String())
	require.Equal(t, "dead", compiled.Get("s3").String())
}

func TestCoroClose(t *testing.T) {
	// Close() while coroutine is suspended stops it.
	src := `
coro := import("coro")
gen := func(yield) {
    yield(1)
    yield(2)
    yield(3)
}
co := coro.new(gen)
v1, _ := co.resume()
co.close()
_, ok := co.resume()  // dead after close
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())

	require.Equal(t, int64(1), compiled.Get("v1").Int64())
	require.Equal(t, false, compiled.Get("ok").Bool())
}

func TestCoroErrorPropagation(t *testing.T) {
	// A runtime error inside the coroutine body propagates to the caller.
	// Calling a non-callable (int) causes "not callable: int" at runtime.
	src := `
coro := import("coro")
boom := func(yield) {
    yield(1)
    x := 5
    x()
}
co := coro.new(boom)
v1, _ := co.resume()
_, _ := co.resume()
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	runErr := compiled.Run()
	require.Error(t, runErr)
}

func TestCoroTypeNameAndString(t *testing.T) {
	src := `
coro := import("coro")
gen := func(yield) {}
co := coro.new(gen)
tn := type_name(co)
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())
	require.Equal(t, "coroutine", compiled.Get("tn").String())
}

func TestCoroYieldUndefined(t *testing.T) {
	// yield with no argument sends undefined.
	src := `
coro := import("coro")
gen := func(yield) { yield() }
co := coro.new(gen)
v, ok := co.resume()
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())

	require.True(t, compiled.Get("v").IsUndefined())
	require.Equal(t, true, compiled.Get("ok").Bool())
}

func TestCoroMultipleCoroutines(t *testing.T) {
	// Two independent coroutines interleaved.
	src := `
coro := import("coro")
mkgen := func(start) {
    return func(yield) {
        yield(start)
        yield(start + 1)
    }
}
c1 := coro.new(mkgen(10))
c2 := coro.new(mkgen(20))
a, _ := c1.resume()
b, _ := c2.resume()
c, _ := c1.resume()
d, _ := c2.resume()
`
	s := tengo.NewScript([]byte(src))
	s.SetImports(stdlib.GetModuleMap("coro"))
	compiled, err := s.Compile()
	require.NoError(t, err)
	require.NoError(t, compiled.Run())

	require.Equal(t, int64(10), compiled.Get("a").Int64())
	require.Equal(t, int64(20), compiled.Get("b").Int64())
	require.Equal(t, int64(11), compiled.Get("c").Int64())
	require.Equal(t, int64(21), compiled.Get("d").Int64())
}
