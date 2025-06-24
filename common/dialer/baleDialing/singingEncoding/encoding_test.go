package singingEncoding_test

import (
	"testing"

	"github.com/sagernet/sing-box/common/dialer/baleDialing/singingEncoding"
)

func TestDecoding01(t *testing.T) {
	myValue := "hello how are you doing"

	enc := singingEncoding.StdEncoding.EncodeToString([]byte(myValue))

	b, err := singingEncoding.StdEncoding.DecodeString(enc)
	if err != nil {
		t.Error(err)
		return
	}

	if string(b) != myValue {
		t.Error("values are not equal")
		return
	}
}
func TestDecoding02(t *testing.T) {
	encoded := "ذڀNυиٕ?r\"\u03a2E>9a&ڿڧٖ:_ټ٥ڎmڗ٬κχ"

	b, err := singingEncoding.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Error(err)
		return
	}

	if len(b) == 0 {
		t.Error("empty value returned")
		return
	}
}
