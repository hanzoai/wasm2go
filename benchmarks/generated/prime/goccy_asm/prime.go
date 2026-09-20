package goccy_asm

type Module struct {
}

func New() *Module {
	m := &Module{}
	return m
}

func (m *Module) IsPrime(l0 int32) int32 {
	return fn0(m, l0)
}

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func ui32(x int32) uint32 { return uint32(x) }

//go:noinline
func wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func wasm_trap_unreachable() { panic("wasm: unreachable") }

func i32_rem_u(x, y uint32) uint32 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x % y
}

func i32_rem_u_s(x, y int32) int32 { return int32(i32_rem_u(uint32(x), uint32(y))) }
