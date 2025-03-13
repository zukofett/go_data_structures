package bitarray

type BitArray uint64

const (
	Bits = 64
)

func (ba BitArray) SetAll() BitArray {
	return ^BitArray(0)
}

func (ba BitArray) ResetAll() BitArray {
	return 0
}

func (ba BitArray) SetOn(at uint) BitArray {
	return ba | maskAt(at)
}

func (ba BitArray) SetOff(at uint) BitArray {
	return ba &^ maskAt(at)
}

func (ba BitArray) Set(at, value uint) BitArray {
	if value != 0 {
		return ba.SetOn(at)
	}
	return ba.SetOff(at)
}

func (ba BitArray) Get(at uint) bool {
	return maskAt(at)&ba != BitArray(0)
}

func (ba BitArray) FlipBit(at uint) BitArray {
	return ba ^ maskAt(at)
}

const (
	// --- Right Bit Masks ---
	rm0 = 0x5555555555555555 // 01010101 ...
	rm1 = 0x3333333333333333 // 00110011 ...
	rm2 = 0x0f0f0f0f0f0f0f0f // 00001111 ...
	rm3 = 0x00ff00ff00ff00ff // etc.
	rm4 = 0x0000ffff0000ffff
	rm5 = 0x00000000ffffffff

	// --- Left Bit Masks ---
	lm0 = 0xaaaaaaaaaaaaaaaa // 10101010 ...
	lm1 = 0xcccccccccccccccc // 11001100 ...
	lm2 = 0xf0f0f0f0f0f0f0f0 // 11110000 ...
	lm3 = 0xff00ff00ff00ff00 // etc.
	lm4 = 0xffff0000ffff0000
	lm5 = 0xffffffff00000000
)

func (ba BitArray) Mirror() BitArray {
	res := (ba&lm5)>>32 | (ba&rm5)<<32
	res = (res&lm4)>>16 | (res&rm4)<<16
	res = (res&lm3)>>8 | (res&rm3)<<8
	res = (res&lm2)>>4 | (res&rm2)<<4
	res = (res&lm1)>>2 | (res&rm1)<<2
	res = (res&lm0)>>1 | (res&rm0)<<1
	return res
}

func (ba BitArray) RotateRight(shift uint) BitArray {
	shift &= (Bits - 1)
    return ba>>shift | ba<<(Bits-shift)
}

func (ba BitArray) RotateLeft(shift uint) BitArray {
	return ba.RotateRight(Bits - shift%Bits)
}

func (ba BitArray) CountOn() int {
	mask := BitArray(0).SetAll()

	ba = ba>>1&(rm0&mask) + ba&(rm0&mask)
	ba = ba>>2&(rm1&mask) + ba&(rm1&mask)
	ba = (ba>>4 + ba) & (rm2 & mask)
	ba += ba >> 8
	ba += ba >> 16
	ba += ba >> 32
	return int(ba) & (1<<7 - 1)
}

func (ba BitArray) CountOff() int {
	return Bits - ba.CountOn()
}

func (ba BitArray) String() string {
	vals := []string{"0", "1"}

	ret := make([]rune, Bits)
	retLen := 0

	ba = ba.Mirror()

	for i := range Bits {
		retLen += copy(ret[retLen:], []rune(vals[1&(ba>>i)]))
	}

	return string(ret)
}

func maskAt(at uint) BitArray {
	return BitArray(1 << at)
}
