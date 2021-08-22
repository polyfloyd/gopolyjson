package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"polyjson"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	interfaces   = flag.String("interfaces", "", "The comma sepparated names of the interfaces denoting polymorphic types")
	discriminant = flag.String("discriminant", "kind", "The name of the JSON field that holds the variant name")
	file         = flag.String("file", "polyjsongen.go", "The name of the file written to the package")
	packagePath  = flag.String("package", "", "The scoped package path")
)

func main() {
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
	for _, name := range strings.Split(*interfaces, ",") {
		typ, err := polyjson.TypeFromInterface(pkg.Syntax, name, *discriminant)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(200)
		}
		types = append(types, *typ)
	}

	outFile := filepath.Join(*packagePath, *file)
	if err := polyjson.WriteMarshalerFile(outFile, pkg.Name, types); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(300)
	}
}
