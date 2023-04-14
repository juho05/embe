package analyzer

import (
	"fmt"

	"github.com/juho05/embe/parser"
)

type Event struct {
	Name         string
	Param        *Param
	ParamOptions []any
}

func (e Event) String() string {
	if e.Param == nil {
		return "event " + e.Name
	}
	return fmt.Sprintf("event %s %s: %s", e.Name, e.Param.Name, e.Param.Type)
}

var Events = make(map[string]Event)

func newEvent(name string, param *Param, options ...any) {
	Events[name] = Event{
		Name:         name,
		Param:        param,
		ParamOptions: options,
	}
}

func init() {
	newEvent("launch", nil)
	newEvent("button", &Param{Name: "name", Type: parser.DTString}, "a", "b")
	newEvent("joystick", &Param{Name: "direction", Type: parser.DTString}, "up", "down", "left", "right", "middle")
	newEvent("tilt", &Param{Name: "direction", Type: parser.DTString}, "left", "right", "forward", "backward")
	newEvent("face", &Param{Name: "direction", Type: parser.DTString}, "up", "down")
	newEvent("wave", &Param{Name: "direction", Type: parser.DTString}, "left", "right")
	newEvent("rotate", &Param{Name: "direction", Type: parser.DTString}, "clockwise", "anticlockwise")
	newEvent("fall", nil)
	newEvent("shake", nil)
	newEvent("light", &Param{Name: "comparison", Type: parser.DTString})
	newEvent("sound", &Param{Name: "comparison", Type: parser.DTString})
	newEvent("shakeval", &Param{Name: "comparison", Type: parser.DTString})
	newEvent("timer", &Param{Name: "comparison", Type: parser.DTString})
	newEvent("receive", &Param{Name: "message", Type: parser.DTString})
}
