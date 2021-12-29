package myerrors

type Type string

const (
	DOMAIN                  Type = "DOMAIN"
	USECASE                 Type = "USECASE"
	REPOSITORY              Type = "REPOSITORY"
	REGISTER_NOT_FOUND      Type = "REGISTER_NOT_FOUND"
	REGISTER_ALREADY_EXISTS Type = "REGISTER_ALREADY_EXISTS"
	UNIDENTIFIED            Type = "UNIDENTIFIED"
	UNAUTHORIZED            Type = "UNAUTHORIZED"
	FORBIDDEN               Type = "FORBIDDEN"
)

type Error struct {
	Message string `json:"message"`
	Type    Type   `json:"type"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(err error, tp Type) *Error {
	return &Error{err.Error(), tp}
}

func NewErrorMessage(message string, tp Type) *Error {
	return &Error{message, tp}
}
