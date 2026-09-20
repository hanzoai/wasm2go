//go:build amd64 && amd64.v2

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $0-24
	NO_LOCAL_POINTERS
	MOVQ m+0(FP), AX
	MOVQ l0+8(FP), BX
	TESTQ	BX, BX
	MOVL	$0, AX
	CMOVQGT	BX, AX
	XORL	CX, CX
	MOVL	$1, DX
	JMP	pc34
pc21:
	LEAQ	(CX)(DX*1), BX
	DECQ	AX
	MOVQ	DX, CX
	MOVQ	BX, DX
pc34:
	TESTQ	AX, AX
	JNE	pc21
	MOVQ	CX, AX
	MOVQ AX, r0+16(FP)
	RET

