package goccy_asm

type Module struct {
}

func New() *Module {
	m := &Module{}
	return m
}

func (m *Module) Fibonacci(l0 int64) int64 {
	return fn0(m, l0)
}

//go:noinline
func wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func wasm_trap_unreachable() { panic("wasm: unreachable") }
