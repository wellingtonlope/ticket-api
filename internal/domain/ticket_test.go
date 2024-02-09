package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenTicket(t *testing.T) {
	exampleDate := time.Now()
	exampleClient, _ := UserRegister("client", "client@mail.com", "password", exampleDate)
	type args struct {
		title       string
		description string
		createdAt   time.Time
		client      *User
	}
	testCases := []struct {
		name          string
		args          args
		assertResult  func(t *testing.T, got *Ticket)
		expectedError error
	}{
		{
			name: "should open a valid ticket",
			args: args{
				title:       "title",
				description: "description",
				createdAt:   exampleDate,
				client:      exampleClient,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.NotNil(t, got)
				assert.Equal(t, "title", got.Title)
				assert.Equal(t, "description", got.Description)
				assert.Equal(t, StatusOpen, got.Status)
				assert.Equal(t, exampleClient.ID, got.Client.ID)
				assert.Equal(t, exampleClient.Name, got.Client.Name)
				assert.Equal(t, exampleClient.Email, got.Client.Email)
				assert.Equal(t, exampleDate, *got.CreatedAt)
				assert.Nil(t, got.UpdatedAt)
			},
			expectedError: nil,
		},
		{
			name: "should return an error if title is empty",
			args: args{
				title:       "",
				description: "description",
				createdAt:   exampleDate,
				client:      exampleClient,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.Nil(t, got)
			},
			expectedError: ErrTicketTitleIsInvalid,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := OpenTicket(tc.args.title, tc.args.description, tc.args.createdAt, *tc.args.client)
			tc.assertResult(t, got)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTicketGet(t *testing.T) {
	exampleDate := time.Now()
	exampleClient, _ := UserRegister("client", "client@mail.com", "password", exampleDate)
	exampleOperator, _ := UserRegister("operator", "operator@mail.com", "password", exampleDate)
	exampleOperator.Profile = ProfileOperator
	exampleTicketOpen, _ := OpenTicket("title", "description", exampleDate, *exampleClient)
	type args struct {
		operator  *User
		updatedAt time.Time
	}
	testCases := []struct {
		name          string
		ticket        Ticket
		args          args
		assertResult  func(t *testing.T, got *Ticket)
		expectedError error
	}{
		{
			name:   "should get a valid ticket",
			ticket: *exampleTicketOpen,
			args: args{
				operator:  exampleOperator,
				updatedAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.Equal(t, exampleOperator.ID, got.Operator.ID)
				assert.Equal(t, exampleOperator.Name, got.Operator.Name)
				assert.Equal(t, exampleOperator.Email, got.Operator.Email)
				assert.Equal(t, StatusInProgress, got.Status)
				assert.Equal(t, exampleDate, *got.UpdatedAt)
			},
			expectedError: nil,
		},
		{
			name:   "should return an error if operator is not an operator",
			ticket: *exampleTicketOpen,
			args: args{
				operator:  exampleClient,
				updatedAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.Equal(t, *exampleTicketOpen, *got)
			},
			expectedError: ErrTicketNoOperator,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ticket.Get(*tc.args.operator, tc.args.updatedAt)
			tc.assertResult(t, &tc.ticket)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTicketClose(t *testing.T) {
	exampleDate := time.Now()
	exampleSolution := "solution"
	exampleClient, _ := UserRegister("client", "client@mail.com", "password", exampleDate)
	exampleOperator, _ := UserRegister("operator", "operator@mail.com", "password", exampleDate)
	exampleOperator.Profile = ProfileOperator
	exampleTicketOpen, _ := OpenTicket("title", "description", exampleDate, *exampleClient)
	exampleTicketGet := func() *Ticket {
		m := *exampleTicketOpen
		_ = m.Get(*exampleOperator, exampleDate)
		return &m
	}()
	type args struct {
		solution  string
		updatedAt time.Time
	}
	testCases := []struct {
		name          string
		ticket        Ticket
		args          args
		assertResult  func(t *testing.T, got *Ticket)
		expectedError error
	}{
		{
			name:   "should close a valid ticket",
			ticket: *exampleTicketGet,
			args: args{
				solution:  exampleSolution,
				updatedAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.Equal(t, StatusClose, got.Status)
				assert.Equal(t, exampleDate, *got.UpdatedAt)
				assert.Equal(t, exampleSolution, got.Solution)
			},
			expectedError: nil,
		},
		{
			name:   "should return an error if no operator",
			ticket: *exampleTicketOpen,
			args: args{
				solution:  exampleSolution,
				updatedAt: exampleDate,
			},
			assertResult: func(t *testing.T, got *Ticket) {
				assert.Equal(t, *exampleTicketOpen, *got)
			},
			expectedError: ErrTicketNoGetToClose,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ticket.Close(tc.args.solution, tc.args.updatedAt)
			tc.assertResult(t, &tc.ticket)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
