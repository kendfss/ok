package vm_test

import (
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/vm"

	"github.com/stretchr/testify/assert"
)

func TestAdd_Execute(t *testing.T) {
	for testName, test := range map[string]struct {
		left, right string
		expected    string
	}{
		"maintain-precision": {"1.2200", "4.7", "5.9200"},
	} {
		t.Run(testName, func(t *testing.T) {
			registers := map[string]*ast.Literal{
				"0": asttest.NewLiteralNumber(test.left),
				"1": asttest.NewLiteralNumber(test.right),
			}
			ins := &vm.Add{Left: "0", Right: "1", Result: "2"}
			assert.NoError(t, ins.Execute(registers, nil, nil))
			assert.Equal(t, test.expected, registers[ins.Result].Value)
		})
	}
}
