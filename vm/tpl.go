package vm

import "text/template"

var tpl = template.New("gen")

// type thunkInfo struct {
// 	Parms  []string
// 	RuleNm string
// 	ExprIx int
// 	Code   string
// }
//
// type program struct {
// 	Init   string
// 	Instrs []ϡinstr
//
// 	Ms []ast.Expression
// 	As []*thunkInfo
// 	Bs []*thunkInfo
// 	Ss []string
//
// 	InstrToRule []int
//
// 	mss map[string]int // reverse map of string to index in Ss
// 	mms map[string]int // reverse map of matcher's raw value to index in Ms
//
// 	ruleNmIx      int
// 	exprIx        int
// 	parmsSet      [][]string  // stack of parms set for code blocks
// 	ruleNmStartIx map[int]int // rule name ix to first rule instr ix
// 	ruleNmEntryIx map[int]int // rule name ix to entry point instr ix
// 	ruleNmToDisNm map[int]int // rule name ix to rule display name ix
// }

const codeTpl = `
{{.Init}}

{{range .As}}
func ({{$.RecvrNm}} *current) on{{.RuleNm}}{{.ExprIx}}() (interface{}, error) {
{{ .Code }}	
}

func (v *ϡvm) callOn{{.RuleNm}}{{.ExprIx}}() (interface{}, error) {
{{range .Parms}}

{{end}}
	v.cur.on{{.RuleNm}}{{.ExprIx}}()
}
{{end}}

{{range .Bs}}
{{end}}
`

func init() {
	template.Must(tpl.Parse(codeTpl))
}
