package mqtt

import "fmt"

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

type Error struct {
	ErrorCode    ErrorCode
	wrappedError error
}

func NewError(code ErrorCode, err error) *Error {
	return &Error{
		ErrorCode:    code,
		wrappedError: err,
	}
}

func (e Error) Error() string {
	if e.wrappedError != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.wrappedError)
	} else {
		return string(e.ErrorCode)
	}
}

func (e Error) Unwrap() error {
	return e.wrappedError
}
