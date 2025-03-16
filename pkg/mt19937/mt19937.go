package mt19937

import (
	"errors"
	"time"
)

var ErrInvalidBytesLength = errors.New("invalid bytes length")

const (
	n         uint32 = 624
	m         uint32 = 397
	f         uint32 = 1812433253
	matrixA   uint32 = 0x9908b0df
	upperMask uint32 = 0x80000000
	lowerMask uint32 = 0x7fffffff
)

type MT19937 struct {
	mt    [n]uint32
	index uint32
}

func New() *MT19937 {
	return NewWithSeed(uint32(time.Now().Unix())) //nolint:gosec // This is obviously in bounds
}

func NewWithSeed(seed uint32) *MT19937 {
	rnd := &MT19937{
		mt:    [n]uint32{},
		index: 0,
	}

	rnd.mt[0] = seed
	for i := uint32(1); i < n; i++ {
		rnd.mt[i] = f*(rnd.mt[i-1]^(rnd.mt[i-1]>>30)) + i
	}
	return rnd
}

func (rnd *MT19937) UInt32() uint32 {
	if rnd.index == 0 {
		rnd.twist()
	}

	y := rnd.mt[rnd.index]
	y ^= y >> 11
	y ^= y << 7 & 0x9d2c5680
	y ^= y << 15 & 0xefc60000
	y ^= y >> 18
	rnd.index = (rnd.index + 1) % n
	return y
}

func (rnd *MT19937) Int31() int32 {
	return int32(rnd.UInt32() >> 1) //nolint:gosec // The shift is taking care of the sign bit
}

func (rnd *MT19937) Bytes(length int) []byte {
	if length <= 0 {
		panic(ErrInvalidBytesLength)
	}

	buf := make([]byte, length)
	for i := 0; i < length; i += 4 {
		val := rnd.Int31()
		for j := 0; j < 4 && i+j < length; j++ {
			// SHOULD BE LITTLE ENDIAN
			buf[i+j] = byte(val & 0xff)
			val >>= 8
		}
	}
	return buf
}

func (rnd *MT19937) twist() {
	for i := range n {
		y := (rnd.mt[i] & upperMask) | (rnd.mt[(i+1)%n] & lowerMask)
		rnd.mt[i] = rnd.mt[(i+m)%n] ^ (y >> 1)
		if y%2 != 0 {
			rnd.mt[i] ^= matrixA
		}
	}
}
