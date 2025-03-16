package fbsutils

import (
	"encoding/base64"

	"github.com/arisu-archive/bluearchive-fbs-utils/pkg/mt19937"
	"github.com/arisu-archive/bluearchive-fbs-utils/pkg/xxhash32"
)

func CreateKey(name string, size int) []byte {
	// Apply xxhash32 to the name
	hashValue := xxhash32.Checksum([]byte(name))

	// Use the generated digest as MT19937 seed to get next 8 bytes
	mt := mt19937.NewWithSeed(hashValue)

	// Convert the next bytes to a string
	return mt.Bytes(size)
}

func CreateTableKey(name string) []byte {
	return CreateKey(name, 8)
}

func CreateZipPassword(name string) []byte {
	return []byte(base64.StdEncoding.EncodeToString(CreateKey(name, 15)))
}
