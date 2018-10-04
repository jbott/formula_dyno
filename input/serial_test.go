package input

import "testing"

func TestDecodeSerialMessage(t *testing.T) {
	data := []byte("12345 43.3213145")

	value, err := DecodeSerialMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	if value != 43.3213145 {
		t.Fatal("Did not decode value correctly")
	}
}
