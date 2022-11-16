package main

import (
	"log"
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

const (
	TRAP_GETC  uint16 = 0x20 // get character from keyboard, not echoed onto the terminal
	TRAP_OUT   uint16 = 0x21 // output a character
	TRAP_PUTS  uint16 = 0x22 // output a word string
	TRAP_IN    uint16 = 0x23 // get character from keyboard, echoed onto the terminal
	TRAP_PUTSP uint16 = 0x24 // output a byte string
	TRAP_HALT  uint16 = 0x25 // halt the problem
)

const (
	MR_KBSR uint16 = 0xFE00 // keyboard status
	MR_KBDR uint16 = 0xFE02 // keyboard data
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

func (m *Memory) memWrite(address uint16, val uint16) {
	m[address] = val
}

func (m *Memory) memRead(address uint16) uint16 {
	if address == MR_KBSR {
		if checkKey != 0 {
			m[MR_KBSR] = 1 << 15
			m[MR_KBDR] = checkKey
		} else {
			m[MR_KBSR] = 0
		}
	}
	return m[address]
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
			// destination register (DR)
			r0 := (instr >> 9) & 0x7
			// first operannd (SR1)
			r1 := (instr >> 6) & 0x7
			// whether we are in immediate mode
			imm_flag := (instr >> 5) & 0x1

			if imm_flag == 1 {
				imm5 := signExtend(instr&0x1F, 5)
				register[r0] = register[r1] + imm5
			} else {
				r2 := instr & 0x7
				register[r0] = register[r1] + register[r2]
			}
			updateFlags(r0)
		case OP_AND:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			imm_flag := (instr >> 5) & 0x1

			if imm_flag == 1 {
				imm5 := signExtend(instr&0x1F, 5)
				register[r0] = register[r1] & imm5
			} else {
				r2 := instr & 0x7
				register[r0] = register[r1] & register[r2]
			}
			updateFlags(r0)
		case OP_NOT:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			register[r0] = ^register[r1]
			updateFlags(r0)
		case OP_BR:
			pcOffset := signExtend(instr&0x1FF, 9)
			condFlag := (instr >> 9) & 0x7
			if condFlag&register[R_COND] != 0 {
				register[R_PC] += pcOffset
			}
		case OP_JMP:
			// Also handles RET
			r1 := (instr >> 6) & 0x7
			register[R_PC] = register[r1]
		case OP_JSR:
			longFlag := (instr > 11) & 0x1
			register[R_R7] = register[R_PC]

			if longFlag == 1 {
				longPcOffset := signExtend(instr&0x7FF, 11)
				register[R_PC] += longPcOffset // JSR
			} else {
				r1 := (instr >> 9) & 0x7
				register[R_PC] = register[R1] // JSSR
			}
		case OP_LD:
			r0 := (instr >> 9) & 0x7
			pcOffset = signExtend(instr&0x1FF, 9)
			register[r0] = memRead(register[R_PC] + pcOffset)
			updateFlags(r0)
		case OP_LDI:
			// destination register
			r0 := (instr >> 9) & 0x7
			// PCoffset
			pcOffset = signExtend(instr&0x1FF, 9)
			// add pcOffset to the current PC, look at that memory location to get the final address
			register[r0] = memRead(memRead(register[R_PC] + pcOffset))
			updateFlags(r0)
		case OP_LDR:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			offset := signExtend(instr&0x3F, 6)
			register[r0] = memRead(register[r1] + offset)
			udpateFlags(r0)
		case OP_LEA:
			r0 := (instr >> 9) & 0x7
			pcOffset = signExtend(instr&0x1FF, 9)
			register[r0] = register[R_PC] + pcOffset
			updateFlags(r0)
		case OP_ST:
			r0 := (instr >> 9) & 0x7
			pcOffset := signExtend(instr&0x1FF, 9)
			memWrite(register[R_PC]+pcOffset, register[r0])
		case OP_STI:
			r0 := (instr >> 9) & 0x7
			pcOffset = signExtend(instr&0x1FF, 9)
			memWrite(memRead(register[R_PC]+pcOffset), register[r0])
		case OP_STR:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			offset = signExtend(instr&0x3F, 6)
			memWrite(register[r1]+offset, register[r0])
		case OP_TRAP:
			register[R_R7] = register[R_PC]

			switch instr & 0xFF {
			case TRAP_GETC:
			case TRAP_OUT:
			case TRAP_PUTS:
			case TRAP_IN:
			case TRAP_PUTSP:
			case TRAP_HALT:
			}
		case OP_RES:
			log.Printf("Invalid instruction")
		case OP_RTI:
			log.Printf("Invalid instruction")
		default:
		}
	}
}
