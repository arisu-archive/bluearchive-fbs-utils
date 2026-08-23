# BlueArchive FlatBuffers Utils

This is a utility for working with FlatBuffers in Go. It provides a set of functions for marshaling and unmarshaling FlatBuffers messages, as well as for converting between FlatBuffers and Go types.

## Installation

```bash
go get github.com/arisu-archive/bluearchive-fbs-utils
```

## Usage

```go
tableKey := []byte{1, 2}

decoded := fbsutils.Decode("QAI=", tableKey)
encoded := fbsutils.Encode(decoded, tableKey)

fmt.Println(decoded, encoded) // A QAI=
```

`Convert` remains available for backward compatibility and behaves like
`Decode`. New code should use the directional `Decode` and `Encode` functions.
For floats, `Encode` reverses values decoded from the protocol's positive wire
domain; legacy decoding leaves non-positive wire values unchanged.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
