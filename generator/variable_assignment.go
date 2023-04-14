package generator

import (
	"github.com/juho05/embe/blocks"
)

type Assignment struct {
	Name         string
	AssignType   blocks.BlockType
	IncreaseType blocks.BlockType
	InputName    string
}

var Assignments = make(map[string]Assignment)

func newAssignment(name string, assignType, increaseType blocks.BlockType, inputName string) {
	Assignments[name] = Assignment{
		Name:         name,
		AssignType:   assignType,
		IncreaseType: increaseType,
		InputName:    inputName,
	}
}

func init() {
	newAssignment("audio.volume", blocks.AudioSetVolume, blocks.AudioAddVolume, "number_1")
	newAssignment("audio.speed", blocks.AudioSetSpeed, blocks.AudioAddSpeed, "number_1")
	newAssignment("lights.back.brightness", blocks.LEDSetBrightness, blocks.LEDAddBrightness, "number_1")
}
