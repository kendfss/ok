package compiler

import (
	"fmt"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/ast/asttest"
	"github.com/elliotchance/ok/vm"
)

func compileMap(compiledFunc *vm.CompiledFunc, n *ast.Map, file *Compiled) (vm.Register, error) {
	// TODO(elliot): Check type is valid for the map.
	// TODO(elliot): Maps with duplicate keys should be an error.

	sizeRegister := compiledFunc.NextRegister()
	compiledFunc.Append(&vm.Assign{
		VariableName: sizeRegister,
		Value:        asttest.NewLiteralNumber(fmt.Sprintf("%d", len(n.Elements))),
	})

	mapRegister := compiledFunc.NextRegister()
	compiledFunc.Append(&vm.MapAlloc{
		// TODO(elliot): This needs to be derived from the actual type.
		Kind: "{}any",

		Size:   sizeRegister,
		Result: mapRegister,
	})

	for _, element := range n.Elements {
		// TODO(elliot): Check keyKind is string.
		keyRegisters, _, err := compileExpr(compiledFunc, element.Key, file)
		if err != nil {
			return "", err
		}

		// TODO(elliot): Check value is the right type for map.
		valueRegisters, _, _ := compileExpr(compiledFunc, element.Value, file)

		compiledFunc.Append(&vm.MapSet{
			Map:   mapRegister,
			Key:   keyRegisters[0],
			Value: valueRegisters[0],
		})
	}

	return mapRegister, nil
}
