package blocks

import "encoding/json"

type Project struct {
	Targets []map[string]json.RawMessage `json:"targets"`
}
