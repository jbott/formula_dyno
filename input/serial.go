package input

import (
	"bufio"
	"errors"
	"io"
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

func WaitForDataLineAndSend(port io.Reader, out chan<- string) {
	rd := bufio.NewReader(port)

	line, err := rd.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	out <- strings.TrimRight(line, "\n")
}

func DecodeSerialMessage(str string) (float64, error) {
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

	out := make(chan string)

	go WaitForDataLineAndSend(port, out)

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
	close(s.stopchan)
}
