package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var funcCalls = map[string]func(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error){
	"audio.stop":       funcAudioStop,
	"audio.playBuzzer": funcAudioPlayBuzzer,

	"time.sleep": funcTimeSleep,

	"mbot.restart": funcMbotRestart,

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
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'audio.playBuzzer' function takes 1 argument: audio.playBuzzer(frequency: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayBuzzerTone, parent)
	g.blocks[parent].Next = &block.ID

	number, err := g.value(block.ID, stmt.Name, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	block.Inputs["number_2"] = number

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

func funcMbotRestart(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
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
