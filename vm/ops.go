package vm

import "strconv"

//+pigeon: ops.go

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

// ϡinstr holds a single instruction: an opcode with its arguments.
type ϡinstr struct {
	op       ϡop
	ruleNmIx int
	args     []uint16
}
