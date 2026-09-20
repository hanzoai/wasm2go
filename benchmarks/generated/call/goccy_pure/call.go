package goccy_pure

import (
	"math"
	"math/bits"
)

type Module struct {
}

func New() *Module {
	m := &Module{}
	return m
}

func (m *Module) Add(l0 int32, l1 int32) int32 {
	return fn0(m, l0, l1)
}
func (m *Module) Sub(l0 int32, l1 int32) int32 {
	return fn1(m, l0, l1)
}
func (m *Module) Mul64(l0 int64, l1 int64) int64 {
	return fn2(m, l0, l1)
}
func (m *Module) DivS(l0 int32, l1 int32) int32 {
	return fn3(m, l0, l1)
}
func (m *Module) Shifts(l0 int32, l1 int32) int32 {
	return fn4(m, l0, l1)
}
func (m *Module) Rotl(l0 int32, l1 int32) int32 {
	return fn5(m, l0, l1)
}
func (m *Module) LtS(l0 int32, l1 int32) int32 {
	return fn6(m, l0, l1)
}
func (m *Module) LtU(l0 int32, l1 int32) int32 {
	return fn7(m, l0, l1)
}

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func ui32(x int32) uint32 { return uint32(x) }

// b2i32 materialises a wasm comparison result — an i32 that is 0 or 1 — from
// the Go bool the comparison expression evaluates to.
//
// It exists as a named helper rather than an inline `func() int32 { ... }()`
// because the gcasm backend requires every direct call left in the compiled
// output to be either a package-local FnN or something the Go inliner removed.
// A func literal is normally inlined at its call site, but the inliner gives up
// once the ENCLOSING function grows past its budget — and a single wasm function
// can translate to tens of thousands of lines of Go, as an interpreter's
// bytecode dispatch loop does. The literal is then outlined into a real closure
// symbol (FnN.funcA.funcB), which reaches the assembler as a direct call gcasm
// cannot marshal. A named helper this small is always inlined, and if it ever
// were not, it would fail loudly at its own symbol rather than as a nested
// closure.
func b2i32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

//go:noinline
func wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func wasm_trap_unreachable() { panic("wasm: unreachable") }

func i32_div_s(x, y int32) int32 {
	if y == -1 && x == math.MinInt32 {
		wasm_trap_int_overflow()
	}
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x / y
}

func i32_rotl(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), int(y&31))) }
