package goccy_pure

func fn0(m *Module, l0 int32, l1 int32) int32 {
	return l0 + l1
}
func fn1(m *Module, l0 int32, l1 int32) int32 {
	return l0 - l1
}
func fn2(m *Module, l0 int64, l1 int64) int64 {
	return l0 * l1
}
func fn3(m *Module, l0 int32, l1 int32) int32 {
	var v3 int32
	_ = v3
	v3 = i32_div_s(l0, l1)
	return v3
}
func fn4(m *Module, l0 int32, l1 int32) int32 {
	return int32(ui32(l0<<(uint(l1)%32)) >> (uint(l1) % 32))
}
func fn5(m *Module, l0 int32, l1 int32) int32 {
	return i32_rotl(l0, l1)
}
func fn6(m *Module, l0 int32, l1 int32) int32 {
	return b2i32(l0 < l1)
}
func fn7(m *Module, l0 int32, l1 int32) int32 {
	return b2i32(ui32(l0) < ui32(l1))
}
