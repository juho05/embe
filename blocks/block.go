package blocks

type Block struct {
	Type     BlockType      `json:"opcode"`
	Next     *string        `json:"next"`
	Parent   *string        `json:"parent"`
	Inputs   map[string]any `json:"inputs"`
	Fields   map[string]any `json:"fields"`
	Shadow   bool           `json:"shadow"`
	TopLevel bool           `json:"topLevel"`
}
