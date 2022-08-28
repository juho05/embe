package generator

import (
	"errors"
	"fmt"
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

	"time.sleep": funcTimeSleep,

	"mbot.restart": funcMBotRestart,

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
	block.Inputs["file_name"], err = g.fieldMenu(menuBlockType, "", "CYBERPI_PLAY_AUDIO_UNTIL_3_FILE_NAME", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v parser.Token) error {
		names := []string{"hi", "bye", "yeah", "wow", "laugh", "hum", "sad", "sigh", "annoyed", "angry", "surprised", "yummy", "curious", "embarrassed", "ready", "sprint", "sleepy", "meow", "start", "switch", "beeps", "buzzing", "jump", "level-up", "low-energy", "prompt", "right", "wrong", "ring", "score", "wake", "warning", "metal-clash", "glass-clink", "inflator", "running-water", "clockwork", "click", "current", "wood-hit", "iron", "drop", "bubble", "wave", "magic", "spitfire", "heartbeat"}
		if !slices.Contains(names, v.Literal.(string)) {
			return g.newError(fmt.Sprintf("Unknown clip name. Available options: %s", strings.Join(names, ", ")), v)
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
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(blocks.AudioPlayMusicInstrumentMenu, "'", "CYBERPI_PLAY_MUSIC_WITH_NOTE_FIELDMENU_1", block.ID, stmt.Name, stmt.Parameters[0], parser.DTString, func(v parser.Token) error {
		if !slices.Contains(names, v.Literal.(string)) {
			return g.newError(fmt.Sprintf("Unknown instrument name. Available options: %s", strings.Join(names, ", ")), v)
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
	block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Name, stmt.Parameters[0], "", func(v parser.Token) error {
		if str, ok := v.Literal.(string); !ok {
			return errWrongType
		} else {
			if str != "all" {
				return g.newError("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", v)
			}
		}
		return nil
	})
	if err == errWrongType {
		block.Inputs["fieldMenu_1"], err = g.fieldMenu(menuBlockType, "\"", menuFieldKey, block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber, func(v parser.Token) error {
			nr := int(v.Literal.(float64))
			if nr != 1 && nr != 2 && nr != 3 && nr != 4 && nr != 5 {
				return g.newError("Unknown LED. Available options: \"all\", 1, 2, 3, 4, 5", v)
			}
			return nil
		})
	}
	return err
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

func funcProgramExit(g *generator, stmt *parser.StmtFuncCall) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'program.exit' function does not take any arguments.", stmt.Name)
	}
	block := g.NewBlock(blocks.ControlStop, false)
	return block, nil
}
