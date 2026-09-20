//go:build arm64

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $16-20
	NO_LOCAL_POINTERS
	MOVD m+0(FP), R0
	MOVW l0+8(FP), R1
	MOVD	32(R0), R2
	MOVD	ZR, R3
pc8:
	LSL	$5, R3, R4
	SUB	R3, R4, R4
	MOVWU	R3, R5
	MOVB	R4, (R2)(R5)
	ADD	$1, R3, R3
	CMPW	R3, R1
	BHI	pc8
	MOVD	ZR, R3
	MOVD	ZR, R0
pc44:
	MOVWU	R3, R4
	MOVW	(R2)(R4), R4
	ADD	R0, R4, R0
	ADD	$4, R3, R3
	CMPW	R3, R1
	BHI	pc44
	MOVW R0, r0+16(FP)
	RET	(R30)

