package goccy_pure

func fn0(m *Module, l0 int64) int64 {
	var v2 int64
	_ = v2
	var v7 int64
	_ = v7
	var v9 int64
	_ = v9
	var __phi9 int64
	_ = __phi9
	var v10 int64
	_ = v10
	var __phi10 int64
	_ = __phi10
	var v11 int64
	_ = v11
	var __phi11 int64
	_ = __phi11
	v2 = int64(0)
	if v2 < l0 {
		v7 = l0
	} else {
		v7 = v2
	}
	__phi9 = v7
	__phi10 = int64(1)
	__phi11 = v2
	v9 = __phi9
	v10 = __phi10
	v11 = __phi11
	for {
		if v9 == int64(0) {
			break
		} else {
			__phi9 = v9 - int64(1)
			__phi10 = v11 + v10
			__phi11 = v10
			v9 = __phi9
			v10 = __phi10
			v11 = __phi11
			continue
		}
		break
	}
	return v11
}
