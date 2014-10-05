// Code written by Chris Stiteler for CS472 (thursday)
// Project 1 - Due 10/02/2014
package main

import (
	"fmt"
)

// the Disassembleable interface is fulfilled by
// both R and I format MIPS instructions, implicitly
type Disassembleable interface {
	disassemble()
}

// rInstruction represents an r-format MIPS instruction
type rInstruction struct {
	bits        uint32
	address     uint32
	instruction string
	opcode      uint32
	src1        uint32
	src2        uint32
	dest        uint32
	function    uint32
}

// iInstruction represents an i-format MIPS instruction
type iInstruction struct {
	bits        uint32
	address     uint32
	instruction string
	opcode      uint32
	src1        uint32
	src2        uint32
	constant    int16
}

type Mask struct {
	bits  uint32
	shift uint8
}

// masks and their subsequent shift values
var opcodeMask = Mask{0xFC000000, 26}
var src1Mask = Mask{0x03E00000, 21}
var src2Mask = Mask{0x001F0000, 16}
var rDestMask = Mask{0x0000F800, 11}
var rFuncMask = Mask{0x0000003F, 0}
var iConstMask = Mask{0x0000FFFF, 0}

// program counter
var pc = uint32(0x7A060)

// processor address size is in bytes
var addressSize = uint32(4)

// our raw hex instruction input
var input = []uint32{
	0x022DA822,
	0x8EF30018,
	0x12A70004,
	0x02689820,
	0xAD930018,
	0x02697824,
	0xAD8FFFF4,
	0x018C6020,
	0x02A4A825,
	0x158FFFF6,
	0x8E59FFF0,
}

func main() {
	// build slice of instructions from input
	instructions := buildInstructions(input)

	// disassemble then print each instruction
	for _, instruct := range instructions {
		instruct.disassemble()
		fmt.Println(instruct)
	}
}

// buildInstructions() builds a slice of Disassemblable{} structs
// for each item, determines instruct. type, opcode, and address
func buildInstructions(input []uint32) []Disassembleable {
	instructions := make([]Disassembleable, len(input))

	// need to get opcode, if 000000, build an rInstruction, else iInstruction
	for i, inputBits := range input {
		var newInstruction Disassembleable
		opcode := maskAndShift(opcodeMask, inputBits)

		if opcode == 0 {
			newInstruction = &rInstruction{bits: inputBits, opcode: opcode,
				address: pc}
		} else {
			newInstruction = &iInstruction{bits: inputBits, opcode: opcode,
				address: pc}
		}

		pc += addressSize
		instructions[i] = newInstruction
	}

	return instructions
}

// r.disassemble() extracts data from rInstruction bits
// also fulfills the Disassembleable interface for rInstructs.
func (r *rInstruction) disassemble() {
	r.src1 = maskAndShift(src1Mask, r.bits)
	r.src2 = maskAndShift(src2Mask, r.bits)
	r.dest = maskAndShift(rDestMask, r.bits)
	r.function = maskAndShift(rFuncMask, r.bits)
	r.instruction = r.getInstruction()
}

// i.dissassemble() extracts data from iInstruction bits
// also fulfills the Disassembleable interface for iInstructs.
func (i *iInstruction) disassemble() {
	i.src1 = maskAndShift(src1Mask, i.bits)
	i.src2 = maskAndShift(src2Mask, i.bits)
	i.constant = maskAndShiftShort(iConstMask, int16(i.bits))
	i.instruction = i.getInstruction()
}

func (r *rInstruction) getInstruction() string {
	switch r.function {
	case 0x20:
		return fmt.Sprintf("add")
	case 0x22:
		return fmt.Sprintf("sub")
	case 0x24:
		return fmt.Sprintf("and")
	case 0x25:
		return fmt.Sprintf("or")
	default:
		return ""
	}
}

func (i *iInstruction) getInstruction() string {
	switch i.opcode {
	case 0x4:
		return fmt.Sprintf("beq")
	case 0x5:
		return fmt.Sprintf("bne")
	case 0x23:
		return fmt.Sprintf("lw")
	case 0x2B:
		return fmt.Sprintf("sw")
	default:
		return ""
	}
}

// String() is a 'toString' method on an rInstruction
func (r *rInstruction) String() string {
	return fmt.Sprintf("%X %s $%d, $%d, $%d", r.address, r.instruction,
		r.dest, r.src1, r.src2)
}

// i.String() is a a 'toString' method on an iInstruction
func (i *iInstruction) String() string {
	// if branch, handle branch format
	if i.instruction == "bne" || i.instruction == "beq" {

		branchAddress := i.getBranchToAddress()
		return fmt.Sprintf("%X %s $%d, $%d address %X", i.address, i.instruction,
			i.src1, i.src2, branchAddress)
	} else {
		// else handle lw/sw format
		return fmt.Sprintf("%X %s $%d, %d ($%d)", i.address, i.instruction,
			i.src2, i.constant, i.src1)
	}
}

// getBranchToAddress() calcs branch address using pc-relative offset
// and the address of the current instruction
func (i *iInstruction) getBranchToAddress() uint32 {
	// shift offset/const 2 bits left to decompress, and account for incrememted pc
	return i.address + addressSize + (uint32(i.constant) << 2)
}

// maskAndShift() returns desired bits in a 32-bit value
// depending on the mask (including a shift value)
func maskAndShift(mask Mask, inputBits uint32) uint32 {
	return (inputBits & mask.bits) >> mask.shift
}

// maskAndShift() returns desired bits in a 16-bit value
// depending on the mask (including a shift value)
func maskAndShiftShort(mask Mask, inputBits int16) int16 {
	return (inputBits & int16(mask.bits)) >> mask.shift
}
