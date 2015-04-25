package gen

import "text/template"

var tpl = template.New("gen")

// TODO : better name...
type genData struct {
	Init string // init code block

	As []struct {
		ReceiverName string
		RuleName     string
		ExprIndex    int
		Code         string
	}
}

const codeTpl = `
{{.Init}}

{{range .As}}
func ({{ .ReceiverName }} *current) on{{ .RuleName }}{{ .ExprIndex }}() (interface{}, error) {
{{ .Code }}	
}

func ({{ .ReceiverName }} *current) callOn{{ .RuleName }}{{ .ExprIndex }}() (interface{}, error) {
{{ .Code }}	
}
{{end}}
`

func init() {
	template.Must(tpl.Parse(codeTpl))
}
