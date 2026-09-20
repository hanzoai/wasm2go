//go:build !arm64 && (!amd64 || !amd64.v2)

package goccy_asm

import "unsafe"

func fn0(m *Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	v6 = int32(0)
	for {
		v11 = v6 * int32(31) & int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(v6))) = uint8(v11)
		v14 = v6 + int32(1)
		if ui32(v14) < ui32(l0) {
			v6 = v14
			continue
		} else {
			break
		}
		break
	}
	v16 = int32(0)
	v19 = v16
	v20 = v16
	for {
		v21 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
		v22 = v20 + v21
		v24 = v19 + int32(4)
		if ui32(v24) < ui32(l0) {
			v19 = v24
			v20 = v22
			continue
		} else {
			break
		}
		break
	}
	return v22
}
