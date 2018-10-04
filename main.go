package main

import (
	"fmt"

	"github.com/jbott/formula_dyno/input"
	"github.com/tarm/serial"
)

func main() {
	serial := input.NewSerial(serial.Config{Name: "/dev/ttyACM0", Baud: 115200})

	// Start all inputs
	serial.Start()
	defer serial.Stop()

	fmt.Printf("time, name, value\n")

	for {
		select {
		case data := <-serial.Output:
			fmt.Printf("%d, torque, %f\n", data.Timestamp.UnixNano(), data.Value)
		}
	}
}
