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

type FuncCall func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error)

var FuncCalls = map[string]FuncCall{
	"audio.stop":           funcAudioStop,
	"audio.playBuzzer":     funcAudioPlayBuzzer,
	"audio.playClip":       funcAudioPlayClip,
	"audio.playInstrument": funcAudioPlayInstrument,
	"audio.playNote":       funcAudioPlayNote,
	"audio.record.start":   funcAudioRecordingStart,
	"audio.record.stop":    funcAudioRecordingStop,
	"audio.record.play":    funcAudioRecordingPlay,

	"lights.back.playAnimation":   funcLEDPlayAnimation,
	"lights.front.setBrightness":  funcLEDSetAmbientBrightness("set"),
	"lights.front.addBrightness":  funcLEDSetAmbientBrightness("add"),
	"lights.front.displayEmotion": funcLEDDisplayEmotion,
	"lights.front.deactivate":     funcLEDDeactivateAmbient,
	"lights.bottom.deactivate":    funcLEDDeactivateFill,
	"lights.bottom.setColor":      funcLEDSetFillColor,
	"lights.back.display":         funcLEDDisplay,
	"lights.back.displayColor":    funcLEDDisplayColor,
	"lights.back.displayColorFor": funcLEDDisplayColorFor,
	"lights.back.deactivate":      funcLEDDeactivate,
	"lights.back.move":            funcLEDMove,

	"display.print":                 funcDisplayPrint(false),
	"display.println":               funcDisplayPrint(true),
	"display.setFontSize":           funcDisplaySetFontSize,
	"display.setColor":              funcDisplaySetColor,
	"display.showLabel":             funcDisplayShowLabel,
	"display.lineChart.addData":     funcDisplayLineChartAddData,
	"display.lineChart.setInterval": funcDisplayLineChartSetInterval,
	"display.barChart.addData":      funcDisplayBarChartAddData,
	"display.table.addData":         funcDisplayTableAddData,
	"display.setOrientation":        funcDisplaySetOrientation,
	"display.clear":                 funcDisplayClear,

	"display.setBackgroundColor": funcDisplaySetBackgroundColor,
	"display.render":             funcDisplayRender,

	"sprite.fromIcon":   funcSpriteFromIcon,
	"sprite.fromText":   funcSpriteFromText,
	"sprite.fromQR":     funcSpriteFromQR,
	"sprite.flipH":      funcSpriteFlip("x"),
	"sprite.flipV":      funcSpriteFlip("y"),
	"sprite.delete":     funcSpriteDelete,
	"sprite.setAnchor":  funcSpriteSetAnchor,
	"sprite.moveLeft":   funcSpriteMove("left"),
	"sprite.moveRight":  funcSpriteMove("x"),
	"sprite.moveUp":     funcSpriteMove("up"),
	"sprite.moveDown":   funcSpriteMove("y"),
	"sprite.moveTo":     funcSpriteMoveTo,
	"sprite.moveRandom": funcSpriteMoveRandom,
	"sprite.rotate":     funcSpriteRotate,
	"sprite.rotateTo":   funcSpriteRotateTo,
	"sprite.setScale":   funcSpriteSetScale,
	"sprite.setColor":   funcSpriteSetColor,
	"sprite.resetColor": funcSpriteResetColor,
	"sprite.show":       funcSpriteShowHide("show"),
	"sprite.hide":       funcSpriteShowHide("hide"),
	"sprite.toFront":    funcSpriteSetLayer("z_max"),
	"sprite.toBack":     funcSpriteSetLayer("z_min"),
	"sprite.layerUp":    funcSpriteChangeLayer("z_up"),
	"sprite.layerDown":  funcSpriteChangeLayer("z_down"),

	"draw.begin":        funcDrawBegin,
	"draw.finish":       funcDrawFinish,
	"draw.clear":        funcDrawClear,
	"draw.setColor":     funcDrawSetColor,
	"draw.setThickness": funcDrawSetThickness,
	"draw.setSpeed":     funcDrawSetSpeed,
	"draw.rotate":       funcDrawRotate,
	"draw.rotateTo":     funcDrawRotateTo,
	"draw.line":         funcDrawLine,
	"draw.circle":       funcDrawCircle,
	"draw.moveUp":       funcDrawMove("up"),
	"draw.moveDown":     funcDrawMove("y"),
	"draw.moveLeft":     funcDrawMove("left"),
	"draw.moveRight":    funcDrawMove("x"),
	"draw.moveTo":       funcDrawMoveTo,
	"draw.moveToCenter": funcDrawMoveToCenter,
	"draw.save":         funcDrawSave,

	"net.broadcast":  funcNetBroadcast,
	"net.setChannel": funcNetSetChannel,
	"net.connect":    funcNetConnect,
	"net.reconnect":  funcNetReconnect,
	"net.disconnect": funcNetDisconnect,

	"sensors.resetAngle":             funcSensorsResetAngle,
	"sensors.resetYawAngle":          funcSensorsResetYawAngle,
	"sensors.defineColor":            funcSensorsDefineColor,
	"sensors.calibrateColors":        funcSensorsCalibrateColors,
	"sensors.enhancedColorDetection": funcSensorsEnhancedColorDetection,

	"motors.run":                  funcMotorsRun("forward"),
	"motors.runBackward":          funcMotorsRun("backward"),
	"motors.moveDistance":         funcMotorsRunDistance("forward"),
	"motors.moveDistanceBackward": funcMotorsRunDistance("backward"),
	"motors.turnLeft":             funcMotorsTurn("cw"),
	"motors.turnRight":            funcMotorsTurn("ccw"),
	"motors.rotateRPM":            funcMotorsRotate("speed"),
	"motors.rotatePower":          funcMotorsRotate("power"),
	"motors.rotateAngle":          funcMotorsRotateAngle,
	"motors.driveRPM":             funcMotorsDrive("speed"),
	"motors.drivePower":           funcMotorsDrive("power"),
	"motors.stop":                 funcMotorsStop,
	"motors.resetAngle":           funcMotorsResetAngle,
	"motors.lock":                 funcMotorsSetLock("1"),
	"motors.unlock":               funcMotorsSetLock("0"),

	"time.wait":       funcTimeWait,
	"time.resetTimer": funcResetTimer,

	"mbot.restart":             funcMBotRestart,
	"mbot.resetParameters":     funcMBotChassisParameters("reset"),
	"mbot.calibrateParameters": funcMBotChassisParameters("calibrate"),

	"script.stop":      funcScriptStop("this script"),
	"script.stopAll":   funcScriptStop("all"),
	"script.stopOther": funcScriptStop("other scripts in sprite"),

	"lists.append":  funcListsAppend,
	"lists.remove":  funcListsRemove,
	"lists.clear":   funcListsClear,
	"lists.insert":  funcListsInsert,
	"lists.replace": funcListsReplace,
}

func funcAudioStop(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioStop, false)
	return block, nil
}

func funcAudioPlayBuzzer(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayBuzzerTone, false)

	number, err := g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	if len(stmt.Parameters) == 2 {
		block.Type = blocks.AudioPlayBuzzerToneWithTime
		block.Inputs["number_1"] = number
		block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		block.Inputs["number_2"] = number
	}

	return block, nil
}

func funcAudioPlayClip(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayClip, false)

	menuBlockType := blocks.AudioPlayClipFileNameMenu
	if len(stmt.Parameters) == 2 {
		untilDone, err := g.literal(stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		} else if untilDone.(bool) {
			block.Type = blocks.AudioPlayClipUntilDone
			menuBlockType = blocks.AudioPlayClipUntilDoneFileNameMenu
		}
	}

	var err error
	block.Inputs["file_name"], err = g.fieldMenu(menuBlockType, "", "CYBERPI_PLAY_AUDIO_UNTIL_3_FILE_NAME", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		names := []string{"hi", "bye", "yeah", "wow", "laugh", "hum", "sad", "sigh", "annoyed", "angry", "surprised", "yummy", "curious", "embarrassed", "ready", "sprint", "sleepy", "meow", "start", "switch", "beeps", "buzzing", "jump", "level-up", "low-energy", "prompt", "right", "wrong", "ring", "score", "wake", "warning", "metal-clash", "glass-clink", "inflator", "running-water", "clockwork", "click", "current", "wood-hit", "iron", "drop", "bubble", "wave", "magic", "spitfire", "heartbeat"}
		if !slices.Contains(names, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown clip name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	return block, err
}

func funcAudioPlayInstrument(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayMusicInstrument, false)

	var err error
	names := []string{"snare", "bass-drum", "side-stick", "crash-cymbal", "open-hi-hat", "closed-hi-hat", "tambourine", "hand-clap", "claves"}
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.AudioPlayMusicInstrumentMenu, "`", "CYBERPI_PLAY_MUSIC_WITH_NOTE_FIELDMENU_1", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		if !slices.Contains(names, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown instrument name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcAudioPlayNote(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioPlayNote, false)

	noteBlock := blocks.NewShadowBlock(blocks.AudioNote, block.ID)
	g.blocks[noteBlock.ID] = noteBlock

	durationParameter := 1
	if len(stmt.Parameters) == 3 {
		durationParameter = 2
		noteName, err := g.literal(stmt.Parameters[0])
		if err != nil {
			g.errors = append(g.errors, err)
		} else {
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
				g.errors = append(g.errors, g.newErrorExpr("Invalid note name.", stmt.Parameters[0]))
			}
			octave, err := g.literal(stmt.Parameters[1])
			if err != nil {
				g.errors = append(g.errors, err)
			} else {
				noteValue := int(octave.(float64))*12 + value
				noteBlock.Fields["NOTE"] = []any{strconv.Itoa(noteValue), nil}
				block.Inputs["number_1"] = []any{1, noteBlock.ID}
			}
		}
	} else if v, ok := stmt.Parameters[0].(*parser.ExprLiteral); ok {
		noteBlock.Fields["NOTE"] = []any{strconv.Itoa(int(v.Token.Literal.(float64))), nil}
		block.Inputs["number_1"] = []any{1, noteBlock.ID}
	} else {
		note, err := g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			g.errors = append(g.errors, err)
		} else {
			noteBlock.Parent = nil
			noteBlock.Fields["NOTE"] = []any{"0", nil}
			block.Inputs["number_1"] = []any{3, note[1].(string), noteBlock.ID}
		}
	}

	var err error
	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[durationParameter])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcAudioRecordingStart(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordStart, false)
	return block, nil
}

func funcAudioRecordingStop(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordStop, false)
	return block, nil
}

func funcAudioRecordingPlay(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.AudioRecordPlay, false)

	if len(stmt.Parameters) == 1 {
		untilDone, err := g.literal(stmt.Parameters[0])
		if err != nil {
			return nil, err
		}
		if untilDone.(bool) {
			block.Type = blocks.AudioRecordPlayUntilDone
		}
	}

	return block, nil
}

func funcLEDPlayAnimation(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDPlayAnimation, false)

	name, err := g.literal(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	names := []string{"rainbow", "spindrift", "meteor_blue", "meteor_green", "flash_red", "flash_orange", "firefly"}
	if !slices.Contains(names, name.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown animation name. Available options: %s", strings.Join(names, ", ")), stmt.Parameters[0])
	}

	block.Fields["LED_animation"] = []any{name.(string), nil}

	return block, nil
}

func funcLEDSetAmbientBrightness(operation string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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
			g.errors = append(g.errors, err)
		}

		block.Inputs["bv"], err = g.value(g.blockID, stmt.Parameters[len(stmt.Parameters)-1])
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

func funcLEDDisplayEmotion(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.UltrasonicShowEmotion, false)

	g.noNext = true
	indexMenu := g.NewBlock(blocks.UltrasonicShowEmotionIndex, true)
	indexMenu.Fields["MBUILD_ULTRASONIC2_GET_DISTANCE_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}

	var err error
	block.Inputs["emotion"], err = g.fieldMenu(blocks.UltrasonicShowEmotionMenu, "", "MBUILD_ULTRASONIC2_SHOW_EMOTION_EMOTION", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		names := []string{"sleepy", "wink", "happy", "dizzy", "thinking"}
		if !slices.Contains(names, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown emotion name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcLEDDeactivateAmbient(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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

func funcLEDDeactivateFill(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorDisableFillColor, false)

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorDisableFillColorIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func funcLEDSetFillColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorSetFillColor, false)

	color, err := g.literal(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	colors := []string{"red", "green", "blue"}
	if !slices.Contains(colors, color.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown color. Available options: %s", strings.Join(colors, ", ")), stmt.Parameters[0])
	}
	block.Fields["fieldMenu_3"] = []any{color, nil}

	g.noNext = true
	indexMenu := g.NewBlock(blocks.SensorColorSetFillColorIndex, true)
	indexMenu.Fields["MBUILD_QUAD_COLOR_SENSOR_GET_STA_WITH_INPUTMENU_INDEX"] = []any{"1", nil}
	block.Inputs["index"] = []any{1, indexMenu.ID}
	return block, nil
}

func funcLEDDisplay(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplay, false)

	names := make([]string, 5)
	for i := range names {
		n, err := g.literal(stmt.Parameters[i])
		if err != nil {
			g.errors = append(g.errors, err)
			continue
		}
		colors := []string{"gray", "red", "orange", "yellow", "green", "cyan", "blue", "magenta", "white"}
		if index := slices.Index(colors, n.(string)); index >= 0 {
			names[i] = strconv.Itoa(index)
		} else {
			g.errors = append(g.errors, g.newErrorExpr(fmt.Sprintf("Unknown color name. Available options: %s", strings.Join(colors, ", ")), stmt.Parameters[i]))
		}
	}

	block.Fields["ledRing"] = []any{strings.Join(names, ""), nil}

	return block, nil
}

var hexColorRegex = regexp.MustCompile("^#[a-fA-F0-9]{6}$")

func funcLEDDisplayColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplaySingleColor, false)
	if len(stmt.Parameters) > 2 {
		block.Type = blocks.LEDDisplaySingleColorWithRGB
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBFieldMenu, stmt, 3, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Parameters[3], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorFieldMenu, stmt, 1, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Parameters[1], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcLEDDisplayColorFor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDDisplaySingleColorWithTime, false)
	if len(stmt.Parameters) > 3 {
		block.Type = blocks.LEDDisplaySingleColorWithRGBAndTime
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBAndTimeFieldMenu, stmt, 4, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Parameters[3], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_5"], err = g.value(block.ID, stmt.Parameters[4])
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithTimeFieldMenu, stmt, 2, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Parameters[1], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[2])
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcLEDDeactivate(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDOff, false)

	err := selectLED(g, block, blocks.LEDOffFieldMenu, stmt, 0, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")

	return block, err
}

func selectLED(g *generator, block *blocks.Block, menuBlockType blocks.BlockType, stmt *parser.StmtCall, paramCountWithoutLED int, menuFieldKey string) error {
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
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		if str, ok := v.(string); !ok {
			return errWrongType
		} else {
			if str != "all" {
				return g.newErrorTk("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", token)
			}
		}
		return nil
	})
	if err == errWrongType {
		block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
			nr := int(v.(float64))
			if nr != 1 && nr != 2 && nr != 3 && nr != 4 && nr != 5 {
				return g.newErrorTk("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", token)
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
			return g.newErrorTk(errorMsg, token)
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
		block.Inputs[orderKey], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, parameters[0], func(v any, token parser.Token) error {
			if str, ok := v.(string); !ok {
				return errWrongType
			} else {
				if str != "all" {
					return g.newErrorTk(errorMsg, token)
				}
			}
			return nil
		})
	} else {
		err = errWrongType
	}
	if err == errWrongType {
		block.Inputs[orderKey], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, parameters[0], func(v any, token parser.Token) error {
			nr := int(v.(float64))
			if nr < 1 || nr > 8 {
				return g.newErrorTk(errorMsg, token)
			}
			return nil
		})
	}
	return err
}

func funcDisplayPrint(newLine bool) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.DisplayPrint, false)
		if newLine {
			block.Type = blocks.DisplayPrintln
		}

		var err error
		block.Inputs["string_2"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcDisplaySetFontSize(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetFont, false)

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.DisplaySetFontMenu, "", "CYBERPI_CONSOLE_SET_FONT_INPUTMENU_1", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		sizes := []int{12, 16, 24, 32}
		if math.Mod(v.(float64), 1.0) != 0 || !slices.Contains(sizes, int(v.(float64))) {
			options := ""
			for _, s := range sizes {
				options = fmt.Sprintf("%s, %d", options, s)
			}
			return g.newErrorTk(fmt.Sprintf("Unknown size. Available options: %s", options), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplaySetColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetBrushColor, false)
	var err error
	if len(stmt.Parameters) == 3 {
		block.Type = blocks.DisplaySetBrushColorRGB
		block.Inputs["number_1"], err = g.valueInRange(block.ID, stmt.Parameters[0], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Parameters[0], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcDisplayShowLabel(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLabelShowSomewhereWithSize, false)
	number, err := g.literal(stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	} else {
		if math.Mod(number.(float64), 1.0) != 0 || number.(float64) < 0 || number.(float64) > 8 {
			return nil, g.newErrorExpr("The label number must lie between 0 and 8.", stmt.Parameters[0])
		}
		block.Fields["fieldMenu_1"] = []any{fmt.Sprintf("%d", int(number.(float64))-1), nil}
	}

	block.Inputs["string_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	sizeIndex := 3
	if len(stmt.Parameters) == 5 {
		sizeIndex = 4
		block.Type = blocks.DisplayLabelShowXYWithSize

		block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[2])
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[3])
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		location, err := g.literal(stmt.Parameters[2])
		if err != nil {
			g.errors = append(g.errors, err)
		} else {
			locations := []string{"top_left", "top_mid", "top_right", "mid_left", "center", "mid_right", "bottom_left", "bottom_mid", "bottom_right"}
			if !slices.Contains(locations, location.(string)) {
				return nil, g.newErrorExpr(fmt.Sprintf("Unknown label location. Available options: %s", strings.Join(locations, ", ")), stmt.Parameters[2])
			}
			block.Fields["fieldMenu_2"] = []any{location, nil}
		}
	}

	block.Inputs["inputMenu_4"], err = g.fieldMenu(blocks.DisplayLabelShowSomewhereWithSizeMenu, "", "CYBERPI_CONSOLE_SET_FONT_INPUTMENU_1", block.ID, stmt.Parameters[sizeIndex], func(v any, token parser.Token) error {
		sizes := []int{12, 16, 24, 32}
		if math.Mod(v.(float64), 1.0) != 0 || !slices.Contains(sizes, int(v.(float64))) {
			options := ""
			for _, s := range sizes {
				options = fmt.Sprintf("%s, %d", options, s)
			}
			return g.newErrorTk(fmt.Sprintf("Unknown size. Available options: %s", options), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayLineChartAddData(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLineChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayLineChartSetInterval(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayLineChartSetInterval, false)

	var err error
	block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayBarChartAddData(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayBarChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayTableAddData(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayTableAddDataAtRowColumn, false)

	var err error
	block.Inputs["string_3"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.DisplayTableAddDataAtRowColumnMenu, "", "CYBERPI_DISPLAY_TABLE_ADD_DATA_AT_ROW_COLUMN_2_FIELDMENU_1", block.ID, stmt.Parameters[1], func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newErrorTk("The value must be an integer.", token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["fieldMenu_2"], err = g.fieldMenu(blocks.DisplayTableAddDataAtRowColumnMenu, "", "CYBERPI_DISPLAY_TABLE_ADD_DATA_AT_ROW_COLUMN_2_FIELDMENU_2", block.ID, stmt.Parameters[2], func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newErrorTk("The value must be an integer.", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplaySetOrientation(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplaySetOrientation, false)

	var err error
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.DisplaySetOrientationMenu, "", "CYBERPI_DISPLAY_ROTATE_TO_2_FIELDMENU_1", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		if math.Mod(v.(float64), 1) != 0 {
			return g.newErrorTk("The value must be an integer.", token)
		}
		value := int(v.(float64))
		if value != -90 && value != 0 && value != 90 && value != 180 {
			return g.newErrorTk("The orientation must be either -90, 0, 90 or 180 degrees.", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayClear(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DisplayClear, false)
	return block, nil
}

func funcDisplaySetBackgroundColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteSetBackgroundFillColor, false)

	var err error
	if len(stmt.Parameters) == 1 {
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Parameters[0], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		block.Type = blocks.SpriteSetBackgroundFillColorRGB

		block.Inputs["number_1"], err = g.valueInRange(block.ID, stmt.Parameters[0], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcDisplayRender(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteScreenRender, false)
	return block, nil
}

func funcSpriteFromIcon(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteDrawPixelWithIcon, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SpriteDrawPixelWithIconInputMenu, "", "CYBERPI_SPRITE_DRAW_PIXEL_WITH_ICON_INPUTMENU_2", block.ID, stmt.Parameters[1], func(v any, token parser.Token) error {
		names := []string{"Music", "Image", "Video", "Clock", "Play", "Pause", "Next", "Prev", "Sound", "Temperature", "Light", "Motion", "Home", "Gear", "List", "Right", "Wrong", "Shut_down", "Refresh", "Trash_can", "Download", "Cloudy", "Rain", "Snow", "Train", "Rocket", "Truck", "Car", "Droplet", "Distance", "Fire", "Magnetic", "Gas", "Vision", "Color", "Overcast", "Sandstorm", "Foggy"}
		if !slices.Contains(names, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown icon name. Available options: %s", strings.Join(names, ", ")), token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcSpriteFromText(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteDrawText, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["string_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteFromQR(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteDrawQR, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["string_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteFlip(axis string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteMirrorWithAxis, false)

		var err error
		block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_3"] = []any{axis, nil}

		return block, nil
	}
}

func funcSpriteDelete(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteDelete, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcSpriteSetAnchor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteSetAlign, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["inputMenu_2"], err = g.fieldMenu(blocks.SpriteSetAlignInputMenu, "", "CYBERPI_SPRITE_SET_ALIGN_INPUTMENU_2", block.ID, stmt.Parameters[1], func(v any, token parser.Token) error {
		locations := []string{"top_left", "top_mid", "top_right", "mid_left", "center", "mid_right", "bottom_left", "bottom_mid", "bottom_right"}
		if !slices.Contains(locations, v.(string)) {
			return g.newErrorTk(fmt.Sprintf("Unknown anchor location. Available options: %s", strings.Join(locations, ", ")), token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteMove(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteMoveXY, false)

		var err error
		block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Fields["fieldMenu_2"] = []any{direction, nil}

		return block, nil
	}
}

func funcSpriteMoveTo(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteMoveTo, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[2])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteMoveRandom(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteMoveRandom, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcSpriteRotate(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteRotate, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteRotateTo(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteRotateTo, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteSetScale(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteSetSize, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcSpriteSetColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteSetColorWithColor, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	if len(stmt.Parameters) == 2 {
		block.Inputs["number_2"], err = g.valueWithRegex(block.ID, stmt.Parameters[1], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		block.Type = blocks.SpriteSetColorWithRGB
		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_4"], err = g.valueInRange(block.ID, stmt.Parameters[3], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcSpriteResetColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SpriteCloseColor, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcSpriteShowHide(showHide string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteShowAndHide, false)

		var err error
		block.Inputs["inputVariable_2"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["string_1"] = []any{showHide, nil}

		return block, nil
	}
}

func funcSpriteSetLayer(layer string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteZMinMax, false)

		var err error
		block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_2"] = []any{layer, nil}

		return block, nil
	}
}

func funcSpriteChangeLayer(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.SpriteZUpDown, false)

		var err error
		block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_2"] = []any{direction, nil}

		return block, nil
	}
}

func funcDrawBegin(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchStart, false)
	return block, nil
}

func funcDrawFinish(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchEnd, false)
	return block, nil
}

func funcDrawClear(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchClear, false)
	return block, nil
}

func funcDrawSetColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchSetColorWithColor, false)

	var err error
	if len(stmt.Parameters) == 1 {
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Parameters[0], hexColorRegex, 9, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			g.errors = append(g.errors, err)
		}
	} else {
		block.Type = blocks.DrawSketchSetColorWithRGB
		block.Inputs["number_1"], err = g.valueInRange(block.ID, stmt.Parameters[0], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcDrawSetThickness(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchSetSize, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDrawSetSpeed(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchSetSpeed, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDrawRotate(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchCW, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDrawRotateTo(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchSetAngle, false)

	var err error
	block.Inputs["angle_1"], err = g.valueWithValidator(block.ID, stmt.Parameters[0], func(v any) bool { return true }, 8, "")
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDrawLine(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchMove, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDrawCircle(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchCircle, false)

	var err error
	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcDrawMove(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.DrawSketchMoveXAndY, false)

		var err error
		block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		block.Fields["fieldMenu_1"] = []any{direction, nil}

		return block, nil
	}
}

func funcDrawMoveTo(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchMoveTo, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["number_2"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcDrawMoveToCenter(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchMoveToCenter, false)
	return block, nil
}

func funcDrawSave(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.DrawSketchSpriteDrawSketch, false)

	var err error
	block.Inputs["string_1"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcLEDMove(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.LEDMove, false)

	var err error
	block.Inputs["led_number"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcNetBroadcast(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetSetWifiBroadcast, false)

	var err error
	block.Inputs["message"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	if len(stmt.Parameters) == 2 {
		block.Type = blocks.NetSetWifiBroadcastWithValue
		block.Inputs["value"], err = g.value(block.ID, stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		}
	}

	return block, nil
}

func funcNetSetChannel(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetSetWifiChannel, false)

	channel, err := g.literal(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	if int(channel.(float64)) != 1 && int(channel.(float64)) != 6 && int(channel.(float64)) != 11 {
		return nil, g.newErrorExpr("Invalid channel. Allowed options: 1, 6, 11", stmt.Parameters[0])
	}
	block.Fields["channel"] = []any{fmt.Sprintf("%d", int(channel.(float64))), nil}

	return block, nil
}

func funcNetConnect(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetConnectWifi, false)

	var err error
	block.Inputs["ssid"], err = g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["wifipassword"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}

	return block, nil
}

func funcNetReconnect(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetWifiReconnect, false)
	return block, nil
}

func funcNetDisconnect(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.NetWifiDisconnect, false)
	return block, nil
}

func funcSensorsResetAngle(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorsResetAxisRotationAngle, false)
	value, err := g.literal(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	axes := []string{"all", "x", "y", "z"}
	if !slices.Contains(axes, value.(string)) {
		return nil, g.newErrorExpr(fmt.Sprintf("Unknown axis. Available options: %s", strings.Join(axes, ", ")), stmt.Parameters[0])
	}
	block.Fields["axis"] = []any{value, nil}
	return block, nil
}

func funcSensorsResetYawAngle(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorsResetYaw, false)
	return block, nil
}

func funcSensorsDefineColor(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorDefineColor, false)

	var err error
	block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Parameters[0], -1, 0, 255)
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Parameters[1], -1, 0, 255)
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Parameters[2], -1, 0, 255)
	if err != nil {
		g.errors = append(g.errors, err)
	}

	if len(stmt.Parameters) == 4 {
		block.Inputs["tolerance"], err = g.value(block.ID, stmt.Parameters[3])
		if err != nil {
			g.errors = append(g.errors, err)
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

func funcSensorsCalibrateColors(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorCalibrate, false)
	block.Fields["index"] = []any{"1", nil}
	return block, nil
}

func funcSensorsEnhancedColorDetection(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.SensorColorDetectionMode, false)

	enable, err := g.literal(stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	if enable.(bool) {
		block.Fields["mode"] = []any{"enhance", nil}
	} else {
		block.Fields["mode"] = []any{"standard", nil}
	}

	return block, nil
}

func funcMotorsRun(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2MoveDirectionWithRPM, false)

		block.Fields["DIRECTION"] = []any{direction, nil}

		var err error
		block.Inputs["POWER"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		if len(stmt.Parameters) == 2 {
			block.Type = blocks.Mbot2MoveDirectionWithTime
			block.Inputs["TIME"], err = g.value(block.ID, stmt.Parameters[1])
			if err != nil {
				g.errors = append(g.errors, err)
			}
		}

		return block, nil
	}
}

func funcMotorsRunDistance(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2MoveMoveWithCmAndInch, false)

		block.Fields["DIRECTION"] = []any{direction, nil}
		block.Fields["fieldMenu_3"] = []any{"cm", nil}

		var err error
		block.Inputs["POWER"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}
		return block, nil
	}
}

func funcMotorsTurn(direction string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2CwAndCcwWithAngle, false)
		block.Fields["fieldMenu_1"] = []any{direction, nil}

		var err error
		block.Inputs["ANGLE"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcMotorsRotate(unit string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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
		block.Inputs[inputField], err = g.fieldMenu(menuType, "", "MBOT2_ENCODER_MOTOR_SET_WITH_TIME_FIELDMENU_1", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
			str := v.(string)
			if str != "ALL" && str != "EM1" && str != "EM2" {
				return g.newErrorTk("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Fields["fieldMenu_4"] = []any{unit, nil}

		block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		if len(stmt.Parameters) == 3 {
			block.Inputs["number_3"], err = g.value(block.ID, stmt.Parameters[2])
			if err != nil {
				g.errors = append(g.errors, err)
			}
		}

		return block, nil
	}
}

func funcMotorsRotateAngle(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorSetWithTimeAngleAndCircle, false)

	var err error
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorSetWithTimeAngleAndCircleMenu, "", "MBOT2_ENCODER_MOTOR_SET_WITH_TIME_FIELDMENU_1", block.ID, stmt.Parameters[0], func(v any, token parser.Token) error {
		str := v.(string)
		if str != "ALL" && str != "EM1" && str != "EM2" {
			return g.newErrorTk("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		g.errors = append(g.errors, err)
	}

	block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Parameters[1])
	return block, err
}

func funcMotorsDrive(unit string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2EncoderMotorDrivePower, false)
		rightPowerKey := "number_2"
		if unit == "speed" {
			block.Type = blocks.Mbot2EncoderMotorDriveSpeed
			rightPowerKey = "RIGHT_POWER"
		}

		var err error
		block.Inputs["LEFT_POWER"], err = g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		block.Inputs[rightPowerKey], err = g.value(block.ID, stmt.Parameters[1])
		if err != nil {
			g.errors = append(g.errors, err)
		}

		return block, nil
	}
}

func funcMotorsStop(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2EncoderMotorStop, false)

	encoderMotor := "ALL"
	if len(stmt.Parameters) == 1 {
		motor, err := g.literal(stmt.Parameters[0])
		if err != nil {
			return nil, err
		}
		encoderMotor = motor.(string)
		if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
			return nil, g.newErrorExpr("Unknown encoder motor. Available options: ALL, EM1, EM2", stmt.Parameters[0])
		}
	}

	block.Fields["fieldMenu_1"] = []any{encoderMotor, nil}

	return block, nil
}

func funcMotorsResetAngle(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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
			},
		}
	}

	var err error
	block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorResetAngleMenu, "", "MBOT2_ENCODER_MOTOR_STOP_FIELDMENU_1", block.ID, motor, func(v any, token parser.Token) error {
		encoderMotor := v.(string)
		if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
			return g.newErrorTk("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcMotorsSetLock(value string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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
				},
			}
		}

		var err error
		block.Inputs["inputMenu_1"], err = g.fieldMenu(blocks.Mbot2EncoderMotorLockUnlockMenu, "", "MBOT2_ENCODER_MOTOR_STOP_FIELDMENU_1", block.ID, motor, func(v any, token parser.Token) error {
			encoderMotor := v.(string)
			if encoderMotor != "ALL" && encoderMotor != "EM1" && encoderMotor != "EM2" {
				return g.newErrorTk("Unknown encoder motor. Available options: ALL, EM1, EM2", token)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return block, nil
	}
}

func funcTimeWait(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ControlWaitUntil, false)

	condition, err := g.value(block.ID, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}

	if stmt.Parameters[0].Type() == parser.DTBool {
		block.Inputs["CONDITION"] = condition
		return block, nil
	} else {
		block.Type = blocks.ControlWait
		seconds, err := g.value(block.ID, stmt.Parameters[0])
		if err != nil {
			return nil, err
		}
		block.Inputs["DURATION"] = seconds
	}

	return block, nil
}

func funcResetTimer(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.Mbot2TimerReset, false)
	return block, nil
}

func funcMBotRestart(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ControlRestart, false)
	return block, nil
}

func funcMBotChassisParameters(parameter string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
		block := g.NewBlock(blocks.Mbot2SetParameters, false)

		if parameter == "calibrate" {
			parameter = "set_auto"
		}

		block.Fields["PARA"] = []any{parameter, nil}

		return block, nil
	}
}

func funcScriptStop(stopOption string) func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
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

func funcListsAppend(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListAdd, false)
	err := selectList(g, block, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	return block, nil
}

func funcListsRemove(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListDelete, false)
	err := selectList(g, block, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	return block, nil
}

func funcListsClear(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListClear, false)
	err := selectList(g, block, stmt.Parameters[0])
	if err != nil {
		return nil, err
	}
	return block, nil
}

func funcListsInsert(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListInsert, false)
	err := selectList(g, block, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Parameters[2])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	return block, nil
}

func funcListsReplace(g *generator, stmt *parser.StmtCall) (*blocks.Block, error) {
	block := g.NewBlock(blocks.ListReplace, false)
	err := selectList(g, block, stmt.Parameters[0])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["INDEX"], err = g.value(block.ID, stmt.Parameters[1])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	block.Inputs["ITEM"], err = g.value(block.ID, stmt.Parameters[2])
	if err != nil {
		g.errors = append(g.errors, err)
	}
	return block, nil
}

func selectList(g *generator, block *blocks.Block, param parser.Expr) error {
	if i, ok := param.(*parser.ExprIdentifier); ok {
		if l, ok := g.definitions.Lists[i.Name.Lexeme]; ok {
			block.Fields["LIST"] = []any{
				l.Name.Lexeme,
				l.ID,
			}
			return nil
		}
		return g.newErrorExpr("Unknown list.", param)
	}
	return g.newErrorExpr("Expected list.", param)
}
