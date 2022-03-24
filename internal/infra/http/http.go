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
	DataFormat          = "2006-01-02T15:04:05"
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
	authMiddleware := AuthMiddleware{
		Authenticator: h.Authenticator,
	}
	userController := UserController{
		UCRegister:        h.UseCases.UserRegister,
		UCLogin:           h.UseCases.UserLogin,
		UCGetAllOperators: h.UseCases.UserGetAllOperators,
		Authenticator:     h.Authenticator,
	}
	ticketController := TicketController{
		UCOpen:             h.UseCases.TicketOpen,
		UCGet:              h.UseCases.TicketGet,
		UCClose:            h.UseCases.TicketClose,
		UCAssignToOperator: h.UseCases.TicketAssignToOperator,
		UCDelete:           h.UseCases.TicketDelete,
		UCGetByID:          h.UseCases.TicketGetByID,
		UCGetAll:           h.UseCases.TicketGetAll,
		UCGetAllByClient:   h.UseCases.TicketGetAllByClient,
		UCGetAllByOperator: h.UseCases.TicketGetAllByOperator,
		UCGetAllOpen:       h.UseCases.TicketGetAllOpen,
		Authenticator:      h.Authenticator,
	}

	h.Server.Register(Route{
		Method:  http.MethodPost,
		Path:    "/users/register",
		Handler: userController.Register,
	})
	h.Server.Register(Route{
		Method:  http.MethodPost,
		Path:    "/users/login",
		Handler: userController.Login,
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/users/operators",
		Handler:     userController.GetAllOperators,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodPost,
		Path:        "/tickets",
		Handler:     ticketController.Open,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodPost,
		Path:        "/tickets/:id/get",
		Handler:     ticketController.Get,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodPost,
		Path:        "/tickets/:id/close",
		Handler:     ticketController.Close,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodPost,
		Path:        "/tickets/:id/assign/:idOperator",
		Handler:     ticketController.AssignToOperator,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodDelete,
		Path:        "/tickets/:id",
		Handler:     ticketController.Delete,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/tickets/:id",
		Handler:     ticketController.GetByID,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/tickets",
		Handler:     ticketController.GetAll,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/tickets/client/:idClient",
		Handler:     ticketController.GetAllByClient,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/tickets/operator/:idOperator",
		Handler:     ticketController.GetAllByOperator,
		Middlewares: []Middleware{authMiddleware.Handle},
	})
	h.Server.Register(Route{
		Method:      http.MethodGet,
		Path:        "/tickets/open",
		Handler:     ticketController.GetAllOpen,
		Middlewares: []Middleware{authMiddleware.Handle},
	})

	return h.Server.Start(port)
}
