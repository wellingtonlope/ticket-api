package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/application/usecase/user"
	"net/http"
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
)

type UserHandler struct {
	userUseCase *user.UserUseCase
}

func initUserHandler(e *echo.Echo, userUseCase *user.UserUseCase) {
	h := UserHandler{userUseCase: userUseCase}
	e.POST("/user/register", h.Register)
	e.POST("/user/login", h.Login)
	e.GET("/user/operators", h.GetAllOperators)
}

func (h *UserHandler) Register(c echo.Context) error {
	request := userRegisterRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.userUseCase.Register(request.Name, request.Email, request.Password)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) Login(c echo.Context) error {
	request := userLoginRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.userUseCase.Login(request.Email, request.Password)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetAllOperators(c echo.Context) error {
	authorization := getAuthorization(c)

	response, myerr := h.userUseCase.GetAllOperator(authorization)

	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}
