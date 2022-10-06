package generator

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

type Param struct {
	Name string
	Type parser.DataType
}

type Signature struct {
	FuncName   string
	Params     []Param
	ReturnType parser.DataType
}

func (s Signature) String() string {
	signature := s.FuncName + "("

	for i, p := range s.Params {
		if i > 0 {
			signature += ", "
		}
		signature = fmt.Sprintf("%s%s: %s", signature, p.Name, p.Type)
	}

	signature += ")"

	if s.ReturnType != "" {
		signature += " : " + string(s.ReturnType)
	}

	return signature
}

type FuncCall struct {
	Name       string
	Signatures []Signature
	Fn         func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error)
}

var FuncCalls = make(map[string]FuncCall)

func newFuncCall(name string, fn func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error), signatures ...[]Param) {
	if len(signatures) == 0 {
		signatures = append(signatures, []Param{})
	}

	call := FuncCall{
		Name:       name,
		Signatures: make([]Signature, len(signatures)),
	}

	for i, s := range signatures {
		call.Signatures[i].FuncName = name
		call.Signatures[i].Params = s
	}

	call.Fn = func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		for _, s := range signatures {
			if len(s) == len(stmt.Parameters) {
				return fn(g, stmt)
			}
		}

		signature := make([]string, len(call.Signatures))
		for i, s := range call.Signatures {
			signature[i] = s.String()
		}

		return nil, g.newError("Wrong argument count.", stmt.Name)
	}

	FuncCalls[name] = call
}

func init() {
	newFuncCall("audio.stop", funcAudioStop)
	newFuncCall("audio.playBuzzer", funcAudioPlayBuzzer, []Param{{Name: "frequency", Type: parser.DTNumber}}, []Param{{Name: "frequency", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.playClip", funcAudioPlayClip, []Param{{Name: "name", Type: parser.DTString}}, []Param{{Name: "name", Type: parser.DTString}, {Name: "block", Type: parser.DTBool}})
	newFuncCall("audio.playInstrument", funcAudioPlayInstrument, []Param{{Name: "name", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.playNote", funcAudioPlayNote, []Param{{Name: "name", Type: parser.DTString}, {Name: "octave", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "note", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("audio.record.start", funcAudioRecordingStart)
	newFuncCall("audio.record.stop", funcAudioRecordingStop)
	newFuncCall("audio.record.play", funcAudioRecordingPlay, []Param{}, []Param{{Name: "block", Type: parser.DTBool}})

	newFuncCall("lights.back.playAnimation", funcLEDPlayAnimation, []Param{{Name: "name", Type: parser.DTString}})
	newFuncCall("lights.front.setBrightness", funcLEDSetAmbientBrightness("set"), []Param{{Name: "value", Type: parser.DTNumber}}, []Param{{Name: "light", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lights.front.addBrightness", funcLEDSetAmbientBrightness("add"), []Param{{Name: "value", Type: parser.DTNumber}}, []Param{{Name: "light", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lights.front.displayEmotion", funcLEDDisplayEmotion, []Param{{Name: "emotion", Type: parser.DTString}})
	newFuncCall("lights.front.deactivate", funcLEDDeactivateAmbient, []Param{}, []Param{{Name: "light", Type: parser.DTNumber}})
	newFuncCall("lights.bottom.deactivate", funcLEDDeactivateFill)
	newFuncCall("lights.bottom.setColor", funcLEDSetFillColor, []Param{{Name: "color", Type: parser.DTString}})
	newFuncCall("lights.back.display", funcLEDDisplay, []Param{{Name: "color1", Type: parser.DTString}, {Name: "color2", Type: parser.DTString}, {Name: "color3", Type: parser.DTString}, {Name: "color4", Type: parser.DTString}, {Name: "color5", Type: parser.DTString}})
	newFuncCall("lights.back.displayColor", funcLEDDisplayColor, []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("lights.back.displayColorFor", funcLEDDisplayColorFor, []Param{{Name: "color", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "color", Type: parser.DTString}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "led", Type: parser.DTNumber}, {Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("lights.back.deactivate", funcLEDDeactivate, []Param{}, []Param{{Name: "led", Type: parser.DTNumber}})
	newFuncCall("lights.back.move", funcLEDMove, []Param{{Name: "n", Type: parser.DTNumber}})

	newFuncCall("display.print", funcDisplayPrint(false), []Param{{Name: "text", Type: parser.DTString}})
	newFuncCall("display.println", funcDisplayPrint(true), []Param{{Name: "text", Type: parser.DTString}})
	newFuncCall("display.setFontSize", funcDisplaySetFontSize, []Param{{Name: "size", Type: parser.DTNumber}})
	newFuncCall("display.setColor", funcDisplaySetColor, []Param{{Name: "color", Type: parser.DTString}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}})
	newFuncCall("display.showLabel", funcDisplayShowLabel, []Param{{Name: "label", Type: parser.DTString}, {Name: "text", Type: parser.DTString}, {Name: "location", Type: parser.DTString}, {Name: "size", Type: parser.DTNumber}}, []Param{{Name: "label", Type: parser.DTString}, {Name: "text", Type: parser.DTString}, {Name: "x", Type: parser.DTNumber}, {Name: "y", Type: parser.DTNumber}, {Name: "size", Type: parser.DTNumber}})
	newFuncCall("display.lineChart.addData", funcDisplayLineChartAddData, []Param{{Name: "value", Type: parser.DTNumber}})
	newFuncCall("display.lineChart.setInterval", funcDisplayLineChartSetInterval, []Param{{Name: "interval", Type: parser.DTNumber}})
	newFuncCall("display.barChart.addData", funcDisplayBarChartAddData, []Param{{Name: "value", Type: parser.DTNumber}})
	newFuncCall("display.table.addData", funcDisplayTableAddData, []Param{{Name: "text", Type: parser.DTString}, {Name: "row", Type: parser.DTNumber}, {Name: "column", Type: parser.DTNumber}})
	newFuncCall("display.setOrientation", funcDisplaySetOrientation, []Param{{Name: "orientation", Type: parser.DTNumber}})
	newFuncCall("display.clear", funcDisplayClear)

	newFuncCall("net.broadcast", funcNetBroadcast, []Param{{Name: "message", Type: parser.DTString}}, []Param{{Name: "message", Type: parser.DTString}, {Name: "value", Type: parser.DTString}})
	newFuncCall("net.setChannel", funcNetSetChannel, []Param{{Name: "channel", Type: parser.DTNumber}})
	newFuncCall("net.connect", funcNetConnect, []Param{{Name: "ssid", Type: parser.DTString}, {Name: "password", Type: parser.DTString}})
	newFuncCall("net.reconnect", funcNetReconnect)
	newFuncCall("net.disconnect", funcNetDisconnect)

	newFuncCall("sensors.resetAngle", funcSensorsResetAngle, []Param{{Name: "axis", Type: parser.DTString}})
	newFuncCall("sensors.resetYawAngle", funcSensorsResetYawAngle)
	newFuncCall("sensors.defineColor", funcSensorsDefineColor, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}}, []Param{{Name: "r", Type: parser.DTNumber}, {Name: "g", Type: parser.DTNumber}, {Name: "b", Type: parser.DTNumber}, {Name: "tolerance", Type: parser.DTNumber}})

	newFuncCall("motors.run", funcMotorsRun("forward"), []Param{{Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.runBackward", funcMotorsRun("backward"), []Param{{Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.moveDistance", funcMotorsRunDistance("forward"), []Param{{Name: "distance", Type: parser.DTNumber}})
	newFuncCall("motors.moveDistanceBackward", funcMotorsRunDistance("backward"), []Param{{Name: "distance", Type: parser.DTNumber}})
	newFuncCall("motors.turnLeft", funcMotorsTurn("cw"), []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.turnRight", funcMotorsTurn("ccw"), []Param{{Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.rotateRPM", funcMotorsRotate("speed"), []Param{{Name: "motor", Type: parser.DTString}, {Name: "rpm", Type: parser.DTNumber}}, []Param{{Name: "motor", Type: parser.DTString}, {Name: "rpm", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.rotatePower", funcMotorsRotate("power"), []Param{{Name: "motor", Type: parser.DTString}, {Name: "power", Type: parser.DTNumber}}, []Param{{Name: "motor", Type: parser.DTString}, {Name: "power", Type: parser.DTNumber}, {Name: "duration", Type: parser.DTNumber}})
	newFuncCall("motors.rotateAngle", funcMotorsRotateAngle, []Param{{Name: "motor", Type: parser.DTString}, {Name: "angle", Type: parser.DTNumber}})
	newFuncCall("motors.driveRPM", funcMotorsDrive("speed"), []Param{{Name: "em1RPM", Type: parser.DTNumber}, {Name: "em2RPM", Type: parser.DTNumber}})
	newFuncCall("motors.drivePower", funcMotorsDrive("power"), []Param{{Name: "em1Power", Type: parser.DTNumber}, {Name: "em2Power", Type: parser.DTNumber}})
	newFuncCall("motors.stop", funcMotorsStop, []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.resetAngle", funcMotorsResetAngle, []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.lock", funcMotorsSetLock("1"), []Param{}, []Param{{Name: "motor", Type: parser.DTString}})
	newFuncCall("motors.unlock", funcMotorsSetLock("0"), []Param{}, []Param{{Name: "motor", Type: parser.DTString}})

	newFuncCall("time.wait", funcTimeWait, []Param{{Name: "duration", Type: parser.DTNumber}}, []Param{{Name: "continueCondition", Type: parser.DTBool}})
	newFuncCall("time.resetTimer", funcResetTimer)

	newFuncCall("mbot.restart", funcMBotRestart)
	newFuncCall("mbot.resetParameters", funcMBotChassisParameters("reset"))
	newFuncCall("mbot.calibrateParameters", funcMBotChassisParameters("calibrate"))

	newFuncCall("script.stop", funcScriptStop("this script"))
	newFuncCall("script.stopAll", funcScriptStop("all"))
	newFuncCall("script.stopOther", funcScriptStop("other scripts in sprite"))

	newFuncCall("lists.append", funcListsAppend, []Param{{Name: "list", Type: parser.DTStringList}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lists.remove", funcListsRemove, []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}})
	newFuncCall("lists.clear", funcListsClear, []Param{{Name: "list", Type: parser.DTStringList}})
	newFuncCall("lists.insert", funcListsInsert, []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
	newFuncCall("lists.replace", funcListsReplace, []Param{{Name: "list", Type: parser.DTStringList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTString}}, []Param{{Name: "list", Type: parser.DTNumberList}, {Name: "index", Type: parser.DTNumber}, {Name: "value", Type: parser.DTNumber}})
}

func funcAudioStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioStop, false)
	return block, nil
}

func funcAudioPlayBuzzer(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayBuzzerTone, false)

	number, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	if len(stmt.Parameters) == 2 {
		block.Type = blocks.AudioPlayBuzzerToneWithTime
		block.Inputs["number_1"] = number
		block.Inputs["number_2"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	} else {
		block.Inputs["number_2"] = number
	}

	return block, nil
}

func funcAudioPlayClip(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayClip, false)

	menuBlockType := blocks.AudioPlayClipFileNameMenu
	if len(stmt.Parameters) == 2 {
		untilDone, err := g.literal(stmt.Name, stmt.Parameters[1], parser.DTBool)
		if err != nil {
			return nil, err
		}
		if untilDone.(bool) {
			block.Type = blocks.AudioPlayClipUntilDone
			menuBlockType = blocks.AudioPlayClipUntilDoneFileNameMenu
		}
	}

	var err error
	block.Inputs["file_name"], err = g.fieldMenu(menuBlockType, "", "CYBERPI_PLAY_AUDIO_UNTIL_3_FILE_NAME", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		names := []string{"hi", "bye", "yeah", "wow", "laugh", "hum", "sad", "sigh", "annoyed", "angry", "surprised", "yummy", "curious", "embarrassed", "ready", "sprint", "sleepy", "meow", "start", "switch", "beeps", "buzzing", "jump", "level-up", "low-energy", "prompt", "right", "wrong", "ring", "score", "wake", "warning", "metal-clash", "glass-clink", "inflator", "running-water", "clockwork", "click", "current", "wood-hit", "iron", "drop", "bubble", "wave", "magic", "spitfire", "heartbeat"}
		if !slices.Contains(names, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown clip name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	return block, err
}

func funcAudioPlayInstrument(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayMusicInstrument, false)

	var err error
	names := []string{"snare", "bass-drum", "side-stick", "crash-cymbal", "open-hi-hat", "closed-hi-hat", "tambourine", "hand-clap", "claves"}
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.AudioPlayMusicInstrumentMenu, "`", "CYBERPI_PLAY_MUSIC_WITH_NOTE_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		if !slices.Contains(names, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown instrument name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcAudioPlayNote(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayNote, false)

	noteBlock := blocks.NewShadowBlock(blocks.AudioNote, block.ID)
	g.blocks[noteBlock.ID] = noteBlock

	durationParameter := 1
	if len(stmt.Parameters) == 3 {
		durationParameter = 2
		noteName, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
		if err != nil {
			return nil, err
		}
		values := map[string]int{
			"c":  0,
			"c#": 1,
			"db": 1,
			"d":  2,
			"d#": 3,
			"eb": 3,
			"e":  4,
			"f":  5,
			"f#": 6,
			"gb": 6,
			"g":  7,
			"g#": 8,
			"ab": 8,
			"a":  9,
			"a#": 10,
			"bb": 10,
			"b":  11,
		}
		value, ok := values[strings.ToLower(noteName.(string))]
		if !ok {
			return nil, g.newError("Invalid note name.", stmt.Parameters[0].(*parser.ExprLiteral).Token)
		}
		octave, err := g.literal(stmt.Name, stmt.Parameters[1], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		noteValue := int(octave.(float64))*12 + value
		noteBlock.Fields["NOTE"] = []any{strconv.Itoa(noteValue), nil}
		block.Inputs["number_1"] = []any{1, noteBlock.ID}
	} else if v, ok := stmt.Parameters[0].(*parser.ExprLiteral); ok {
		if v.Token.DataType != parser.DTNumber {
			return nil, g.newError(fmt.Sprintf("The value must be of type %s.", parser.DTNumber), v.Token)
		}
		noteBlock.Fields["NOTE"] = []any{strconv.Itoa(int(v.Token.Literal.(float64))), nil}
		block.Inputs["number_1"] = []any{1, noteBlock.ID}
	} else {
		note, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		noteBlock.Parent = nil
		noteBlock.Fields["NOTE"] = []any{"0", nil}
		block.Inputs["number_1"] = []any{3, note[1].(string), noteBlock.ID}
	}

	var err error
	block.Inputs["number_2"], err = g.value(block.ID, stmt.Name, stmt.Parameters[durationParameter], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcAudioRecordingStart(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordStart, false)
	return block, nil
}

func funcAudioRecordingStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordStop, false)
	return block, nil
}

func funcAudioRecordingPlay(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordPlay, false)

	if len(stmt.Parameters) == 1 {
		untilDone, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTBool)
		if err != nil {
			return nil, err
		}
		if untilDone.(bool) {
			block.Type = blocks.AudioRecordPlayUntilDone
		}
	}

	return block, nil
}

func funcLEDPlayAnimation(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDPlayAnimation, false)

	name, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	names := []string{"rainbow", "spindrift", "meteor_blue", "meteor_green", "flash_red", "flash_orange", "firefly"}
	if !slices.Contains(names, name.(string)) {
		return nil, g.newError(fmt.Sprintf("Unknown animation name. Available options: %s", strings.Join(names, ", ")), stmt.Parameters[0].(*parser.ExprLiteral).Token)
	}

	block.Fields["LED_animation"] = []any{name.(string), nil}

	return block, nil
}

func funcLEDSetAmbientBrightness(operation string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		blockType := blocks.UltrasonicSetBrightness
		indexType := blocks.UltrasonicSetBrightnessIndex
		orderType := blocks.UltrasonicSetBrightnessOrder
		if operation == "add" {
			blockType = blocks.UltrasonicAddBrightness
			indexType = blocks.UltrasonicAddBrightnessIndex
			orderType = blocks.UltrasonicAddBrightnessOrder
		}
		block := g.NewBlock(blockType, false)

		err := selectAmbientLight(g, block, orderType, stmt.Name, stmt.Parameters, 1, "order", "MBUILD_ULTRASONIC2_SET_BRI_ORDER", true)
		if err != nil {
			return nil, err
		}

		block.Inputs["bv"], err = g.value(g.blockID, stmt.Name, stmt.Parameters[len(stmt.Parameters)-1], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		g.noNext = true
		indexMenu := g.NewBlock(indexType, true)
		indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
		block.Inputs["index"] = []any{1, indexMenu.ID}

		return block, nil
	}
}

func funcLEDDisplayEmotion(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.UltrasonicShowEmotion, false)

	g.noNext = true
	indexMenu := g.NewBlock(blocks.UltrasonicShowEmotionIndex, true)
	indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	var err error
	block.Inputs["emotion"], err = g.fieldMenu(blocks.UltrasonicShowEmotionMenu, "", "MBUILD_ULTRASONIC2_SHOW_EMOTION_EMOTION", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		names := []string{"sleepy", "wink", "happy", "dizzy", "thinking"}
		if !slices.Contains(names, v.(string)) {
			return g.newError(fmt.Sprintf("Unknown emotion name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcLEDDeactivateAmbient(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.UltrasonicOffLED, false)

	err := selectAmbientLight(g, block, blocks.UltrasonicOffLEDInput, stmt.Name, stmt.Parameters, 0, "inputMenu_3", "MBUILD_ULTRASONIC2_SET_BRI_ORDER", true)
	if err != nil {
		return nil, err
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.UltrasonicOffLEDIndex, true)
	indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	return block, nil
}

func funcLEDDeactivateFill(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorDisableFillColor, false)

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorDisableFillColorIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func funcLEDSetFillColor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorSetFillColor, false)

	color, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}
	colors := []string{"red", "green", "blue"}
	if !slices.Contains(colors, color.(string)) {
		return nil, g.newError(fmt.Sprintf("Unknown color. Available options: %s", strings.Join(colors, ", ")), parameterToken(stmt.Parameters[0]))
	}
	block.Fields["fieldMenu_3"] = []any{color, nil}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorSetFillColorIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func funcLEDDisplay(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplay, false)

	names := make([]string, 5)
	for i := range names {
		n, err := g.literal(stmt.Name, stmt.Parameters[i], parser.DTString)
		if err != nil {
			return nil, err
		}
		colors := []string{"gray", "red", "orange", "yellow", "green", "cyan", "blue", "magenta", "white"}
		if index := slices.Index(colors, n.(string)); index >= 0 {
			names[i] = strconv.Itoa(index)
		} else {
			return nil, g.newError(fmt.Sprintf("Unknown color name. Available options: %s", strings.Join(colors, ", ")), stmt.Parameters[i].(*parser.ExprLiteral).Token)
		}
	}

	block.Fields["ledRing"] = []any{strings.Join(names, ""), nil}

	return block, nil
}

var hexColorRegex = regexp.MustCompile("^#[a-fA-F0-9]{6}$")

func funcLEDDisplayColor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplaySingleColor, false)
	if len(stmt.Parameters) > 2 {
		block.Type = blocks.LEDDisplaySingleColorWithRGB
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBFieldMenu, stmt, 3, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorFieldMenu, stmt, 1, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString, hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcLEDDisplayColorFor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplaySingleColorWithTime, false)
	if len(stmt.Parameters) > 3 {
		block.Type = blocks.LEDDisplaySingleColorWithRGBAndTime
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBAndTimeFieldMenu, stmt, 4, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_5"], err = g.value(block.ID, stmt.Name, stmt.Parameters[4], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithTimeFieldMenu, stmt, 2, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString, hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			return nil, err
		}
		block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcLEDDeactivate(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDOff, false)

	err := selectLED(g, block, blocks.LEDOffFieldMenu, stmt, 0, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")

	return block, err
}

func selectLED(g *generator, block *blocks.Block, menuBlockType blocks.BlockType, stmt *parser.StmtFuncCall, paramCountWithoutLED int, menuFieldKey string) error {
	if len(stmt.Parameters) == paramCountWithoutLED {
		stmt.Parameters = append([]parser.Expr{&parser.ExprLiteral{
			Token: parser.Token{
				Type:     parser.TkLiteral,
				Literal:  "all",
				DataType: parser.DTString,
			},
		}}, stmt.Parameters...)
	}

	errWrongType := errors.New("wrong-type")
	var err error
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Name, stmt.Parameters[0], "", func(v any, token parser.Token) error {
		if str, ok := v.(string); !ok {
			return errWrongType
		} else {
			if str != "all" {
				return g.newError("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", token)
			}
		}
		return nil
	})
	if err == errWrongType {
		block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, func(v any, token parser.Token) error {
			nr := int(v.(float64))
			if nr != 1 && nr != 2 && nr != 3 && nr != 4 && nr != 5 {
				return g.newError("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", token)
			}
			return nil
		})
	}
	return err
}

func selectAmbientLight(g *generator, block *blocks.Block, menuBlockType blocks.BlockType, token parser.Token, parameters []parser.Expr, paramCountWithoutLight int, orderKey string, menuFieldKey string, allowAll bool) error {
	errorMsg := "Unknown light. Available options: 1, 2, 3, 4, 5, 6, 7, 8"
	if allowAll {
		errorMsg = "Unknown light. Available options: \"all\", 1, 2, 3, 4, 5, 6, 7, 8"
	}

	if len(parameters) == paramCountWithoutLight {
		if !allowAll {
			return g.newError(errorMsg, token)
		}
		parameters = append([]parser.Expr{&parser.ExprLiteral{
			Token: parser.Token{
				Type:     parser.TkLiteral,
				Literal:  "all",
				DataType: parser.DTString,
			},
		}}, parameters...)
	}

	errWrongType := errors.New("wrong-type")
	var err error
	if allowAll {
		block.Inputs[orderKey], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, token, parameters[0], "", func(v any, token parser.Token) error {
			if str, ok := v.(string); !ok {
				return errWrongType
			} else {
				if str != "all" {
					return g.newError(errorMsg, token)
				}
			}
			return nil
		})
	} else {
		err = errWrongType
	}
	if err == errWrongType {
		block.Inputs[orderKey], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, token, parameters[0], parser.DTNumber, func(v any, token parser.Token) error {
			nr := int(v.(float64))
			if nr < 1 || nr > 8 {
				return g.newError(errorMsg, token)
			}
			return nil
		})
	}
	return err
}

func funcDisplayPrint(newLine bool) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.DisplayPrint, false)
		if newLine {
			block.Type = blocks.DisplayPrintln
		}

		var err error
		block.Inputs["string_2"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString)
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcDisplaySetFontSize(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetFont, false)

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.DisplaySetFontMenu, "", "CYBERPI_CONSOLE_SET_FONT_INPUTMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, func(v any, token parser.Token) error {
		sizes := []int{12, 16, 24, 32}
		if math.Mod(v.(float64), 1.0) != 0 || !slices.Contains(sizes, int(v.(float64))) {
			options := ""
			for _, s := range sizes {
				options = fmt.Sprintf("%s, %d", options, s)
			}
			return g.newError(fmt.Sprintf("Unknown size. Available options: %s", options), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplaySetColor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetBrushColor, false)
	var err error
	if len(stmt.Parameters) == 3 {
		block.Type = blocks.DisplaySetBrushColorRGB
		block.Inputs["number_1"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, -1, 0, 255)
		if err != nil {
			return nil, err
		}
	} else {
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcDisplayShowLabel(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLabelShowSomewhereWithSize, false)
	number, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	if math.Mod(number.(float64), 1.0) != 0 || number.(float64) < 0 || number.(float64) > 8 {
		return nil, g.newError("The label number must lie between 0 and 8.", stmt.Parameters[0].(*parser.ExprLiteral).Token)
	}
	block.Fields["fieldMenu_1"] = []any{fmt.Sprintf("%d", int(number.(float64))-1), nil}

	block.Inputs["string_2"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString)
	if err != nil {
		return nil, err
	}

	sizeIndex := 3
	if len(stmt.Parameters) == 5 {
		sizeIndex = 4
		block.Type = blocks.DisplayLabelShowXYWithSize

		block.Inputs["number_2"], err = g.value(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	} else {
		location, err := g.literal(stmt.Name, stmt.Parameters[2], parser.DTString)
		if err != nil {
			return nil, err
		}
		locations := []string{"top_left", "top_mid", "top_right", "mid_left", "center", "mid_right", "bottom_left", "bottom_mid", "bottom_right"}
		if !slices.Contains(locations, location.(string)) {
			return nil, g.newError(fmt.Sprintf("Unknown label location. Available options: %s", strings.Join(locations, ", ")), stmt.Parameters[2].(*parser.ExprLiteral).Token)
		}
		block.Fields["fieldMenu_2"] = []any{location, nil}
	}

	block.Inputs["inputMenu_4"], err = g.fieldMenu(blocks.DisplayLabelShowSomewhereWithSizeMenu, "", "CYBERPI_CONSOLE_SET_FONT_INPUTMENU_1", block.ID, stmt.Name, stmt.Parameters[sizeIndex], parser.DTNumber, func(v any, token parser.Token) error {
		sizes := []int{12, 16, 24, 32}
		if math.Mod(v.(float64), 1.0) != 0 || !slices.Contains(sizes, int(v.(float64))) {
			options := ""
			for _, s := range sizes {
				options = fmt.Sprintf("%s, %d", options, s)
			}
			return g.newError(fmt.Sprintf("Unknown size. Available options: %s", options), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayLineChartAddData(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLineChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayLineChartSetInterval(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLineChartSetInterval, false)

	var err error
	block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayBarChartAddData(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayBarChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayTableAddData(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayTableAddDataAtRowColumn, false)

	var err error
	block.Inputs["string_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.DisplayTableAddDataAtRowColumnMenu, "", "CYBERPI_DISPLAY_TABLE_ADD_DATA_AT_ROW_COLUMN_2_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newError("The value must be an integer.", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	block.Inputs["fieldMenu_2"], err = g.fieldMenu(blocks.DisplayTableAddDataAtRowColumnMenu, "", "CYBERPI_DISPLAY_TABLE_ADD_DATA_AT_ROW_COLUMN_2_FIELDMENU_2", block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newError("The value must be an integer.", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplaySetOrientation(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetOrientation, false)

	var err error
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.DisplaySetOrientationMenu, "", "CYBERPI_DISPLAY_ROTATE_TO_2_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newError("The value must be an integer.", token)
		}
		value := int(v.(float64))
		if value != -90 && value != 0 && value != 90 && value != 180 {
			return g.newError("The orientation must be either -90, 0, 90 or 180 degrees.", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayClear(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayClear, false)
	return block, nil
}

func funcLEDMove(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDMove, false)

	var err error
	block.Inputs["led_number"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcNetBroadcast(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetSetWifiBroadcast, false)

	var err error
	block.Inputs["message"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	if len(stmt.Parameters) == 2 {
		block.Type = blocks.NetSetWifiBroadcastWithValue
		block.Inputs["value"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString)
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcNetSetChannel(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetSetWifiChannel, false)

	channel, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	if int(channel.(float64)) != 1 && int(channel.(float64)) != 6 && int(channel.(float64)) != 11 {
		return nil, g.newError("Invalid channel. Allowed options: 1, 6, 11", parameterToken(stmt.Parameters[0]))
	}
	block.Fields["channel"] = []any{fmt.Sprintf("%d", int(channel.(float64))), nil}

	return block, nil
}

func funcNetConnect(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetConnectWifi, false)

	var err error
	block.Inputs["ssid"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	block.Inputs["wifipassword"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcNetReconnect(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetWifiReconnect, false)
	return block, nil
}

func funcNetDisconnect(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetWifiDisconnect, false)
	return block, nil
}

func funcSensorsResetAngle(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorsResetAxisRotationAngle, false)
	value, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}
	axes := []string{"all", "x", "y", "z"}
	if !slices.Contains(axes, value.(string)) {
		return nil, g.newError(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(axes, ", ")), parameterToken(stmt.Parameters[0]))
	}
	block.Fields["axis"] = []any{value, nil}
	return block, nil
}

func funcSensorsResetYawAngle(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorsResetYaw, false)
	return block, nil
}

func funcSensorsDefineColor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorDefineColor, false)

	var err error
	block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, -1, 0, 255)
	if err != nil {
		return nil, err
	}
	block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, -1, 0, 255)
	if err != nil {
		return nil, err
	}
	block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, -1, 0, 255)
	if err != nil {
		return nil, err
	}

	if len(stmt.Parameters) == 4 {
		block.Inputs["tolerance"], err = g.value(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	} else {
		block.Inputs["tolerance"] = []any{1, []any{4, "50"}}
	}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorDefineColorIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func funcMotorsRun(direction string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2MoveDirectionWithRPM, false)

		block.Fields["DIRECTION"] = []any{direction, nil}

		var err error
		block.Inputs["POWER"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		if len(stmt.Parameters) == 2 {
			block.Type = blocks.Mbot2MoveDirectionWithTime
			block.Inputs["TIME"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
			if err != nil {
				return nil, err
			}
		}

		return block, nil
	}
}

func funcMotorsRunDistance(direction string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2MoveMoveWithCmAndInch, false)

		block.Fields["DIRECTION"] = []any{direction, nil}
		block.Fields["fieldMenu_3"] = []any{"cm", nil}

		var err error
		block.Inputs["POWER"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		return block, nil
	}
}

func funcMotorsTurn(direction string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2CwAndCcwWithAngle, false)
		block.Fields["fieldMenu_1"] = []any{direction, nil}

		var err error
		block.Inputs["ANGLE"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcMotorsRotate(unit string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		blockType := blocks.Mbot2EncoderMotorSet
		menuType := blocks.Mbot2EncoderMotorSetMenu
		inputField := "inputMenu_1"
		if len(stmt.Parameters) == 3 {
			inputField = "fieldMenu_1"
			blockType = blocks.Mbot2EncoderMotorSetWithTime
			menuType = blocks.Mbot2EncoderMotorSetWithTimeMenu
		}

		block := g.NewBlock(blockType, false)

		var err error
		block.Inputs[inputField], err = g.fieldMenu(menuType, "", "MBOT2_ENCODER_MOTOR_SET_WITH_TIME_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
			str := v.(string)
			if str != "ALL" && str != "EM1" && str != "EM2" {
				return g.newError("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_4"] = []any{unit, nil}

		block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		if len(stmt.Parameters) == 3 {
			block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber)
			if err != nil {
				return nil, err
			}
		}

		return block, nil
	}
}

func funcMotorsRotateAngle(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorSetWithTimeAngleAndCircle, false)

	var err error
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorSetWithTimeAngleAndCircleMenu, "", "MBOT2_ENCODER_MOTOR_SET_WITH_TIME_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
		str := v.(string)
		if str != "ALL" && str != "EM1" && str != "EM2" {
			return g.newError("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	return block, err
}

func funcMotorsDrive(unit string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2EncoderMotorDrivePower, false)
		rightPowerKey := "number_2"
		if unit == "speed" {
			block.Type = blocks.Mbot2EncoderMotorDriveSpeed
			rightPowerKey = "RIGHT_POWER"
		}

		var err error
		block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		block.Inputs[rightPowerKey], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcMotorsStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorStop, false)

	encoderMotor := "ALL"
	if len(stmt.Parameters) == 1 {
		motor, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
		if err != nil {
			return nil, err
		}
		encoderMotor = motor.(string)
		if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
			return nil, g.newError("Unknown encoder motor. Available options: ALL, EM1, EM2", parameterToken(stmt.Parameters[0]))
		}
	}

	block.Fields["fieldMenu_1"] = []any{encoderMotor, nil}

	return block, nil
}

func funcMotorsResetAngle(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorResetAngle, false)

	var motor parser.Expr
	if len(stmt.Parameters) == 1 {
		motor = stmt.Parameters[0]
	} else {
		motor = &parser.ExprLiteral{
			Token: parser.Token{
				Type:     parser.TkLiteral,
				Literal:  "ALL",
				DataType: parser.DTString,
				Line:     stmt.Name.Line,
			},
		}
	}

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorResetAngleMenu, "", "MBOT2_ENCODER_MOTOR_STOP_FIELDMENU_1", block.ID, stmt.Name, motor, parser.DTString, func(v any, token parser.Token) error {
		encoderMotor := v.(string)
		if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
			return g.newError("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcMotorsSetLock(value string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2EncoderMotorLockUnlock, false)

		block.Fields["fieldMenu_2"] = []any{value, nil}

		var motor parser.Expr
		if len(stmt.Parameters) == 1 {
			motor = stmt.Parameters[0]
		} else {
			motor = &parser.ExprLiteral{
				Token: parser.Token{
					Type:     parser.TkLiteral,
					Literal:  "ALL",
					DataType: parser.DTString,
					Line:     stmt.Name.Line,
				},
			}
		}

		var err error
		block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorLockUnlockMenu, "", "MBOT2_ENCODER_MOTOR_STOP_FIELDMENU_1", block.ID, stmt.Name, motor, parser.DTString, func(v any, token parser.Token) error {
			encoderMotor := v.(string)
			if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
				return g.newError("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcTimeWait(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ControlWaitUntil, false)

	condition, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTBool)
	if err == nil {
		block.Inputs["CONDITION"] = condition
		return block, nil
	} else {
		block.Type = blocks.ControlWait
		seconds, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		block.Inputs["DURATION"] = seconds
	}

	return block, nil
}

func funcResetTimer(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2TimerReset, false)
	return block, nil
}

func funcMBotRestart(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ControlRestart, false)
	return block, nil
}

func funcMBotChassisParameters(parameter string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2SetParameters, false)

		if parameter == "calibrate" {
			parameter = "set_auto"
		}

		block.Fields["PARA"] = []any{parameter, nil}

		return block, nil
	}
}

func funcScriptStop(stopOption string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.ControlStop, false)
		block.NoNext = true

		block.Fields["STOP_OPTION"] = []any{stopOption, nil}

		hasNext := "false"
		if stopOption == "other scripts in sprite" {
			hasNext = "true"
			block.NoNext = false
		}
		block.Mutation = map[string]any{
			"tagName":  "mutation",
			"children": []any{},
			"hasnext":  hasNext,
		}

		return block, nil
	}
}

func funcListsAppend(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListAdd, false)
	valueType, err := selectList(g, block, stmt.Name, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], valueType)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcListsRemove(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListDelete, false)
	_, err := selectList(g, block, stmt.Name, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcListsClear(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListClear, false)
	_, err := selectList(g, block, stmt.Name, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcListsInsert(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListInsert, false)
	valueType, err := selectList(g, block, stmt.Name, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Name, stmt.Parameters[2], valueType)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcListsReplace(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListReplace, false)
	valueType, err := selectList(g, block, stmt.Name, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Name, stmt.Parameters[2], valueType)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func selectList(g *generator, block *blocks.Block, token parser.Token, param parser.Expr) (parser.DataType, error) {
	if i, ok := param.(*parser.ExprIdentifier); ok {
		token = i.Name
		if l, ok := g.lists[i.Name.Lexeme]; ok {
			block.Fields["LIST"] = []any{
				l.Name.Lexeme,
				l.ID,
			}
			return parser.DataType(strings.TrimSuffix(string(l.DataType), "[]")), nil
		}
		return "", g.newError("Unknown list.", token)
	}
	if l, ok := param.(*parser.ExprLiteral); ok {
		token = l.Token
	}
	return "", g.newError("Expected list.", token)
}
