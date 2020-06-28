package vm_test

import (
	"errors"
	"testing"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/number"
	"github.com/elliotchance/ok/vm"

	"github.com/stretchr/testify/assert"
)

func TestDivide_Execute(t *testing.T) {
	for testName, test := range map[string]struct {
		left, right string
		expected    string
		err         error
	}{
		"success":     {"1.2200", "4.7", "0.25957446808510638298", nil},
		"divide-zero": {"1.2200", "0", "0", errors.New("division by zero")},
	} {
		t.Run(testName, func(t *testing.T) {
			registers := map[string]*ast.Literal{
				"0": ast.NewLiteralNumber(test.left),
				"1": ast.NewLiteralNumber(test.right),
			}
			ins := &vm.Divide{Left: "0", Right: "1", Result: "2"}
			assert.Equal(t, test.err, ins.Execute(registers, nil, nil))
			actual := number.NewNumber(registers[ins.Result].Value)
			assert.Equal(t, test.expected, number.Format(actual, -1))
		})
	}
}
