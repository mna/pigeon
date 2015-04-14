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

### Rule name reference

The `parser.rules` map of names to rule nodes would not be required, a rule reference would be simply a jump to the opcode instruction of the start of that rule.

The `parser.rstack` slice of rules serves only to get the rule's identifier (or display name) in error messages. In the VM implementation, a simple mapping of instruction index to rule identifier (or display name) saves memory and achieves the same purpose.

### Variable sets (scope of labels)

The `parser.vstack` field holds the stack of variable sets - a mapping of label name to value that applies (is in scope) at a given position in the parser. In the VM implementation, a counter would keep track of the current scope depth, and the variable sets would be lazily created only when the first label in a scope is encountered. It would be stored in a `map[int]map[string]interface{}`, where the `int` is the scope depth.

On scope exit, if the value for that scope is not nil or an empty map, then the key would be deleted to avoid corruption if the parser goes back to that scope level.

### Memoization

