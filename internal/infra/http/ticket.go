package http

import (
	"encoding/json"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/security"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/ticket"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"net/http"
	"time"
)

type (
	TicketUserResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	TicketResponse struct {
		ID          string              `json:"id"`
		Title       string              `json:"title"`
		Description string              `json:"description"`
		Solution    string              `json:"solution,omitempty"`
		Status      string              `json:"status"`
		Client      *TicketUserResponse `json:"client"`
		Operator    *TicketUserResponse `json:"operator,omitempty"`
		CreatedAt   string              `json:"created_at"`
		UpdatedAt   string              `json:"updated_at,omitempty"`
	}
	OpenRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	CloseRequest struct {
		Solution string `json:"solution"`
	}
)

func ticketUserResponseFromUserOutput(output *ticket.UserOutput) *TicketUserResponse {
	if output == nil {
		return nil
	}

	return &TicketUserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	}
}

func ticketResponseFromOutput(output ticket.Output) *TicketResponse {
	var createdAt, updatedAt string
	if output.CreatedAt != nil {
		createdAt = output.CreatedAt.Format(DataFormat)
	}
	if output.UpdatedAt != nil {
		updatedAt = output.UpdatedAt.Format(DataFormat)
	}
	return &TicketResponse{
		ID:          output.ID,
		Title:       output.Title,
		Description: output.Description,
		Solution:    output.Solution,
		Status:      output.Status,
		Client:      ticketUserResponseFromUserOutput(output.Client),
		Operator:    ticketUserResponseFromUserOutput(output.Operator),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func ticketsResponseFromOutputs(outputs []ticket.Output) *[]TicketResponse {
	responses := make([]TicketResponse, 0, len(outputs))
	for _, output := range outputs {
		responses = append(responses, *ticketResponseFromOutput(output))
	}
	return &responses
}

type TicketController struct {
	UCOpen             ticket.Open
	UCGet              ticket.Get
	UCClose            ticket.Close
	UCAssignToOperator ticket.AssignToOperator
	UCDelete           ticket.Delete
	UCGetByID          ticket.GetByID
	UCGetAll           ticket.GetAll
	UCGetAllByClient   ticket.GetAllByClient
	UCGetAllByOperator ticket.GetAllByOperator
	UCGetAllOpen       ticket.GetAllOpen
	Authenticator      security.Authenticator
}

func (c *TicketController) Open(request Request) Response {
	var openRequest OpenRequest
	err := json.Unmarshal([]byte(request.Body), &openRequest)
	if err != nil {
		return Response{
			HttpCode: http.StatusBadRequest,
			Body:     wrapError(ErrInvalidJsonBody),
		}
	}

	input := ticket.OpenInput{
		Title:       openRequest.Title,
		Description: openRequest.Description,
		CreatedAt:   time.Now(),
		LoggedUser:  *request.LoggedUser,
	}
	output, err := c.UCOpen.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case domain.ErrTicketTitleIsInvalid:
			httpStatus = http.StatusBadRequest
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketResponseFromOutput(*output)),
	}
}

func (c *TicketController) Get(request Request) Response {
	input := ticket.GetInput{
		TicketID:   request.Params["id"],
		UpdatedAt:  time.Now(),
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGet.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		case repository.ErrTicketNotFound:
			httpStatus = http.StatusNotFound
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketResponseFromOutput(*output)),
	}
}

func (c *TicketController) Close(request Request) Response {
	var closeRequest CloseRequest
	err := json.Unmarshal([]byte(request.Body), &closeRequest)

	input := ticket.CloseInput{
		TicketID:   request.Params["id"],
		Solution:   closeRequest.Solution,
		UpdatedAt:  time.Now(),
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCClose.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case domain.ErrTicketNoGetToClose:
			httpStatus = http.StatusBadRequest
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		case repository.ErrTicketNotFound:
			httpStatus = http.StatusNotFound
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketResponseFromOutput(*output)),
	}
}

func (c *TicketController) AssignToOperator(request Request) Response {
	input := ticket.AssignToOperatorInput{
		TicketID:   request.Params["id"],
		OperatorID: request.Params["idOperator"],
		UpdatedAt:  time.Now(),
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCAssignToOperator.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case domain.ErrTicketNoOperator:
			httpStatus = http.StatusBadRequest
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		case repository.ErrTicketNotFound, repository.ErrUserNotFound:
			httpStatus = http.StatusNotFound
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketResponseFromOutput(*output)),
	}
}

func (c *TicketController) Delete(request Request) Response {
	input := ticket.DeleteInput{
		TicketID:   request.Params["id"],
		LoggedUser: *request.LoggedUser,
	}
	_, err := c.UCDelete.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		case repository.ErrTicketNotFound:
			httpStatus = http.StatusNotFound
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusNoContent,
		Body:     "",
	}
}

func (c *TicketController) GetByID(request Request) Response {
	input := ticket.GetByIDInput{
		TicketID:   request.Params["id"],
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGetByID.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case repository.ErrTicketNotFound:
			httpStatus = http.StatusNotFound
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketResponseFromOutput(*output)),
	}
}

func (c *TicketController) GetAll(request Request) Response {
	input := ticket.GetAllInput{
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGetAll.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketsResponseFromOutputs(*output)),
	}
}

func (c *TicketController) GetAllByClient(request Request) Response {
	input := ticket.GetAllByClientInput{
		ClientID:   request.Params["idClient"],
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGetAllByClient.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketsResponseFromOutputs(*output)),
	}
}

func (c *TicketController) GetAllByOperator(request Request) Response {
	input := ticket.GetAllByOperatorInput{
		OperatorID: request.Params["idOperator"],
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGetAllByOperator.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketsResponseFromOutputs(*output)),
	}
}

func (c *TicketController) GetAllOpen(request Request) Response {
	input := ticket.GetAllOpenInput{
		LoggedUser: *request.LoggedUser,
	}
	output, err := c.UCGetAllOpen.Handle(input)

	if err != nil {
		httpStatus := http.StatusInternalServerError
		switch err {
		case security.ErrForbidden:
			httpStatus = http.StatusForbidden
		}
		return Response{
			HttpCode: httpStatus,
			Body:     wrapError(err),
		}
	}

	return Response{
		HttpCode: http.StatusOK,
		Body:     wrapBody(ticketsResponseFromOutputs(*output)),
	}
}
