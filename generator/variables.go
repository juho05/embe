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

var Variables = map[string]Var{
	"audio.volume": {blockType: blocks.AudioGetVolume, dataType: parser.DTNumber},
	"audio.speed":  {blockType: blocks.AudioGetSpeed, dataType: parser.DTNumber},

	"led.brightness": {blockType: blocks.LEDGetBrightness, dataType: parser.DTNumber},

	"mbot.battery": {blockType: blocks.SensorBatteryLevelMacAddressAndSoOn, dataType: parser.DTNumber, fields: map[string]any{"fieldMenu_1": []any{"battery", nil}}},

	"sensors.wavingAngle":     {blockType: blocks.SensorWaveAngle, dataType: parser.DTNumber},
	"sensors.wavingSpeed":     {blockType: blocks.SensorWaveSpeed, dataType: parser.DTNumber},
	"sensors.shakingStrength": {blockType: blocks.SensorShakingStrength, dataType: parser.DTNumber},
	"sensors.brightness":      {blockType: blocks.SensorBrightness, dataType: parser.DTNumber},
	"sensors.loudness":        {blockType: blocks.SensorLoudness, dataType: parser.DTNumber},
	"sensors.distance": {blockType: blocks.SensorUltrasonicDistance, dataType: parser.DTNumber, fn: func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicDistanceMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	}},
	"sensors.outOfRange": {blockType: blocks.SensorUltrasonicOutOfRange, dataType: parser.DTBool, fn: func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicOutOfRangeMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	}},
}
