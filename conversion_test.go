package fbsutils

import "testing"

func TestDecodeString(t *testing.T) {
	t.Parallel()

	key := []byte{1, 2}
	const encoded = "QAI="

	got := Decode(encoded, key)
	if got != "A" {
		t.Errorf("Decode(%q, %v) = %q, want %q", encoded, key, got, "A")
	}
}

func TestEncodeString(t *testing.T) {
	t.Parallel()

	key := []byte{1, 2}
	const decoded = "A"

	got := Encode(decoded, key)
	if got != "QAI=" {
		t.Errorf("Encode(%q, %v) = %q, want %q", decoded, key, got, "QAI=")
	}
}

func TestStringUnicode(t *testing.T) {
	t.Parallel()

	key := []byte{0}
	const (
		decoded = "学生🙂"
		encoded = "ZlsfdT3YQt4="
	)

	if got := Encode(decoded, key); got != encoded {
		t.Errorf("Encode(%q, %v) = %q, want %q", decoded, key, got, encoded)
	}
	if got := Decode(encoded, key); got != decoded {
		t.Errorf("Decode(%q, %v) = %q, want %q", encoded, key, got, decoded)
	}
}

func TestEncodeFloat32(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float32
		key   []byte
		want  float32
	}{
		{name: "positive modulus", value: 5, key: []byte{4}, want: 200000},
		{name: "negative modulus", value: -4, key: []byte{3}, want: 120000},
		{name: "empty key", value: 5, want: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := Encode(test.value, test.key)
			if got != test.want {
				t.Errorf("Encode(%v, %v) = %v, want %v", test.value, test.key, got, test.want)
			}
		})
	}
}

func TestEncodeFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		key   []byte
		want  float64
	}{
		{name: "positive modulus", value: 5, key: []byte{4}, want: 200000},
		{name: "negative modulus", value: -4, key: []byte{3}, want: 120000},
		{name: "empty key", value: 5, want: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := Encode(test.value, test.key)
			if got != test.want {
				t.Errorf("Encode(%v, %v) = %v, want %v", test.value, test.key, got, test.want)
			}
		})
	}
}

func TestDecodeFloats(t *testing.T) {
	t.Parallel()

	if got := Decode(float32(200000), []byte{4}); got != 5 {
		t.Errorf("Decode(float32(200000), %v) = %v, want %v", []byte{4}, got, float32(5))
	}
	if got := Decode(float32(120000), []byte{3}); got != -4 {
		t.Errorf("Decode(float32(120000), %v) = %v, want %v", []byte{3}, got, float32(-4))
	}
	if got := Decode(float64(200000), []byte{4}); got != 5 {
		t.Errorf("Decode(float64(200000), %v) = %v, want %v", []byte{4}, got, float64(5))
	}
	if got := Decode(float64(120000), []byte{3}); got != -4 {
		t.Errorf("Decode(float64(120000), %v) = %v, want %v", []byte{3}, got, float64(-4))
	}
}

func TestEncodeIntegers(t *testing.T) {
	t.Parallel()

	key32 := []byte{1, 0, 0, 0}
	if got := Encode(int32(2), key32); got != 3 {
		t.Errorf("Encode(int32(2), %v) = %d, want %d", key32, got, 3)
	}
	if got := Encode(uint32(2), key32); got != 3 {
		t.Errorf("Encode(uint32(2), %v) = %d, want %d", key32, got, 3)
	}

	key64 := []byte{1, 0, 0, 0, 0, 0, 0, 0}
	if got := Encode(int64(2), key64); got != 3 {
		t.Errorf("Encode(int64(2), %v) = %d, want %d", key64, got, 3)
	}
	if got := Encode(uint64(2), key64); got != 3 {
		t.Errorf("Encode(uint64(2), %v) = %d, want %d", key64, got, 3)
	}

	key8 := []byte{1}
	if got := Encode(uint8(2), key8); got != 3 {
		t.Errorf("Encode(uint8(2), %v) = %d, want %d", key8, got, 3)
	}
	if got := Encode(uint8(2), nil); got != 2 {
		t.Errorf("Encode(uint8(2), nil) = %d, want %d", got, 2)
	}
	if got := Decode(uint8(2), nil); got != 2 {
		t.Errorf("Decode(uint8(2), nil) = %d, want %d", got, 2)
	}
}

func TestEncodeSmallIntegers(t *testing.T) {
	t.Parallel()

	// int8(-5) = 0xFB; 0xFB ^ 0x01 = 0xFA = int8(-6).
	if got := Encode(int8(-5), []byte{1}); got != -6 {
		t.Errorf("Encode(int8(-5), %v) = %d, want %d", []byte{1}, got, -6)
	}
	// int16(258) = LE [0x02 0x01]; XOR [0xFF 0x00] = [0xFD 0x01] = int16(509).
	key16 := []byte{0xFF, 0x00}
	if got := Encode(int16(258), key16); got != 509 {
		t.Errorf("Encode(int16(258), %v) = %d, want %d", key16, got, 509)
	}
	// Single-byte key must cycle across both bytes:
	// int16(-2) = LE [0xFE 0xFF]; XOR [0x01 0x01] = [0xFF 0xFE] = int16(-257).
	if got := Encode(int16(-2), []byte{1}); got != -257 {
		t.Errorf("Encode(int16(-2), %v) = %d, want %d", []byte{1}, got, -257)
	}
	if got := Encode(uint16(2), []byte{1, 0}); got != 3 {
		t.Errorf("Encode(uint16(2), %v) = %d, want %d", []byte{1, 0}, got, 3)
	}
}

func TestSmallIntegerRoundTrip(t *testing.T) {
	t.Parallel()

	key := []byte{0x5A, 0xC3}
	if got := Decode(Encode(int8(-77), key), key); got != -77 {
		t.Errorf("Decode(Encode(int8(-77))) = %d, want %d", got, -77)
	}
	if got := Decode(Encode(int16(-12345), key), key); got != -12345 {
		t.Errorf("Decode(Encode(int16(-12345))) = %d, want %d", got, -12345)
	}
	if got := Decode(Encode(uint16(54321), key), key); got != 54321 {
		t.Errorf("Decode(Encode(uint16(54321))) = %d, want %d", got, 54321)
	}
}

func TestSmallIntegerGuards(t *testing.T) {
	t.Parallel()

	key := []byte{1}
	if got := Encode(int8(0), key); got != 0 {
		t.Errorf("Encode(int8(0), %v) = %d, want 0", key, got)
	}
	if got := Encode(int16(0), key); got != 0 {
		t.Errorf("Encode(int16(0), %v) = %d, want 0", key, got)
	}
	if got := Encode(uint16(0), key); got != 0 {
		t.Errorf("Encode(uint16(0), %v) = %d, want 0", key, got)
	}
	if got := Encode(int8(7), nil); got != 7 {
		t.Errorf("Encode(int8(7), nil) = %d, want 7", got)
	}
	if got := Encode(int16(7), nil); got != 7 {
		t.Errorf("Encode(int16(7), nil) = %d, want 7", got)
	}
	if got := Encode(uint16(7), nil); got != 7 {
		t.Errorf("Encode(uint16(7), nil) = %d, want 7", got)
	}
}

func TestBooleanIdentity(t *testing.T) {
	t.Parallel()

	if got := Encode(true, []byte{4}); !got {
		t.Errorf("Encode(true, %v) = %t, want true", []byte{4}, got)
	}
	if got := Decode(false, []byte{4}); got {
		t.Errorf("Decode(false, %v) = %t, want false", []byte{4}, got)
	}
}

func TestConvertCompatibility(t *testing.T) {
	t.Parallel()

	key := []byte{1, 2}
	const encoded = "QAI="

	if got := Convert(encoded, key); got != "A" {
		t.Errorf("Convert(%q, %v) = %q, want %q", encoded, key, got, "A")
	}
	if got := Convert(true, key); !got {
		t.Errorf("Convert(true, %v) = %t, want true", key, got)
	}
}
