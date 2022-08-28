package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Assignment struct {
	AssignType   blocks.BlockType
	IncreaseType blocks.BlockType
	DataType     parser.DataType
	InputName    string
}

var Assignments = map[string]Assignment{
	"audio.volume": {
		AssignType:   blocks.AudioSetVolume,
		IncreaseType: blocks.AudioAddVolume,
		DataType:     parser.DTNumber,
		InputName:    "number_1",
	},
	"audio.speed": {
		AssignType:   blocks.AudioSetSpeed,
		IncreaseType: blocks.AudioAddSpeed,
		DataType:     parser.DTNumber,
		InputName:    "number_1",
	},

	"led.brightness": {
		AssignType:   blocks.LEDSetBrightness,
		IncreaseType: blocks.LEDAddBrightness,
		DataType:     parser.DTNumber,
		InputName:    "number_1",
	},
}
