package testdata

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	input := Square{
		TopLeft: [2]int{2, 3},
		Width:   2,
		Height:  2,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	output, err := UnmarshalShapeJSON(b)
	if err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalNested(t *testing.T) {
	input := Union{
		A: Square{
			TopLeft: [2]int{2, 3},
			Width:   2,
			Height:  2,
		},
		B: Circle{
			Center: [2]int{6, 3},
			Radius: 4,
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	output, err := UnmarshalShapeJSON(b)
	if err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalStructField(t *testing.T) {
	input := Area{
		Color: "red",
		Shape: Square{
			TopLeft: [2]int{2, 3},
			Width:   2,
			Height:  2,
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var output Area
	if err := json.Unmarshal(b, &output); err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalStructFieldNull(t *testing.T) {
	input := Area{
		Color: "red",
		Shape: nil,
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var output Area
	if err := json.Unmarshal(b, &output); err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalStructFieldSlice(t *testing.T) {
	input := Pattern{
		Size: 12,
		Shapes: []Shape{
			Square{
				TopLeft: [2]int{2, 3},
				Width:   2,
				Height:  2,
			},
			Circle{
				Center: [2]int{6, 3},
				Radius: 4,
			},
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var output Pattern
	if err := json.Unmarshal(b, &output); err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalStructFieldMap(t *testing.T) {
	input := NamedPattern{
		Sizes: map[string]int{"foo": 12},
		Shapes: map[string]Shape{
			"foo": Square{
				TopLeft: [2]int{2, 3},
				Width:   2,
				Height:  2,
			},
			"bar": Circle{
				Center: [2]int{6, 3},
				Radius: 4,
			},
		},
	}
	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	var output NamedPattern
	if err := json.Unmarshal(b, &output); err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", input)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}

func TestMarshalStructFieldSkipField(t *testing.T) {
	expect := ShapeShifter{
		From:   nil,
		To:     nil,
		SkipMe: nil,
		Err:    "oh no!",
	}
	b := `{
		"SkipMe": {"kind": "Circle", "Center": [1, 2], "Radius": 12},
		"Err": "oh no!"
	}`
	var output ShapeShifter
	if err := json.Unmarshal([]byte(b), &output); err != nil {
		t.Logf("JSON: %s", b)
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expect, output) {
		t.Logf("JSON: %s", b)
		t.Logf("Exp: %#v", expect)
		t.Logf("Got: %#v", output)
		t.Fatal("Mismatched data")
	}
}
