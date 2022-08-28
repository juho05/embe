package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var Events = map[string]func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error){
	"start":    eventStart,
	"button":   eventButton,
	"joystick": eventDirectionKey,
	"tilt":     eventAction(blocks.EventDetectAttitude, "tilt", "left", "right", "forward", "backward"),
	"face":     eventAction(blocks.EventDetectAttitude, "face", "up", "down"),
	"wave":     eventAction(blocks.EventDetectAction, "wave", "left", "right", "up", "down"),
	"rotate":   eventAction(blocks.EventDetectAction, "", "clockwise", "anticlockwise"),
	"fall":     eventActionSingle(blocks.EventDetectAction, "freefall"),
	"shake":    eventActionSingle(blocks.EventDetectAction, "shake"),
	"light":    eventSensor("light_sensor"),
	"sound":    eventSensor("microphone"),
	"shakeval": eventSensor("shake_val"),
	"timer":    eventSensor("timer"),
}

func eventStart(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	if err := assertNoEventParameter(g, stmt); err != nil {
		return nil, err
	}
	return blocks.NewBlockTopLevel(blocks.EventLaunch), nil
}

func eventButton(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	if err := assertEventParameter(g, stmt, parser.DTString); err != nil {
		return nil, err
	}
	param := stmt.Parameter.Literal.(string)
	if param != "a" && param != "b" {
		return nil, g.newError(`Unknown button. Available options: "a", "b".`, stmt.Parameter)
	}
	block := blocks.NewBlockTopLevel(blocks.EventButtonPress)
	block.Fields["fieldMenu_2"] = []any{param, nil}
	return block, nil
}

func eventDirectionKey(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	if err := assertEventParameter(g, stmt, parser.DTString); err != nil {
		return nil, err
	}
	param := stmt.Parameter.Literal.(string)
	if param != "left" && param != "right" && param != "up" && param != "down" && param != "middle" {
		return nil, g.newError(`Unknown button. Available options: "left", "right", "up", "down", "middle".`, stmt.Parameter)
	}
	block := blocks.NewBlockTopLevel(blocks.EventDirectionKeyPress)
	block.Fields["fieldMenu_2"] = []any{param, nil}
	return block, nil
}

func eventAction(blockType blocks.BlockType, prefix string, options ...string) func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtEvent) (*blocks.Block, error) {
		if err := assertEventParameter(g, stmt, parser.DTString); err != nil {
			return nil, err
		}
		param := stmt.Parameter.Literal.(string)
		valid := false
		for _, o := range options {
			if param == o {
				valid = true
				break
			}
		}
		if !valid {
			return nil, g.newError(fmt.Sprintf(`Unknown argument value. Available options: "%s".`, strings.Join(options, "\", \"")), stmt.Parameter)
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
		if err := assertEventParameter(g, stmt, parser.DTString); err != nil {
			return nil, err
		}
		param := stmt.Parameter.Literal.(string)
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

func assertNoEventParameter(g *generator, stmt *parser.StmtEvent) error {
	if (stmt.Parameter != parser.Token{}) {
		return g.newError(fmt.Sprintf("The '%s' event does not take any arguments.", stmt.Name.Lexeme), stmt.Parameter)
	}
	return nil
}

func assertEventParameter(g *generator, stmt *parser.StmtEvent, dataType parser.DataType) error {
	if (stmt.Parameter == parser.Token{}) {
		return g.newError(fmt.Sprintf("The '%s' event takes a value of type %s as an argument.", stmt.Name.Lexeme, dataType), stmt.Name)
	}
	if stmt.Parameter.DataType != dataType {
		return g.newError(fmt.Sprintf("Wrong data type. Expected '%s'.", dataType), stmt.Parameter)
	}
	return nil
}
