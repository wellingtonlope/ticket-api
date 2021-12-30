package mongodb

import (
	"context"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TicketRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *TicketRepositoryMongo) Insert(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error) {
	_, err := r.Collection.InsertOne(context.Background(), ticket)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return ticket, nil
}

func (r *TicketRepositoryMongo) GetById(id string) (*domain.Ticket, *myerrors.Error) {
	ticket := domain.Ticket{}

	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&ticket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, myerrors.NewErrorMessage("ticket not found", myerrors.REGISTER_NOT_FOUND)
		}
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return &ticket, nil
}

func (r *TicketRepositoryMongo) Update(ticket *domain.Ticket) (*domain.Ticket, *myerrors.Error) {
	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": ticket.ID}, bson.M{"$set": ticket})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return ticket, nil
}

func (r *TicketRepositoryMongo) GetAll() (*[]domain.Ticket, *myerrors.Error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := domain.Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, myerrors.NewError(err, myerrors.REPOSITORY)
		}
		tickets = append(tickets, ticket)
	}
	return &tickets, nil
}

func (r *TicketRepositoryMongo) GetAllOpen() (*[]domain.Ticket, *myerrors.Error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"status": domain.STATUS_OPEN})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := domain.Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, myerrors.NewError(err, myerrors.REPOSITORY)
		}
		tickets = append(tickets, ticket)
	}
	return &tickets, nil
}

func (r *TicketRepositoryMongo) GetAllByOperator(operator *domain.User) (*[]domain.Ticket, *myerrors.Error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"operator.id": operator.ID})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := domain.Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, myerrors.NewError(err, myerrors.REPOSITORY)
		}
		tickets = append(tickets, ticket)
	}
	return &tickets, nil
}

func (r *TicketRepositoryMongo) GetAllByClient(client *domain.User) (*[]domain.Ticket, *myerrors.Error) {
	var tickets []domain.Ticket

	cur, err := r.Collection.Find(context.Background(), bson.M{"client.id": client.ID})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		ticket := domain.Ticket{}
		if err = cur.Decode(&ticket); err != nil {
			return nil, myerrors.NewError(err, myerrors.REPOSITORY)
		}
		tickets = append(tickets, ticket)
	}
	return &tickets, nil
}

func (r *TicketRepositoryMongo) Delete(id string) *myerrors.Error {
	_, myerr := r.GetById(id)
	if myerr != nil {
		return myerr
	}

	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return nil
}
