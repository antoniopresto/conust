package conust

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestCodec(t *testing.T) {
	codecTests := []struct {
		name    string
		input   string
		encoded string
		decoded string
	}{
		{name: "empty", input: "", encoded: "", decoded: ""},

		{name: "zero 1", input: "0", encoded: "5", decoded: "0"},
		{name: "zero 2", input: "+000", encoded: "5", decoded: "0"},
		{name: "zero 3", input: "-000", encoded: "5", decoded: "0"},
		{name: "zero 4", input: "000.0000", encoded: "5", decoded: "0"},

		{
			name:    "all digits",
			input:   "1234567890abcdefghij.klmnopqrstuvwxyz",
			encoded: "7k1234567890abcdefghijklmnopqrstuvwxyz",
			decoded: "1234567890abcdefghij.klmnopqrstuvwxyz",
		},
		{
			name:    "negative all digits",
			input:   "-1234567890abcdefghij.klmnopqrstuvwxyz",
			encoded: "3fyxwvutsrqzponmlkjihgfedcba9876543210~",
			decoded: "-1234567890abcdefghij.klmnopqrstuvwxyz",
		},

		{name: "holes in the middle", input: "005f002k00.0i0k0", encoded: "785f002k000i0k", decoded: "5f002k00.0i0k"},

		{name: "one", input: "1", encoded: "711", decoded: "1"},
		{name: "ugly one", input: "+00001", encoded: "711", decoded: "1"},
		{name: "negative one", input: "-1", encoded: "3yy~", decoded: "-1"},
		{name: "ugly negative one", input: "-000001", encoded: "3yy~", decoded: "-1"},
		{name: "ugly positive int", input: "+00000123000", encoded: "76123", decoded: "123000"},
		{name: "ugly negative int", input: "-00000123000", encoded: "3tyxw~", decoded: "-123000"},
		{name: "fractional", input: "54321.12345", encoded: "755432112345", decoded: "54321.12345"},
		{name: "negative fractional", input: "-54321.12345", encoded: "3uuvwxyyxwvu~", decoded: "-54321.12345"},
		{
			name:    "ugly fractional",
			input:   "+00054321000.00012345000",
			encoded: "785432100000012345",
			decoded: "54321000.00012345",
		},
		{
			name:    "ugly negative fractional",
			input:   "-00054321000.00012345000",
			encoded: "3ruvwxyzzzzzzyxwvu~",
			decoded: "-54321000.00012345",
		},
		{name: "cowboy hat", input: "cowboy.hat", encoded: "76cowboyhat", decoded: "cowboy.hat"},
		{name: "negative cowboy hat", input: "-cowboy.hat", encoded: "3tnb3ob1ip6~", decoded: "-cowboy.hat"},
		{
			name:    "maximum int length",
			input:   "12345678901234567890123456789012345.1",
			encoded: "7z1123456789012345678901234567890123451",
			decoded: "12345678901234567890123456789012345.1",
		},
		{
			name:    "maximum negative int length",
			input:   "-12345678901234567890123456789012345.1",
			encoded: "30yyxwvutsrqzyxwvutsrqzyxwvutsrqzyxwvuy~",
			decoded: "-12345678901234567890123456789012345.1",
		},
		{
			name:    "maximum fracleading zero count",
			input:   "0.000000000000000000000000000000000004325430",
			encoded: "60y432543",
			decoded: "0.00000000000000000000000000000000000432543",
		},

		{
			name:    "example 1",
			input:   "12000000000000000000000000000000000000",
			encoded: "7z412",
			decoded: "12000000000000000000000000000000000000",
		},
		{name: "example 2", input: "1200", encoded: "7412", decoded: "1200"},
		{name: "example 3", input: "12", encoded: "7212", decoded: "12"},
		{name: "example 4", input: "1.2", encoded: "7112", decoded: "1.2"},
		{name: "example 5", input: "0.12", encoded: "6z12", decoded: "0.12"},
		{name: "example 6", input: "0.0012", encoded: "6x12", decoded: "0.0012"},
		{
			name:    "example 6.2",
			input:   "0.0000000000000000000000000000000000012",
			encoded: "60y12",
			decoded: "0.0000000000000000000000000000000000012",
		},
		{
			name:    "example 6.3",
			input:   "-0.0000000000000000000000000000000000012",
			encoded: "4z1yx~",
			decoded: "-0.0000000000000000000000000000000000012",
		},
		{name: "example 7", input: "-0.0012", encoded: "42yx~", decoded: "-0.0012"},
		{name: "example 8", input: "-0.12", encoded: "40yx~", decoded: "-0.12"},
		{name: "example 9", input: "-1.2", encoded: "3yyx~", decoded: "-1.2"},
		{name: "example 10", input: "-12", encoded: "3xyx~", decoded: "-12"},
		{name: "example 11", input: "-1200", encoded: "3vyx~", decoded: "-1200"},
		{
			name:    "example 12",
			input:   "-12000000000000000000000000000000000000",
			encoded: "30vyx~",
			decoded: "-12000000000000000000000000000000000000",
		},
	}
	codec := new(Codec)
	for _, i := range codecTests {
		t.Run(i.name, func(t *testing.T) {
			encoded, _ := codec.EncodeToken(i.input)

			if i.encoded != encoded {
				t.Fatalf("Encoding expected: %v, got %v\n", i.encoded, encoded)
			}

			decoded, _ := codec.DecodeToken(encoded)
			if i.decoded != decoded {
				t.Fatalf("Decoding expected: %v, got %v\n", i.decoded, decoded)
			}
		})
	}
}

func TestSortedness(t *testing.T) {
	step := 0.01
	prev := LessThanAny
	c := new(Codec)
	for i := -111111.0; i <= 111111.0; i++ {
		str := fmt.Sprintf("%3f", i*step)
		encoded, ok := c.EncodeToken(str)
		if !ok {
			t.Fatal("Encoding failed for", i)
		}
		if prev >= encoded {
			t.Fatal("at", i*step, " ", prev, "is not smaller than", encoded)
		}
		prev = encoded
	}
}

func BenchmarkEncoding(b *testing.B) {
	step := 0.001
	c := new(Codec)
	to := float64(b.N / 2)
	from := -1 * to
	for i := from; i <= to; i++ {
		str := strconv.FormatFloat(i*step, 'f', -1, 64)
		encoded, ok := c.EncodeToken(str)
		if !ok {
			b.Fatal("Encoding failed for", i)
		}

		_, ok = c.DecodeToken(encoded)
		if !ok {
			b.Fatal("Decoding failed for", encoded, "in", i)
		}
	}
}

func TestEncodeMixedText(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		ok     bool
		output string
	}{
		{name: "empty", input: "", ok: true, output: ""},
		{name: "no numbers", input: "quick brown fox", ok: true, output: "quick brown fox"},
		{name: "only numbers", input: "423", ok: true, output: "73423"},
		{name: "only numbers end with space", input: "423 ", ok: true, output: "73423 "},
		{name: "only numbers starts with space", input: " 423", ok: true, output: " 73423"},
		{name: "sign is like text", input: "-423", ok: true, output: "- 73423"},
		{name: "mixed 1.1", input: "300Z", ok: true, output: "733 Z"},
		{name: "mixed 1.2", input: "300 Z", ok: true, output: "733 Z"},
		{name: "mixed 2.1", input: "A300Z", ok: true, output: "A 733 Z"},
		{name: "mixed 2.2", input: "A 300 Z", ok: true, output: "A 733 Z"},
		{name: "mixed 3.1", input: "A300", ok: true, output: "A 733"},
		{name: "mixed 3.2", input: "A 300", ok: true, output: "A 733"},
		{
			name:   "mixed 4",
			input:  "If 2x + 3y = 8 and 4x + 12y = 28, what is x and y?",
			ok:     true,
			output: "If 712 x + 713 y = 718 and 714 x + 7212 y = 7228 , what is x and y?",
		},
		{name: "mixed c1", input: "SomeCam300D", ok: true, output: "SomeCam 733 D"},
		{name: "mixed c2", input: "SomeCam600D", ok: true, output: "SomeCam 736 D"},
		{name: "mixed c3", input: "SomeCam1000D", ok: true, output: "SomeCam 741 D"},
		{name: "mixed c4", input: "SomeCam1100D", ok: true, output: "SomeCam 7411 D"},
	}
	c := new(Codec)
	for _, i := range testCases {
		t.Run(i.name, func(t *testing.T) {
			encoded, ok := c.EncodeMixedText(i.input)

			if ok != i.ok {
				t.Fatalf("ok expected %v got %v", i.ok, ok)
			}

			if encoded != i.output {
				t.Fatalf("output expected %s got %s", i.output, encoded)
			}
		})
	}
}

func BenchmarkEncodeMixedText(b *testing.B) {
	var charPool = [...]byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
		' ', ' ', ' ', ' ', ' ', ' ',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	var genText = make([]byte, 1024)
	var textLen = len(genText)

	rand.Seed(42)
	c := new(Codec)

	for i := 0; i < textLen; i++ {
		genText[i] = charPool[rand.Intn(len(charPool))]
	}

	for i := 0; i < b.N; i++ {
		start := rand.Intn(textLen)
		end := rand.Intn(textLen-start) + start
		fragment := string(genText[start:end])
		_, ok := c.EncodeMixedText(fragment)

		if !ok {
			b.Fatal("Encoding failed for " + fragment)
		}
	}
}

func ExampleCodec_EncodeToken() {
	c := new(Codec)

	out, ok := c.EncodeToken("86400")
	fmt.Printf("%q, %v\n", out, ok)

	out, ok = c.EncodeToken("-3.14")
	fmt.Printf("%q, %v\n", out, ok)

	out, ok = c.EncodeToken("-base36.number")
	fmt.Printf("%q, %v\n", out, ok)

	// Output:
	// "75864", true
	// "3ywyv~", true
	// "3top7lwtc5dol8~", true
}

func ExampleCodec_DecodeToken() {
	c := new(Codec)

	out, ok := c.DecodeToken("42yx~")
	fmt.Printf("%q, %v\n", out, ok)

	out, ok = c.DecodeToken("6w125")
	fmt.Printf("%q, %v\n", out, ok)

	// Output:
	// "-0.0012", true
	// "0.000125", true
}
func ExampleCodec_EncodeMixedText() {
	c := new(Codec)

	out, ok := c.EncodeMixedText("SomeCam 40d")
	fmt.Printf("%q, %v\n", out, ok)

	out, ok = c.EncodeMixedText("SomeCam350d")
	fmt.Printf("%q, %v\n", out, ok)

	out, ok = c.EncodeMixedText("SomeCam1100 d")
	fmt.Printf("%q, %v\n", out, ok)

	// Output:
	// "SomeCam 724 d", true
	// "SomeCam 7335 d", true
	// "SomeCam 7411 d", true
}
