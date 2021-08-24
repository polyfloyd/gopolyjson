package polyjson

import (
	"fmt"
	"go/ast"
	"regexp"
)

type Type struct {
	Name         string
	Variants     []string
	Discriminant string
}

type (
	Struct struct {
		Name              string
		PolymorphicFields []StructField
	}
	StructField struct {
		Name     string // Name of the field in the struct.
		JSONName string // Name of the field as encoded in JSON.
		Type     string // Name of the polymorphic type.
		Kind     string // How the unmarshaling logic should behave. "Scalar" | "Slice" | "Map"
	}
)

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
		return nil, fmt.Errorf("interface %q has no function that can be used to locate variants", interfaceName)
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
		return nil, fmt.Errorf("interface %q has no implementors", interfaceName)
	}

	return &Type{
		Name:         interfaceName,
		Variants:     variantNames,
		Discriminant: discriminant,
	}, nil
}

func PolymorphicStructFields(files []*ast.File, polymorphicTypes []string) ([]Struct, error) {
	var structs []Struct
	for typeSpec := range iterTypeSpecs(files) {
		struc, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		var polymorphicFields []StructField
		for _, field := range struc.Fields.List {
			if ident, ok := field.Type.(*ast.Ident); ok && containsString(polymorphicTypes, ident.Name) {
				for _, name := range field.Names {
					polymorphicFields = append(polymorphicFields, StructField{
						Name:     name.Name,
						JSONName: jsonFieldName(field),
						Type:     ident.Name,
						Kind:     "Scalar",
					})
				}
				continue
			}
			if arr, ok := field.Type.(*ast.ArrayType); ok {
				if ident, ok := arr.Elt.(*ast.Ident); ok && containsString(polymorphicTypes, ident.Name) {
					for _, name := range field.Names {
						polymorphicFields = append(polymorphicFields, StructField{
							Name:     name.Name,
							JSONName: jsonFieldName(field),
							Type:     ident.Name,
							Kind:     "Slice",
						})
					}
				}
			}
			if mp, ok := field.Type.(*ast.MapType); ok {
				if keyIdent, ok := mp.Key.(*ast.Ident); ok && keyIdent.Name == "string" {
					if ident, ok := mp.Value.(*ast.Ident); ok && containsString(polymorphicTypes, ident.Name) {
						for _, name := range field.Names {
							polymorphicFields = append(polymorphicFields, StructField{
								Name:     name.Name,
								JSONName: jsonFieldName(field),
								Type:     ident.Name,
								Kind:     "Map",
							})
						}
					}
				}
			}
		}

		if len(polymorphicFields) > 0 {
			structs = append(structs, Struct{
				Name:              typeSpec.Name.Name,
				PolymorphicFields: polymorphicFields,
			})
		}
	}

	return structs, nil
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

func jsonFieldName(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}
	m := regexp.MustCompile(`json:"(.+),?.*"`).FindStringSubmatch(field.Tag.Value)
	if m == nil {
		return ""
	}
	return m[1]
}

func containsString(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
