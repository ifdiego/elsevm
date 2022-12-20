package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/eiannone/keyboard"
)

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

const (
	FL_POS uint16 = 1 << 0 // P
	FL_ZRO uint16 = 1 << 1 // Z
	FL_NEG uint16 = 1 << 2 // N
)

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
	MR_KBSR uint16 = 0xFE00 // keyboard status
	MR_KBDR uint16 = 0xFE02 // keyboard data
)

const (
	TRAP_GETC  uint16 = 0x20 // get character from keyboard, not echoed onto the terminal
	TRAP_OUT   uint16 = 0x21 // output a character
	TRAP_PUTS  uint16 = 0x22 // output a word string
	TRAP_IN    uint16 = 0x23 // get character from keyboard, echoed onto the terminal
	TRAP_PUTSP uint16 = 0x24 // output a byte string
	TRAP_HALT  uint16 = 0x25 // halt the problem
)

type Memory [math.MaxUint16 + 1]uint16

var memory Memory
var register [R_COUNT]uint16
var nativeEndian binary.ByteOrder

func keyboardRead() uint16 {
	symbol, controlKey, err := keyboard.GetSingleKey()

	if controlKey == keyboard.KeyEsc || controlKey == keyboard.KeyCtrlC {
		log.Println("Pressed escaping")
		os.Exit(0)
	}

	if err != nil {
		log.Printf("Error, %s", err)
	}

	return uint16(symbol)
}

func signExtend(x uint16, bit_count int) uint16 {
	if (x >> (bit_count - 1) & 1) != 0 {
		x |= 0xFFFF << bit_count
	}
	return x
}

func swap16(x uint16) uint16 {
	if nativeEndian == binary.BigEndian {
		return (x << 8) | (x >> 8)
	} else {
		return x
	}
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

func readImageFile(file string) {
	// the origin tells us where in memory to place the image
	var origin uint16

	data, _ := readImage(file)
	buffer := bytes.NewBuffer(data)
	origin = swap16(binary.BigEndian.Uint16(buffer.Next(2)))

	bufferLen := buffer.Len()

	for i := 0; i < bufferLen; i++ {
		b := buffer.Next(2)
		if len(b) == 0 {
			break
		}
		memory[origin] = swap16(binary.BigEndian.Uint16(b))
		origin++
	}

	log.Printf("Program has been read into memory, contains %d bytes, %d words", bufferLen, bufferLen/2)
}

func readImage(imagePath string) ([]byte, int64) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var size int64 = info.Size()
	data := make([]byte, size)

	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	return data, size
}

func (m *Memory) memWrite(address uint16, val uint16) {
	m[address] = val
}

func (m *Memory) memRead(address uint16) uint16 {
	if address == MR_KBSR {
		checkKey := keyboardRead()

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
	imageFilePath := flag.String("image", "empty image", "go -image program.obj")
	flag.Parse()
	readImageFile(*imageFilePath)

	// since exactly one condition flag should be set at anu given time, set the Z flag
	register[R_COND] = FL_ZRO

	// set the PC to starting position
	// 0x3000 is the default
	var PC_START uint16 = 0x3000
	register[R_PC] = PC_START
	log.Println("Computer starting...")

	for {
		instr := memory.memRead(register[R_PC])
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
			longFlag := (instr >> 11) & 0x1
			register[R_R7] = register[R_PC]

			if longFlag == 1 {
				longPcOffset := signExtend(instr&0x7FF, 11)
				register[R_PC] += longPcOffset // JSR
			} else {
				r1 := (instr >> 9) & 0x7
				register[R_PC] = register[r1] // JSSR
			}
		case OP_LD:
			r0 := (instr >> 9) & 0x7
			pcOffset := signExtend(instr&0x1FF, 9)
			register[r0] = memory.memRead(register[R_PC] + pcOffset)
			updateFlags(r0)
		case OP_LDI:
			// destination register (DR)
			r0 := (instr >> 9) & 0x7
			// PCoffset 9
			pcOffset := signExtend(instr&0x1FF, 9)
			// add pcOffset to the current PC, look at that memory location to get the final address
			register[r0] = memory.memRead(memory.memRead(register[R_PC] + pcOffset))
			updateFlags(r0)
		case OP_LDR:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			offset := signExtend(instr&0x3F, 6)
			register[r0] = memory.memRead(register[r1] + offset)
			updateFlags(r0)
		case OP_LEA:
			r0 := (instr >> 9) & 0x7
			pcOffset := signExtend(instr&0x1FF, 9)
			register[r0] = register[R_PC] + pcOffset
			updateFlags(r0)
		case OP_ST:
			r0 := (instr >> 9) & 0x7
			pcOffset := signExtend(instr&0x1FF, 9)
			memory.memWrite(register[R_PC]+pcOffset, register[r0])
		case OP_STI:
			r0 := (instr >> 9) & 0x7
			pcOffset := signExtend(instr&0x1FF, 9)
			memory.memWrite(memory.memRead(register[R_PC]+pcOffset), register[r0])
		case OP_STR:
			r0 := (instr >> 9) & 0x7
			r1 := (instr >> 6) & 0x7
			offset := signExtend(instr&0x3F, 6)
			memory.memWrite(register[r1]+offset, register[r0])
		case OP_TRAP:
			register[R_R7] = register[R_PC]

			switch instr & 0xFF {
			case TRAP_GETC:
				// read a single ASCII char
				register[R_R0] = keyboardRead()
			case TRAP_OUT:
				fmt.Printf("%c", register[R_R0])
			case TRAP_PUTS:
				// one char per word
				for c := register[R_R0]; memory[c] != 0x00; c++ {
					fmt.Printf("%c", memory[c])
				}
			case TRAP_IN:
				fmt.Println("Enter a character: ")
				symbol := keyboardRead()
				fmt.Printf("%c", symbol)
				r0 := symbol
				updateFlags(r0)
			case TRAP_PUTSP:
				// one char per byte (two bytes per word)
				// here we need to swap back to
				// big endian format
				for c := register[R_R0]; memory[c] != 0x00; c++ {
					char1 := memory[c]
					fmt.Printf("%c", char1&0xFF)
					char2 := char1 & 0xFF >> 8
					if char2 != 0 {
						fmt.Printf("%c", char2)
					}
				}
			case TRAP_HALT:
				log.Printf("Computer halting...")
				os.Exit(0)
			}
		case OP_RES:
			log.Printf("Invalid instruction")
		case OP_RTI:
			log.Printf("Invalid instruction")
		default:
		}
	}
}
