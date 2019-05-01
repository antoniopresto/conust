// Package conust transforms numbers into string tokens for which the simple string comparison
// produces the same result as the numeric comparison of the original numbers would.
// The input is not limited to decimal numbers, other bases up to base 36 are accepted with the
// restriction that it must be lowercased.
// The expected input format is ^[+-]?[0-9a-z]+(\.[0-9a-z]+)?$ Failing to satisfy this may result in panics.
//
// Transforming tokens back into numbers is also possible. This operation requires that the
// tokens are as they were generated by the encoder, modifications to them might also case panics
// when decoding.
//
// The conversion adds a few characters to the length of the original numeric string, but at
// the same time it can save some space by storing only the significant portion of the number
// omitting trailing and leading zeros of it in the output.
package conust

// [48 49 50 51 52 53 54 55 56 57 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122]
var digits36 = [...]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z'}

var digits36Reversed = [...]byte{
	'z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r', 'q',
	'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g',
	'f', 'e', 'd', 'c', 'b', 'a', '9', '8', '7', '6',
	'5', '4', '3', '2', '1', '0'}

const maxDigitValue = 35
const maxMagnitudeDigitValue = 34

const digit0 byte = '0'
const digit1 byte = '1'
const digit9 byte = '9'
const digitA byte = 'a'
const digitZ byte = 'z'
const minusByte byte = '-'
const plusByte byte = '+'

const signNegativeMagPositive byte = '3'
const signNegativeMagNegative byte = '4'
const zeroOutput = "5"
const signPositiveMagNegative byte = '6'
const signPositiveMagPositive byte = '7'

// LessThanAny is a string which is less than any encoded value
// You can use this constant as the lower boundary for generated tokens.
const LessThanAny = "2"

// GreaterThanAny is a string which is greater than any encoded value.
// You can use this constant as the upper boundary for generated tokens.
const GreaterThanAny = "8"

const zeroInput = "0"

const decimalPoint byte = '.'
const negativeNumberTerminator byte = '~'
const inTextSeparator byte = ' '

// Codec can transform strings to and from the Conust format.
//
// It has Encode and Decode functions to transform simple numbers to and from the Conust format.
//
// There is also EncodeInText, a convenience function, that encodes each group of decimal numbers
// and returns the resulting string. So that for example the strings "Item 20" and "Item 100" become
// "Item 722" and "Item 731" which sort as the numeric value in them would naturally imply.
type Codec interface {
	// Encode turns the input into the alphanumerically sortable Conust string.
	// If the input hase a base higher than 10 and contains letter characters, it must be lowercased.
	Encode(in string) (out string, ok bool)
	// Decode turns a Conust string back into its normal representation
	Decode(in string) (out string, ok bool)
	// EncodeInText is a convinience function that replaces all groups of decimal numbers of the input
	// with conust strings and returns the resulting string
	EncodeInText(in string) (out string, ok bool)
}

func isDigit(digit byte) bool {
	return (digit >= digit0 && digit <= digit9) ||
		(digit >= digitA && digit <= digitZ)
}

func digitToInt(digit byte) int {
	if digit < digitA {
		return int(digit - digit0)
	}
	return 10 + int(digit-digitA)
}

func reversedDigitToInt(digit byte) int {
	if digit < digitA {
		return 26 + int(digit9-digit)
	}
	return int(digitZ - digit)
}

func intToDigit(i int) byte {
	return digits36[i]
}

func intToReversedDigit(i int) byte {
	return digits36Reversed[i]
}

func flipDigit(digit byte) byte {
	return intToReversedDigit(digitToInt(digit))
}
