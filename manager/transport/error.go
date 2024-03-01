// SPDX-License-Identifier: Apache-2.0

package transport

import "fmt"

// ErrorCode represents an OCPP error code.
type ErrorCode string

const (
	ErrorFormatViolation               ErrorCode = "FormatViolation"
	ErrorGenericError                  ErrorCode = "GenericError"
	ErrorInternalError                 ErrorCode = "InternalError"
	ErrorMessageTypeNotSupported       ErrorCode = "MessageTypeNotSupported"
	ErrorNotImplemented                ErrorCode = "NotImplemented"
	ErrorNotSupported                  ErrorCode = "NotSupported"
	ErrorOccurrenceConstraintViolation ErrorCode = "OccurrenceConstraintViolation"
	ErrorPropertyConstraintViolation   ErrorCode = "PropertyConstraintViolation"
	ErrorProtocolError                 ErrorCode = "ProtocolError"
	ErrorRpcFrameworkError             ErrorCode = "RpcFrameworkError"
	ErrorSecurityError                 ErrorCode = "SecurityError"
	ErrorTypeConstraintViolation       ErrorCode = "TypeConstraintViolation"
)

// An Error allows a function to return an error that includes an OCPP ErrorCode.
type Error struct {
	ErrorCode    ErrorCode
	WrappedError error
}

func NewError(code ErrorCode, err error) *Error {
	return &Error{
		ErrorCode:    code,
		WrappedError: err,
	}
}

func (e Error) Error() string {
	if e.WrappedError != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.WrappedError)
	} else {
		return string(e.ErrorCode)
	}
}

func (e Error) Unwrap() error {
	return e.WrappedError
}
