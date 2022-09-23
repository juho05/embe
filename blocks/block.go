package blocks

import "github.com/google/uuid"

type Block struct {
	ID       string         `json:"-"`
	NoNext   bool           `json:"-"`
	Type     BlockType      `json:"opcode"`
	Next     *string        `json:"next"`
	Parent   *string        `json:"parent"`
	Inputs   map[string]any `json:"inputs"`
	Fields   map[string]any `json:"fields"`
	Mutation map[string]any `json:"mutation,omitempty"`
	Shadow   bool           `json:"shadow"`
	TopLevel bool           `json:"topLevel"`
	X        int            `json:"x,omitempty"`
	Y        int            `json:"y,omitempty"`
}

var topLevelX = -520

func NewStage() {
	topLevelX = -520
}

func NewBlock(blockType BlockType, parent string) *Block {
	var p *string
	p = &parent
	if parent == "" {
		p = nil
	}
	return &Block{
		ID:     uuid.NewString(),
		Type:   blockType,
		Parent: p,
		Inputs: make(map[string]any),
		Fields: make(map[string]any),
	}
}

func NewShadowBlock(blockType BlockType, parent string) *Block {
	var p *string
	p = &parent
	if parent == "" {
		p = nil
	}
	return &Block{
		ID:     uuid.NewString(),
		Type:   blockType,
		Parent: p,
		Inputs: make(map[string]any),
		Fields: make(map[string]any),
		Shadow: true,
	}
}

func NewBlockTopLevel(blockType BlockType) *Block {
	topLevelX += 550
	return &Block{
		ID:       uuid.NewString(),
		Type:     blockType,
		Inputs:   make(map[string]any),
		Fields:   make(map[string]any),
		TopLevel: true,
		X:        topLevelX,
		Y:        80,
	}
}
