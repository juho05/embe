package generator

import (
	"github.com/Bananenpro/embe/blocks"
	"github.com/Bananenpro/embe/parser"
)

var events = map[string]func(g *generator, stmt *parser.StmtEvent) (blocks.Block, error){
	"start": eventStart,
}

func eventStart(g *generator, stmt *parser.StmtEvent) (blocks.Block, error) {
	if (stmt.Parameter != parser.Token{}) {
		return blocks.Block{}, g.newError("The 'start' event does not take any arguments.", stmt.Parameter)
	}
	block := blocks.NewBlockTopLevel(blocks.WhenLaunch)
	return block, nil
}
