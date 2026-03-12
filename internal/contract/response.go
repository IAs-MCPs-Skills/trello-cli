package contract

import "encoding/json"

type successEnvelope struct {
	OK   bool `json:"ok"`
	Data any  `json:"data"`
}

// Success builds a JSON success envelope.
func Success(data any) ([]byte, error) {
	return json.Marshal(successEnvelope{OK: true, Data: data})
}
