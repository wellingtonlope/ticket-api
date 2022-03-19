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
		Title       string
		Description string
	}
)

func ticketUserResponseFromUserOutput(output *ticket.TicketUserOutput) *TicketUserResponse {
	if output == nil {
		return nil
	}

	return &TicketUserResponse{
		ID:    output.ID,
		Name:  output.Name,
		Email: output.Email,
	}
}

func ticketResponseFromOutput(output *ticket.TicketOutput) *TicketResponse {
	var updatedAt string
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
		CreatedAt:   output.CreatedAt.Format(DataFormat),
		UpdatedAt:   updatedAt,
	}
}

type TicketController struct {
	UCOpen        *ticket.Open
	UCGet         *ticket.Get
	Authenticator security.Authenticator
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
		Body:     wrapBody(ticketResponseFromOutput(output)),
	}
}

func (c *TicketController) Get(request Request) Response {
	id := request.Params["id"]

	input := ticket.GetInput{
		TicketID:   id,
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
		Body:     wrapBody(ticketResponseFromOutput(output)),
	}
}
