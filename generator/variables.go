package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Var struct {
	blockType blocks.BlockType
	dataType  parser.DataType
	fields    map[string]any
	fn        func(g *generator, parent *blocks.Block)
}

var variables = map[string]Var{
	"audio.volume": {blockType: blocks.GetVolume, dataType: parser.DTNumber},
	"audio.speed":  {blockType: blocks.GetSpeed, dataType: parser.DTNumber},

	"mbot.tiltedForward":  {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"tiltforward", nil}}},
	"mbot.tiltedBackward": {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"tiltback", nil}}},
	"mbot.tiltedLeft":     {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"tiltleft", nil}}},
	"mbot.tiltedRight":    {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"tiltright", nil}}},
	"mbot.faceUp":         {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"faceup", nil}}},
	"mbot.faceDown":       {blockType: blocks.DetectAttitude, dataType: parser.DTBool, fields: map[string]any{"tilt": []any{"facedown", nil}}},

	"mbot.battery": {blockType: blocks.BatteryLevelMacAddressAndSoOn, dataType: parser.DTNumber, fields: map[string]any{"fieldMenu_1": []any{"battery", nil}}},

	"sensors.brightness": {blockType: blocks.Brightness, dataType: parser.DTNumber},
	"sensors.loudness":   {blockType: blocks.Loudness, dataType: parser.DTNumber},
	"sensors.distance": {blockType: blocks.UltrasonicDistance, dataType: parser.DTNumber, fn: func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.UltrasonicDistanceMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	}},
	"sensors.outOfRange": {blockType: blocks.UltrasonicOutOfRange, dataType: parser.DTBool, fn: func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.UltrasonicOutOfRangeMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	}},
}
