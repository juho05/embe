package blocks

type Block struct {
	ID       string         `json:"-"`
	Type     BlockType      `json:"opcode"`
	Next     *string        `json:"next"`
	Parent   *string        `json:"parent"`
	Inputs   map[string]any `json:"inputs"`
	Fields   map[string]any `json:"fields"`
	Shadow   bool           `json:"shadow"`
	TopLevel bool           `json:"topLevel"`
}

func NewBlock(id string, blockType BlockType, topLevel bool, parent *string) Block {
	return Block{
		ID:       id,
		Type:     blockType,
		Inputs:   make(map[string]any),
		Fields:   make(map[string]any),
		TopLevel: topLevel,
	}
}
