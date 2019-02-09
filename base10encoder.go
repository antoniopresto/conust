package conust

import (
	"strings"
)

type base10Encoder struct {
	isEmpty                  bool
	isZero                   bool
	isPositive               bool
	intStart                 int
	intSignificantPartEnd    int
	intEnd                   int
	fracLeadingZeroCount     int
	fracSignificantPartStart int
	fracEnd                  int
}

// NewBase10Encoder returns an encoder that outputs base(10) Conust strings
func NewBase10Encoder() Encoder {
	return base10Encoder{}
}

//FromString turns input number into a base(10) Conust string
func (e base10Encoder) FromString(s string) (out string, ok bool) {
	// TODO
	return "", false
}

func (e base10Encoder) analyzeInput(s string) bool {
	length := len(s)
	// empty input results in empty but ok output
	if length == 0 {
		e.isEmpty = true
		return true
	}

	i := 0

	// determine sign
	switch s[0] {
	case minusByte:
		e.isPositive = false
		i++
	case plusByte:
		e.isPositive = true
		i++
	default:
		e.isPositive = true
	}
	// a sign only is bad input
	if i >= length {
		return false
	}

	// skip leading zeroes
	for ; i < length && s[i] == digit0; i++ {
	}

	// if there were only zeroes, the result is zeroOutput
	if i >= length {
		e.isZero = true
		return true
	}

	// determine integer part bounds
	e.intStart = i
	trailingZeroCount := 0
	for ; i < length; i++ {
		if s[i] == positiveIntegerTerminator {
			break
		}
		if s[i] == digit0 {
			trailingZeroCount++
			continue
		}
		if s[i] > digit9 || s[i] < digit0 {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.intSignificantPartEnd = (i - 1) - trailingZeroCount
	e.intEnd = (i - 1)

	// init fraction variables
	e.fracSignificantPartStart = -1
	e.fracEnd = -1
	e.fracLeadingZeroCount = 0

	// if no fraction present, end processing
	if i >= length-1 {
		return true
	}

	// skip over decimal separator
	i++

	// process fraction part
	e.fracSignificantPartStart = i
	for ; i < length && s[i] == digit0; i++ {
	}

	// fraction contains only zeroes
	if i >= length {
		e.fracSignificantPartStart = -1
		e.fracEnd = -1
		return true
	}

	e.fracLeadingZeroCount = i - e.fracSignificantPartStart
	e.fracSignificantPartStart = i
	trailingZeroCount = 0
	for ; i < length; i++ {
		if s[i] == digit0 {
			trailingZeroCount++
			continue
		}
		if s[i] > digit9 || s[i] < digit0 {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.fracEnd = (i - 1) - trailingZeroCount
	return true
}

// FromI32 turns input number into a base(10) Conust string
func (e base10Encoder) FromInt32(i int32) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int32Preproc(i))
}

// FromInt64 turns input number into a base(10) Conust string
func (e base10Encoder) FromInt64(i int64) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int64Preproc(i))
}

// FromFloat32 turns input number into a base(10) Conust string
func (e base10Encoder) FromFloat32(f float32) string {
	// TODO
	return ""
}

// FromFloat64 turns input number into a base(10) Conust string
func (e base10Encoder) FromFloat64(f float64) string {
	// TODO
	return ""
}

func (e base10Encoder) fromIntString(positive bool, absNumber string) string {
	var b strings.Builder
	b.Grow(builderInitialCap)
	if positive {
		b.WriteByte(signPositive)
		e.encode(&b, true, absNumber)
	} else {
		b.WriteByte(signNegative)
		e.encode(&b, false, absNumber)
	}
	return b.String()
}

func (e base10Encoder) encode(b *strings.Builder, positive bool, number string) {
	if positive {
		b.WriteByte(intToDigit(len(number)))
		b.WriteString(strings.TrimRight(number, trailing0))
	} else {
		b.WriteByte(intToReversedDigit36(len(number)))
		number = strings.TrimRight(number, trailing0)
		for j := 0; j < len(number); j++ {
			b.WriteByte(flipDigit10(number[j]))
		}
		b.WriteByte(negativeIntegerTerminator)
	}
}
