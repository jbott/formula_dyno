package input

import (
	"bytes"
	"strings"
	"testing"
)

const (
	sample_data = "12345 43.3213145\n"
)

func TestDecodeSerialMessage(t *testing.T) {
	// Trim newline from constant string so it should match the
	// WaitForDataLineAndSend func
	trimmed_data := strings.TrimRight(sample_data, "\n")

	value, err := DecodeSerialMessage(trimmed_data)
	if err != nil {
		t.Fatal(err)
	}

	if value != 43.3213145 {
		t.Fatal("Did not decode value correctly")
	}
}

func TestSerialRecieve(t *testing.T) {
	rd := bytes.NewReader([]byte(sample_data))

	out := make(chan string)

	go WaitForDataLineAndSend(rd, out)

	// Block until we read data from goroutine
	output_data := <-out

	// Trim newline from constant string so it should match the
	// WaitForDataLineAndSend func
	trimmed_data := strings.TrimRight(sample_data, "\n")

	if strings.Compare(trimmed_data, string(output_data)) != 0 {
		t.Fatal("Sample data and output data don't match!")
	}
}
