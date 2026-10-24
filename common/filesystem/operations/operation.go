package operations

import "encoding/json"

type Operation struct {
	Action  string          `json:"action"`
	Request json.RawMessage `json:"request"`
}
