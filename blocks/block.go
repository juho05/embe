package blocks

import "github.com/google/uuid"

type Block struct {
	ID       string         `json:"-"`
	Type     BlockType      `json:"opcode"`
	Next     *string        `json:"next"`
	Parent   *string        `json:"parent"`
	Inputs   map[string]any `json:"inputs"`
	Fields   map[string]any `json:"fields"`
	Shadow   bool           `json:"shadow"`
	TopLevel bool           `json:"topLevel"`
	X        int            `json:"x"`
	Y        int            `json:"y"`
}

func NewBlock(blockType BlockType, parent *string) Block {
	return Block{
		ID:     uuid.NewString(),
		Type:   blockType,
		Parent: parent,
		Inputs: make(map[string]any),
		Fields: make(map[string]any),
		Y:      40,
	}
}

var topLevelX = -400

func NewBlockTopLevel(blockType BlockType) Block {
	topLevelX += 430
	return Block{
		ID:       uuid.NewString(),
		Type:     blockType,
		Inputs:   make(map[string]any),
		Fields:   make(map[string]any),
		TopLevel: true,
		X:        topLevelX,
		Y:        40,
	}
}
