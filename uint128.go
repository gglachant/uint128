// Copyright 2018 Weborama. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package uint128 // import "github.com/weborama/uint128"

//go:generate go run make_tables.go

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	numBits     = 128
	numHalfBits = 64
)

// Uint128 defines an unsigned integer of 128 bits.
type Uint128 struct {
	H, L uint64
}

// Zero is a 0 valued Uint128.
func Zero() Uint128 {
	return Uint128{H: 0, L: 0}
}

// MaxUint128 returns the maximum value of an Uint128.
func MaxUint128() Uint128 {
	return Uint128{0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}
}

// Cmp compares two Uint128 and returns one of the following values:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func (x Uint128) Cmp(y Uint128) int {
	return Cmp(x, y)
}

// IsZero returns true if x is zero.
func (x Uint128) IsZero() bool {
	return IsZero(x)
}

// ShiftLeft shifts x to the left by the provided number of bits.
func (x Uint128) ShiftLeft(bits uint) Uint128 {
	return ShiftLeft(x, bits)
}

// ShiftRight shifts x to the right by the provided number of bits.
func (x Uint128) ShiftRight(bits uint) Uint128 {
	return ShiftRight(x, bits)
}

// And returns the logical AND of x and y.
func (x Uint128) And(y Uint128) Uint128 {
	return And(x, y)
}

// AndNot returns the logical AND NOT of x and y.
func (x Uint128) AndNot(y Uint128) Uint128 {
	return AndNot(x, y)
}

// Not returns the logical AND of x and y.
func (x Uint128) Not() Uint128 {
	return Not(x)
}

// Xor returns the logical XOR of x and y.
func (x Uint128) Xor(y Uint128) Uint128 {
	return Xor(x, y)
}

// Or returns the logical OR of x and y.
func (x Uint128) Or(y Uint128) Uint128 {
	return Or(x, y)
}

// Add adds x and y.
func (x Uint128) Add(y Uint128) Uint128 {
	return Add(x, y)
}

// Add128 returns the sum with carry of x, y and carry: sum = x + y + carry.
// The carry input must be 0 or 1; otherwise the behavior is undefined.
// The carryOut output is guaranteed to be 0 or 1.
func (x Uint128) Add128(y, carry Uint128) (sum, carryOut Uint128) {
	sum, carryOut = Add128(x, y, carry)

	return
}

// Incr increments x by one.
func (x Uint128) Incr() Uint128 {
	return Incr(x)
}

// Sub subtracts x and y.
func (x Uint128) Sub(y Uint128) Uint128 {
	return Sub(x, y)
}

// Decr decrements x by one.
func (x Uint128) Decr() Uint128 {
	return Decr(x)
}

// Mul multiplies x and y.
// Overflow is not checked.
func (x Uint128) Mul(y Uint128) Uint128 {
	return Mul(x, y)
}

// Div divides x by y.
func (x Uint128) Div(y Uint128) Uint128 {
	return Div(x, y)
}

// Mod returns x % y.
func (x Uint128) Mod(y Uint128) Uint128 {
	return Mod(x, y)
}

// RotateLeft rotates x left by k bits.
func (x Uint128) RotateLeft(k uint) Uint128 {
	return RotateLeft(x, k)
}

// RotateRight rotates x right by k bits.
func (x Uint128) RotateRight(k uint) Uint128 {
	return RotateRight(x, k)
}

// NewFromString creates a new Uint128 from its hexadecimal string representation.
// XXX: Do a proper job of it.
func NewFromString(str string) (x Uint128, err error) {
	x = Uint128{0, 0}
	// nolint: gomnd // Number of characters in a hexadecimal representation of an uint128
	if len(str) > 32 {
		return x, fmt.Errorf("s:%s length greater than 32", str) // nolint: goerr113
	}

	b, err := hex.DecodeString(fmt.Sprintf("%032s", str))
	if err != nil {
		return x, fmt.Errorf("hex.DecodeString(): %w", err)
	}

	rdr := bytes.NewReader(b)
	err = binary.Read(rdr, binary.BigEndian, &x)

	return
}

// HexString returns a Hexadecimal string representation of an Uint128.
// HexString returns a 32-character zero-padded lowercase hexadecimal string representation of an Uint128.
func (x Uint128) HexString() string {
	return fmt.Sprintf("%016x%016x", x.H, x.L)
}

// String returns a "0x" prefixed 32-character zero-padded lowercase hexadecimal string representation of an Uint128.
func (x Uint128) String() string {
	return "0x" + x.HexString()
}

// Format is a custom formatter for Uint128.
func (x Uint128) Format(fmtState fmt.State, c rune) {
	switch c {
	case 'v':
		if fmtState.Flag('#') { // %#v
			fmt.Fprint(fmtState, x.String())
		} else if fmtState.Flag('+') { // %+v
			// Format H and L as decimal strings individually first
			hStr := fmt.Sprintf("%d", x.H)
			lStr := fmt.Sprintf("%d", x.L)
			fmt.Fprintf(fmtState, "(H:%s, L:%s)", hStr, lStr)
		} else { // %v
			// Format H and L as decimal strings individually first
			hStr := fmt.Sprintf("%d", x.H)
			lStr := fmt.Sprintf("%d", x.L)
			fmt.Fprintf(fmtState, "(H:%s, L:%s)", hStr, lStr)
		}
	case 's':
		fmt.Fprint(fmtState, x.String())
	case 'b':
		binStr := fmt.Sprintf("%064b%064b", x.H, x.L)
		if fmtState.Flag('#') { // %#b
			fmt.Fprint(fmtState, "0b"+binStr)
		} else {
			fmt.Fprint(fmtState, binStr)
		}
	case 'x':
		hexStr := x.HexString()
		if fmtState.Flag('#') { // %#x
			fmt.Fprint(fmtState, "0x"+hexStr)
		} else {
			fmt.Fprint(fmtState, hexStr)
		}
	case 'X':
		hexStr := strings.ToUpper(x.HexString())
		if fmtState.Flag('#') { // %#X
			fmt.Fprint(fmtState, "0X"+hexStr)
		} else {
			fmt.Fprint(fmtState, hexStr)
		}
	case 'd', 'o':
		fmt.Fprintf(fmtState, "%%!%c(NOT_IMPLEMENTED)", c)
	case 'T':
		fmt.Fprintf(fmtState, "%T", x)
	default:
		// For unknown verbs, pass it to Sprintf but use default %v for the Uint128 part
		// This is a bit tricky, might be better to just error or use a default.
		// For now, let's try to be somewhat helpful.
		// Constructing format string like "%" + string(c)
		// fmt.Fprintf(fmtState, "%"+string(c), x.String()) // This might recurse or not be what user expects for all verbs
		// Fallback to simple string for other verbs for now.
		fmt.Fprintf(fmtState, "%%!%c(Uint128=%s)", c, x.String())
	}
}
