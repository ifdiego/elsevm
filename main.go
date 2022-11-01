package main

import (
	"math"
)

type Memory [math.MaxUint16 + 1]uint16

var memory Memory

const (
	R_R0 uint16 = iota
	R_R1
	R_R2
	R_R3
	R_R4
	R_R5
	R_R6
	R_R7
	R_PC // program counter
	R_COND
	R_COUNT
)

var register [R_COUNT]uint16

const (
	OP_BR   uint16 = iota // branch
	OP_ADD                // add
	OP_LD                 // load
	OP_ST                 // store
	OP_JSR                // jump register
	OP_AND                // bitwise and
	OP_LDR                // load register
	OP_STR                // store register
	OP_RTI                // unused
	OP_NOT                // bitwise not
	OP_LDI                // load indirect
	OP_STI                // store indirect
	OP_JMP                // jump
	OP_RES                // reserved (unused)
	OP_LEA                // load effective address
	OP_TRAP               // execute trap
)

const (
	FL_POS uint16 = 1 << 0 // P
	FL_ZRO uint16 = 1 << 1 // Z
	FL_NEG uint16 = 1 << 2 // N
)

func signExtend(x uint16, bit_count int) uint16 {
	if (x >> (bit_count - 1) & 1) != 0 {
		x |= 0xFFFF << bit_count
	}
	return x
}

func updateFlags(r uint16) {
	if register[r] == 0 {
		register[R_COND] = FL_ZRO
	} else if (register[r] >> 15) == 1 { // a 1 in the left-most bit indicates negative
		register[R_COND] = FL_NEG
	} else {
		register[R_COND] = FL_POS
	}
}

func main() {
	// set the PC to starting position
	// 0x3000 is the default
	var PC_START uint16 = 0x3000
	register[R_PC] = PC_START

	for {
		instr := memRead(register[R_PC])
		op := instr >> 12
		register[R_PC]++

		switch op {
		case OP_ADD:
		case OP_AND:
		case OP_NOT:
		case OP_BR:
		case OP_JMP:
		case OP_JSR:
		case OP_LD:
		case OP_LDI:
		case OP_LDR:
		case OP_LEA:
		case OP_ST:
		case OP_STI:
		case OP_STR:
		case OP_TRAP:
		case OP_RES:
		case OP_RTI:
		default:
		}
	}
}
