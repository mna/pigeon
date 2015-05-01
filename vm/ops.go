package vm

import (
	"errors"
	"strconv"
)

//+ϡ following code is part of the generated parser

// ϡop represents an opcode.
type ϡop byte

// list of opcodes in the pigeon VM.
const (
	ϡopExit ϡop = iota
	ϡopCall
	ϡopCallA
	ϡopCallB
	ϡopCumulOrF
	ϡopJump
	ϡopJumpIfF
	ϡopJumpIfT
	ϡopMatch
	ϡopNilIfF
	ϡopNilIfT
	ϡopPop
	ϡopPopVJumpIfF
	ϡopPush
	ϡopRestore
	ϡopRestoreIfF
	ϡopReturn
	ϡopStoreIfT
	ϡopTakeLOrJump
	ϡopmax // must always be after the last valid opcode

	// ϡopPlaceholder is an (invalid) opcode used by the Generator
	// to insert opcodes that need the index of the starting instruction
	// of a rule that hasn't been generated yet.
	//
	// It must be placed after ϡopmax (because it is invalid in the
	// final program) and it has one argument, the index in the strings
	// array of the identifier of the rule.
	ϡopPlaceholder
)

// ϡlookupOp translates an opcode to a string.
var ϡlookupOp = []string{
	ϡopExit: "exit", ϡopCall: "call", ϡopCallA: "callA",
	ϡopCallB: "callB", ϡopCumulOrF: "cumulOrF",
	ϡopJump: "jump", ϡopJumpIfF: "jumpIfF", ϡopJumpIfT: "jumpIfT",
	ϡopMatch: "match", ϡopNilIfF: "nilIfF", ϡopNilIfT: "nilIfT",
	ϡopPop: "pop", ϡopPopVJumpIfF: "popVJumpIfF",
	ϡopPush: "push", ϡopRestore: "restore", ϡopRestoreIfF: "restoreIfF",
	ϡopReturn: "return", ϡopStoreIfT: "storeIfT", ϡopTakeLOrJump: "takeLOrJump",
}

// String returns the string representation of the opcode.
func (op ϡop) String() string {
	if 0 <= op && int(op) < len(ϡlookupOp) {
		return ϡlookupOp[op]
	}
	return "ϡop(" + strconv.Itoa(int(op)) + ")"
}

// ϡinstr encodes an opcode with its arguments as a 64-bits unsigned
// integer. The bits are used as follows:
//
// o : 6 bits = opcode (max=63)
// n : 10 bits = for PUSHL, number of values in array (max=1023)
// l : 16 bits = instruction index (max=65535)
//
// So a single PUSH instruction can encode 2 indices (first arg is the stack ID).
// The 64-bit value looks like this:
// oooooonn nnnnnnnn llllllll llllllll llllllll llllllll llllllll llllllll
//
// And if a PUSH (L) instruction has more than 2 indices, it can store 4 full
// indices per subsequent values (4 * 16 bits = 64 bits).
type ϡinstr uint64

// limits and masks.
const (
	ϡiBits = 64
	ϡlBits = 16
	ϡnBits = 10
	ϡoBits = 6
	ϡlPerI = ϡiBits / ϡlBits

	ϡlMask = 1<<ϡlBits - 1
	ϡnMask = 1<<ϡnBits - 1
	ϡoMask = 1<<ϡoBits - 1
)

// decode decodes the instruction and returns the 5 parts:
// the opcode, the number of L array values, and the 3 instruction
// indices.
func (i ϡinstr) decode() (op ϡop, n, ix0, ix1, ix2 int) {
	ix2 = int(i & ϡlMask)
	i >>= ϡlBits
	ix1 = int(i & ϡlMask)
	i >>= ϡlBits
	ix0 = int(i & ϡlMask)
	i >>= ϡlBits
	n = int(i & ϡnMask)
	i >>= ϡnBits
	op = ϡop(i & ϡoMask)
	return
}

// decodeLs decodes the instruction as a list of L instruction
// indices (as a follow-up value to a PUSHL opcode).
func (i ϡinstr) decodeLs() (ix0, ix1, ix2, ix3 int) {
	ix3 = int(i & ϡlMask)
	i >>= ϡlBits
	ix2 = int(i & ϡlMask)
	i >>= ϡlBits
	ix1 = int(i & ϡlMask)
	i >>= ϡlBits
	ix0 = int(i & ϡlMask)
	return
}

// ϡencodeInstr encodes the provided operation and its arguments into
// a list of instruction values. It may return an error if any part
// of the instruction overflows the allowed values.
func ϡencodeInstr(op ϡop, args ...int) ([]ϡinstr, error) {
	var is []ϡinstr

	if op >= ϡopmax && op != ϡopPlaceholder {
		return nil, errors.New("invalid op value")
	}
	if len(args) > ϡnMask {
		return nil, errors.New("too many arguments")
	}

	// first instruction contains opcode
	is = append(is, ϡinstr(op)<<(ϡiBits-ϡoBits))
	n := uint(len(args))
	if n == 0 {
		return is, nil
	}
	off := uint(ϡiBits - ϡoBits - ϡnBits)
	is[0] |= ϡinstr(n) << off

	ix := 0
	for i, arg := range args {
		if arg > ϡlMask {
			return nil, errors.New("argument value too big")
		}

		mod := uint((i + 1) % ϡlPerI)
		if mod == 0 {
			is = append(is, 0)
			ix++
		}

		is[ix] |= ϡinstr(arg) << (off - (mod * ϡlBits))
	}

	return is, nil
}
