Go PolyJSON
===========

[![Build Status](https://github.com/polyfloyd/gopolyjson/workflows/CI/badge.svg)](https://github.com/polyfloyd/gopolyjson/actions)

Go Code generator of JSON marshalers/unmarshalers for polymorphic
data structures.


## Usage
```
go install -v github.com/polyfloyd/gopolyjson/cmd/polyjson@latest
```

The program targets a single package and generates marshalers for one or more
polymorphic types along with unmarshalers for structures that use such
polymorphic types in fields.

The recommended usage is to define an interface denoting the polymorphic types
with at least a method that:
* Takes no arguments
* Returns nothing
* Is private (starts with a lowercase letter)

The program will automatically discover associated variants that implement this
interface.
```go
package shapes

type Shape interface {
	iShape()
}

func (Triangle) iShape() {}
func (Square) iShape()   {}

type Square struct {
	TopLeft       [2]int
	Width, Height int
}

type Triangle struct {
	P0 [2]int
	P1 [2]int
	P2 [2]int
}
```
Invoke the generator by specifying the name of the interface with `-type`:
```
polyjson -type Shape -package ./shapes
```
This will place a `polyjsongen.go` file in the package specified containing the
(un)marshalers.

Now, you can encode and decode your data like this:
```go
inputShape := Square{TopLeft: [2]int{1, 2}, Width: 4, Height: 4}

b, err := json.Marshal(inputShape)
fmt.Printf("%s\n", b) // {"kind": "Square", "TopLeft": [1, 2], "Width": 4, "Height": 4}

decodedShape, err := UnmarshalShapeJSON(b)
fmt.Printf("%T\n", decodedShape) // Square
````


## Reference

### Discriminant
Each struct that is encoded to JSON that is a variant of a polymorphic type
gains an additional field that holds the name of the type. It is placed in
between the other fields of the struct. The default name of this field is
`kind` and it can be altered by setting the `-discriminant` flag.

### Specifying types
The `-type` argument accepts a polymorphic type specification in the following forms:
* `-type Shape` Interface only
* `-type Shape:Triangle,Square` Interface with explicitly named variants
* `-type Shape:Triangle=triangle_shape,Square=square_shape` Interface with
  explicitly named variants and their JSON counterparts

### Usage in Go
To encode a struct, just use `json.Marshal` like you normally would.

To decode a polymorphic type, use the generated `Unmarshal<type>JSON`, e.g.
`UnmarshalShapeJSON`. This will probe the JSON object for the variant kind and
unmarshal into the associated struct.

### Interoperability with Rust and Serde
PolyJSON was initially made for exchanging data between Go and Rust services.
Rust in combination with [serde_json](https://crates.io/crates/serde-json) can
understand the format generated by this tool like so:

```rust
#[derive(Serialize, Deserialize)]
#[serde(tag = "kind")]
pub enum Shape {
	Square {
		top_left: [i32; 2],
		width: i32,
		height: i32,
	},
	Triangle {
		p0: [i32; 2],
		p1: [i32; 2],
		p2: [i32; 2],
	},
}
```
