# pigeon - moving to a VM implementation

The original recursive parser had issues with pathological input, where it could generate stack overflows (e.g. in `test/linear/linear_test.go`, with a 1MB input file). It could also benefit from a different approach with less function call (and possibly allocation) overhead.

The transition to a Virtual Machine (VM) based implementation could be relatively simple. By representing the various expressions and matchers with relatively high-level opcodes, it should be possible to avoid excessive dispatch overhead while avoiding the problems inherent to the recursive implementation.

## Overview

### Matchers

The parser generator would translate all literal matchers in the AST to a list of `Matcher` interfaces:

```
type Matcher interface {
    Match(savepointReader) bool
}

// interface name and methods TBD.
type savepointReader interface {
    current() savepoint
    advance()
}
```

The `AnyMatcher`, `LitMatcher` and `CharClassMatcher` nodes would map to such a `Matcher` interface implementation. Identical literals in the grammar would map to the same `Matcher` value, the same index in the list of matchers.

### Code blocks

Code blocks would still get generated as methods on the `*current` type, but the thunks would be added to a list - actually two separate lists:

* `athunks` : list of action method thunks, signature `func() (interface{}, error)`
* `bthunks` : list of predicate method thunks, signature `func() (bool, error)`

CALL opcodes would have an index argument indicating which thunk to call (e.g. `CALLA 2` or `CALLB 0`).

### Rule name reference

The `parser.rules` map of names to rule nodes would not be required, a rule reference would be simply a jump to the opcode instruction of the start of that rule.

The `parser.rstack` slice of rules serves only to get the rule's identifier (or display name) in error messages. In the VM implementation, a simple mapping of instruction index to rule identifier (or display name) saves memory and achieves the same purpose. The exact way to do the mapping is TBD.

### Variable sets (scope of labels)

The `parser.vstack` field holds the stack of variable sets - a mapping of label name to value that applies ("is in scope") at a given position in the parser. In the VM implementation, a counter would keep track of the current scope depth, and the variable sets would be lazily created only when the first label in a scope is encountered. It would be stored in a `[]map[string]interface{}`, where the index is the scope depth.

On scope exit, if the value for that scope is not nil or an empty map, then the map would be deleted to avoid corruption if the parser goes back to that scope level.

### Memoization

Memoization remains an option on the parser/VM. When an expression returns to its caller index, the values it produced will be stored, along with the starting parser position, the ending position and the index of the first instruction of this expression. Anytime a JUMP would occur to that expression for the same parser position, the VM would bypass the JUMP and instead put the memoized values on the stack directly, advance the parser at the saved ending position and resume execution at the caller's return instruction.

### Error reporting

One of the goals of this rewrite as a VM is to provide better error reporting, namely using the [farthest failure position][ffp] heuristic. The VM will track the FFP along with the instruction of the expression that failed the farthest so better error messages can be returned (the position, rule identifier or display name, and possibly the expected terminal or list of terminals).

Panic recovery would work the same as now, with an option to disable it to get the stack trace.

### Debugging

The debug option would be supported as it is now, although the output will likely be quite different. Exact logging TBD.

### API

The API covered by the API stability guarantee in the doc will remain stable. Internal symbols not part of this API should use a prefix-naming-scheme to avoid clashes with user-defined code (e.g. π?). The accepted PEG syntax remains exactly the same, with the same semantics.

## Opcodes

Each rule and expression execution (a rule is really a specific kind of expression, the RuleRefExpr, and the starting rule is a RuleRefExpr where the identifier is that of the first rule in the grammar) perform the following steps:

* Pop the return instruction index from the stack, store it locally;
* Push the current parser position to the stack;
* Execute the match/expression;
* Pop the current parser position from the stack;
* Push the boolean match result;
* Push the result value or nil;
* Jump to the return instruction index.

The VM defines the following registers (?):

* pc : the program counter, increments on each instruction, `JUMP` sets it directly.
* pt : the point, a position. `PUSH pt` pushes its value onto the stack, `POP pt` pops a value from the stack and assigns it to the register and `SET pt` sets the parser's current position to the value of the pt register.

The following opcodes are required:

* `JUMP i` where `i` is the instruction index. Runs the instruction at `i` as the next instruction.
* `EXIT` pops the return value and the match boolean from the stack and terminates the VM execution, returning those two values.
* `PUSH p` pushes the current parser position to the stack.


## Examples

N means a number, M means a matcher index, I means an instruction index, V means a value generated from a match, B means a boolean (match or not), P means a position.

VM has the following typed registers:

* R int : a return index. POP R stores the value in this register, GOTO R jumps to the instruction at index R.

### E1

Grammar:

```
A <- 'a'
```

* Matchers: 'a'
* Thunks: none

Opcodes:

0: PUSH N (2) : push the return index
1: GOTO N (3) : goto instruction N
2: EXIT : pop V and B, exit VM, return V and B
3: [Rule A] POP R : save the return index, stack is now empty
4:          MATCH M : run the matcher, if not a match, restore P, push B and V
5:          GOTO R : return the the saved return index

### E2

Grammar:

```
A <- 'a' 'b'
```

* Matchers: 'a', 'b'
* Thunks: none

Opcodes:

0: PUSH N (2) : push the return index
1: GOTO N (3) : goto instruction N
2: EXIT : pop V and B, exit VM, return V and B
3: [Rule A] 
4:          


[ffp]: http://arxiv.org/abs/1405.6646
