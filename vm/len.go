package vm

import (
	"fmt"
	"strings"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
)

// Len is used to determine the size of an array or map.
type Len struct {
	Argument, Result Register
}

// Execute implements the Instruction interface for the VM.
func (ins *Len) Execute(registers map[Register]*ast.Literal, _ *int, _ *VM) error {
	r := registers[ins.Argument]
	var result int
	switch {
	case strings.HasPrefix(r.Kind, "[]"):
		result = len(r.Array)

	case strings.HasPrefix(r.Kind, "{}"):
		result = len(r.Map)

	case r.Kind == "string":
		result = len([]rune(r.Value))

	default:
		result = len(r.Value)
	}

	registers[ins.Result] = asttest.NewLiteralNumber(fmt.Sprintf("%d", result))

	return nil
}

// String is the human-readable description of the instruction.
func (ins *Len) String() string {
	return fmt.Sprintf("%s = len(%s)", ins.Result, ins.Argument)
}
