package gen

import (
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	data := genData{
		Init:         "this is init code",
		ReceiverName: "cc",
		As: []struct {
			RuleName  string
			ExprIndex int
			Args      []string
			Code      string
		}{
			{"ruleA", 1, nil, "this is code!"},
		},
	}
	if err := tpl.Execute(os.Stdout, data); err != nil {
		t.Error(err)
	}
}
