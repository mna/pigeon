package vm

import (
	"io/ioutil"
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
	if err := tpl.Execute(ioutil.Discard, data); err != nil {
		t.Error(err)
	}
}
