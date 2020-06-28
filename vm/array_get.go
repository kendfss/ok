package vm

import (
	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/number"
)

// ArrayGet gets a value from the array by its index.
type ArrayGet struct {
	Array, Index, Result string
}

// Execute implements the Instruction interface for the VM.
func (ins *ArrayGet) Execute(registers map[string]*ast.Literal, _ *int, _ *VM) error {
	index := number.Int64(number.NewNumber(registers[ins.Index].Value))
	registers[ins.Result] = registers[ins.Array].Array[index]

	return nil
}
