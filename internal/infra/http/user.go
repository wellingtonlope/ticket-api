package http

import (
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
	"github.com/wellingtonlope/ticket-api/internal/domain"
)

type (
	userRegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	userTokenResponse struct {
		Token string `json:"token"`
	}
	UserResponse struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}
)

func userResponseFromUserOutput(user *user.UserOutput) *UserResponse {
	var updatedAt string
	if user.UpdatedAt != nil {
		updatedAt = user.UpdatedAt.Format("2006-01-02 15:04:05")
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedAt,
	}
}

func usersResponseFromUsersOutput(outputs *[]user.UserOutput) *[]UserResponse {
	response := make([]UserResponse, 0, len(*outputs))
	for _, output := range *outputs {
		response = append(response, *userResponseFromUserOutput(&output))
	}
	return &response
}

type UserHandler struct {
	userRegister        *user.Register
	userLogin           *user.Login
	userGetAllOperators *user.GetAllOperators
	authenticator       *jwt.Authenticator
}

func initUserHandler(e *echo.Echo, useCases *usecase.AllUseCases, authenticator *jwt.Authenticator) {
	handler := &UserHandler{
		userRegister:        useCases.UserRegister,
		userLogin:           useCases.UserLogin,
		userGetAllOperators: useCases.UserGetAllOperators,
		authenticator:       authenticator,
	}

	authMiddleware := NewAuthMiddleware(authenticator)

	e.POST("/users/register", handler.register)
	e.POST("/users/login", handler.login)
	e.GET("/users/operators", handler.getAllOperators, authMiddleware.Handle)
}

func (h *UserHandler) register(c echo.Context) error {
	request := userRegisterRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapError(err))
	}

	input := user.RegisterInput{
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: time.Now(),
	}

	output, err := h.userRegister.Handle(input)

	if err != nil {
		switch err {
		case user.ErrUserAlreadyExists:
		case domain.ErrNameIsInvalid:
		case domain.ErrEmailIsInvalid:
		case domain.ErrPasswordIsInvalid:
			return c.JSON(http.StatusBadRequest, wrapError(err))
		default:
			return c.JSON(http.StatusInternalServerError, wrapError(err))
		}
	}

	token, err := h.authenticator.Generate(output.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, wrapError(err))
	}

	return c.JSON(http.StatusCreated, userTokenResponse{Token: token})
}

func (h *UserHandler) login(c echo.Context) error {
	request := userLoginRequest{}
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, wrapError(err))
	}

	input := user.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	}

	output, err := h.userLogin.Handle(input)
	if err != nil {
		switch err {
		case user.ErrUserEmailPasswordWrong:
			return c.JSON(http.StatusBadRequest, wrapError(err))
		default:
			return c.JSON(http.StatusInternalServerError, wrapError(err))
		}
	}

	token, err := h.authenticator.Generate(output.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, wrapError(err))
	}

	return c.JSON(http.StatusOK, userTokenResponse{Token: token})
}

func (h *UserHandler) getAllOperators(c echo.Context) error {
	loggedUser := c.Get(ContextUser).(*domain.User)
	input := user.GetAllOperatorsInput{
		LoggedUser: *loggedUser,
	}

	output, err := h.userGetAllOperators.Handle(input)
	if err != nil {
		switch err {
		case security.ErrForbidden:
			return c.JSON(http.StatusForbidden, wrapError(err))
		default:
			return c.JSON(http.StatusInternalServerError, wrapError(err))
		}
	}

	return c.JSON(http.StatusOK, usersResponseFromUsersOutput(output))
}
