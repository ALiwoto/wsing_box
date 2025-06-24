package singingEncoding

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
)

// Define the ranges of characters we'll use for encoding
var singingChars = [][]rune{
	{'a', 'z'}, // Latin lowercase
	{'A', 'Z'}, // Latin uppercase
	{'0', '9'}, // Digits
	{'!', '/'}, // Some ASCII symbols
	{':', '@'}, // More ASCII symbols
	{'[', '`'}, // More ASCII symbols
	{'{', '~'}, // More ASCII symbols
	{'α', 'ω'}, // Greek lowercase
	{'Α', 'Ω'}, // Greek uppercase
	{'а', 'я'}, // Cyrillic lowercase
	{'ا', 'ی'}, // Persian characters
}

// This will hold all the individual characters we can use for encoding
var allSingingChars []rune

// init populates allSingingChars with all valid characters from our ranges
func init() {
	// Create a map to ensure uniqueness
	uniqueChars := make(map[rune]bool)

	// Process each range
	for _, charRange := range singingChars {
		if len(charRange) != 2 {
			for _, current := range charRange {
				uniqueChars[current] = true
				// println(int(current), string(current))
			}
		}
		start, end := charRange[0], charRange[1]
		for r := start; r <= end; r++ {
			uniqueChars[r] = true
			// println(int(r), string(r))
		}
	}

	// Convert the map keys to our slice
	allSingingChars = make([]rune, 0, len(uniqueChars))
	for r := range uniqueChars {
		allSingingChars = append(allSingingChars, r)
	}
}

// SingingEncoding represents an encoding defined by a given alphabet
type SingingEncoding struct {
	encoder []rune
	size    int
}

// NewEncoding returns a new SingingEncoding that uses the given alphabet
// which consists of the provided character ranges
func NewEncoding(charRanges [][]rune) *SingingEncoding {
	// Create a map to ensure uniqueness
	uniqueChars := make(map[rune]bool)

	// Process each range
	for _, charRange := range charRanges {
		start, end := charRange[0], charRange[1]
		for r := start; r <= end; r++ {
			uniqueChars[r] = true
		}
	}

	// Convert the map keys to our slice
	encoder := make([]rune, 0, len(uniqueChars))
	for r := range uniqueChars {
		encoder = append(encoder, r)
	}

	return &SingingEncoding{
		encoder: encoder,
		size:    len(encoder),
	}
}

// StdEncoding is the standard singing encoding using Latin, digits, and Persian characters
var StdEncoding = NewEncoding(singingChars)

// EncodeToString returns the singing encoding of src.
func (enc *SingingEncoding) EncodeToString(src []byte) string {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write(src); err != nil {
		// Handle error - in practice you might want to return an error instead
		return ""
	}
	if err := gw.Close(); err != nil {
		// Handle error
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (enc *SingingEncoding) EncodeToStringOld(src []byte) string {
	if len(src) == 0 {
		return ""
	}

	// gzip algorithm
	buff := bytes.NewBuffer([]byte{})
	w, _ := gzip.NewWriterLevel(buff, gzip.BestCompression)
	w.Write(src)
	w.Close()

	newSrc, err := io.ReadAll(buff)
	if err == nil {
		src = newSrc
	}

	// Convert byte array to a big integer
	num := new(big.Int).SetBytes(src)

	// Calculate the encoded string
	var result []rune
	zero := big.NewInt(0)
	base := big.NewInt(int64(enc.size))
	mod := new(big.Int)

	// Convert the big integer to our encoding
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		index := mod.Int64()
		result = append(result, enc.encoder[index])
	}

	// The result is built in reverse order, so we need to reverse it
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// DecodeString returns the bytes represented by the singing encoded string s.
func (enc *SingingEncoding) DecodeString(s string) ([]byte, error) {
	compressed, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	gr, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	decompressed, err := io.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

func (enc *SingingEncoding) DecodeStringOld(s string) ([]byte, error) {
	if len(s) == 0 {
		return []byte{}, nil
	}

	// Create a reverse lookup map for decoding
	decoderMap := make(map[rune]int)
	for i, r := range enc.encoder {
		decoderMap[r] = i
	}

	// Convert the encoded string back to a big integer
	num := big.NewInt(0)
	base := big.NewInt(int64(enc.size))

	for _, r := range s {
		value, ok := decoderMap[r]
		if !ok {
			return nil, fmt.Errorf("singing encoding: invalid character %q", r)
		}

		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(value)))
	}

	// gzip algorithm
	buff := bytes.NewBuffer(num.Bytes())
	reader, err := gzip.NewReader(buff)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(reader)

	// Convert big integer back to bytes
	// return num.Bytes(), nil
}
