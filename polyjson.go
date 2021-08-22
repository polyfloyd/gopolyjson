package polyjson

import (
	"fmt"
	"go/ast"
)

type Type struct {
	Name         string
	Variants     []string
	Discriminant string
}

func TypeFromInterface(files []*ast.File, interfaceName, discriminant string) (*Type, error) {
	// Find the interface denoted by the type name.
	var iface *ast.TypeSpec
	for typeSpec := range iterTypeSpecs(files) {
		if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
			if interfaceName == typeSpec.Name.Name {
				iface = typeSpec
				// No break, consume the iterator.
			}
		}
	}
	if iface == nil {
		return nil, fmt.Errorf("unable to locate interface declaration for %q", interfaceName)
	}

	// Find the function that is expected to be implemented by type variants.
	// It must adhere to these properties:
	// * No parameters
	// * No results
	// * Private (name starts with lowercase)
	var typeFunctionName string
	ifaceT := iface.Type.(*ast.InterfaceType)
	for _, methodField := range ifaceT.Methods.List {
		if fun, ok := methodField.Type.(*ast.FuncType); ok {
			if len(fun.Params.List) != 0 || fun.Results != nil {
				continue
			}
		}
		if len(methodField.Names) == 0 {
			continue
		}
		name := methodField.Names[0].Name
		if name[0] < 'a' || 'z' < name[0] {
			continue
		}
		typeFunctionName = name
	}
	if typeFunctionName == "" {
		return nil, fmt.Errorf("interface %q has no function that can be used to locate variants", iface.Name.Name)
	}

	var variantNames []string
	for typeSpec := range iterFuncDecls(files) {
		if typeSpec.Name.Name != typeFunctionName {
			continue
		}
		if len(typeSpec.Recv.List) != 1 {
			continue
		}
		if ident, ok := typeSpec.Recv.List[0].Type.(*ast.Ident); ok {
			variantNames = append(variantNames, ident.Name)
		}
	}
	if len(variantNames) == 0 {
		return nil, fmt.Errorf("interface %q has implementors")
	}

	return &Type{
		Name:         interfaceName,
		Variants:     variantNames,
		Discriminant: discriminant,
	}, nil
}

func iterTypeSpecs(files []*ast.File) <-chan *ast.TypeSpec {
	out := make(chan *ast.TypeSpec)
	go func() {
		defer close(out)
		for _, file := range files {
			for _, decl := range file.Decls {
				if d, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range d.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							out <- typeSpec
						}
					}
				}
			}
		}
	}()
	return out
}

func iterFuncDecls(files []*ast.File) <-chan *ast.FuncDecl {
	out := make(chan *ast.FuncDecl)
	go func() {
		defer close(out)
		for _, file := range files {
			for _, decl := range file.Decls {
				if d, ok := decl.(*ast.FuncDecl); ok {
					out <- d
				}
			}
		}
	}()
	return out
}
