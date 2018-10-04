package input

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

const (
	BUFFER_SIZE = 128
)

type Serial struct {
	DataEventEmitter

	config   serial.Config
	stopchan chan struct{}
}

func NewSerial(cfg serial.Config) Serial {
	s := Serial{
		DataEventEmitter: DataEventEmitter{
			Output: make(chan DataEvent),
		},

		config:   cfg,
		stopchan: make(chan struct{}),
	}

	return s
}

type LineEvent struct {
	Timestamp time.Time
	Data      string
}

func WaitForDataLineAndSend(port io.Reader, out chan<- LineEvent) {
	rd := bufio.NewReader(port)

	for {
		line, err := rd.ReadString('\n')
		if err == nil || err == io.EOF {
			// Send an event
			now := time.Now()
			out <- LineEvent{Timestamp: now, Data: strings.TrimRight(line, "\r\n")}

			if err == io.EOF {
				// Done reading, exit
				return
			}
		} else {
			log.Fatal(err)
		}
	}
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

	// This will ensure our WaitForDataLineAndSend goroutine closes when this
	// function exits
	defer port.Close()

	out := make(chan LineEvent)

	go WaitForDataLineAndSend(port, out)

	for {
		select {
		case line_event := <-out:
			value, err := DecodeSerialMessage(line_event.Data)
			if err == nil {
				s.Output <- DataEvent{
					Timestamp: time.Now(),
					Value:     value,
				}
			}

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
