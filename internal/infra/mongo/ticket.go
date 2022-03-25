package mongo

import (
	"context"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketRepository struct {
	Collection *mongo.Collection
}

type TicketUser struct {
	ID    string `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type Ticket struct {
	ID          primitive.ObjectID  `bson:"_id"`
	Title       string              `bson:"title"`
	Description string              `bson:"description,omitempty"`
	Solution    string              `bson:"solution,omitempty"`
	Status      string              `bson:"status"`
	Client      *TicketUser         `bson:"client"`
	Operator    *TicketUser         `bson:"operator,omitempty"`
	CreatedAt   *primitive.DateTime `bson:"created_at"`
	UpdatedAt   *primitive.DateTime `bson:"updated_at,omitempty"`
}

func domainToTicket(ticketDomain domain.Ticket) *Ticket {
	objectId, err := primitive.ObjectIDFromHex(ticketDomain.ID)
	if err != nil {
		objectId = primitive.NewObjectID()
	}

	ticket := Ticket{
		ID:          objectId,
		Title:       ticketDomain.Title,
		Description: ticketDomain.Description,
		Solution:    ticketDomain.Solution,
		Status:      string(ticketDomain.Status),
	}
	if ticketDomain.Client != nil {
		ticket.Client = &TicketUser{
			ID:    ticketDomain.Client.ID,
			Name:  ticketDomain.Client.Name,
			Email: ticketDomain.Client.Email.String(),
		}
	}
	if ticketDomain.Operator != nil {
		ticket.Operator = &TicketUser{
			ID:    ticketDomain.Operator.ID,
			Name:  ticketDomain.Operator.Name,
			Email: ticketDomain.Operator.Email.String(),
		}
	}
	if ticketDomain.CreatedAt != nil {
		createdAt := primitive.NewDateTimeFromTime(*ticketDomain.CreatedAt)
		ticket.CreatedAt = &createdAt
	}
	if ticketDomain.UpdatedAt != nil {
		updatedAt := primitive.NewDateTimeFromTime(*ticketDomain.UpdatedAt)
		ticket.UpdatedAt = &updatedAt
	}
	return &ticket
}

func (t Ticket) toDomain() *domain.Ticket {
	ticketDomain := domain.Ticket{
		ID:          t.ID.Hex(),
		Title:       t.Title,
		Description: t.Description,
		Solution:    t.Solution,
		Status:      domain.Status(t.Status),
	}
	if t.Client != nil {
		email, _ := domain.NewEmail(t.Client.Email)
		ticketDomain.Client = &domain.TicketUser{
			ID:    t.Client.ID,
			Name:  t.Client.Name,
			Email: *email,
		}
	}
	if t.Operator != nil {
		email, _ := domain.NewEmail(t.Client.Email)
		ticketDomain.Operator = &domain.TicketUser{
			ID:    t.Operator.ID,
			Name:  t.Operator.Name,
			Email: *email,
		}
	}
	if t.CreatedAt != nil {
		createdAt := t.CreatedAt.Time()
		ticketDomain.CreatedAt = &createdAt
	}
	if t.UpdatedAt != nil {
		updatedAt := t.UpdatedAt.Time()
		ticketDomain.UpdatedAt = &updatedAt
	}
	return &ticketDomain
}

func (r *TicketRepository) Insert(domainTicket domain.Ticket) (*domain.Ticket, error) {
	ticket := domainToTicket(domainTicket)
	_, err := r.Collection.InsertOne(context.Background(), ticket)
	if err != nil {
		return nil, err
	}

	return ticket.toDomain(), nil
}

func (r *TicketRepository) GetByID(id string) (*domain.Ticket, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}
	ticket := Ticket{}

	err = r.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrTicketNotFound
		}
		return nil, err
	}

	return ticket.toDomain(), nil
}

func (r *TicketRepository) Update(domainTicket domain.Ticket) (*domain.Ticket, error) {
	ticket := domainToTicket(domainTicket)
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": ticket.ID}, bson.M{"$set": ticket})
	if err != nil {
		return nil, err
	}

	return ticket.toDomain(), nil
}

func (r *TicketRepository) GetAll() (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket.toDomain())
	}
	return &tickets, nil
}

func (r *TicketRepository) GetAllOpen() (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"status": domain.STATUS_OPEN})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket.toDomain())
	}
	return &tickets, nil
}

func (r *TicketRepository) GetAllByOperatorID(operatorID string) (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"operator.id": operatorID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket.toDomain())
	}
	return &tickets, nil
}

func (r *TicketRepository) GetAllByClientID(clientID string) (*[]domain.Ticket, error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"client.id": clientID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket.toDomain())
	}
	return &tickets, nil
}

func (r *TicketRepository) DeleteByID(id string) error {
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repository.ErrUserNotFound
	}

	_, err = r.Collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
