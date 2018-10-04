package main

import (
	"fmt"

	"github.com/jbott/formula_dyno/input"
	"github.com/tarm/serial"
)

func main() {
	fmt.Println("vim-go")

	input.NewSerial(serial.Config{Name: "COM45", Baud: 115200})
}
