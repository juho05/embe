package generator

import (
	"fmt"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Var struct {
	Name     string
	DataType parser.DataType

	blockType blocks.BlockType
	fields    map[string]any
	fn        func(g *generator, parent *blocks.Block)
}

func (v Var) String() string {
	return fmt.Sprintf("var %s: %s", v.Name, v.DataType)
}

var Variables = make(map[string]Var)

func newVar(name string, blockType blocks.BlockType, dataType parser.DataType, fields map[string]any, fn func(g *generator, parent *blocks.Block)) {
	Variables[name] = Var{
		Name:      name,
		blockType: blockType,
		DataType:  dataType,
		fields:    fields,
		fn:        fn,
	}
}

func init() {
	newVar("audio.volume", blocks.AudioGetVolume, parser.DTNumber, nil, nil)
	newVar("audio.speed", blocks.AudioGetSpeed, parser.DTNumber, nil, nil)

	newVar("lights.back.brightness", blocks.LEDGetBrightness, parser.DTNumber, nil, nil)

	newVar("time.timer", blocks.Mbot2TimerGet, parser.DTNumber, nil, nil)

	newVar("mbot.battery", blocks.SensorBatteryLevelMacAddressAndSoOn, parser.DTNumber, map[string]any{"fieldMenu_1": []any{"battery", nil}}, nil)
	newVar("mbot.mac", blocks.SensorBatteryLevelMacAddressAndSoOn, parser.DTString, map[string]any{"fieldMenu_1": []any{"mac", nil}}, nil)
	newVar("mbot.hostname", blocks.Mbot2Hostname, parser.DTString, nil, nil)

	newVar("sensors.wavingAngle", blocks.SensorWaveAngle, parser.DTNumber, nil, nil)
	newVar("sensors.wavingSpeed", blocks.SensorWaveSpeed, parser.DTNumber, nil, nil)
	newVar("sensors.shakingStrength", blocks.SensorShakingStrength, parser.DTNumber, nil, nil)
	newVar("sensors.brightness", blocks.SensorBrightness, parser.DTNumber, nil, nil)
	newVar("sensors.loudness", blocks.SensorLoudness, parser.DTNumber, nil, nil)
	newVar("sensors.distance", blocks.SensorUltrasonicDistance, parser.DTNumber, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicDistanceMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	})
	newVar("sensors.outOfRange", blocks.SensorUltrasonicOutOfRange, parser.DTBool, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicOutOfRangeMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	})
	newVar("sensors.lineDeviation", blocks.SensorColorGetOffTrack, parser.DTNumber, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		indexMenu := g.NewBlock(blocks.SensorColorGetOffTrackIndex, true)
		indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, indexMenu.ID}
	})

	newVar("net.connected", blocks.NetWifiIsConnected, parser.DTBool, nil, nil)
}
