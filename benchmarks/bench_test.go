package benchmarks

import (
	"context"
	"encoding/binary"
	_ "embed"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	fib_asm "github.com/goccy/wasm2go/benchmarks/generated/fib/goccy_asm"
	fib_pure "github.com/goccy/wasm2go/benchmarks/generated/fib/goccy_pure"
	fib_ncruces "github.com/goccy/wasm2go/benchmarks/generated/fib/ncruces"

	prime_asm "github.com/goccy/wasm2go/benchmarks/generated/prime/goccy_asm"
	prime_pure "github.com/goccy/wasm2go/benchmarks/generated/prime/goccy_pure"
	prime_ncruces "github.com/goccy/wasm2go/benchmarks/generated/prime/ncruces"

	mem_asm "github.com/goccy/wasm2go/benchmarks/generated/mem/goccy_asm"
	mem_pure "github.com/goccy/wasm2go/benchmarks/generated/mem/goccy_pure"
	mem_ncruces "github.com/goccy/wasm2go/benchmarks/generated/mem/ncruces"

	call_asm "github.com/goccy/wasm2go/benchmarks/generated/call/goccy_asm"
	call_pure "github.com/goccy/wasm2go/benchmarks/generated/call/goccy_pure"
	call_ncruces "github.com/goccy/wasm2go/benchmarks/generated/call/ncruces"
)

//go:embed testdata/fib.wasm
var fibWasm []byte

//go:embed testdata/primes.wasm
var primesWasm []byte

//go:embed testdata/mem.wasm
var memWasm []byte

//go:embed testdata/arith.wasm
var arithWasm []byte

// Native reference implementations

func fibNative(n int64) int64 {
	if n <= 0 {
		return 0
	}
	var v1 int64 = 1
	var v2 int64 = 0
	for n > 0 {
		n--
		t3 := v2
		v2 = v1
		v1 = t3 + v2
	}
	return v2
}

func primeNative(n int32) int32 {
	if n < 2 {
		return 0
	}
	if n == 2 {
		return 1
	}
	if n%2 == 0 {
		return 0
	}
	for i := int32(3); i < n; i += 2 {
		if n%i == 0 {
			return 0
		}
	}
	return 1
}

func memNative(size int32) int32 {
	mem := make([]byte, 65536)
	for i := int32(0); i < size; i++ {
		mem[i] = byte((i * 31) & 0xFF)
	}
	var acc int32
	for i := int32(0); i < size; i += 4 {
		acc += int32(binary.LittleEndian.Uint32(mem[i:]))
	}
	return acc
}

func setupWazero(ctx context.Context, wasmBytes []byte, useInterpreter bool, fnName string) (wazero.Runtime, api.Module, api.Function, func()) {
	var r wazero.Runtime
	if useInterpreter {
		r = wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	} else {
		r = wazero.NewRuntime(ctx)
	}
	mod, err := r.Instantiate(ctx, wasmBytes)
	if err != nil {
		panic(err)
	}
	fn := mod.ExportedFunction(fnName)
	if fn == nil {
		panic("exported function not found: " + fnName)
	}
	cleanup := func() {
		_ = r.Close(ctx)
	}
	return r, mod, fn, cleanup
}

func TestParity(t *testing.T) {
	ctx := context.Background()

	// 1. Fibonacci parity
	_, _, wazeroFibComp, cleanup1 := setupWazero(ctx, fibWasm, false, "fibonacci")
	defer cleanup1()
	_, _, wazeroFibInterp, cleanup2 := setupWazero(ctx, fibWasm, true, "fibonacci")
	defer cleanup2()

	mFibNcruces := fib_ncruces.New()
	mFibPure := fib_pure.New()
	mFibAsm := fib_asm.New()

	for _, n := range []int64{0, 1, 2, 5, 10, 20, 35} {
		want := fibNative(n)
		if got := mFibNcruces.Xfibonacci(n); got != want {
			t.Fatalf("Fib ncruces mismatch for n=%d: got %d, want %d", n, got, want)
		}
		if got := mFibPure.Fibonacci(n); got != want {
			t.Fatalf("Fib pure mismatch for n=%d: got %d, want %d", n, got, want)
		}
		if got := mFibAsm.Fibonacci(n); got != want {
			t.Fatalf("Fib asm mismatch for n=%d: got %d, want %d", n, got, want)
		}
		resComp, _ := wazeroFibComp.Call(ctx, uint64(n))
		if int64(resComp[0]) != want {
			t.Fatalf("Fib wazero-comp mismatch for n=%d: got %d, want %d", n, resComp[0], want)
		}
		resInterp, _ := wazeroFibInterp.Call(ctx, uint64(n))
		if int64(resInterp[0]) != want {
			t.Fatalf("Fib wazero-interp mismatch for n=%d: got %d, want %d", n, resInterp[0], want)
		}
	}

	// 2. Prime parity
	_, _, wazeroPrimeComp, cleanup3 := setupWazero(ctx, primesWasm, false, "is_prime")
	defer cleanup3()
	_, _, wazeroPrimeInterp, cleanup4 := setupWazero(ctx, primesWasm, true, "is_prime")
	defer cleanup4()

	mPrimeNcruces := prime_ncruces.New()
	mPrimePure := prime_pure.New()
	mPrimeAsm := prime_asm.New()

	for _, n := range []int32{0, 1, 2, 3, 4, 17, 25, 97, 100, 997, 1000} {
		want := primeNative(n)
		if got := mPrimeNcruces.Xis_prime(n); got != want {
			t.Fatalf("Prime ncruces mismatch for n=%d: got %d, want %d", n, got, want)
		}
		if got := mPrimePure.IsPrime(n); got != want {
			t.Fatalf("Prime pure mismatch for n=%d: got %d, want %d", n, got, want)
		}
		if got := mPrimeAsm.IsPrime(n); got != want {
			t.Fatalf("Prime asm mismatch for n=%d: got %d, want %d", n, got, want)
		}
		resComp, _ := wazeroPrimeComp.Call(ctx, uint64(uint32(n)))
		if int32(resComp[0]) != want {
			t.Fatalf("Prime wazero-comp mismatch for n=%d: got %d, want %d", n, resComp[0], want)
		}
		resInterp, _ := wazeroPrimeInterp.Call(ctx, uint64(uint32(n)))
		if int32(resInterp[0]) != want {
			t.Fatalf("Prime wazero-interp mismatch for n=%d: got %d, want %d", n, resInterp[0], want)
		}
	}

	// 3. Memory parity
	_, _, wazeroMemComp, cleanup5 := setupWazero(ctx, memWasm, false, "bench_mem")
	defer cleanup5()
	_, _, wazeroMemInterp, cleanup6 := setupWazero(ctx, memWasm, true, "bench_mem")
	defer cleanup6()

	mMemNcruces := mem_ncruces.New()
	mMemPure := mem_pure.New()
	mMemAsm := mem_asm.New()

	for _, sz := range []int32{64, 256, 1024, 4096, 16384} {
		want := memNative(sz)
		if got := mMemNcruces.Xbench_mem(sz); got != want {
			t.Fatalf("Mem ncruces mismatch for sz=%d: got %d, want %d", sz, got, want)
		}
		if got := mMemPure.BenchMem(sz); got != want {
			t.Fatalf("Mem pure mismatch for sz=%d: got %d, want %d", sz, got, want)
		}
		if got := mMemAsm.BenchMem(sz); got != want {
			t.Fatalf("Mem asm mismatch for sz=%d: got %d, want %d", sz, got, want)
		}
		resComp, _ := wazeroMemComp.Call(ctx, uint64(uint32(sz)))
		if int32(resComp[0]) != want {
			t.Fatalf("Mem wazero-comp mismatch for sz=%d: got %d, want %d", sz, resComp[0], want)
		}
		resInterp, _ := wazeroMemInterp.Call(ctx, uint64(uint32(sz)))
		if int32(resInterp[0]) != want {
			t.Fatalf("Mem wazero-interp mismatch for sz=%d: got %d, want %d", sz, resInterp[0], want)
		}
	}

	// 4. Call parity
	_, _, wazeroCallComp, cleanup7 := setupWazero(ctx, arithWasm, false, "add")
	defer cleanup7()
	_, _, wazeroCallInterp, cleanup8 := setupWazero(ctx, arithWasm, true, "add")
	defer cleanup8()

	mCallNcruces := call_ncruces.New()
	mCallPure := call_pure.New()
	mCallAsm := call_asm.New()

	for _, pair := range [][2]int32{{0, 0}, {1, 2}, {100, 200}, {-50, 50}, {12345, 67890}} {
		want := pair[0] + pair[1]
		if got := mCallNcruces.Xadd(pair[0], pair[1]); got != want {
			t.Fatalf("Call ncruces mismatch for %v: got %d, want %d", pair, got, want)
		}
		if got := mCallPure.Add(pair[0], pair[1]); got != want {
			t.Fatalf("Call pure mismatch for %v: got %d, want %d", pair, got, want)
		}
		if got := mCallAsm.Add(pair[0], pair[1]); got != want {
			t.Fatalf("Call asm mismatch for %v: got %d, want %d", pair, got, want)
		}
		resComp, _ := wazeroCallComp.Call(ctx, uint64(uint32(pair[0])), uint64(uint32(pair[1])))
		if int32(resComp[0]) != want {
			t.Fatalf("Call wazero-comp mismatch for %v: got %d, want %d", pair, resComp[0], want)
		}
		resInterp, _ := wazeroCallInterp.Call(ctx, uint64(uint32(pair[0])), uint64(uint32(pair[1])))
		if int32(resInterp[0]) != want {
			t.Fatalf("Call wazero-interp mismatch for %v: got %d, want %d", pair, resInterp[0], want)
		}
	}
}

var sinkI64 int64
var sinkI32 int32

// --- Benchmark Fibonacci (Iterative, n=35) ---

func BenchmarkFib_Native(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkI64 = fibNative(35)
	}
}

func BenchmarkFib_Goccy_ASM(b *testing.B) {
	m := fib_asm.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI64 = m.Fibonacci(35)
	}
}

func BenchmarkFib_Goccy_Pure(b *testing.B) {
	m := fib_pure.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI64 = m.Fibonacci(35)
	}
}

func BenchmarkFib_Ncruces(b *testing.B) {
	m := fib_ncruces.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI64 = m.Xfibonacci(35)
	}
}

func BenchmarkFib_Wazero_Compiler(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, fibWasm, false, "fibonacci")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, 35)
		sinkI64 = int64(res[0])
	}
}

func BenchmarkFib_Wazero_Interpreter(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, fibWasm, true, "fibonacci")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, 35)
		sinkI64 = int64(res[0])
	}
}

// --- Benchmark Primes (Trial division, n=10,000) ---

const benchPrimeTarget int32 = 9973

func BenchmarkPrime_Native(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkI32 = primeNative(benchPrimeTarget)
	}
}

func BenchmarkPrime_Goccy_ASM(b *testing.B) {
	m := prime_asm.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.IsPrime(benchPrimeTarget)
	}
}

func BenchmarkPrime_Goccy_Pure(b *testing.B) {
	m := prime_pure.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.IsPrime(benchPrimeTarget)
	}
}

func BenchmarkPrime_Ncruces(b *testing.B) {
	m := prime_ncruces.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.Xis_prime(benchPrimeTarget)
	}
}

func BenchmarkPrime_Wazero_Compiler(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, primesWasm, false, "is_prime")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(benchPrimeTarget)))
		sinkI32 = int32(res[0])
	}
}

func BenchmarkPrime_Wazero_Interpreter(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, primesWasm, true, "is_prime")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(benchPrimeTarget)))
		sinkI32 = int32(res[0])
	}
}

// --- Benchmark Memory (Read/Write 16KB Wasm buffer) ---

const benchMemSize int32 = 16384

func BenchmarkMem_Native(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkI32 = memNative(benchMemSize)
	}
}

func BenchmarkMem_Goccy_ASM(b *testing.B) {
	m := mem_asm.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.BenchMem(benchMemSize)
	}
}

func BenchmarkMem_Goccy_Pure(b *testing.B) {
	m := mem_pure.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.BenchMem(benchMemSize)
	}
}

func BenchmarkMem_Ncruces(b *testing.B) {
	m := mem_ncruces.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.Xbench_mem(benchMemSize)
	}
}

func BenchmarkMem_Wazero_Compiler(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, memWasm, false, "bench_mem")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(benchMemSize)))
		sinkI32 = int32(res[0])
	}
}

func BenchmarkMem_Wazero_Interpreter(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, memWasm, true, "bench_mem")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(benchMemSize)))
		sinkI32 = int32(res[0])
	}
}

// --- Benchmark Call Latency (Trivial add(a, b) invocation) ---

func addNative(a, b int32) int32 {
	return a + b
}

func BenchmarkCall_Native(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkI32 = addNative(int32(i), 42)
	}
}

func BenchmarkCall_Goccy_ASM(b *testing.B) {
	m := call_asm.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.Add(int32(i), 42)
	}
}

func BenchmarkCall_Goccy_Pure(b *testing.B) {
	m := call_pure.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.Add(int32(i), 42)
	}
}

func BenchmarkCall_Ncruces(b *testing.B) {
	m := call_ncruces.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkI32 = m.Xadd(int32(i), 42)
	}
}

func BenchmarkCall_Wazero_Compiler(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, arithWasm, false, "add")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(i)), 42)
		sinkI32 = int32(res[0])
	}
}

func BenchmarkCall_Wazero_Interpreter(b *testing.B) {
	ctx := context.Background()
	_, _, fn, cleanup := setupWazero(ctx, arithWasm, true, "add")
	defer cleanup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, _ := fn.Call(ctx, uint64(uint32(i)), 42)
		sinkI32 = int32(res[0])
	}
}

