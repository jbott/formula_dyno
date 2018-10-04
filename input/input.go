package input

import (
	"time"
)

type DataEvent struct {
	Time time.Time
}

type Input interface {
	Start() error
	Stop() error
	Channel() (<-chan DataEvent, error)
}
