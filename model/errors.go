package model

type PaymentError struct {
	Error   error
	Code    int
	Message string
}

func NewError(err error, code int, message string) PaymentError {
	return PaymentError{
		Error:   err,
		Code:    code,
		Message: message,
	}
}

var (
	InternalError          = PaymentError{Code: 500, Message: "INTERNAL_SERVER_ERROR"}
	InvalidParams          = PaymentError{Code: 400, Message: "INVALID_PARAMS"}
	PermissionRestrictions = PaymentError{Code: 400, Message: "PERMISSION_RESTRICTIONS"}
	NotEnoughMoney         = PaymentError{Code: 409, Message: "Not enough money to withdrawal"}
)
