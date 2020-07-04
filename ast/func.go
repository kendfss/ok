package ast

import "strings"

// Argument is used to define a name and type for a function argument.
type Argument struct {
	Name string
	Type string
}

// Func represents the definition of a function.
type Func struct {
	// Name is the name of the function being declared.
	Name string

	// Arguments may contain zero or more elements. They will always be in the
	// order in which their are declared.
	Arguments []*Argument

	// Returns may contain zero or more types. They will always be in the order
	// in which they are declared.
	Returns []string

	// Statements can have zero or more elements for each of the ordered
	// discreet statements in the function.
	Statements []Node

	Pos string
}

// String returns the signature, like `Foo(x number, y number) (string, string)`.
func (f *Func) String() string {
	var args []string
	for _, arg := range f.Arguments {
		args = append(args, arg.Name+" "+arg.Type)
	}

	returnSignature := ""
	if len(f.Returns) == 1 {
		returnSignature = " " + f.Returns[0]
	}
	if len(f.Returns) > 1 {
		returnSignature = " (" + strings.Join(f.Returns, ", ") + ")"
	}

	return "func " + f.Name + "(" + strings.Join(args, ", ") + ")" + returnSignature
}

// Position returns the position.
func (f *Func) Position() string {
	return f.Pos
}
