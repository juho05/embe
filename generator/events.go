package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Event struct {
	Name  string
	Fn    func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error)
	Param *Param
}

func (e Event) String() string {
	if e.Param == nil {
		return "@" + e.Name
	}
	return fmt.Sprintf("@%s %s: %s", e.Name, e.Param.Name, e.Param.Type)
}

var Events = make(map[string]Event)

func newEvent(name string, fn func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error), param *Param) {
	Events[name] = Event{
		Name:  name,
		Fn:    fn,
		Param: param,
	}
}

func init() {
	newEvent("start", eventStart, nil)
	newEvent("button", eventButton, &Param{Name: "name", Type: parser.DTString})
	newEvent("joystick", eventDirectionKey, &Param{Name: "direction", Type: parser.DTString})
	newEvent("tilt", eventAction(blocks.EventDetectAttitude, "tilt", "left", "right", "forward", "backward"), &Param{Name: "direction", Type: parser.DTString})
	newEvent("face", eventAction(blocks.EventDetectAttitude, "face", "up", "down"), &Param{Name: "direction", Type: parser.DTString})
	newEvent("wave", eventAction(blocks.EventDetectAction, "wave", "left", "right"), &Param{Name: "direction", Type: parser.DTString})
	newEvent("rotate", eventAction(blocks.EventDetectAction, "", "clockwise", "anticlockwise"), &Param{Name: "direction", Type: parser.DTString})
	newEvent("fall", eventActionSingle(blocks.EventDetectAction, "freefall"), nil)
	newEvent("shake", eventActionSingle(blocks.EventDetectAction, "shake"), nil)
	newEvent("light", eventSensor("light_sensor"), &Param{Name: "comparison", Type: parser.DTString})
	newEvent("sound", eventSensor("microphone"), &Param{Name: "comparison", Type: parser.DTString})
	newEvent("shakeval", eventSensor("shakeval"), &Param{Name: "comparison", Type: parser.DTString})
	newEvent("timer", eventSensor("timer"), &Param{Name: "comparison", Type: parser.DTString})
	newEvent("receive", eventReceive, &Param{Name: "message", Type: parser.DTString})
}

func eventStart(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	if err := assertNoEventParameter(g, stmt); err != nil {
		return nil, err
	}
	return blocks.NewBlockTopLevel(blocks.EventLaunch), nil
}

func eventButton(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	param, err := getParameter(g, stmt, parser.DTString, []string{"a", "b"})
	if err != nil {
		return nil, err
	}
	block := blocks.NewBlockTopLevel(blocks.EventButtonPress)
	block.Fields["fieldMenu_2"] = []any{param, nil}
	return block, nil
}

func eventDirectionKey(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	param, err := getParameter(g, stmt, parser.DTString, []string{"left", "right", "up", "down", "middle"})
	if err != nil {
		return nil, err
	}
	block := blocks.NewBlockTopLevel(blocks.EventDirectionKeyPress)
	block.Fields["fieldMenu_2"] = []any{param, nil}
	return block, nil
}

func eventAction(blockType blocks.BlockType, prefix string, options ...string) func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
		param, err := getParameter(g, stmt, parser.DTString, options)
		if err != nil {
			return nil, err
		}
		if param == "backward" {
			param = "back"
		}
		block := blocks.NewBlockTopLevel(blockType)
		block.Fields["tilt"] = []any{"is_" + prefix + param, nil}
		return block, nil
	}
}

func eventActionSingle(blockType blocks.BlockType, name string) func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
		if err := assertNoEventParameter(g, stmt); err != nil {
			return nil, err
		}
		block := blocks.NewBlockTopLevel(blockType)
		block.Fields["tilt"] = []any{"is_" + name, nil}
		return block, nil
	}
}

func eventSensor(sensor string) func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
		param, err := getParameter[string](g, stmt, parser.DTString, nil)
		if err != nil {
			return nil, err
		}
		parts := strings.SplitAfter(param, ">")
		if len(parts) == 1 {
			parts = strings.SplitAfter(param, "<")
		}
		if len(parts) != 2 {
			return nil, g.newError(`Invalid argument. Expected format: "< NUMBER" or "> NUMBER", e.g "< 12.3".`, stmt.Parameter)
		}
		parts[0] = strings.TrimSpace(parts[0])
		parts[1] = strings.TrimSpace(parts[1])
		if parts[0] != "<" && parts[0] != ">" {
			return nil, g.newError(`Invalid argument. Expected format: "< NUMBER" or "> NUMBER", e.g "< 12.3".`, stmt.Parameter)
		}
		num, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, g.newError(`Invalid argument. Expected format: "< NUMBER" or "> NUMBER", e.g "< 12.3".`, stmt.Parameter)
		}

		block := blocks.NewBlockTopLevel(blocks.EventSensorValueBiggerOrSmaller)
		block.Inputs["number_3"] = []any{1, []any{4, fmt.Sprintf("%v", num)}}
		block.Fields["fieldMenu_2"] = []any{sensor, nil}

		operator := "greater"
		if parts[0] == "<" {
			operator = "smaller"
		}
		block.Fields["fieldMenu_3"] = []any{operator, nil}
		return block, nil
	}
}

func eventReceive(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	param, err := getParameter[string](g, stmt, parser.DTString, nil)
	if err != nil {
		return nil, err
	}
	block := blocks.NewBlockTopLevel(blocks.EventReceivedMessage)
	block.Inputs["message"] = []any{1, []any{10, param}}
	return block, nil
}

func assertNoEventParameter(g *generator, stmt *parser.StmtEvent) error {
	if (stmt.Parameter != parser.Token{}) {
		return g.newError(fmt.Sprintf("The '%s' event does not take any arguments.", stmt.Name.Lexeme), stmt.Parameter)
	}
	return nil
}

func getParameter[T comparable](g *generator, stmt *parser.StmtEvent, dataType parser.DataType, options []T) (T, error) {
	var value T
	if (stmt.Parameter == parser.Token{}) {
		return value, g.newError(fmt.Sprintf("The '%s' event takes a value of type %s as an argument.", stmt.Name.Lexeme, dataType), stmt.Name)
	}
	if stmt.Parameter.Type == parser.TkIdentifier {
		if constant, ok := g.constants[stmt.Parameter.Lexeme]; ok {
			if constant.Type != dataType {
				return value, g.newError(fmt.Sprintf("Wrong data type. Expected '%s'.", dataType), stmt.Parameter)
			}
			value = constant.Value.Literal.(T)
		} else {
			return value, g.newError("Unknown constant.", stmt.Parameter)
		}
	} else {
		if stmt.Parameter.DataType != dataType {
			return value, g.newError(fmt.Sprintf("Wrong data type. Expected '%s'.", dataType), stmt.Parameter)
		}
		value = stmt.Parameter.Literal.(T)
	}

	if options != nil {
		valid := false
		for _, o := range options {
			if value == o {
				valid = true
				break
			}
		}
		if !valid {
			strOptions := make([]string, len(options))
			for i, o := range options {
				strOptions[i] = fmt.Sprintf("%v", o)
			}
			return value, g.newError(fmt.Sprintf("Invalid value. Available options: %s", strings.Join(strOptions, ", ")), stmt.Parameter)
		}
	}

	return value, nil
}
