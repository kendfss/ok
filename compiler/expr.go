package compiler

import (
	"fmt"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/vm"
)

// compileExpr return the (result register, result type, error)
func compileExpr(compiledFunc *vm.CompiledFunc, expr ast.Node, file *Compiled) ([]vm.Register, []string, error) {
	switch e := expr.(type) {
	case *ast.Assign:
		err := compileAssign(compiledFunc, e, file)

		return nil, nil, err

	case *ast.Literal:
		returns := compiledFunc.NextRegister()
		compiledFunc.Append(&vm.Assign{
			VariableName: returns,
			Value:        e,
		})

		return []vm.Register{returns}, []string{e.Kind}, nil

	case *ast.Func:
		cf, err := CompileFunc(e, file)
		if err != nil {
			return nil, nil, err
		}

		file.FuncDefs[e.Name] = e
		file.Funcs[e.Name] = cf

		// TODO(elliot): Doesn't return true function type.
		fnType := "func (number, number) number"

		returns := compiledFunc.NextRegister()
		compiledFunc.Append(&vm.Assign{
			VariableName: returns,
			Value: &ast.Literal{
				Kind:  fnType,
				Value: e.Name,
			},
		})

		return []vm.Register{returns}, []string{fnType}, nil

	case *ast.Array:
		returns, kind, err := compileArray(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		return []vm.Register{returns}, []string{kind}, nil

	case *ast.Map:
		returns, err := compileMap(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		// TODO(elliot): Doesn't return type.
		return []vm.Register{returns}, []string{"{}"}, nil

	case *ast.Call:
		results, resultKinds, err := compileCall(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		// TODO(elliot): Doesn't return kind.
		return results, resultKinds, nil

	case *ast.Identifier:
		// TODO(elliot): Doesn't check that the upper scope variable exists or
		//  fetches the correct type.
		if e.Name[0] == '^' {
			return []vm.Register{vm.Register(e.Name)}, []string{"number"}, nil
		}

		if v, ok := compiledFunc.Variables[e.Name]; ok || e.Name[0] == '^' {
			return []vm.Register{vm.Register(e.Name)}, []string{v}, nil
		}

		return nil, nil, fmt.Errorf("%s undefined variable: %s",
			e.Pos, e.Name)

	case *ast.Binary:
		result, ty, err := compileBinary(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		return []vm.Register{result}, []string{ty}, nil

	case *ast.Group:
		return compileExpr(compiledFunc, e.Expr, file)

	case *ast.Unary:
		result, ty, err := compileUnary(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		return []vm.Register{result}, []string{ty}, nil

	case *ast.Key:
		result, ty, err := compileKey(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		return []vm.Register{result}, []string{ty}, nil

	case *ast.Interpolate:
		result, err := compileInterpolate(compiledFunc, e, file)
		if err != nil {
			return nil, nil, err
		}

		return []vm.Register{result}, []string{"string"}, nil
	}

	panic(expr)
}
