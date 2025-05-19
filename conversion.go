package fbsutils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"reflect"
	"unicode/utf16"
)

type FlatData interface {
	InitKey(key []byte)
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
	FlatDataName() string
}

type FlatBuffer struct {
	TableKey []byte `json:"-"`
}

func (f *FlatBuffer) InitKey(key []byte) {
	f.TableKey = key
}

func Convert[T any](value T, tableKey []byte) T {
	switch v := any(value).(type) {
	case int32:
		return any(ConvertInt32(v, tableKey)).(T)
	case int64:
		return any(ConvertInt64(v, tableKey)).(T)
	case uint32:
		return any(ConvertUInt32(v, tableKey)).(T)
	case uint64:
		return any(ConvertUInt64(v, tableKey)).(T)
	case uint8:
		return any(ConvertUbyte(v, tableKey)).(T)
	case float32:
		return any(ConvertFloat32(v, tableKey)).(T)
	case float64:
		return any(ConvertFloat64(v, tableKey)).(T)
	case string:
		if v == "" {
			return value
		}
		return any(ConvertString(v, tableKey)).(T)
	default:
		// If type not supported, return the original value
		panic("type not supported:" + reflect.TypeOf(value).String())
	}
}

// XorBytes performs XOR operation between value and key bytes.
func XorBytes(value, key []byte) []byte {
	if len(value) == 0 || len(key) == 0 {
		return value
	}

	result := make([]byte, len(value))
	for i := range value {
		result[i] = value[i] ^ key[i%len(key)]
	}
	return result
}

// ConvertInt32 converts an int32 value using XOR.
func ConvertInt32(value int32, key []byte) int32 {
	if value == 0 {
		return 0
	}
	var result int32
	err := binary.Read(
		bytes.NewReader(XorBytes(binary.LittleEndian.AppendUint32(nil, uint32(value)), key)),
		binary.LittleEndian,
		&result,
	)
	if err != nil {
		return value
	}
	return result
}

// ConvertInt64 converts an int64 value using XOR.
func ConvertInt64(value int64, key []byte) int64 {
	if value == 0 {
		return 0
	}
	var result int64
	err := binary.Read(
		bytes.NewReader(XorBytes(binary.LittleEndian.AppendUint64(nil, uint64(value)), key)),
		binary.LittleEndian,
		&result,
	)
	if err != nil {
		return value
	}
	return result
}

// ConvertUInt32 converts a uint32 value using XOR.
func ConvertUInt32(value uint32, key []byte) uint32 {
	if value == 0 {
		return 0
	}
	var result uint32
	err := binary.Read(
		bytes.NewReader(XorBytes(binary.LittleEndian.AppendUint32(nil, value), key)),
		binary.LittleEndian,
		&result,
	)
	if err != nil {
		return value
	}
	return result
}

// ConvertUInt64 converts a uint64 value using XOR.
func ConvertUInt64(value uint64, key []byte) uint64 {
	if value == 0 {
		return 0
	}
	var result uint64
	err := binary.Read(
		bytes.NewReader(XorBytes(binary.LittleEndian.AppendUint64(nil, value), key)),
		binary.LittleEndian,
		&result,
	)
	if err != nil {
		return value
	}
	return result
}

// ConvertUbyte converts a uint8 value using XOR.
func ConvertUbyte(value uint8, key []byte) uint8 {
	if value == 0 {
		return 0
	}
	return value ^ key[0]
}

func calculateModulus(key []byte) int {
	if len(key) == 0 {
		return 1
	}
	modulus := int(key[0] % 10)
	if modulus <= 1 {
		modulus = 7
	}
	if key[0]&1 == 1 {
		modulus = -modulus
	}

	return modulus
}

// ConvertFloat32 converts a float32 value using XOR.
func ConvertFloat32(value float32, key []byte) float32 {
	modulus := calculateModulus(key)
	if value > 0 && modulus != 1 {
		return float32(value) / float32(modulus) / 10000
	}
	return value
}

// ConvertFloat64 converts a float64 value using XOR.
func ConvertFloat64(value float64, key []byte) float64 {
	modulus := calculateModulus(key)
	if value > 0 && modulus != 1 {
		return float64(value) / float64(modulus) / 10000
	}
	return value
}

// ConvertString converts a base64 encoded string.
func ConvertString(value string, key []byte) string {
	if len(value) == 0 {
		return ""
	}

	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return value
	}

	xorred := XorBytes(raw, key)
	// Convert UTF-16 bytes to runes
	runes := make([]uint16, len(xorred)/2)
	for i := 0; i < len(runes); i++ {
		runes[i] = uint16(xorred[i*2]) | (uint16(xorred[i*2+1]) << 8)
	}
	return string(utf16.Decode(runes))
}
