package mongodb

import (
	"context"
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *UserRepositoryMongo) Insert(user *domain.User) (*domain.User, *myerrors.Error) {
	_, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return user, nil
}

func (r *UserRepositoryMongo) GetById(id string) (*domain.User, *myerrors.Error) {
	user := domain.User{}

	err := r.Collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, myerrors.NewErrorMessage("user not found", myerrors.REGISTER_NOT_FOUND)
		}
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return &user, nil
}

func (r *UserRepositoryMongo) GetByEmail(email string) (*domain.User, *myerrors.Error) {
	user := domain.User{}

	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, myerrors.NewErrorMessage("user not found", myerrors.REGISTER_NOT_FOUND)
		}
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	return &user, nil
}

func (r *UserRepositoryMongo) GetAllOperator() (*[]domain.User, *myerrors.Error) {
	var users []domain.User
	cur, err := r.Collection.Find(context.Background(), bson.M{"profile": domain.PROFILE_OPERATOR})
	if err != nil {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		user := domain.User{}
		if err = cur.Decode(&user); err != nil {
			return nil, myerrors.NewError(err, myerrors.REPOSITORY)
		}
		users = append(users, user)
	}
	return &users, nil
}

func (r *UserRepositoryMongo) Delete(id string) *myerrors.Error {
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
