package generator

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var funcCalls = map[string]func(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error){
	"audio.stop":            funcAudioStop,
	"audio.playBuzzer":      funcAudioPlayBuzzer,
	"audio.playNote":        funcAudioPlayNote,
	"audio.playInstrument":  funcAudioPlayInstrument,
	"audio.playClip":        funcAudioPlayClip,
	"audio.recording.start": funcAudioRecordingStart,
	"audio.recording.stop":  funcAudioRecordingStop,
	"audio.recording.play":  funcAudioRecordingPlay,

	"led.display":       funcLEDDisplay,
	"led.move":          funcLEDMove,
	"led.playAnimation": funcLEDPlayAnimation,

	"time.sleep": funcTimeSleep,

	"mbot.restart": funcMBotRestart,

	"program.exit": funcProgramExit,
}

func funcAudioStop(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.stop' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.StopAudio, parent)
	g.blocks[parent].Next = &block.ID
	return block, nil
}

func funcAudioPlayBuzzer(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playBuzzer' function takes 1-2 arguments: audio.playBuzzer(frequency: number, duration?: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayBuzzerTone, parent)
	g.blocks[parent].Next = &block.ID

	number, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	if len(stmt.Parameters) == 2 {
		block.Type = blocks.PlayBuzzerToneWithTime
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

func funcAudioPlayClip(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 && len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playClip' function takes 1-2 arguments: audio.playClip(name: string, untilDone?: boolean)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayClip, parent)
	g.blocks[parent].Next = &block.ID

	menuBlock := blocks.NewShadowBlock(blocks.PlayClipFileNameMenu, block.ID)
	g.blocks[menuBlock.ID] = menuBlock
	block.Inputs["file_name"] = []any{1, menuBlock.ID}
	if len(stmt.Parameters) == 2 {
		untilDone, err := g.literal(stmt.Name, stmt.Parameters[1], parser.DTBool)
		if err != nil {
			return nil, err
		}
		if untilDone.(bool) {
			block.Type = blocks.PlayClipUntilDone
			menuBlock.Type = blocks.PlayClipUntilDoneFileNameMenu
		}
	}

	fileName, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	names := []string{"hi", "bye", "yeah", "wow", "laugh", "hum", "sad", "sigh", "annoyed", "angry", "surprised", "yummy", "curious", "embarrassed", "ready", "sprint", "sleepy", "meow", "start", "switch", "beeps", "buzzing", "jump", "level-up", "low-energy", "prompt", "right", "wrong", "ring", "score", "wake", "warning", "metal-clash", "glass-clink", "inflator", "running-water", "clockwork", "click", "current", "wood-hit", "iron", "drop", "bubble", "wave", "magic", "spitfire", "heartbeat"}
	if !slices.Contains(names, fileName.(string)) {
		return nil, g.newError(fmt.Sprintf("Unknown clip name. Available options: %s", strings.Join(names, ", ")), stmt.Parameters[0].(*parser.ExprLiteral).Token)
	}

	menuBlock.Fields["CYBERPI_PLAY_AUDIO_UNTIL_3_FILE_NAME"] = []any{fileName, nil}

	return block, nil
}

func funcAudioPlayInstrument(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 2 {
		return nil, g.newError("The 'audio.playInstrument' function takes 2 arguments: audio.playInstrument(name: string, duration: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayMusicInstrument, parent)
	g.blocks[parent].Next = &block.ID

	menuBlock := blocks.NewShadowBlock(blocks.PlayMusicInstrumentMenu, block.ID)
	g.blocks[menuBlock.ID] = menuBlock
	block.Inputs["fieldMenu_1"] = []any{1, menuBlock.ID}

	instrumentName, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTString)
	if err != nil {
		return nil, err
	}

	names := []string{"snare", "bass-drum", "side-stick", "crash-cymbal", "open-hi-hat", "closed-hi-hat", "tambourine", "hand-clap", "claves"}
	if !slices.Contains(names, instrumentName.(string)) {
		return nil, g.newError(fmt.Sprintf("Unknown instrument name. Available options: %s", strings.Join(names, ", ")), stmt.Parameters[0].(*parser.ExprLiteral).Token)
	}

	menuBlock.Fields["CYBERPI_PLAY_MUSIC_WITH_NOTE_FIELDMENU_1"] = []any{fmt.Sprintf("'%v'", instrumentName), nil}

	block.Inputs["number_3"], err = g.value(block.ID, stmt.Name, stmt.Parameters[1], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcAudioPlayNote(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 2 && len(stmt.Parameters) != 3 {
		return nil, g.newError("The 'audio.playNote' function takes 2-3 arguments: audio.playNote(note: number, duration: number) audio.playNote(name: string, octave: number, duration: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayNote, parent)
	g.blocks[parent].Next = &block.ID

	noteBlock := blocks.NewShadowBlock(blocks.Note, block.ID)
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

func funcAudioRecordingStart(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.record.start' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.RecordStart, parent)
	g.blocks[parent].Next = &block.ID
	return block, nil
}

func funcAudioRecordingStop(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.record.stop' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.RecordStop, parent)
	g.blocks[parent].Next = &block.ID
	return block, nil
}

func funcAudioRecordingPlay(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) > 1 {
		return nil, g.newError("The 'audio.record.play' function takes 0-1 arguments: audio.record.play(untilDone?: boolean)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayRecord, parent)
	g.blocks[parent].Next = &block.ID

	if len(stmt.Parameters) == 1 {
		untilDone, err := g.literal(stmt.Name, stmt.Parameters[0], parser.DTBool)
		if err != nil {
			return nil, err
		}
		if untilDone.(bool) {
			block.Type = blocks.PlayRecordUntilDone
		}
	}

	return block, nil
}

func funcLEDPlayAnimation(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'led.playAnimation' function takes 1 argument: led.playAnimation(name: string)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayLEDAnimation, parent)
	g.blocks[parent].Next = &block.ID

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

func funcLEDDisplay(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 5 {
		return nil, g.newError("The 'led.display' function takes 5 arguments: led.display(color1: string, color2: string, color3: string, color4: string, color5: string)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.LEDDisplay, parent)
	g.blocks[parent].Next = &block.ID

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

func funcLEDMove(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'led.move' function takes 1 argument: led.scroll(amount: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.LEDMove, parent)
	g.blocks[parent].Next = &block.ID

	var err error
	block.Inputs["led_number"], err = g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func funcTimeSleep(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'time.sleep' function takes 1 argument: time.sleep(seconds: number) or time.sleep(continueCondition: boolean)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.WaitUntil, parent)
	g.blocks[parent].Next = &block.ID

	condition, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTBool)
	if err == nil {
		block.Inputs["CONDITION"] = condition
		return block, nil
	} else {
		block.Type = blocks.Wait
		seconds, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
		if err != nil {
			return nil, err
		}
		block.Inputs["DURATION"] = seconds
	}

	return block, nil
}

func funcMBotRestart(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'mbot.restart' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.Restart, parent)
	g.blocks[parent].Next = &block.ID
	return block, nil
}

func funcProgramExit(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'program.exit' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.Stop, parent)
	g.blocks[parent].Next = &block.ID
	return block, nil
}
