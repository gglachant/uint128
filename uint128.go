// Copyright 2018 Weborama. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package uint128 // import "github.com/weborama/uint128"

//go:generate go run make_tables.go

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrStringTooLong = errors.New("string length greater than 32")

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

// MulWithOverflow multiplies x and y, returning the product and a boolean indicating if overflow occurred.
func (x Uint128) MulWithOverflow(y Uint128) (Uint128, bool) {
	return MulWithOverflow(x, y)
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
func NewFromString(str string) (x Uint128, err error) {
	x = Uint128{0, 0}
	//nolint:mnd // Number of characters in a hexadecimal representation of an uint128.
	if len(str) > 32 {
		return x, fmt.Errorf("s:%s: %w", str, ErrStringTooLong)
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
func (x Uint128) Format(fmtState fmt.State, verb rune) {
	switch verb {
	case 'v':
		x.formatV(fmtState)
	case 's':
		x.formatS(fmtState)
	case 'b':
		x.formatB(fmtState)
	case 'x':
		x.formatX(fmtState, false)
	case 'X':
		x.formatX(fmtState, true)
	case 'd', 'o':
		x.formatNotImplemented(fmtState, verb)
	case 'T':
		x.formatType(fmtState)
	default:
		x.formatUnknown(fmtState, verb)
	}
}

func (x Uint128) formatV(fmtState fmt.State) {
	switch {
	case fmtState.Flag('#'):
		fmt.Fprint(fmtState, x.String())
	case fmtState.Flag('+'):
		fmt.Fprint(fmtState, x.formatHL())
	default:
		fmt.Fprint(fmtState, x.formatHL())
	}
}

func (x Uint128) formatS(fmtState fmt.State) {
	fmt.Fprint(fmtState, x.String())
}

func (x Uint128) formatB(fmtState fmt.State) {
	binStr := fmt.Sprintf("%064b%064b", x.H, x.L)
	if fmtState.Flag('#') {
		fmt.Fprint(fmtState, "0b"+binStr)
	} else {
		fmt.Fprint(fmtState, binStr)
	}
}

func (x Uint128) formatX(fmtState fmt.State, upper bool) {
	hexStr := x.HexString()
	if upper {
		hexStr = strings.ToUpper(hexStr)
	}

	if fmtState.Flag('#') {
		prefix := "0x"
		if upper {
			prefix = "0X"
		}

		fmt.Fprint(fmtState, prefix+hexStr)
	} else {
		fmt.Fprint(fmtState, hexStr)
	}
}

func (x Uint128) formatNotImplemented(fmtState fmt.State, verb rune) {
	fmt.Fprintf(fmtState, "%%!%c(NOT_IMPLEMENTED)", verb)
}

func (x Uint128) formatType(fmtState fmt.State) {
	fmt.Fprintf(fmtState, "%T", x)
}

func (x Uint128) formatUnknown(fmtState fmt.State, verb rune) {
	fmt.Fprintf(fmtState, "%%!%c(Uint128=%s)", verb, x.String())
}

func (x Uint128) formatHL() string {
	hStr := strconv.FormatUint(x.H, 10)
	lStr := strconv.FormatUint(x.L, 10)

	return fmt.Sprintf("(H:%s, L:%s)", hStr, lStr)
}
