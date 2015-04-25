package gen

import "text/template"

var tpl = template.New("gen")

// TODO : better name...
type genData struct {
	Instrs       []uint64
	Init         string // init code block
	ReceiverName string

	As []struct {
		RuleName  string
		ExprIndex int
		Args      []string
		Code      string
	}
	Bs []struct {
		RuleName  string
		ExprIndex int
		Args      []string
		Code      string
	}
}

const codeTpl = `
{{.Init}}

{{range .As}}
func ({{ $.ReceiverName }} *current) on{{ .RuleName }}{{ .ExprIndex }}() (interface{}, error) {
{{ .Code }}	
}

func (v *ϡvm) callOn{{ .RuleName }}{{ .ExprIndex }}() (interface{}, error) {
{{range .Args}}

{{end}}
}
{{end}}
`

func init() {
	template.Must(tpl.Parse(codeTpl))
}
