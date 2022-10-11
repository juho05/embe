package analyzer

import (
	"fmt"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Assignment struct {
	Name         string
	DataType     parser.DataType
	AssignType   blocks.BlockType
	IncreaseType blocks.BlockType
	InputName    string
}

func (a Assignment) String() string {
	return fmt.Sprintf("var %s: %s", a.Name, a.DataType)
}

var Assignments = make(map[string]Assignment)

func newAssignment(name string, dataType parser.DataType, assignType, increaseType blocks.BlockType, inputName string) {
	Assignments[name] = Assignment{
		Name:         name,
		DataType:     dataType,
		AssignType:   assignType,
		IncreaseType: increaseType,
		InputName:    inputName,
	}
}

func init() {
	newAssignment("audio.volume", parser.DTNumber, blocks.AudioSetVolume, blocks.AudioAddVolume, "number_1")
	newAssignment("audio.speed", parser.DTNumber, blocks.AudioSetSpeed, blocks.AudioAddSpeed, "number_1")
	newAssignment("lights.back.brightness", parser.DTNumber, blocks.LEDSetBrightness, blocks.LEDAddBrightness, "number_1")
}
