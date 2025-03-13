package bitarray

import (
	"math"
	"testing"
)

func TestSetAll(t *testing.T) {
	t.Run("full bitarray", func(t *testing.T) {
		ba := BitArray(math.MaxUint)

		got := uint(ba.SetAll())
		var want uint = math.MaxUint

		if got != want {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
	})

	t.Run("empty bitarray", func(t *testing.T) {
		var ba BitArray

		got := uint(ba.SetAll())
		var want uint = math.MaxUint

		if got != want {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
	})
}

func TestResetAll(t *testing.T) {
	t.Run("full bitarray", func(t *testing.T) {
		ba := BitArray(math.MaxUint)

		got := uint(ba.ResetAll())
		var want uint = 0

		if got != want {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
	})

	t.Run("empty bitarray", func(t *testing.T) {
		var ba BitArray

		got := uint(ba.ResetAll())
		var want uint = 0

		if got != want {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
	})
}

func TestSetOn(t *testing.T) {
	var ba BitArray
	var want uint = 1
	for i := uint(0); i < Bits; i++ {
		want |= 1 << i
		got := ba.SetOn(i)

		if got != BitArray(want) {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
		ba = got
	}
}

func TestSetOff(t *testing.T) {
	ba := BitArray(math.MaxUint)
	var want uint = math.MaxUint
	for i := uint(0); i < Bits; i++ {
		want <<= 1
		got := ba.SetOff(i)

		if got != BitArray(want) {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
		ba = got
	}
}

func TestSetAt(t *testing.T) {
	var ba BitArray
	var want uint
	for i := uint(0); i < Bits; i++ {
		got := ba.SetAt(i, 1)
		want |= 1 << i

		if got != BitArray(want) {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
		ba = got
	}

	for i := Bits-1; i >= 0; i-- {
		got := ba.SetAt(uint(i), 0)
		want ^= 1<<i

		if got != BitArray(want) {
			t.Fatalf("got: %d; want: %d;", got, want)
		}
		ba = got
	}
}

func TestCountOn(t *testing.T) {
	var ba BitArray
	for i := 0; i <= Bits; i++ {
		testCountOn(t, ba, i)
		ba = ba<<1 | 1
	}

	for i := Bits; i >= 0; i-- {
		testCountOn(t, ba, i)
		ba = ba << 1
	}

	for i := 0; i < 256; i++ {
		for k := 0; k < 64-8; k++ {
			testCountOn(t, BitArray(i)<<k, tab[i].pop)
		}
	}
}

func testCountOn(t *testing.T, ba BitArray, want int) {
	got := ba.CountOn()
	if got != want {
		t.Fatalf("CountOn(%#016x) == %d; want %d", ba, got, want)
	}
}

func TestRotateLeft(t *testing.T) {
	m := BitArray(deBruijn)
	for k := uint(0); k < 128; k++ {
		x := m
		got := x.RotateLeft(k)
		want := x<<(k&0x3f) | x>>(64-k&0x3f)
		if got != want {
			t.Fatalf("RotateLeft(\n%#016x, \n%d) == \n%#016x; want \n%#016x", x, k, got, want)
		}
		got = got.RotateRight(k)
		if got != x {
			t.Fatalf("RotateRight(\n%#08x, \n%d) == \n%#08x; \nwant \n%#08x", want, k, got, x)
		}
	}
}

func TestMirror(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		testReverse(t, BitArray(1)<<i, BitArray(1)<<(63-i))
	}

	for _, test := range []struct {
		x, r BitArray
	}{
		{0, 0},
		{0x1, 0x8 << 60},
		{0x2, 0x4 << 60},
		{0x3, 0xc << 60},
		{0x4, 0x2 << 60},
		{0x5, 0xa << 60},
		{0x6, 0x6 << 60},
		{0x7, 0xe << 60},
		{0x8, 0x1 << 60},
		{0x9, 0x9 << 60},
		{0xa, 0x5 << 60},
		{0xb, 0xd << 60},
		{0xc, 0x3 << 60},
		{0xd, 0xb << 60},
		{0xe, 0x7 << 60},
		{0xf, 0xf << 60},
		{0x5686487, 0xe12616a000000000},
		{0x0123456789abcdef, 0xf7b3d591e6a2c480},
	} {
		testReverse(t, test.x, test.r)
		testReverse(t, test.r, test.x)
	}
}

func testReverse(t *testing.T, ba, want BitArray) {
	got := ba.Mirror()
	if got != want {
		t.Fatalf("Reverse(%#08x) == %#016x; want %#016x", ba, got, want)
	}

}

// ----------------------------------------------------------------------------
// Testing support

const deBruijn = 0x03f79d71b4ca8b09

type entry = struct {
	nlz, ntz, pop int
}

// tab contains results for all uint8 values
var tab [256]entry

func init() {
	tab[0] = entry{8, 8, 0}
	for i := 1; i < len(tab); i++ {
		// nlz
		x := i // x != 0
		n := 0
		for x&0x80 == 0 {
			n++
			x <<= 1
		}
		tab[i].nlz = n

		// ntz
		x = i // x != 0
		n = 0
		for x&1 == 0 {
			n++
			x >>= 1
		}
		tab[i].ntz = n

		// pop
		x = i // x != 0
		n = 0
		for x != 0 {
			n += int(x & 1)
			x >>= 1
		}
		tab[i].pop = n
	}
}
