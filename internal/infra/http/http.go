package http

import (
	"encoding/json"
	"errors"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase"
	"net/http"
)

const (
	AuthorizationHeader = "Authorization"
	ContentTypeHeader   = "Content-Type"
	ContentTypeJSON     = "application/json"
)

type Handler func(r Request) Response
type Middleware func(h Handler) Handler
type (
	Request struct {
		LoggedUser *security.User
		Header     map[string]string
		Params     map[string]string
		Query      map[string]string
		Body       string
	}
	Response struct {
		HttpCode int
		Body     string
	}
	Route struct {
		Method      string
		Path        string
		Handler     Handler
		Middlewares []Middleware
	}
	ErrorResponse struct {
		Message string `json:"message"`
	}
)

var (
	ErrInvalidJsonBody = errors.New("invalid json body")
)

type Server interface {
	Register(r Route)
	Start(port int) error
}

type Http struct {
	Server        Server
	UseCases      *usecase.AllUseCases
	Authenticator security.Authenticator
}

func wrapError(err error) string {
	response := ErrorResponse{
		Message: err.Error(),
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func wrapBody(body interface{}) string {
	bytes, err := json.Marshal(body)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (h *Http) Start(port int) error {
	controller := UserController{
		UCRegister:        h.UseCases.UserRegister,
		UCLogin:           h.UseCases.UserLogin,
		UCGetAllOperators: h.UseCases.UserGetAllOperators,
		Authenticator:     h.Authenticator,
	}
	authMiddleware := AuthMiddleware{
		Authenticator: h.Authenticator,
	}

	h.Server.Register(Route{
		Method:  http.MethodPost,
		Path:    "/users/register",
		Handler: controller.Register,
	})
	h.Server.Register(Route{
		Method:  http.MethodPost,
		Path:    "/users/login",
		Handler: controller.Login,
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/users/operators",
		Handler:     controller.GetAllOperators,
		Middlewares: []Middleware{authMiddleware.Handle},
	})

	return h.Server.Start(port)
}
