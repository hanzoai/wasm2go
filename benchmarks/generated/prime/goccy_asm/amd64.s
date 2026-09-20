//go:build amd64 && amd64.v2

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $16-20
	NO_LOCAL_POINTERS
	MOVQ m+0(FP), AX
	MOVLQSX l0+8(FP), BX
	CMPL	BX, $2
	JCS	pc61
	JEQ	pc50
	BTL	$0, BX
	JCC	pc42
	MOVL	BX, l0+8(FP)
	MOVL	$3, CX
	JMP	pc72
pc42:
	XORL	AX, AX
	MOVL AX, r0+16(FP)
	RET
pc50:
	MOVL	$1, AX
	MOVL AX, r0+16(FP)
	RET
pc61:
	XORL	AX, AX
	MOVL AX, r0+16(FP)
	RET
pc69:
	ADDL	$2, CX
pc72:
	CMPL	BX, CX
	JLS	pc123
	XCHGL	AX, AX
	XCHGL	AX, AX
	TESTL	CX, CX
	JNE	pc103
	MOVL	CX, 4(SP)
	CALL ·wasm_trap_div_zero(SB)
	MOVL	4(SP), AX
	TESTL	AX, AX
	MOVL	AX, CX
	MOVL	l0+8(FP), BX
pc103:
	JEQ	pc134
	MOVL	BX, AX
	XORL	DX, DX
	DIVL	CX
	TESTL	DX, DX
	JNE	pc69
	XORL	AX, AX
	MOVL AX, r0+16(FP)
	RET
pc123:
	MOVL	$1, AX
	MOVL AX, r0+16(FP)
	RET
pc134:
	CALL	·gcasmTrapDivZero(SB)
	XCHGL	AX, AX

