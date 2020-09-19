package vm

import (
	"encoding/gob"
	"os"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/types"
)

// File is the root structure that will be serialized into the okc file.
type File struct {
	// Imports lists all the packages that this package relies on.
	Imports map[string]map[string]*types.Type

	Funcs      map[string]*CompiledFunc
	FuncDefs   map[string]*ast.Func
	Tests      []*CompiledTest
	Interfaces map[string]map[string]*types.Type
	Constants  map[string]*ast.Literal
}

func (f *File) ResolveType(t *types.Type) *types.Type {
	// TODO(elliot): Remove me in the future.
	if t.Name == "error.Error" {
		return types.ErrorInterface
	}

	if t.Kind == types.KindUnresolvedInterface {
		return types.NewInterface(t.Name, f.Interfaces[t.Name])
	}

	return t
}

// Store will create or replace the okc file for the provided package name.
func Store(file *File, packageName string) error {
	err := os.MkdirAll(Directory, 0755)
	if err != nil {
		return err
	}

	filePath := PathForPackage(packageName)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(file)
	if err != nil {
		return err
	}

	return nil
}

func Load(packageName string) (*File, error) {
	// Ignore packages that are build in (standard library).
	if p, ok := Packages[packageName]; ok {
		return p, nil
	}

	filePath := PathForPackage(packageName)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(f)
	var okcFile File
	err = decoder.Decode(&okcFile)
	if err != nil {
		return nil, err
	}

	return &okcFile, nil
}
