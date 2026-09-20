package goccy_pure

func fn0(m *Module, l0 int32) int32 {
	var v12 int32
	_ = v12
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	if ui32(l0) < ui32(int32(2)) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return int32(0)
L2:
	;
	goto L3
L3:
	;
	if l0 == int32(2) {
		goto L4
	} else {
		goto L5
	}
L4:
	;
	return int32(1)
L5:
	;
	goto L6
L6:
	;
	v12 = i32_rem_u_s(l0, int32(2))
	if v12 == int32(0) {
		goto L7
	} else {
		goto L8
	}
L7:
	;
	return int32(0)
L8:
	;
	goto L9
L9:
	;
	v19 = int32(3)
	goto L10
L10:
	;
	if ui32(l0) <= ui32(v19) {
		goto L12
	} else {
		goto L13
	}
L11:
	;
	return int32(1)
L12:
	;
	goto L11
L13:
	;
	v21 = i32_rem_u_s(l0, v19)
	if v21 == int32(0) {
		goto L14
	} else {
		goto L15
	}
L14:
	;
	return int32(0)
L15:
	;
	goto L16
L16:
	;
	v19 = v19 + int32(2)
	goto L10
}
