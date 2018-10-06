package input

import (
	"time"

	"github.com/brutella/can"
)

type Can struct {
	DataEventEmitter

	bus      *can.Bus
	stopchan chan struct{}
}

func NewCan(bus *can.Bus) Can {
	c := Can{
		DataEventEmitter: DataEventEmitter{
			Output: make(chan DataEvent),
		},

		bus:      bus,
		stopchan: make(chan struct{}),
	}

	return c
}

func (c *Can) Run() {
	out := make(chan can.Frame)

	c.bus.SubscribeFunc(func(frm can.Frame) {
		out <- frm
	})

	// Kick off publish function
	go c.bus.ConnectAndPublish()

	for {
		select {
		case frm := <-out:
			c.Output <- DataEvent{
				Timestamp: time.Now(),
				IntValue:  frm.ID,
				RawData:   frm.Data[:frm.Length],
			}

		case <-c.stopchan:
			// Done!
			c.bus.Disconnect()
			return
		}
	}
}

func (c *Can) Start() {
	go c.Run()
}

func (c *Can) Stop() {
	close(c.stopchan)
}
