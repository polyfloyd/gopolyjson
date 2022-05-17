package main

import (
	"reflect"
	"testing"

	"github.com/polyfloyd/go-polyjson"
)

func TestTypeFlag(t *testing.T) {
	testCases := []struct {
		Flag   string
		Expect *polyjson.TypeFromInterfaceArgs
	}{
		{
			Flag: "Shape",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{},
			},
		},
		{
			Flag: "Shape:",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{},
			},
		},
		{
			Flag: "Shape:Foo",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{"Foo": "Foo"},
			},
		},
		{
			Flag: "Shape:Foo=foo",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{"Foo": "foo"},
			},
		},
		{
			Flag: "Shape:Foo=foo,Bar=bar",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{"Foo": "foo", "Bar": "bar"},
			},
		},
		{
			Flag: "Shape:Foo,Bar,",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{"Foo": "Foo", "Bar": "Bar"},
			},
		},
		{
			Flag: "Shape:Foo,Bar=bar",
			Expect: &polyjson.TypeFromInterfaceArgs{
				Interface:    "Shape",
				VariantRemap: map[string]string{"Foo": "Foo", "Bar": "bar"},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.Flag, func(t *testing.T) {
			var types typesFlags
			if err := types.Set(test.Flag); err != nil {
				if test.Expect != nil {
					t.Fatal(err)
				}
				return
			}
			if test.Expect == nil {
				t.Fatal("An error was expected")
			}
			if !reflect.DeepEqual(types[0], *test.Expect) {
				t.Logf("Exp: %#v", *test.Expect)
				t.Logf("Got: %#v", types[0])
				t.Fatal("Mismatched data")
			}
		})
	}
}
