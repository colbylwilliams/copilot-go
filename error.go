package copilot

import (
	"encoding/json"
	"errors"
)

// Error represents an error that occurred during the agent request.
type Error struct {
	// Type is a string that specifies the error's type. type can have a value of reference, function or agent.
	Type ErrorType `json:"type"`
	// Code is string controlled by the agent describing the nature of an error.
	Code string `json:"code"`
	// Message is string that specifies the error message shown to the user.
	Message string `json:"message"`
	// Identifier is string that serves as a unique identifier to link the error with other resources such as references or function calls.
	Identifier string `json:"identifier"`
}

func (a *Error) Error() string {
	return a.Message
}

// ErrorType is a string that specifies the error's type. type can have a value of reference, function or agent.
type ErrorType struct {
	name string
}

func (a ErrorType) String() string {
	return a.name
}

var (
	ErrorTypeReference = ErrorType{"reference"}
	ErrorTypeFunction  = ErrorType{"function"}
	ErrorTypeAgent     = ErrorType{"agent"}
)

func (a *ErrorType) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	if name != ErrorTypeAgent.name && name != ErrorTypeFunction.name && name != ErrorTypeReference.name {
		return errors.New("invalid agent error type, got " + name)
	}

	a.name = name
	return nil
}

func (a ErrorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.name)
}
