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

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	OK    bool        `json:"ok"`
	Error errorDetail `json:"error"`
}

// ErrorEnvelope builds a JSON error envelope from a code and message.
func ErrorEnvelope(code, message string) ([]byte, error) {
	return json.Marshal(errorEnvelope{
		OK:    false,
		Error: errorDetail{Code: code, Message: message},
	})
}

// ErrorFromContractError builds a JSON error envelope from a ContractError.
func ErrorFromContractError(ce *ContractError) ([]byte, error) {
	return ErrorEnvelope(ce.Code, ce.Message)
}
