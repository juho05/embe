package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var funcCalls = map[string]func(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error){
	"audio.stop":       funcAudioStop,
	"audio.playBuzzer": funcAudioPlayBuzzer,
}

func funcAudioStop(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 0 {
		return nil, g.newError("The 'audio.stop' function does not take any arguments.", stmt.Name)
	}
	block := blocks.NewBlock(blocks.StopAudio, parent)
	if parent != "" {
		g.blocks[parent].Next = &block.ID
	}
	return block, nil
}

func funcAudioPlayBuzzer(g *generator, stmt *parser.StmtFuncCall, parent string) (*blocks.Block, error) {
	if len(stmt.Parameters) != 1 {
		return nil, g.newError("The 'audio.playBuzzer' function takes 1 argument: audio.playBuzzer(frequency: number)", stmt.Name)
	}
	block := blocks.NewBlock(blocks.PlayBuzzerTone, parent)
	if parent != "" {
		g.blocks[parent].Next = &block.ID
	}

	number, err := g.value(block.ID, stmt.Parameters[0], parser.DTNumber)
	if err != nil {
		return nil, err
	}
	block.Inputs["number_2"] = number

	return block, nil
}
