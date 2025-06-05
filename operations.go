// Copyright 2018 Weborama. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package uint128 // import "github.com/weborama/uint128"

import "math/bits"

// Cmp compares two Uint128 and returns one of the following values:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
//
//nolint:varnamelen // idiomatic parameter naming.
func Cmp(x, y Uint128) int {
	if x.H < y.H {
		return -1
	} else if x.H > y.H {
		return 1
	}

	if x.L < y.L {
		return -1
	} else if x.L > y.L {
		return 1
	}

	return 0
}

// IsZero returns true if x is zero.
func IsZero(x Uint128) bool {
	return x.H == 0 && x.L == 0
}

// ShiftLeft shifts x to the left by the provided number of bits.
//
//nolint:varnamelen // idiomatic parameter naming.
func ShiftLeft(x Uint128, bits uint) Uint128 {
	switch {
	case bits >= numBits:
		x.H = 0
		x.L = 0
	case bits >= numHalfBits:
		x.H = x.L << (bits - numHalfBits)
		x.L = 0
	default:
		x.H <<= bits
		x.H |= x.L >> (numHalfBits - bits)
		x.L <<= bits
	}

	return x
}

// ShiftRight shifts x to the right by the provided number of bits.
//
//nolint:varnamelen // idiomatic parameter naming.
func ShiftRight(x Uint128, bits uint) Uint128 {
	switch {
	case bits >= numBits:
		x.H = 0
		x.L = 0
	case bits >= numHalfBits:
		x.L = x.H >> (bits - numHalfBits)
		x.H = 0
	default:
		x.L >>= bits
		x.L |= x.H << (numHalfBits - bits)
		x.H >>= bits
	}

	return x
}

// And returns the logical AND of x and y.
func And(x, y Uint128) Uint128 {
	x.H &= y.H
	x.L &= y.L

	return x
}

// AndNot returns the logical AND NOT of x and y.
func AndNot(x, y Uint128) Uint128 {
	x.H &^= y.H
	x.L &^= y.L

	return x
}

// Not returns the logical NOT of x and o.
func Not(x Uint128) Uint128 {
	x.H = ^x.H
	x.L = ^x.L

	return x
}

// Xor returns the logical XOR of x and y.
func Xor(x, y Uint128) Uint128 {
	x.H ^= y.H
	x.L ^= y.L

	return x
}

// Or returns the logical OR of x and y.
func Or(x, y Uint128) Uint128 {
	x.H |= y.H
	x.L |= y.L

	return x
}

// Add adds x and y.
func Add(x, y Uint128) Uint128 {
	sum, _ := Add128(x, y, Zero())

	return sum
}

// Add128 returns the sum with carry of x, y and carry: sum = x + y + carry.
// The carry input must be 0 or 1; otherwise the behavior is undefined.
// The carryOut output is guaranteed to be 0 or 1.
func Add128(x, y, carry Uint128) (sum, carryOut Uint128) {
	sum.L, carryOut.L = bits.Add64(x.L, y.L, carry.L)
	sum.H, carryOut.L = bits.Add64(x.H, y.H, carryOut.L)

	return
}

// Incr increments x by one.
func Incr(x Uint128) Uint128 {
	return Add(x, Uint128{H: 0, L: 1})
}

// Sub subtracts x and y.
//
//nolint:varnamelen // idiomatic parameter naming.
func Sub(x, y Uint128) Uint128 {
	pL := x.L
	x.L -= y.L
	x.H -= y.H

	if x.L > pL {
		x.H--
	}

	return x
}

// Decr decrements x by one.
//
//nolint:varnamelen // idiomatic parameter naming.
func Decr(x Uint128) Uint128 {
	pL := x.L
	x.L--

	if x.L > pL {
		x.H--
	}

	return x
}

// Len returns the minimum number of bits required to represent x; the result is 0 for x == 0.
func Len(x Uint128) int {
	if x.H == 0 {
		return bits.Len64(x.L)
	}

	return numHalfBits + bits.Len64(x.H)
}

// Mul multiplies x and y.
// Overflow is not checked.
//
//nolint:varnamelen // idiomatic parameter naming.
func Mul(x, y Uint128) Uint128 {
	// x = x1*2^64 + x0
	// y = y1*2^64 + y0
	// x*y = (x1*y1)*2^128 + (x1*y0)*2^64 + (x0*y1)*2^64 + x0*y0
	// (x1*y1)*2^128 will overflow, so we ignore it.
	h, l := bits.Mul64(x.L, y.L)
	h += x.H*y.L + x.L*y.H

	return Uint128{H: h, L: l}
}

// MulWithOverflow multiplies x and y, returning the product and a boolean indicating if overflow occurred.
//
//nolint:varnamelen // idiomatic parameter naming.
func MulWithOverflow(x, y Uint128) (Uint128, bool) {
	// x = x1*2^64 + x0
	// y = y1*2^64 + y0
	// x*y = (x1*y1)*2^128 + (x1*y0)*2^64 + (x0*y1)*2^64 + x0*y0
	// If x1*y1 != 0, then overflow definitely occurred (bits above 128)
	// Also, if the high 64 bits of the sum overflowed, that's overflow
	x0, x1 := x.L, x.H
	y0, y1 := y.L, y.H

	// Compute the partial products
	h, l := bits.Mul64(x0, y0)
	m1 := x1 * y0
	m2 := x0 * y1

	// Add the cross terms to the high part
	h, carry1 := bits.Add64(h, m1, 0)
	h, carry2 := bits.Add64(h, m2, 0)

	// If x1*y1 != 0, that's overflow
	// If carry1 or carry2 overflowed, that's overflow
	overflow := (x1 != 0 && y1 != 0) || carry1 != 0 || carry2 != 0

	return Uint128{H: h, L: l}, overflow
}

// Div divides x by y.
func Div(x, y Uint128) Uint128 {
	return divMod(x, y, true)
}

// Mod returns x % y.
func Mod(x, y Uint128) Uint128 {
	return divMod(x, y, false)
}

// divMod implements 128-bit division and modulo.
// This is a basic binary restoring division algorithm.
//
//nolint:varnamelen // idiomatic parameter naming.
func divMod(x, y Uint128, returnDiv bool) Uint128 {
	switch {
	case IsZero(y):
		panic("division by zero")
	case IsZero(x):
		if returnDiv {
			return Zero()
		}

		return Zero() // Mod(0, y) is 0
	case Cmp(x, y) < 0:
		if returnDiv {
			return Zero()
		}

		return x // Mod(x,y) is x if x < y
	case Cmp(x, y) == 0:
		if returnDiv {
			return Uint128{H: 0, L: 1}
		}

		return Zero() // Mod(x,x) is 0
	}

	// At this point, x > y and y != 0
	quotient, remainder := divModMainLoop(x, y)

	if returnDiv {
		return quotient
	}

	return remainder
}

// divModMainLoop performs the main binary restoring division loop.
//
//nolint:varnamelen // idiomatic parameter naming.
func divModMainLoop(x, y Uint128) (Uint128, Uint128) {
	remainder := Zero()
	quotient := Zero()

	var i uint
	for i = range numBits {
		shift := numBits - 1 - i
		remainder = ShiftLeft(remainder, 1)

		if (ShiftRight(x, shift)).L&1 == 1 {
			remainder.L |= 1
		}

		if Cmp(remainder, y) >= 0 {
			remainder = Sub(remainder, y)
			quotient = Or(quotient, ShiftLeft(Uint128{H: 0, L: 1}, shift))
		}
	}

	return quotient, remainder
}

// RotateLeft rotates x left by k bits.
//
//nolint:varnamelen // idiomatic parameter naming.
func RotateLeft(x Uint128, k uint) Uint128 {
	k %= numBits
	if k == 0 {
		return x
	}
	// result is (x << k) | (x >> (numBits - k))
	shiftedLeft := ShiftLeft(x, k)
	shiftedRight := ShiftRight(x, numBits-k)

	return Or(shiftedLeft, shiftedRight)
}

// RotateRight rotates x right by k bits.
//
//nolint:varnamelen // idiomatic parameter naming.
func RotateRight(x Uint128, k uint) Uint128 {
	k %= numBits
	if k == 0 {
		return x
	}
	// result is (x >> k) | (x << (numBits - k))
	shiftedRight := ShiftRight(x, k)
	shiftedLeft := ShiftLeft(x, numBits-k)

	return Or(shiftedRight, shiftedLeft)
}

// LeadingZeros returns the number of leading zero bits in x; the result is 128 for x == 0.
func LeadingZeros(x Uint128) int {
	return numBits - Len(x)
}

// OnesCount returns the number of one bits ("population count") in x.
func OnesCount(x Uint128) int {
	return bits.OnesCount64(x.H) + bits.OnesCount64(x.L)
}

// TrailingZeros returns the number of trailing zero bits in x; the result is 128 for x == 0.
func TrailingZeros(x Uint128) int {
	if x.L == 0 {
		return bits.TrailingZeros64(x.H) + numHalfBits
	}

	return bits.TrailingZeros64(x.L)
}

// Reverse returns the value of x with its bits in reversed order.
func Reverse(x Uint128) (y Uint128) {
	y.L, y.H = bits.Reverse64(x.H), bits.Reverse64(x.L)

	return
}

// ReverseBytes returns the value of x with its bytes in reversed order.
func ReverseBytes(x Uint128) (y Uint128) {
	y.L, y.H = bits.ReverseBytes64(x.H), bits.ReverseBytes64(x.L)

	return
}
