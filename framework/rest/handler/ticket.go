package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/application/usecase/ticket"
	"net/http"
)

type (
	ticketOpenRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	ticketGetRequest struct {
		Id string `param:"id"`
	}
	ticketCloseRequest struct {
		Id       string `param:"id"`
		Solution string `json:"solution"`
	}
	ticketAssignToOperatorRequest struct {
		Id         string `param:"id"`
		IdOperator string `param:"idOperator"`
	}
	ticketDeleteRequest struct {
		Id string `param:"id"`
	}
	ticketGetByIdRequest struct {
		Id string `param:"id"`
	}
	ticketGetAllByClientRequest struct {
		IdClient string `param:"idClient"`
	}
	ticketGetAllByOperatorRequest struct {
		IdOperator string `param:"idOperator"`
	}
)

type TicketHandler struct {
	ticketUseCase *ticket.TicketUseCase
}

func initTicketHandler(e *echo.Echo, ticketUseCase *ticket.TicketUseCase) {
	h := TicketHandler{ticketUseCase: ticketUseCase}
	e.POST("/ticket", h.Open)
	e.POST("/ticket/:id/get", h.Get)
	e.POST("/ticket/:id/close", h.Close)
	e.POST("/ticket/:id/assign/:idOperator", h.AssignToOperator)
	e.DELETE("/ticket/:id", h.Delete)
	e.GET("/ticket/:id", h.GetById)
	e.GET("/ticket", h.GetAll)
	e.GET("/ticket/client/:idClient", h.GetAllByClient)
	e.GET("/ticket/operator/:idOperator", h.GetAllByOperator)
	e.GET("/ticket/open", h.GetAllOpen)
}

func (h *TicketHandler) Open(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketOpenRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.Open(request.Title, request.Description, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *TicketHandler) Get(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketGetRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.Get(request.Id, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) Close(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketCloseRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.Close(request.Id, request.Solution, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) AssignToOperator(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketAssignToOperatorRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.AssignToOperator(request.Id, request.IdOperator, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) Delete(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketDeleteRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	myerr := h.ticketUseCase.Delete(request.Id, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TicketHandler) GetById(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketGetByIdRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.GetById(request.Id, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) GetAll(c echo.Context) error {
	authorization := getAuthorization(c)

	response, myerr := h.ticketUseCase.GetAll(authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) GetAllByClient(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketGetAllByClientRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.GetAllByClient(request.IdClient, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) GetAllByOperator(c echo.Context) error {
	authorization := getAuthorization(c)

	request := ticketGetAllByOperatorRequest{}
	err := c.Bind(&request)
	if err != nil {
		return handlerError(c, myerrors.NewError(err, myerrors.UNIDENTIFIED))
	}

	response, myerr := h.ticketUseCase.GetAllByOperator(request.IdOperator, authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TicketHandler) GetAllOpen(c echo.Context) error {
	authorization := getAuthorization(c)

	response, myerr := h.ticketUseCase.GetAllOpen(authorization)
	if myerr != nil {
		return handlerError(c, myerr)
	}

	return c.JSON(http.StatusOK, response)
}
