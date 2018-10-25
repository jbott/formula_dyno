package output

import (
	"github.com/jbott/formula_dyno/input"
)

type DataEventConsumer struct {
	DataIn chan DataEvent
	input.Input
}
