// Copyright 2018 Weborama. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package uint128_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/weborama/uint128"
)

func TestNewFromString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input     string
		expected  uint128.Uint128
		expectErr bool
	}{
		// Exactly 32 hex chars
		{"00000000000000000000000000000001", uint128.Uint128{H: 0, L: 1}, false},
		// Exactly 32 invalid hex chars
		{"0000000000000000000000000000000R", uint128.Uint128{}, true},
		// Fewer than 32 hex chars, should pad
		{"abc", uint128.Uint128{H: 0, L: 0xabc}, false},
		// 32 chars, nonzero high
		{"00000000000000010000000000000000", uint128.Uint128{H: 1, L: 0}, false},
		// Too long
		{"000000000000000000000000000000000", uint128.Uint128{}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()

			result, err := uint128.NewFromString(testCase.input)
			if testCase.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q, got none", testCase.input)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error for input %q: %v", testCase.input, err)
			}

			if result != testCase.expected {
				t.Errorf("input %q: got %v, want %v", testCase.input, result, testCase.expected)
			}
		})
	}
}

func TestUint128Operations(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label                 string
		input                 uint128.Uint128
		expectedLen           int
		expectedLeadingZeros  int
		expectedOnesCount     int
		expectedTrailingZeros int
		expectedReverse       uint128.Uint128
		expectedReverseBytes  uint128.Uint128
		expectedIncr          uint128.Uint128
		expectedDecr          uint128.Uint128
	}{
		{
			"{H: 0x0, L: 0x456}",
			uint128.Uint128{H: 0x0, L: 0x456},
			11, 117, 5, 1,
			uint128.Uint128{H: 0x6a20000000000000, L: 0x0},
			uint128.Uint128{H: 0x5604000000000000, L: 0x0},
			uint128.Uint128{H: 0x0, L: 0x456},
			uint128.Uint128{H: 0x0, L: 0x456},
		},
		{
			"{H: 0x1, L: 0x456}",
			uint128.Uint128{H: 0x1, L: 0x456},
			65, 63, 6, 1,
			uint128.Uint128{H: 0x6a20000000000000, L: 0x8000000000000000},
			uint128.Uint128{H: 0x5604000000000000, L: 0x100000000000000},
			uint128.Uint128{H: 0x1, L: 0x456},
			uint128.Uint128{H: 0x1, L: 0x456},
		},
		{
			"{H: 0x1, L: 0x0}",
			uint128.Uint128{H: 0x1, L: 0x0},
			65, 63, 1, 64,
			uint128.Uint128{H: 0x0, L: 0x8000000000000000},
			uint128.Uint128{H: 0x0, L: 0x100000000000000},
			uint128.Uint128{H: 0x1, L: 0x1},
			uint128.Uint128{H: 0x1, L: 0xFFFFFFFFFFFFFFFF},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			var rval int

			var ruint uint128.Uint128

			rval = uint128.Len(testCase.input)
			if rval != testCase.expectedLen {
				t.Fatalf("Len - Expected:%d Got:%d", testCase.expectedLen, rval)
			}

			rval = uint128.LeadingZeros(testCase.input)
			if rval != testCase.expectedLeadingZeros {
				t.Fatalf("LeadingZeros - Expected:%d Got:%d", testCase.expectedLeadingZeros, rval)
			}

			rval = uint128.OnesCount(testCase.input)
			if rval != testCase.expectedOnesCount {
				t.Fatalf("OnesCount - Expected:%d Got:%d", testCase.expectedOnesCount, rval)
			}

			rval = uint128.TrailingZeros(testCase.input)
			if rval != testCase.expectedTrailingZeros {
				t.Fatalf("TrailingZeros - Expected:%d Got:%d", testCase.expectedTrailingZeros, rval)
			}

			ruint = uint128.Reverse(testCase.input)
			if ruint != testCase.expectedReverse {
				t.Fatalf("Reverse - Expected:%v Got:%v", testCase.expectedReverse, ruint)
			}

			ruint = uint128.ReverseBytes(testCase.input)
			if ruint != testCase.expectedReverseBytes {
				t.Fatalf("ReverseBytes - Expected:%v Got:%v", testCase.expectedReverseBytes, ruint)
			}
		})
	}
}

func TestAdd128(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label         string
		x, y, carry   uint128.Uint128
		sum, carryOut uint128.Uint128
	}{
		{
			label: "zero sum should result zero",
			x:     uint128.Zero(),
			y:     uint128.Zero(),
			sum:   uint128.Zero(),
		},
		{
			label: "zero+1 sum should result 1",
			x:     uint128.Zero(),
			y:     uint128.Uint128{H: 0, L: 1},
			sum:   uint128.Uint128{H: 0, L: 1},
		},
		{
			label: "1+maxuint64 sum should overflow from Low to High",
			x:     uint128.Uint128{H: 0, L: 1},
			y:     uint128.Uint128{H: 0, L: math.MaxUint64},
			sum:   uint128.Uint128{H: 1, L: 0},
		},
		{
			label:    "1+maxuint128 sum should overflow",
			x:        uint128.Uint128{H: 0, L: 1},
			y:        uint128.MaxUint128(),
			sum:      uint128.Uint128{H: 0, L: 0},
			carryOut: uint128.Uint128{H: 0, L: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			sum, carryOut := testCase.x.Add128(testCase.y, testCase.carry)

			assert(t, sum, testCase.sum, "unexpected sum")
			assert(t, carryOut, testCase.carryOut, "unexpected carryOut")
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label string
		x, y  uint128.Uint128
		sum   uint128.Uint128
	}{
		{
			label: "zero sum should result zero",
			x:     uint128.Zero(),
			y:     uint128.Zero(),
			sum:   uint128.Zero(),
		},
		{
			label: "zero+1 sum should result 1",
			x:     uint128.Zero(),
			y:     uint128.Uint128{H: 0, L: 1},
			sum:   uint128.Uint128{H: 0, L: 1},
		},
		{
			label: "1+maxuint64 sum should overflow from Low to High",
			x:     uint128.Uint128{H: 0, L: 1},
			y:     uint128.Uint128{H: 0, L: math.MaxUint64},
			sum:   uint128.Uint128{H: 1, L: 0},
		},
		{
			label: "1+maxuint128 sum should overflow",
			x:     uint128.Uint128{H: 0, L: 1},
			y:     uint128.MaxUint128(),
			sum:   uint128.Uint128{H: 0, L: 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			sum := testCase.x.Add(testCase.y)

			assert(t, sum, testCase.sum, "unexpected sum")
		})
	}
}

func TestIncr(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label string
		x     uint128.Uint128
		incr  uint128.Uint128
	}{
		{
			label: "zero++ sum should result 1",
			x:     uint128.Zero(),
			incr:  uint128.Uint128{H: 0, L: 1},
		},
		{
			label: "maxuint64++ sum should overflow from Low to High",
			x:     uint128.Uint128{H: 0, L: math.MaxUint64},
			incr:  uint128.Uint128{H: 1, L: 0},
		},
		{
			label: "maxuint128++ sum should overflow",
			x:     uint128.MaxUint128(),
			incr:  uint128.Uint128{H: 0, L: 0},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			incr := testCase.x.Incr()

			assert(t, incr, testCase.incr, "unexpected incr")
		})
	}
}

func assert(t *testing.T, got, expected uint128.Uint128, message string) {
	t.Helper()

	if got != expected {
		t.Errorf("error: %s (got: %v, expected: %v)", message, got, expected)
	}
}

// assertFormat is a helper for TestFormat.
func assertFormat(t *testing.T, val any, formatSpec, expected, testName string) {
	t.Helper()

	got := fmt.Sprintf(formatSpec, val)
	if got != expected {
		t.Errorf("%s: for format %q of %#v, got %q, want %q",
			testName, formatSpec, val, got, expected)
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	valZero := uint128.Zero()
	valABC := uint128.Uint128{L: 0xabc}                         // Decimal 2748
	valMixed := uint128.Uint128{H: 0x12, L: 0xdef0000000000000} // H=18, L=16063007335082819584
	valMax := uint128.MaxUint128()

	// Helper to generate expected binary strings
	binStr := func(u uint128.Uint128) string {
		return fmt.Sprintf("%064b%064b", u.H, u.L)
	}

	testCases := []struct {
		name     string
		val      any
		format   string
		expected string
	}{
		// Note: no testing of %p verb (pointer) as this is not deterministic.

		// Zero value
		{"Zero %s", valZero, "%s", "0x00000000000000000000000000000000"},
		{"Zero %v", valZero, "%v", "(H:0, L:0)"},
		{"Zero %+v", valZero, "%+v", "(H:0, L:0)"},
		{"Zero %#v", valZero, "%#v", "0x00000000000000000000000000000000"},
		{"Zero %x", valZero, "%x", "00000000000000000000000000000000"},
		{"Zero %#x", valZero, "%#x", "0x00000000000000000000000000000000"},
		{"Zero %X", valZero, "%X", "00000000000000000000000000000000"},
		{"Zero %#X", valZero, "%#X", "0X00000000000000000000000000000000"},
		{"Zero %b", valZero, "%b", binStr(valZero)},
		{"Zero %#b", valZero, "%#b", "0b" + binStr(valZero)},
		{"Zero %d", valZero, "%d", "%!d(NOT_IMPLEMENTED)"},
		{"Zero %o", valZero, "%o", "%!o(NOT_IMPLEMENTED)"},
		{"Zero %T", valZero, "%T", "uint128.Uint128"}, // Note that this is covered by the fmt package.
		{"Zero %?", valZero, "%?", fmt.Sprintf("%%!?(Uint128=%s)", valZero.String())},

		// valABC
		{"ABC %s", valABC, "%s", "0x00000000000000000000000000000abc"},
		{"ABC %v", valABC, "%v", "(H:0, L:2748)"},
		{"ABC %+v", valABC, "%+v", "(H:0, L:2748)"},
		{"ABC %#v", valABC, "%#v", "0x00000000000000000000000000000abc"},
		{"ABC %x", valABC, "%x", "00000000000000000000000000000abc"},
		{"ABC %#x", valABC, "%#x", "0x00000000000000000000000000000abc"},
		{"ABC %X", valABC, "%X", "00000000000000000000000000000ABC"},
		{"ABC %#X", valABC, "%#X", "0X00000000000000000000000000000ABC"},
		{"ABC %b", valABC, "%b", binStr(valABC)},
		{"ABC %#b", valABC, "%#b", "0b" + binStr(valABC)},
		{"ABC %d", valABC, "%d", "%!d(NOT_IMPLEMENTED)"},
		{"ABC %o", valABC, "%o", "%!o(NOT_IMPLEMENTED)"},
		{"ABC %T", valABC, "%T", "uint128.Uint128"}, // Note that this is covered by the fmt package.
		{"ABC %?", valABC, "%?", fmt.Sprintf("%%!?(Uint128=%s)", valABC.String())},

		// valMixed
		{"Mixed %s", valMixed, "%s", "0x0000000000000012def0000000000000"},
		{"Mixed %v", valMixed, "%v", "(H:18, L:16064339870830559232)"},
		{"Mixed %+v", valMixed, "%+v", "(H:18, L:16064339870830559232)"},
		{"Mixed %#v", valMixed, "%#v", "0x0000000000000012def0000000000000"},
		{"Mixed %x", valMixed, "%x", "0000000000000012def0000000000000"},
		{"Mixed %#x", valMixed, "%#x", "0x0000000000000012def0000000000000"},
		{"Mixed %X", valMixed, "%X", "0000000000000012DEF0000000000000"},
		{"Mixed %#X", valMixed, "%#X", "0X0000000000000012DEF0000000000000"},
		{"Mixed %b", valMixed, "%b", binStr(valMixed)},
		{"Mixed %#b", valMixed, "%#b", "0b" + binStr(valMixed)},
		{"Mixed %d", valMixed, "%d", "%!d(NOT_IMPLEMENTED)"},
		{"Mixed %o", valMixed, "%o", "%!o(NOT_IMPLEMENTED)"},
		{"Mixed %T", valMixed, "%T", "uint128.Uint128"}, // Note that this is covered by the fmt package.
		{"Mixed %?", valMixed, "%?", fmt.Sprintf("%%!?(Uint128=%s)", valMixed.String())},

		// valMax
		{"Max %s", valMax, "%s", "0xffffffffffffffffffffffffffffffff"},
		{"Max %v", valMax, "%v", "(H:18446744073709551615, L:18446744073709551615)"},
		{"Max %+v", valMax, "%+v", "(H:18446744073709551615, L:18446744073709551615)"},
		{"Max %#v", valMax, "%#v", "0xffffffffffffffffffffffffffffffff"},
		{"Max %x", valMax, "%x", "ffffffffffffffffffffffffffffffff"},
		{"Max %#x", valMax, "%#x", "0xffffffffffffffffffffffffffffffff"},
		{"Max %X", valMax, "%X", "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},
		{"Max %#X", valMax, "%#X", "0XFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},
		{"Max %b", valMax, "%b", binStr(valMax)},
		{"Max %#b", valMax, "%#b", "0b" + binStr(valMax)},
		{"Max %d", valMax, "%d", "%!d(NOT_IMPLEMENTED)"},
		{"Max %o", valMax, "%o", "%!o(NOT_IMPLEMENTED)"},
		{"Max %T", valMax, "%T", "uint128.Uint128"}, // Note that this is covered by the fmt package.
		{"Max %?", valMax, "%?", fmt.Sprintf("%%!?(Uint128=%s)", valMax.String())},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assertFormat(t, testCase.val, testCase.format, testCase.expected, testCase.name)
		})
	}
}

func TestDivMod(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label        string
		x, y         uint128.Uint128
		expectedDiv  uint128.Uint128
		expectedMod  uint128.Uint128
		expectPanic  bool
		specificTest string // "Div" or "Mod" for panic tests
	}{
		{
			label:       "Div by zero",
			x:           uint128.Uint128{L: 1},
			y:           uint128.Zero(),
			expectPanic: true,
		},
		{
			label:       "Mod by zero",
			x:           uint128.Uint128{L: 1},
			y:           uint128.Zero(),
			expectPanic: true,
		},
		{
			label:       "Div 0 by non-zero",
			x:           uint128.Zero(),
			y:           uint128.Uint128{L: 5},
			expectedDiv: uint128.Zero(),
			expectedMod: uint128.Zero(),
		},
		{
			label:       "Div by one",
			x:           uint128.Uint128{L: 100},
			y:           uint128.Uint128{L: 1},
			expectedDiv: uint128.Uint128{L: 100},
			expectedMod: uint128.Zero(),
		},
		{
			label:       "Div number by itself",
			x:           uint128.Uint128{L: 77},
			y:           uint128.Uint128{L: 77},
			expectedDiv: uint128.Uint128{L: 1},
			expectedMod: uint128.Zero(),
		},
		{
			label:       "Simple Div and Mod",
			x:           uint128.Uint128{L: 10},
			y:           uint128.Uint128{L: 3},
			expectedDiv: uint128.Uint128{L: 3},
			expectedMod: uint128.Uint128{L: 1},
		},
		{
			label:       "Div with high part non-zero",
			x:           uint128.Uint128{H: 2, L: 0}, // 2 * 2^64
			y:           uint128.Uint128{L: 2},
			expectedDiv: uint128.Uint128{H: 1, L: 0}, // 1 * 2^64
			expectedMod: uint128.Zero(),
		},
		{
			label:       "Mod with remainder non-zero",
			x:           uint128.Uint128{H: 1, L: 5}, // 1 * 2^64 + 5
			y:           uint128.Uint128{L: 2},
			expectedDiv: uint128.Uint128{H: 0, L: (1 << 63) + 2}, // equivalent of (2^64 + 5) / 2 = 2^63 + 2
			expectedMod: uint128.Uint128{L: 1},
		},
		{
			label:       "X < Y",
			x:           uint128.Uint128{L: 5},
			y:           uint128.Uint128{L: 10},
			expectedDiv: uint128.Zero(),
			expectedMod: uint128.Uint128{L: 5},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			if testCase.expectPanic {
				// Test Div for panic
				t.Run("DivPanic", func(t *testing.T) {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("Div did not panic as expected")
						}
					}()

					_ = testCase.x.Div(testCase.y)
				})
				// Test Mod for panic
				t.Run("ModPanic", func(t *testing.T) {
					defer func() {
						if r := recover(); r == nil {
							t.Errorf("Mod did not panic as expected")
						}
					}()

					_ = testCase.x.Mod(testCase.y)
				})
			} else {
				divResult := testCase.x.Div(testCase.y)
				assert(t, divResult, testCase.expectedDiv, "unexpected Div result")

				modResult := testCase.x.Mod(testCase.y)
				assert(t, modResult, testCase.expectedMod, "unexpected Mod result")
			}
		})
	}
}

func TestMul(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label    string
		x, y     uint128.Uint128
		expected uint128.Uint128
	}{
		{
			label:    "0 * 0 = 0",
			x:        uint128.Zero(),
			y:        uint128.Zero(),
			expected: uint128.Zero(),
		},
		{
			label:    "1 * 0 = 0",
			x:        uint128.Uint128{H: 0, L: 1},
			y:        uint128.Zero(),
			expected: uint128.Zero(),
		},
		{
			label:    "0 * 1 = 0",
			x:        uint128.Zero(),
			y:        uint128.Uint128{H: 0, L: 1},
			expected: uint128.Zero(),
		},
		{
			label:    "1 * 1 = 1",
			x:        uint128.Uint128{H: 0, L: 1},
			y:        uint128.Uint128{H: 0, L: 1},
			expected: uint128.Uint128{H: 0, L: 1},
		},
		{
			label:    "2 * 3 = 6",
			x:        uint128.Uint128{H: 0, L: 2},
			y:        uint128.Uint128{H: 0, L: 3},
			expected: uint128.Uint128{H: 0, L: 6},
		},
		{
			label:    "math.MaxUint64 * 2",
			x:        uint128.Uint128{H: 0, L: math.MaxUint64},
			y:        uint128.Uint128{H: 0, L: 2},
			expected: uint128.Uint128{H: 1, L: math.MaxUint64 - 1},
		},
		{
			label:    "Test high part non-zero",
			x:        uint128.Uint128{H: 1, L: 0}, // 2^64
			y:        uint128.Uint128{H: 0, L: 2},
			expected: uint128.Uint128{H: 2, L: 0}, // 2^65
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			result := testCase.x.Mul(testCase.y)
			assert(t, result, testCase.expected, "unexpected product")
		})
	}
}

func TestMulWithOverflow(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		label    string
		x, y     uint128.Uint128
		expected uint128.Uint128
		overflow bool
	}{
		{
			label:    "0 * 0 = 0 (no overflow)",
			x:        uint128.Zero(),
			y:        uint128.Zero(),
			expected: uint128.Zero(),
			overflow: false,
		},
		{
			label:    "1 * 0 = 0 (no overflow)",
			x:        uint128.Uint128{H: 0, L: 1},
			y:        uint128.Zero(),
			expected: uint128.Zero(),
			overflow: false,
		},
		{
			label:    "math.MaxUint64 * 2 (no overflow)",
			x:        uint128.Uint128{H: 0, L: math.MaxUint64},
			y:        uint128.Uint128{H: 0, L: 2},
			expected: uint128.Uint128{H: 1, L: math.MaxUint64 - 1},
			overflow: false,
		},
		{
			label:    "High part non-zero, no overflow",
			x:        uint128.Uint128{H: 1, L: 0},
			y:        uint128.Uint128{H: 0, L: 2},
			expected: uint128.Uint128{H: 2, L: 0},
			overflow: false,
		},
		{
			label:    "Overflow: high-high multiplication",
			x:        uint128.Uint128{H: 1, L: 0},
			y:        uint128.Uint128{H: 1, L: 0},
			expected: uint128.Uint128{H: 0, L: 0}, // lower 128 bits
			overflow: true,
		},
		{
			label: "Overflow: cross terms",
			x:     uint128.Uint128{H: math.MaxUint64, L: math.MaxUint64},
			y:     uint128.Uint128{H: 1, L: 1},
			expected: func() uint128.Uint128 {
				r, _ := uint128.Uint128{H: math.MaxUint64, L: math.MaxUint64}.MulWithOverflow(uint128.Uint128{H: 1, L: 1})

				return r
			}(),
			overflow: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.label, func(t *testing.T) {
			t.Parallel()

			result, overflow := testCase.x.MulWithOverflow(testCase.y)
			if result != testCase.expected {
				t.Errorf("unexpected product: got %v, want %v", result, testCase.expected)
			}

			if overflow != testCase.overflow {
				t.Errorf("unexpected overflow: got %v, want %v", overflow, testCase.overflow)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		val      uint128.Uint128
		k        uint
		leftExp  uint128.Uint128
		rightExp uint128.Uint128
	}{
		{
			name:     "Rotate 0 by 0",
			val:      uint128.Zero(),
			k:        0,
			leftExp:  uint128.Zero(),
			rightExp: uint128.Zero(),
		},
		{
			name:     "Rotate 1 by 0",
			val:      uint128.Uint128{L: 1},
			k:        0,
			leftExp:  uint128.Uint128{L: 1},
			rightExp: uint128.Uint128{L: 1},
		},
		{
			name:     "Rotate 1 by 1",
			val:      uint128.Uint128{L: 1},
			k:        1,
			leftExp:  uint128.Uint128{L: 2},
			rightExp: uint128.Uint128{H: 1 << 63}, // Wraps around
		},
		{
			name:     "Rotate by 64",
			val:      uint128.Uint128{L: 0x123456789abcdef0},
			k:        64,
			leftExp:  uint128.Uint128{H: 0x123456789abcdef0, L: 0},
			rightExp: uint128.Uint128{H: 0x123456789abcdef0, L: 0},
		},
		{
			name:     "Rotate by 128",
			val:      uint128.Uint128{H: 0x1, L: 0x2},
			k:        128,
			leftExp:  uint128.Uint128{H: 0x1, L: 0x2},
			rightExp: uint128.Uint128{H: 0x1, L: 0x2},
		},
		{
			name:     "Rotate by 129 (equiv 1)",
			val:      uint128.Uint128{L: 1},
			k:        129,
			leftExp:  uint128.Uint128{L: 2},
			rightExp: uint128.Uint128{H: 1 << 63},
		},
		{
			name:     "Rotate value with H and L by 1",
			val:      uint128.Uint128{H: 1, L: 1},
			k:        1,
			leftExp:  uint128.Uint128{H: 2, L: 2},
			rightExp: uint128.Uint128{H: 1 << 63, L: (1 >> 1) | (1 << 63)},
		},
		{
			name:     "Rotate value with H and L by 32",
			val:      uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB},
			k:        32,
			leftExp:  uint128.Or(uint128.ShiftLeft(uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB}, 32), uint128.ShiftRight(uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB}, 128-32)),
			rightExp: uint128.Or(uint128.ShiftRight(uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB}, 32), uint128.ShiftLeft(uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB}, 128-32)),
		},
		{
			name:     "Rotate value with H and L by 63",
			val:      uint128.Uint128{H: 0x1, L: 0x8000000000000000},
			k:        63,
			leftExp:  uint128.Or(uint128.ShiftLeft(uint128.Uint128{H: 0x1, L: 0x8000000000000000}, 63), uint128.ShiftRight(uint128.Uint128{H: 0x1, L: 0x8000000000000000}, 128-63)),
			rightExp: uint128.Or(uint128.ShiftRight(uint128.Uint128{H: 0x1, L: 0x8000000000000000}, 63), uint128.ShiftLeft(uint128.Uint128{H: 0x1, L: 0x8000000000000000}, 128-63)),
		},
		{
			name:     "Rotate value with H and L by 65",
			val:      uint128.Uint128{H: 0x1, L: 0x2},
			k:        65,
			leftExp:  uint128.Or(uint128.ShiftLeft(uint128.Uint128{H: 0x1, L: 0x2}, 65), uint128.ShiftRight(uint128.Uint128{H: 0x1, L: 0x2}, 128-65)),
			rightExp: uint128.Or(uint128.ShiftRight(uint128.Uint128{H: 0x1, L: 0x2}, 65), uint128.ShiftLeft(uint128.Uint128{H: 0x1, L: 0x2}, 128-65)),
		},
		{
			name:     "Rotate 0xFF...FF by 1",
			val:      uint128.MaxUint128(),
			k:        1,
			leftExp:  uint128.MaxUint128(),
			rightExp: uint128.MaxUint128(),
		},
		{
			name:     "Rotate 0x80...00 by 1 (left)",
			val:      uint128.Uint128{H: 1 << 63},
			k:        1,
			leftExp:  uint128.Uint128{L: 1},
			rightExp: uint128.Uint128{H: 1 << 62},
		},
		{
			name:     "Rotate 0x0...01 by 1 (right)",
			val:      uint128.Uint128{L: 1},
			k:        1,
			leftExp:  uint128.Uint128{L: 2},
			rightExp: uint128.Uint128{H: 1 << 63},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			leftActual := testCase.val.RotateLeft(testCase.k)
			assert(t, leftActual, testCase.leftExp, "unexpected RotateLeft result")

			rightActual := testCase.val.RotateRight(testCase.k)
			assert(t, rightActual, testCase.rightExp, "unexpected RotateRight result")
		})
	}
}

func TestCmp(t *testing.T) {
	t.Parallel()

	var (
		zero   = uint128.Zero()
		one    = uint128.Uint128{L: 1}
		two    = uint128.Uint128{L: 2}
		maxL   = uint128.Uint128{L: math.MaxUint64}
		minH   = uint128.Uint128{H: 1}
		max128 = uint128.MaxUint128()
	)

	testCases := []struct {
		name     string
		x, y     uint128.Uint128
		expected int
	}{
		{"zero == zero", zero, zero, 0},
		{"zero < one", zero, one, -1},
		{"one > zero", one, zero, 1},
		{"one == one", one, one, 0},
		{"one < two", one, two, -1},
		{"two > one", two, one, 1},
		{"maxL < minH", maxL, minH, -1},
		{"minH > maxL", minH, maxL, 1},
		{"maxL == maxL", maxL, maxL, 0},
		{"minH == minH", minH, minH, 0},
		{"max128 == max128", max128, max128, 0},
		{"zero < max128", zero, max128, -1},
		{"max128 > zero", max128, zero, 1},
		{"one < maxL", one, maxL, -1},
		{"maxL > one", maxL, one, 1},
		{"one < minH", one, minH, -1},
		{"minH > one", minH, one, 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := testCase.x.Cmp(testCase.y)
			if got != testCase.expected {
				t.Errorf("error: %s Cmp(%v, %v) (got: %d, expected: %d)", testCase.name, testCase.x, testCase.y, got, testCase.expected)
			}

			gotPkg := uint128.Cmp(testCase.x, testCase.y)
			if gotPkg != testCase.expected {
				t.Errorf("error: %s uint128.Cmp(%v, %v) (got: %d, expected: %d)", testCase.name, testCase.x, testCase.y, gotPkg, testCase.expected)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	var (
		zero   = uint128.Zero()
		one    = uint128.Uint128{L: 1}
		minH   = uint128.Uint128{H: 1}
		max128 = uint128.MaxUint128()
	)

	testCases := []struct {
		name     string
		val      uint128.Uint128
		expected bool
	}{
		{"Zero is zero", zero, true},
		{"One is not zero", one, false},
		{"MinH is not zero", minH, false},
		{"Max is not zero", max128, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := testCase.val.IsZero()
			if got != testCase.expected {
				t.Errorf("error: %s IsZero() (got: %t, expected: %t)", testCase.name, got, testCase.expected)
			}

			gotPkg := uint128.IsZero(testCase.val)
			if gotPkg != testCase.expected {
				t.Errorf("error: %s uint128.IsZero() (got: %t, expected: %t)", testCase.name, gotPkg, testCase.expected)
			}
		})
	}
}

func expectedShiftLeft(val, zero uint128.Uint128, shift uint) uint128.Uint128 {
	switch {
	case shift == 0:
		return val
	case shift >= 128:
		return zero
	case shift >= 64:
		return uint128.Uint128{H: val.L << (shift - 64), L: 0}
	default:
		return uint128.Uint128{
			H: (val.H << shift) | (val.L >> (64 - shift)),
			L: val.L << shift,
		}
	}
}

func expectedShiftRight(val, zero uint128.Uint128, shift uint) uint128.Uint128 {
	switch {
	case shift == 0:
		return val
	case shift >= 128:
		return zero
	case shift >= 64:
		return uint128.Uint128{
			L: val.H >> (shift - 64),
			H: 0,
		}
	default:
		return uint128.Uint128{
			L: (val.L >> shift) | (val.H << (64 - shift)),
			H: val.H >> shift,
		}
	}
}

func TestShift(t *testing.T) {
	t.Parallel()

	var (
		u1L    = uint128.Uint128{L: 1}
		u1H    = uint128.Uint128{H: 1}
		uH1L1  = uint128.Uint128{H: 1, L: 1}
		max128 = uint128.MaxUint128()
		zero   = uint128.Zero()
	)

	shifts := []uint{0, 1, 32, 63, 64, 65, 127, 128, 129}

	testCases := []struct {
		name string
		val  uint128.Uint128
	}{
		{"Shift 1L", u1L},
		{"Shift 1H", u1H},
		{"Shift H1L1", uH1L1},
		{"Shift Max", max128},
		{"Shift Zero", zero},
	}

	for _, testCase := range testCases {
		for _, shift := range shifts {
			t.Run(testCase.name+fmt.Sprintf(" by %d", shift), func(t *testing.T) {
				t.Parallel()

				// ShiftLeft
				expectedL := expectedShiftLeft(testCase.val, zero, shift)
				gotL := testCase.val.ShiftLeft(shift)
				assert(t, gotL, expectedL, fmt.Sprintf("ShiftLeft by %d", shift))

				// ShiftRight
				expectedR := expectedShiftRight(testCase.val, zero, shift)
				gotR := testCase.val.ShiftRight(shift)
				assert(t, gotR, expectedR, fmt.Sprintf("ShiftRight by %d", shift))
			})
		}
	}
}

func TestLogicalOps(t *testing.T) {
	t.Parallel()

	var (
		zero   = uint128.Zero()
		val1   = uint128.Uint128{H: 0xF0F0F0F0F0F0F0F0, L: 0xABABABABABABABAB}
		val2   = uint128.Uint128{H: 0x0F0F0F0F0F0F0F0F, L: 0xBABABABABABABABA}
		max128 = uint128.MaxUint128()
	)

	// Not
	t.Run("Not Zero", func(t *testing.T) { t.Parallel(); assert(t, zero.Not(), max128, "Not Zero") })
	t.Run("Not Max", func(t *testing.T) { t.Parallel(); assert(t, max128.Not(), zero, "Not Max") })
	t.Run("Not val1", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: ^val1.H, L: ^val1.L}
		assert(t, val1.Not(), expected, "Not val1")
	})

	// And
	t.Run("val1 And val2", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: val1.H & val2.H, L: val1.L & val2.L}

		assert(t, val1.And(val2), expected, "val1 And val2")
	})
	t.Run("val1 And zero", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.And(zero), zero, "val1 And zero")
	})
	t.Run("val1 And max", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.And(max128), val1, "val1 And max")
	})

	// Or
	t.Run("val1 Or val2", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: val1.H | val2.H, L: val1.L | val2.L}

		assert(t, val1.Or(val2), expected, "val1 Or val2")
	})
	t.Run("val1 Or zero", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Or(zero), val1, "val1 Or zero")
	})
	t.Run("zero Or val2", func(t *testing.T) {
		t.Parallel()

		assert(t, zero.Or(val2), val2, "zero Or val2")
	})
	t.Run("val1 Or max", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Or(max128), max128, "val1 Or max")
	})

	// Xor
	t.Run("val1 Xor val2", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: val1.H ^ val2.H, L: val1.L ^ val2.L}
		assert(t, val1.Xor(val2), expected, "val1 Xor val2")
	})
	t.Run("val1 Xor zero", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Xor(zero), val1, "val1 Xor zero")
	})
	t.Run("val1 Xor val1", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Xor(val1), zero, "val1 Xor val1")
	})
	t.Run("val1 Xor max", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: val1.H ^ max128.H, L: val1.L ^ max128.L}
		assert(t, val1.Xor(max128), expected, "val1 Xor max")
	})

	// AndNot
	t.Run("val1 AndNot val2", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: val1.H &^ val2.H, L: val1.L &^ val2.L}
		assert(t, val1.AndNot(val2), expected, "val1 AndNot val2")
	})
	t.Run("val1 AndNot zero", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.AndNot(zero), val1, "val1 AndNot zero")
	})
	t.Run("zero AndNot val1", func(t *testing.T) {
		t.Parallel()

		assert(t, zero.AndNot(val1), zero, "zero AndNot val1")
	})
	t.Run("val1 AndNot val1", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.AndNot(val1), zero, "val1 AndNot val1")
	})
	t.Run("max AndNot val1", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: max128.H &^ val1.H, L: max128.L &^ val1.L}
		assert(t, max128.AndNot(val1), expected, "max AndNot val1")
	})
}

func TestSubDecr(t *testing.T) {
	t.Parallel()

	var (
		zero   = uint128.Zero()
		oneL   = uint128.Uint128{L: 1}
		twoL   = uint128.Uint128{L: 2}
		maxL   = uint128.Uint128{L: math.MaxUint64}
		oneH   = uint128.Uint128{H: 1}
		val1   = uint128.Uint128{H: 10, L: 20}
		val2   = uint128.Uint128{H: 3, L: 5}
		max128 = uint128.MaxUint128()
	)

	// Sub
	t.Run("Sub zero from zero", func(t *testing.T) {
		t.Parallel()

		assert(t, zero.Sub(zero), zero, "Sub zero from zero")
	})
	t.Run("Sub one from one", func(t *testing.T) {
		t.Parallel()

		assert(t, oneL.Sub(oneL), zero, "Sub one from one")
	})
	t.Run("Sub one from two", func(t *testing.T) {
		t.Parallel()

		assert(t, twoL.Sub(oneL), oneL, "Sub one from two")
	})
	t.Run("Sub val2 from val1", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Sub(val2), uint128.Uint128{H: 7, L: 15}, "Sub val2 from val1")
	})
	t.Run("Sub one from zero (underflow)", func(t *testing.T) {
		t.Parallel()

		assert(t, zero.Sub(oneL), max128, "Sub one from zero")
	})
	t.Run("Sub oneH from zero (underflow)", func(t *testing.T) {
		t.Parallel()

		expected := uint128.Uint128{H: math.MaxUint64, L: 0}
		assert(t, zero.Sub(oneH), expected, "Sub oneH from zero")
	})
	t.Run("Sub maxL from oneH (borrow)", func(t *testing.T) {
		t.Parallel()

		assert(t, oneH.Sub(maxL), uint128.Uint128{H: 0, L: 1}, "Sub maxL from oneH")
	})
	t.Run("Sub oneL from oneH (borrow)", func(t *testing.T) {
		t.Parallel()

		assert(t, oneH.Sub(oneL), uint128.Uint128{H: 0, L: math.MaxUint64}, "Sub oneL from oneH")
	})
	t.Run("Sub max from max", func(t *testing.T) {
		t.Parallel()

		assert(t, max128.Sub(max128), zero, "Sub max from max")
	})
	t.Run("Sub zero from max", func(t *testing.T) {
		t.Parallel()

		assert(t, max128.Sub(zero), max128, "Sub zero from max")
	})

	// Decr
	t.Run("Decr one", func(t *testing.T) {
		t.Parallel()

		assert(t, oneL.Decr(), zero, "Decr one")
	})
	t.Run("Decr zero", func(t *testing.T) {
		t.Parallel()

		assert(t, zero.Decr(), max128, "Decr zero")
	})
	t.Run("Decr oneH", func(t *testing.T) {
		t.Parallel()

		assert(t, oneH.Decr(), maxL, "Decr oneH")
	})
	t.Run("Decr val1", func(t *testing.T) {
		t.Parallel()

		assert(t, val1.Decr(), uint128.Uint128{H: 10, L: 19}, "Decr val1")
	})
	t.Run("Decr max", func(t *testing.T) {
		t.Parallel()

		assert(t, max128.Decr(), uint128.Uint128{H: math.MaxUint64, L: math.MaxUint64 - 1}, "Decr val1")
	})
}
