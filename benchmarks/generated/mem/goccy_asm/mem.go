package goccy_asm

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Module struct {
	memory      []byte
	maxMem      uint64
	M           unsafe.Pointer
	memMu       *sync.Mutex
	memSize     *atomic.Uint64
	dataEnd     uint32
	memShared   bool
	threads     *threadPool
	threadStart func(*Module, int32, int32)
}

func New() *Module {
	m := &Module{}
	m.memory = make([]byte, 65536, 81920)
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	m.memSize.Store(65536)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = 4294967296
	return m
}

const InitialMemoryBytes = 65536

func NewWithMemory(memory []byte, memSize uint64) *Module {
	m := &Module{}
	m.memory = memory
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	if memSize > 4294836224 {
		panic("wasm2go: memory size exceeds the implementation limit (4294836224 bytes)")
	}
	m.memSize.Store(memSize)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = uint64(len(memory))
	return m
}
func NewFromSnapshot(memory []byte, memSize uint64, globals []uint64) *Module {
	m := &Module{}
	m.memory = memory
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	if memSize > 4294836224 {
		panic("wasm2go: memory size exceeds the implementation limit (4294836224 bytes)")
	}
	m.memSize.Store(memSize)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = uint64(len(memory))
	restoreGlobals(m, globals)
	return m
}

func (m *Module) BenchMem(l0 int32) int32 {
	return fn0(m, l0)
}
func (m *Module) Memory() []byte {
	return m.memory
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

// accessMemory runs f with the module's current linear memory while
// holding the same lock memoryGrow takes to mutate the memory slice
// header or relocate its backing array. It is the ONE safe way to
// touch linear memory from OUTSIDE the module's execution goroutine —
// e.g. a watchdog goroutine raising CPython's eval-breaker bit while
// an evaluation is running. For the duration of f the memory can
// neither be resliced nor relocated, so f's writes land in the array
// the guest observes; a grow that raced in just before blocks until f
// returns and then copies f's writes forward with the rest of the
// contents. Determinism notes for callers:
//
//   - f MUST NOT call back into the module or into memoryGrow — that
//     would self-deadlock.
//   - f should be short: a running guest blocks inside memory.grow
//     until f returns (ordinary guest loads/stores do not block).
//   - Bytes the guest reads or writes concurrently with f (that is
//     the point of an eval-breaker-style flag) are exchanged with
//     plain single-word accesses; keep such shared words
//     word-aligned and word-sized.
func accessMemory(m *Module, f func(mem []byte)) {
	m.memMu.Lock()
	defer m.memMu.Unlock()
	f(m.memory)
}

var spinAgents int32
var spinOversubscribed uint32

type threadPool struct {
	nextTID atomic.Int32
	wg      sync.WaitGroup

	parkMu sync.Mutex
	parked map[uint64][]chan struct{}
}

// wake releases up to count waiters on ea and reports how many it woke.
func (p *threadPool) wake(ea uint64, count int32) int32 {
	p.parkMu.Lock()
	defer p.parkMu.Unlock()
	waiters := p.parked[ea]
	n := int32(len(waiters))
	if count >= 0 && count < n {
		n = count
	}
	for _, ch := range waiters[:n] {
		close(ch)
	}
	if int(n) == len(waiters) {
		delete(p.parked, ea)
	} else {
		p.parked[ea] = waiters[n:]
	}
	return n
}

// saveGlobals returns the module's mutable globals, in a form that can be handed back
// to restoreGlobals. It is how a snapshot of an instance captures the state that does not
// live in linear memory.
func saveGlobals(m *Module) []uint64 {
	g := make([]uint64, 0)
	return g
}

// restoreGlobals puts a snapshot's globals back. A snapshot from a different module (or a
// different build of the same one) has a different global count; rather than
// index out of bounds, take what fits and leave the rest at their declared
// initializers.
func restoreGlobals(m *Module, g []uint64) {
	if len(g) != 0 {
		return
	}
}
