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

var FuncCalls = map[string]func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error){
	"audio.stop":            funcAudioStop,
	"audio.playBuzzer":      funcAudioPlayBuzzer,
	"audio.playNote":        funcAudioPlayNote,
	"audio.playInstrument":  funcAudioPlayInstrument,
	"audio.playClip":        funcAudioPlayClip,
	"audio.recording.start": funcAudioRecordingStart,
	"audio.recording.stop":  funcAudioRecordingStop,
	"audio.recording.play":  funcAudioRecordingPlay,

	"led.display":         funcLEDDisplay,
	"led.displayColor":    funcLEDDisplayColor,
	"led.displayColorFor": funcLEDDisplayColorFor,
	"led.deactivate":      funcLEDDeactivate,
	"led.move":            funcLEDMove,
	"led.playAnimation":   funcLEDPlayAnimation,

	"display.print":                 funcDisplayPrint(false),
	"display.println":               funcDisplayPrint(true),
	"display.setFontSize":           funcDisplaySetFontSize,
	"display.setColor":              funcDisplaySetColor,
	"display.showLabel":             funcDisplayShowLabel,
	"display.lineChart.addData":     funcDisplayLineChartAddData,
	"display.lineChart.setInterval": funcDisplayLineChartSetInterval,
	"display.barChart.addData":      funcDisplayBarChartAddData,
	"display.table.addData":         funcDisplayTableAddData,
	"display.clear":                 funcDisplayClear,
	"display.setOrientation":        funcDisplaySetOrientation,

	"net.broadcast":  funcNetBroadcast,
	"net.setChannel": funcNetSetChannel,
	"net.connect":    funcNetConnect,
	"net.reconnect":  funcNetReconnect,
	"net.disconnect": funcNetDisconnect,

	"motors.run":                 funcMotorsRun("forward"),
	"motors.runBackward":         funcMotorsRun("backward"),
	"motors.runDistance":         funcMotorsRunDistance("forward"),
	"motors.runDistanceBackward": funcMotorsRunDistance("backward"),
	"motors.turn":                funcMotorsTurn,
	"motors.rotateRPM":           funcMotorsRotate("speed"),
	"motors.rotatePower":         funcMotorsRotate("power"),
	"motors.rotateAngle":         funcMotorsRotateAngle,
	"motors.stop":                funcMotorsStop,
	"motors.resetAngle":          funcMotorsResetAngle,
	"motors.lock":                funcMotorsSetLock("1"),
	"motors.unlock":              funcMotorsSetLock("0"),

	"time.sleep": funcTimeSleep,

	"mbot.restart":             funcMBotRestart,
	"mbot.calibrateParameters": funcMBotChassisParameters("calibrate"),
	"mbot.resetParameters":     funcMBotChassisParameters("reset"),

	"program.exit": funcProgramExit,
}

func funcAudioStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.stop' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.AudioStop, false)
	return block, nil
}

func funcAudioPlayBuzzer(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playBuzzer' function takes 1-2 arguments: audio.playBuzzer(frequency: number, duration?: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playClip' function takes 1-2 arguments: audio.playClip(name: string, untilDone?: boolean)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playInstrument' function takes 2 arguments: audio.playInstrument(name: string, duration: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.AudioPlayMusicInstrument, false)

	var err error
	names := []string{"snare", "bass-drum", "side-stick", "crash-cymbal", "open-hi-hat", "closed-hi-hat", "tambourine", "hand-clap", "claves"}
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.AudioPlayMusicInstrumentMenu, "'", "CYBERPI_PLAY_MUSIC_WITH_NOTE_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v any, token parser.Token) error {
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
	if len(stmt.Parameters) != 2 && len(stmt.Parameters) != 3 {
		return nil, g.newError("The 'audio.playNote' function takes 2-3 arguments: audio.playNote(note: number, duration: number) audio.playNote(name: string, octave: number, duration: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.record.start' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.AudioRecordStart, false)
	return block, nil
}

func funcAudioRecordingStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.record.stop' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.AudioRecordStop, false)
	return block, nil
}

func funcAudioRecordingPlay(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) > 1 {
		return nil, g.newError("The 'audio.record.play' function takes 0-1 arguments: audio.record.play(untilDone?: boolean)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'led.playAnimation' function takes 1 argument: led.playAnimation(name: string)", stmt.Name)
	}
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

func funcLEDDisplay(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 5 {
		return nil, g.newError("The 'led.display' function takes 5 arguments: led.display(color1: string, color2: string, color3: string, color4: string, color5: string)", stmt.Name)
	}
	block := g.NewBlock(blocks.LEDDisplay, false)

	names := make([]string, 5)
	for i := range names {
		n, err := g.literal(stmt.Name, stmt.Parameters[i], parser.DTString)
		if err != nil {
			return nil, err
		}
		colors := []string{"gray", "red", "orange", "yellow", "green", "cyan", "blue", "magenta"}
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
	if len(stmt.Parameters) < 1 || len(stmt.Parameters) > 4 {
		return nil, g.newError("The 'led.displayColor' function takes 1-4 arguments: led.displayColor(led?: number, color: string) or led.displayColor(led?: number, r: number, g: number, b: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.LEDDisplaySingleColor, false)
	if len(stmt.Parameters) > 2 {
		block.Type = blocks.LEDDisplaySingleColorWithRGB
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBFieldMenu, stmt, 3, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorFieldMenu, stmt, 1, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString, 9, hexColorRegex, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcLEDDisplayColorFor(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) < 2 || len(stmt.Parameters) > 5 {
		return nil, g.newError("The 'led.displayColorFor' function takes 2-5 arguments: led.displayColorFor(led?: number, color: string, duration: number) or led.displayColor(led?: number, r: number, g: number, b: number, duration: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.LEDDisplaySingleColorWithTime, false)
	if len(stmt.Parameters) > 3 {
		block.Type = blocks.LEDDisplaySingleColorWithRGBAndTime
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithRGBAndTimeFieldMenu, stmt, 4, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}

		block.Inputs["r"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["g"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["b"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[3], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_5"], err = g.value(block.ID, stmt.Name, stmt.Parameters[4], parser.DTNumber)
		if err != nil {
			return nil, err
		}
	} else {
		err := selectLED(g, block, blocks.LEDDisplaySingleColorWithTime, stmt, 2, "CYBERPI_LED_SHOW_SINGLE_WITH_COLOR_AND_TIME_2_FIELDMENU_1")
		if err != nil {
			return nil, err
		}
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[1], parser.DTString, 9, hexColorRegex, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
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
	if len(stmt.Parameters) > 1 {
		return nil, g.newError("The 'led.deactivate' function takes 0-1 arguments: led.deactivate(led?: number)", stmt.Name)
	}
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

func funcDisplayPrint(newLine bool) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		if len(stmt.Parameters) != 1 {
			return nil, g.newError(fmt.Sprintf("The '%s' function takes 1 argument: %s(text: string)", stmt.Name.Lexeme, stmt.Name.Lexeme), stmt.Name)
		}
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
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'display.setFontSize' function takes 1 argument: display.setFontSize(size: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 3 {
		return nil, g.newError("The 'display.setColor' function takes 1 or 3 arguments: display.setColor(color: string) or display.setColor(r: number, g: number, b: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.DisplaySetBrushColor, false)
	var err error
	if len(stmt.Parameters) == 3 {
		block.Type = blocks.DisplaySetBrushColorRGB
		block.Inputs["number_1"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_2"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
		block.Inputs["number_3"], err = g.valueInRange(block.ID, stmt.Name, stmt.Parameters[2], parser.DTNumber, 4, 0, 255)
		if err != nil {
			return nil, err
		}
	} else {
		block.Inputs["color_1"], err = g.valueWithRegex(block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, 9, hexColorRegex, "The value must be a valid hex color (\"#000000\" - \"#ffffff\").")
		if err != nil {
			return nil, err
		}
	}

	return block, nil
}

func funcDisplayShowLabel(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 4 && len(stmt.Parameters) != 5 {
		return nil, g.newError("The 'display.showLabel' function takes 3 or 4 arguments: display.showLabel(label: number, text: string, location: string, size: number) or display.showLabel(label: number, text: string, x: number, y: number, size: number)", stmt.Name)
	}

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
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'display.lineChart.addData' function takes 1 argument: display.lineChart.addData(value: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.DisplayLineChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayLineChartSetInterval(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'display.lineChart.setInterval' function takes 1 argument: display.lineChart.setInterval(interval: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.DisplayLineChartSetInterval, false)

	var err error
	block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayBarChartAddData(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'display.barChart.addData' function takes 1 argument: display.barChart.addData(value: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.DisplayBarChartAddData, false)

	var err error
	block.Inputs["number_1"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcDisplayTableAddData(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 3 {
		return nil, g.newError("The 'display.table.addData' function takes 3 arguments: display.table.addData(text: string, row: number, column: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'display.setOrientation' function takes 1 argument: display.setOrientation(orientation: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'display.clear' function takes no arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.DisplayClear, false)
	return block, nil
}

func funcLEDMove(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'led.move' function takes 1 argument: led.scroll(amount: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.LEDMove, false)

	var err error
	block.Inputs["led_number"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcNetBroadcast(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'net.broadcast' function takes 1-2 arguments: net.broadcast(message: string, value?: string)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'net.setChannel' function takes 1 argument: net.setChannel(channel: number)", stmt.Name)
	}
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
	if len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'net.connect' function takes 2 arguments: net.connect(ssid: string, password: string)", stmt.Name)
	}
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
	if len(stmt.Parameters) > 0 {
		return nil, g.newError("The 'net.reconnect' function takes no arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.NetWifiReconnect, false)
	return block, nil
}

func funcNetDisconnect(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) > 0 {
		return nil, g.newError("The 'net.disconnect' function takes no arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.NetWifiDisconnect, false)
	return block, nil
}

func funcMotorsRun(direction string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
			return nil, g.newError("The 'motors.run' function takes 1-2 arguments: motors.run(rpm: number, duration: number)", stmt.Name)
		}
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
		if len(stmt.Parameters) != 1 {
			return nil, g.newError("The 'motors.runDistance' function takes 1 argument: motors.runDistance(distance: number)", stmt.Name)
		}
		block := g.NewBlock(blocks.Mbot2MoveDirectionWithRPM, false)

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

func funcMotorsTurn(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'motors.turn' function takes 1 argument: motors.turn(angle: number)", stmt.Name)
	}
	block := g.NewBlock(blocks.Mbot2CwAndCcwWithAngle, false)
	block.Fields["fieldMenu_1"] = []any{"ccw", nil}

	var err error
	block.Inputs["ANGLE"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcMotorsRotate(unit string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		if len(stmt.Parameters) != 2 && len(stmt.Parameters) != 3 {
			return nil, g.newError("The 'motors.rotate' function takes 1 argument: motors.rotate(motor: string, rpm: number, duration?: number)", stmt.Name)
		}
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
	if len(stmt.Parameters) != 2 {
		return nil, g.newError("The `motors.rotateAngle` function takes 1-2 arguments: motors.rotateAngle(motor: string, angle: number)", stmt.Name)
	}
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

func funcMotorsStop(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) > 1 {
		return nil, g.newError("The `motors.stop` function takes 0-1 arguments: motors.stop(motor?: string)", stmt.Name)
	}
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
	if len(stmt.Parameters) > 1 {
		return nil, g.newError("The `motors.resetAngle` function takes 0-1 arguments: motors.resetAngle(motor?: string)", stmt.Name)
	}
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
		if len(stmt.Parameters) > 1 {
			return nil, g.newError("The `motors.lock` function takes 0-1 arguments: motors.lock(motor?: string)", stmt.Name)
		}
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

func funcTimeSleep(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'time.sleep' function takes 1 argument: time.sleep(seconds: number) or time.sleep(continueCondition: boolean)", stmt.Name)
	}
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

func funcMBotRestart(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'mbot.restart' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.ControlRestart, false)
	return block, nil
}

func funcMBotChassisParameters(parameter string) func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	return func(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
		if len(stmt.Parameters) > 0 {
			return nil, g.newError(fmt.Sprintf("The `mbot.%sParameters` function takes no arguments.", parameter), stmt.Name)
		}
		block := g.NewBlock(blocks.Mbot2SetParameters, false)

		if parameter == "calibrate" {
			parameter = "set_auto"
		}

		block.Fields["PARA"] = []any{parameter, nil}

		return block, nil
	}
}

func funcProgramExit(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'program.exit' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.ControlStop, false)
	return block, nil
}
