package vm

import (
	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/types"
)

// A Symbol contains a literal value that can be referenced by instructions.
type Symbol struct {
	Type  string
	Value string `json:",omitempty"`

	// Interface will only be set when Func is provided.
	Interface string        `json:",omitempty"`
	Func      *CompiledFunc `json:",omitempty"`
}

func (s *Symbol) Literal() *ast.Literal {
	return &ast.Literal{
		Kind:  types.TypeFromString(s.Type),
		Value: s.Value,
	}
}
