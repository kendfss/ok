package compiler

import (
	"path"
	"strings"

	"github.com/elliotchance/ok/ast"
	"github.com/elliotchance/ok/parser"
	"github.com/elliotchance/ok/types"
	"github.com/elliotchance/ok/util"
	"github.com/elliotchance/ok/vm"
)

// Compile will return the compiled file. If there are any dependent packages
// they will also be compiled.
func Compile(
	rootPath,
	pkgPath string,
	includeTests bool,
	anonFunctionName int,
) (*vm.File, []error) {
	p := parser.NewParser(anonFunctionName)
	p.ParseDirectory(path.Join(rootPath, pkgPath), includeTests)
	if errs := p.Errors(); len(errs) > 0 {
		return nil, errs
	}

	packageName := util.PackageNameFromPath(rootPath, pkgPath)

	funcs := map[string]*ast.Func{}
	imports := map[string]*types.Type{}

	for _, pkgName := range p.Imports() {
		// TODO(elliot): Check import location exists.

		if _, ok := vm.Packages[pkgName]; ok {
			imports[pkgName] = vm.Packages[pkgName].Interface()
		} else {
			subFile, errs := Compile(rootPath, pkgName, false, anonFunctionName+10000)
			if len(errs) > 0 {
				return nil, errs
			}

			imports[pkgName] = subFile.Interface()
		}
	}

	for fnName, fn := range p.Funcs() {
		funcs[fnName] = fn
	}

	// The parser cannot resolve imported types, so we can fill those in now.
	//
	// TODO(elliot): This is pretty similar to the same implementation in the
	//  parser's resolveInterfaces().
	var err error
	for _, fn := range funcs {
		// Argument types
		for _, arg := range fn.Arguments {
			parts := strings.Split(arg.Type.Name, ".")
			if len(parts) == 2 {
				arg.Type = imports[parts[0]].Properties[parts[1]].Returns[0]
			}
		}

		// Return types
		for i := range fn.Returns {
			parts := strings.Split(fn.Returns[i].Name, ".")
			if len(parts) == 2 {
				fn.Returns[i] = imports[parts[0]].Properties[parts[1]].Returns[0]
			}
		}

		// TODO(elliot): This doesn't handle resolving in statements, see
		//  resolveInterfacesInStatements().
	}

	okcFile, err := compile(funcs, p.Tests(), p.Constants, imports)
	if err != nil {
		return nil, []error{err}
	}

	packageAlias := strings.ReplaceAll(packageName, "/", "__")

	file := &vm.File{
		Constants: p.Constants,
		Imports:   imports,
		Types:     map[vm.TypeRegister]*types.Type{},
		Symbols:   okcFile.Symbols,
	}

	for _, fn := range funcs {
		file.AddSymbolFunc(&vm.CompiledFunc{
			Type:       fn.Type(),
			Name:       fn.Name,
			UniqueName: fn.UniqueName,
		})
	}

	compiledPackageFn, err := CompileFunc(p.Package(packageAlias), file, nil)
	if err != nil {
		return nil, []error{err}
	}
	okcFile.PackageFunc = compiledPackageFn

	err = vm.Store(okcFile, packageName)
	if err != nil {
		return nil, []error{err}
	}

	return okcFile, nil
}
