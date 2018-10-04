package input

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

const (
	BUFFER_SIZE = 128
)

type Serial struct {
	config   serial.Config
	stopchan chan struct{}
}

func NewSerial(cfg serial.Config) Serial {
	s := Serial{
		config:   cfg,
		stopchan: make(chan struct{}),
	}

	return s
}

func WaitForDataAndSend(port *serial.Port, out chan<- []byte) {
	buf := make([]byte, BUFFER_SIZE)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	out <- buf[:n]
}

func DecodeSerialMessage(data []byte) (float64, error) {
	str := string(data)

	parts := strings.SplitN(str, " ", 2)
	if len(parts) != 2 {
		return 0.0, errors.New("Serial data malformed")
	}

	val, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0.0, err
	}

	return val, nil
}

func (s *Serial) Run() error {
	port, err := serial.OpenPort(&s.config)
	if err != nil {
		return err
	}

	defer port.Close()

	out := make(chan []byte)

	go WaitForDataAndSend(port, out)

	for {
		select {
		case data := <-out:
			// Convert data []byte into a DataEvent and send it
			log.Print(data)

		case <-s.stopchan:
			// Done!
			return nil
		}
	}
}

func (s *Serial) Start() {
	go s.Run()
}

func (s *Serial) Stop() {
}
