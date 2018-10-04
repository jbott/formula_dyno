package input

import (
	"time"
)

type DataEvent struct {
	Timestamp time.Time
	Value     float64
}

type Input interface {
	Start() error
	Stop() error
}

type DataEventEmitter struct {
	Output chan DataEvent
	Input
}
