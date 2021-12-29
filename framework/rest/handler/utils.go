package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"net/http"
)

func handlerError(c echo.Context, err *myerrors.Error) error {
	switch err.Type {
	case myerrors.DOMAIN, myerrors.REGISTER_ALREADY_EXISTS, myerrors.USECASE, myerrors.UNIDENTIFIED:
		return c.JSON(http.StatusBadRequest, err)
	case myerrors.REGISTER_NOT_FOUND:
		return c.JSON(http.StatusNotFound, err)
	case myerrors.REPOSITORY:
		return c.JSON(http.StatusInternalServerError, err)
	case myerrors.UNAUTHORIZED:
		return c.JSON(http.StatusUnauthorized, err)
	case myerrors.FORBIDDEN:
		return c.JSON(http.StatusForbidden, err)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}
