package dispatcher

import (
	"encoding/json"
)

type Response struct {
	Error error
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error"`
	}{
		Error: r.Error.Error(),
	})
}
