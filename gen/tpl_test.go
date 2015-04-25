package gen

import (
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	data := genData{
		Init: "this is init code",
		As: []struct {
			ReceiverName string
			RuleName     string
			ExprIndex    int
			Code         string
		}{
			{"c", "ruleA", 1, "this is code!"},
		},
	}
	if err := tpl.Execute(os.Stdout, data); err != nil {
		t.Error(err)
	}
}
