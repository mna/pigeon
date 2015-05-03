# pigeon - moving to a VM implementation

The original recursive parser had issues with pathological input, where it could generate stack overflows (e.g. in `test/linear/linear_test.go`, with a 1MB input file). It could also benefit from a different approach with less function call (and possibly allocation) overhead.

The transition to a Virtual Machine (VM) based implementation could be relatively simple. By representing the various expressions and matchers with relatively high-level opcodes, it should be possible to avoid excessive dispatch overhead while avoiding the problems inherent to the recursive implementation.

The goals of this reimplementation, in no particular order, are:

* Avoid stack overflows due to too many recursive calls.
* More memory-efficient implementation.
* Better performance.
* Better error reporting with the "farthest failure position" technique.
* Less code, though not at the expense of readability (should still be simple code).
* Better isolation of implementation details vs exposed API, using a prefix to avoid clashes with user code for internal symbols.

## Overview

### Matchers

The parser generator would translate all literal matchers in the AST to a list of `Matcher` interfaces:

```
type Matcher interface {
    Match(peekReader) bool
}

// the parser would implement that interface
type peekReader interface {
    peek() savepoint // position and current rune
    read()           // advance to next rune
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

The API covered by the API stability guarantee in the doc will remain stable. Internal symbols not part of this API will use a prefix-naming-scheme to avoid clashes with user-defined code. The prefix will be U+03E1 `ϡ` ([see here][bird]). This is an interesting choice because:

* it looks a little bit like a bird. If you squint.
* my font correctly prints it (DejaVu Sans Mono)
* it is a valid letter that can start a Go symbol
* it is considered lowercase/not exported
* it is highly unlikely to clash with user code. You have to make a conscious effort to use a symbol that starts with this prefix, so accidental use of internal symbols is highly unlikely.
* it doesn't have any controversial meaning (looks like it [hardly has any meaning that we know of][sampi]).

The accepted PEG syntax remains exactly the same, with the same semantics.

### Code generation

Use go's `text/template` package and data structures to generate the code, instead of a string with fmt verbs.

## Opcodes

✓ CALL : pop I (I1), push the next instruction index to the I stack, jump to I1. Starts a new variable stack (?).
✓ CALLA N : pop V stack value and discard, pop P stack value and use to construct the current value, call action thunk at index N, push return value on the V strack.
✓ CALLB N : call boolean thunk at index N, push FAIL on the V stack if the thunk returned FALSE, TRUE otherwise.
✓ CUMULORF : pop 2 values from V (V and V-1), add V to the V-1 array (V-1 may be fail, replace with an array if that's the case), push to V. If V is FAIL, push FAIL instead of the cumulative array.
✓ EXIT : pop V, exit VM and return value V and true if V is not FAIL, return nil and false otherwise.
✓ JUMP N : inconditional jump to integer N.
✓ JUMPIFF N : jump to integer N if top V stack value is FAIL.
✓ JUMPIFT N : jump to integer N if top V stack value is not FAIL.
✓ MATCH N : save the start position, run the matcher at index N, if matcher returns true, push the slice of bytes from the start to the current parser position on stack V, otherwise push FAIL.
✓ NILIFF : pop top V stack value, push NIL if V is FAIL, FAIL otherwise.
✓ NILIFT : pop top V stack value, push NIL if V is not FAIL, FAIL otherwise.
✓ POPL : pop the top value from the L stack, discard.
✓ POPP : pop the top value from the P stack, discard.
✓ POPVJUMPIFF N : if top V stack value is FAIL, pop it and jump to integer N.
✓ PUSHI N : push integer N on the I stack.
✓ PUSHL n N... : push an array of n integers on the L stack.
✓ PUSHP : push the current parser position on the P stack.
✓ PUSHVE : push empty slice of interface{} on the V stack (typed nil).
✓ PUSHVF : push value FAIL on the V stack.
✓ PUSHVN : push value nil on the V stack.
✓ RESTORE : pop P stack value, restore the parser's position.
✓ RESTOREIFF : pop P, restore the parser's position if peek of top V stack value is FAIL, otherwise discard P.
✓ RETURN : pop I, jump to this instruction.
✓ STOREIFT N : pop top V stack value, if V is not FAIL store it in the current variable stack under label at index N, push V back on the V stack.
✓ TAKELORJUMP N : pop L, take one value off of the array of integers and push that value on the I stack, push L back. If L is empty, don't push anything to I, jump to N.

## Examples

Value may be the sentinel value MatchFailed, indicating no match. VM has four distinct stacks:

* Position stack (P) P[...]
* Instruction index stack (I) I[...]
* Value stack (V) V[...]
* Loop stack (L) L[...]

It also has four distinct lists:

* Matchers (M)
* Action thunks (A)
* Predicate thunks (B)
* Strings (S) : used for labels and to map instruction index to rule names.

The following statement always holds:

* A Matcher always consumes one `I` value and always produces one `V` value.

### Bootstrap sequence

0: PUSHI N : push N on instruction index stack, N = 3 I[3]
.: PUSHA
1: CALL : pop I, push next instruction index to I, jump to I I[2]
2: EXIT : pop V, decompose and return v, b (if V is MatchFailed, return nil, false, otherwise return V, true).

### E1 - Matcher

Grammar:

```
A <- 'a'
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
3: [Rule A, 'a'] PUSHP : save current parser position P[pa] I[2]
.:               MATCH 0 : run the matcher at 0, push V P[pa] I[2] V[va]
.:               RESTOREIFF : pop P, restore position if peek V is MatchFailed P[] I[2] V[va]
.:               RETURN : pop I, return P[] I[] V[va]

### E2 - Sequence

Grammar:

```
A <- 'a' 'b'
```

* M: 'a', 'b'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, Seq] PUSHP : P[ps] I[2]
04:               PUSHV fail : P[ps] I[2] V[f]
05:               PUSHL : P[ps] I[2] V[f] L[[Ia Ib]]
06:               TAKELORJUMP 11 : pop L, take one value off and push to I, push L. Jump to N if array is empty. P[ps] I[2 Ia] V[f] L[[Ib]]
07:               CALL : P[ps] I[2 8] V[f] L[[Ib]]
08:               CUMULORF : pop 2 values from V, cumulate in an array, push to V. If top V value is FAIL, pop 2 and push FAIL.
09:               JUMPIFF 11 : jump to N if top V stack is FAIL
10:               JUMP 06 : back to TAKE instruction
11:               POPL : P[ps] I[2] V[v] L[]
12:               RESTOREIFF : pop P, restore position if peek V is FAIL P[] I[2] V[v] L[]
13:               RETURN : pop I, return P[] I[] V[v] L[]

### E3 - Choice

Grammar:

```
A <- 'a' / 'b'
```

* M: 'a', 'b'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, Choice] PUSHL : P[] I[2] L[[Ia Ib]]
04:                  TAKELORJUMP 09 : P[] I[2 Ia] L[[Ib]]
05:                  CALL : P[] I[2 7] L[[Ib]]
06:                  JUMPIFT 10 : jump to N if top V stack is not FAIL
07:                  POPV (remove the FAIL)
08:                  JUMP 04
09:                  PUSHV FAIL (return fail)
10:                  POPL : P[] I[2] V[v] L[]
11:                  RETURN : P[] I[] V[v] L[]

### E4 - Repetition (`*`)

Grammar:

```
A <- 'a'*
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, ZoM] PUSHV empty : P[] I[2] V[ve]
04:               PUSHI Ia : P[] I[2 Ia] V[ve]
05:               CALL : P[] I[2 6] V[ve]
06:               POPVJUMPIFF 09 : if top V value is FAIL, pop it and jump to N
07:               CUMULORF : P[] I[2] V[ve]
08:               JUMP 04
09:               RETURN : P[] I[] V[ve]

### E5 - Repetition (`+`)

Grammar:

```
A <- 'a'+
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, ZoM] PUSHV fail : P[] I[2] V[f]
04:               PUSHI Ia : P[] I[2 Ia] V[f]
05:               CALL : P[] I[2 6] V[f]
06:               POPVJUMPIFF 09 : if top V value is FAIL, pop it and jump to N
07:               CUMULORF : P[] I[2] V[vf]
08:               JUMP 04
09:               RETURN : P[] I[] V[vf]

### E6 - Optional (`?`)

Grammar:

```
A <- 'a'?
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, ZoO] PUSHI Ia : P[] I[2 Ia] V[]
04:               CALL : P[] I[2 5] V[]
05:               POPVJUMPIFF 07 : if top V value is FAIL, pop it and jump to N
06:               RETURN : P[] I[] V[v]
07:               PUSHV nil : P[] I[2] V[n]
08:               RETURN : P[] I[] V[n]

### E7 - Rule reference

Grammar:

```
A <- B
B <- 'a'
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, Ref] PUSHI Ib : P[] I[2 Ib] V[]
04:               CALL : P[] I[2 5] V[]
05:               RETURN : P[] I[2] V[v]

### E8 - Predicates (and / not)

Grammar:

```
A <- !'a'
```

* M: 'a'
* A, B: none

Opcodes:

(bootstrap)
03: [Rule A, Not] PUSHP : P[p] I[2] V[]
04:               PUSHI Ia : P[p] I[2 Ia] V[]
05:               CALL : P[p] I[2 6] V[]
06:               NILIFF : pop V, push NIL if v is FAIL, FAIL otherwise P[p] I[2] V[b]
                  (NILIFT for and predicate)
07:               RESTORE : P[] I[2] V[b]
08:               RETURN : P[] I[2] V[b]

### E9 - Boolean thunks (and / not)

Grammar:

```
A <- !{return true, nil}
```

* M, A: none
* B: f1 {return true, nil}

Opcodes:

(bootstrap)
03: [Rule A, Not] CALLB 0 : call boolean thunk at index 0, push return value on V P[] I[2] V[]
04:               NILIFF : pop V, push NIL if v is FAIL, FAIL otherwise P[] I[2] V[b]
                  (NILIFT for and predicate)
05:               RETURN : P[] I[2] V[b]

### E10 - Labeled and action

Grammar:

```
A <- label:'a' { return nil, nil }
```

* M: 'a'
* A: f1 { return nil, nil }
* B: none

Opcodes:

(bootstrap)
03: [Rule A, Act] PUSHP : P[p] I[2] V[]
04:               PUSHI Il : P[p] I[2 Il] V[]
05:               CALL : P[p] I[2 6] V[]
06:               JUMPIFF 09 : P[p] I[2] V[v]
07:               CALLA 0 : pop V and discard, pop P, call action thunk at index 0, push return value on V P[] I[2] V[v]
08:               RETURN : P[] I[2] V[v]
09:               POPP : P[] I[2] V[f]
10:               RETURN : P[] I[2] V[f]

11: [Rule A, Lab] PUSHI Ia : P[p] I[2 6 Ia] V[]
12:               CALL : P[p] I[2 6 13] V[]
13:               STOREIFT lbl : store value in label on vstack if V is not FAIL. P[p] I[2 6] V[v]
14:               RETURN : P[p] I[2 6] V[v]

[ffp]: http://arxiv.org/abs/1405.6646
[bird]: http://unicode-table.com/en/03E1/
[sampi]: http://simple.wikipedia.org/wiki/Sampi_%28letter%29
