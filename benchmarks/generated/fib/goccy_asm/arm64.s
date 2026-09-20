//go:build arm64

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $16-24
	NO_LOCAL_POINTERS
	MOVD m+0(FP), R0
	MOVD l0+8(FP), R1
	CMP	$0, R1
	CSEL	GT, R1, ZR, R1
	MOVD	ZR, R2
	MOVD	$1, R3
	JMP	pc36
pc20:
	ADD	R3, R2, R0
	SUB	$1, R1, R1
	MOVD	R3, R2
	MOVD	R0, R3
pc36:
	CBNZ	R1, pc20
	MOVD	R2, R0
	MOVD R0, r0+16(FP)
	RET	(R30)

