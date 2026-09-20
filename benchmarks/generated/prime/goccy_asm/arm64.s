//go:build arm64

#include "textflag.h"
#include "funcdata.h"

TEXT ·fn0(SB), $32-20
	NO_LOCAL_POINTERS
	MOVD m+0(FP), R0
	MOVW l0+8(FP), R1
	CMPW	$2, R1
	BLO	pc84
	BEQ	pc68
	TBZ	$0, R1, pc52
	MOVW	R1, l0+8(FP)
	MOVD	$3, R2
	JMP	pc104
pc52:
	MOVD	ZR, R0
	MOVW R0, r0+16(FP)
	RET	(R30)
pc68:
	MOVD	$1, R0
	MOVW R0, r0+16(FP)
	RET	(R30)
pc84:
	MOVD	ZR, R0
	MOVW R0, r0+16(FP)
	RET	(R30)
pc100:
	ADD	$2, R2, R2
pc104:
	CMPW	R2, R1
	BLS	pc172
	HINT	$0
	HINT	$0
	CBNZW	R2, pc140
	MOVW	R2, 20(RSP)
	CALL ·wasm_trap_div_zero(SB)
	MOVW	l0+8(FP), R1
	MOVW	20(RSP), R2
pc140:
	CBZW	R2, pc188
	UDIVW	R2, R1, R3
	MSUBW	R2, R1, R3, R3
	CBNZW	R3, pc100
	MOVD	ZR, R0
	MOVW R0, r0+16(FP)
	RET	(R30)
pc172:
	MOVD	$1, R0
	MOVW R0, r0+16(FP)
	RET	(R30)
pc188:
	CALL	·gcasmTrapDivZero(SB)
	HINT	$0

