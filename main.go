package main

import (
	"fmt"
	"log"

	"github.com/brutella/can"
	"github.com/jbott/formula_dyno/input"
	"github.com/tarm/serial"
)

func main() {
	serial := input.NewSerial(serial.Config{Name: "/dev/ttyACM0", Baud: 115200})

	bus, err := can.NewBusForInterfaceWithName("can1")
	if err != nil {
		log.Fatal(err)
	}
	can := input.NewCan(bus)

	// Start all inputs
	serial.Start()
	defer serial.Stop()

	can.Start()
	defer can.Stop()

	fmt.Printf("time, name, value\n")

	for {
		select {
		case data := <-serial.Output:
			fmt.Printf("%d, torque, %f\n", data.Timestamp.UnixNano(), data.Value)
		case data := <-can.Output:
			fmt.Printf("%d, can_%d, %v\n", data.Timestamp.UnixNano(), data.IntValue, data.RawData)
		}
	}
}
