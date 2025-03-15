# BlueArchive FlatBuffers Utils

This is a utility for working with FlatBuffers in Go. It provides a set of functions for marshaling and unmarshaling FlatBuffers messages, as well as for converting between FlatBuffers and Go types.

## Installation

```bash
go get github.com/arisu-archive/bluearchive-fbs-utils
```

## Usage

```go
package main

import (
	"fmt"

	fbsutils "github.com/arisu-archive/bluearchive-fbs-utils"
)

func main() {
	tableKey := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	converted := fbsutils.Convert(1, tableKey)
	fmt.Println(converted)
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
