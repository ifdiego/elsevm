package main

import (
	"fmt"
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

func main() {
	fmt.Println(math.MaxUint16) // 65535
}
