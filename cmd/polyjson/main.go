package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	polyjson "github.com/polyfloyd/gopolyjson"
)

func main() {
	var typeArgs typesFlags
	flag.Var(&typeArgs, "type", "Specify a type and optional variant mappings. This flag can be specified multiple times for each type. Pattern: <interface>:[variant2[=jsonVariant2],]..")

	discriminant := flag.String("discriminant", "kind", "The name of the JSON field that holds the variant name")

	file := flag.String("file", "polyjsongen.go", "The name of the file written to the package")
	packagePath := flag.String("package", "", "The scoped package path")
	flag.Parse()

	if *packagePath == "" {
		fmt.Fprintf(os.Stderr, "a package path is required\n")
		os.Exit(1)
	}

	cfg := &packages.Config{Mode: packages.NeedSyntax | packages.NeedName}
	pkgs, err := packages.Load(cfg, *packagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(100)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(100)
	}
	pkg := pkgs[0]

	var types []polyjson.Type
	var typeNames []string
	for _, typeArg := range typeArgs {
		typeArg.Discriminant = *discriminant
		typ, err := polyjson.TypeFromInterface(pkg.Syntax, typeArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(200)
		}
		types = append(types, *typ)
		typeNames = append(typeNames, typ.Name)
	}

	structs, err := polyjson.PolymorphicStructFields(pkg.Syntax, typeNames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(300)
	}

	outFile := filepath.Join(*packagePath, *file)
	if err := polyjson.WriteMarshalerFile(outFile, pkg.Name, types, structs); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(300)
	}
}

type typesFlags []polyjson.TypeFromInterfaceArgs

func (*typesFlags) String() string {
	return "Type specifications"
}

func (f *typesFlags) Set(value string) error {
	iface := strings.Split(value, ":")
	if len(iface) == 1 {
		*f = append(*f, polyjson.TypeFromInterfaceArgs{
			Interface:    iface[0],
			VariantRemap: map[string]string{},
		})
		return nil
	}

	variants := map[string]string{}
	for _, variantGroup := range strings.Split(iface[1], ",") {
		if variantGroup == "" {
			continue
		}
		ss := strings.Split(variantGroup, "=")
		name, jsonName := ss[0], ss[0]
		if len(ss) == 2 {
			jsonName = ss[1]
		}
		variants[name] = jsonName
	}

	*f = append(*f, polyjson.TypeFromInterfaceArgs{
		Interface:    iface[0],
		VariantRemap: variants,
	})
	return nil
}
