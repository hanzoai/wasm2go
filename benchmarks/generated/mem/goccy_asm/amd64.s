//go:build amd64 && amd64.v2

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $0-20
	NO_LOCAL_POINTERS
	MOVQ m+0(FP), AX
	MOVLQSX l0+8(FP), BX
	MOVQ	32(AX), CX
	XORL	DX, DX
pc6:
	MOVL	DX, SI
	SHLL	$5, DX
	SUBL	SI, DX
	MOVL	SI, DI
	MOVB	DL, (CX)(DI*1)
	LEAL	1(SI), DX
	CMPL	BX, DX
	JHI	pc6
	XORL	DX, DX
	XORL	SI, SI
pc29:
	MOVL	DX, DI
	ADDL	(CX)(DI*1), SI
	ADDL	$4, DX
	CMPL	BX, DX
	JHI	pc29
	MOVL	SI, AX
	MOVL AX, r0+16(FP)
	RET

