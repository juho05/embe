package generator

import (
	"github.com/juho05/embe/blocks"
)

type Var struct {
	Name string

	blockType blocks.BlockType
	fields    map[string]any
	fn        func(g *generator, parent *blocks.Block)
}

var Variables = make(map[string]Var)

func newVar(name string, blockType blocks.BlockType, fields map[string]any, fn func(g *generator, parent *blocks.Block)) {
	Variables[name] = Var{
		Name:      name,
		blockType: blockType,
		fields:    fields,
		fn:        fn,
	}
}

func init() {
	newVar("audio.volume", blocks.AudioGetVolume, nil, nil)
	newVar("audio.speed", blocks.AudioGetSpeed, nil, nil)

	newVar("lights.back.brightness", blocks.LEDGetBrightness, nil, nil)

	newVar("time.timer", blocks.Mbot2TimerGet, nil, nil)

	newVar("mbot.battery", blocks.SensorBatteryLevelMacAddressAndSoOn, map[string]any{"fieldMenu_1": []any{"battery", nil}}, nil)
	newVar("mbot.mac", blocks.SensorBatteryLevelMacAddressAndSoOn, map[string]any{"fieldMenu_1": []any{"mac", nil}}, nil)
	newVar("mbot.hostname", blocks.Mbot2Hostname, nil, nil)

	newVar("sensors.wavingAngle", blocks.SensorWaveAngle, nil, nil)
	newVar("sensors.wavingSpeed", blocks.SensorWaveSpeed, nil, nil)
	newVar("sensors.shakingStrength", blocks.SensorShakingStrength, nil, nil)
	newVar("sensors.brightness", blocks.SensorBrightness, nil, nil)
	newVar("sensors.loudness", blocks.SensorLoudness, nil, nil)
	newVar("sensors.distance", blocks.SensorUltrasonicDistance, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicDistanceMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	})
	newVar("sensors.outOfRange", blocks.SensorUltrasonicOutOfRange, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		block := g.NewBlock(blocks.SensorUltrasonicOutOfRangeMenu, true)
		block.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, block.ID}
	})
	newVar("sensors.lineDeviation", blocks.SensorColorGetOffTrack, nil, func(g *generator, parent *blocks.Block) {
		g.noNext = true
		indexMenu := g.NewBlock(blocks.SensorColorGetOffTrackIndex, true)
		indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
		parent.Inputs["index"] = []any{1, indexMenu.ID}
	})

	newVar("net.connected", blocks.NetWifiIsConnected, nil, nil)

	newVar("draw.positionX", blocks.DrawSketchGetXYAngleAndSize, map[string]any{"fieldMenu_1": []any{"x", nil}}, nil)
	newVar("draw.positionY", blocks.DrawSketchGetXYAngleAndSize, map[string]any{"fieldMenu_1": []any{"y", nil}}, nil)
	newVar("draw.rotation", blocks.DrawSketchGetXYAngleAndSize, map[string]any{"fieldMenu_1": []any{"angle", nil}}, nil)
	newVar("draw.thickness", blocks.DrawSketchGetXYAngleAndSize, map[string]any{"fieldMenu_1": []any{"size", nil}}, nil)
}
