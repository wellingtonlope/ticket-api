package mongo

import (
	"context"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

type User struct {
	ID        primitive.ObjectID  `bson:"_id"`
	Name      string              `bson:"name"`
	Email     string              `bson:"email"`
	Password  string              `bson:"password"`
	Profile   string              `bson:"profile"`
	CreatedAt *primitive.DateTime `bson:"created_at"`
	UpdatedAt *primitive.DateTime `bson:"updated_at,omitempty"`
}

func domainToUser(userDomain domain.User) *User {
	objectId, err := primitive.ObjectIDFromHex(userDomain.ID)
	if err != nil {
		objectId = primitive.NewObjectID()
	}

	user := User{
		ID:       objectId,
		Name:     userDomain.Name,
		Email:    userDomain.Email.String(),
		Password: userDomain.Password.String(),
		Profile:  string(userDomain.Profile),
	}
	if userDomain.CreatedAt != nil {
		createdAt := primitive.NewDateTimeFromTime(*userDomain.CreatedAt)
		user.CreatedAt = &createdAt
	}
	if userDomain.UpdatedAt != nil {
		updatedAt := primitive.NewDateTimeFromTime(*userDomain.UpdatedAt)
		user.UpdatedAt = &updatedAt
	}
	return &user
}

func (u User) toDomain() *domain.User {
	email, _ := domain.NewEmail(u.Email)
	password, _ := domain.NewPasswordHashed(u.Password)
	userDomain := domain.User{
		ID:       u.ID.Hex(),
		Name:     u.Name,
		Email:    *email,
		Password: *password,
		Profile:  domain.Profile(u.Profile),
	}
	if u.CreatedAt != nil {
		createdAt := u.CreatedAt.Time()
		userDomain.CreatedAt = &createdAt
	}
	if u.UpdatedAt != nil {
		updatedAt := u.CreatedAt.Time()
		userDomain.UpdatedAt = &updatedAt
	}
	return &userDomain
}

func (r *UserRepository) Insert(domainUser domain.User) (*domain.User, error) {
	user := domainToUser(domainUser)
	_, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return user.toDomain(), nil
}

func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, repository.ErrUserNotFound
	}
	user := User{}

	err = r.Collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return user.toDomain(), nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := User{}

	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return user.toDomain(), nil
}

func (r *UserRepository) GetAllOperator() (*[]domain.User, error) {
	var users []domain.User
	cur, err := r.Collection.Find(context.Background(), bson.M{"profile": domain.ProfileOperator})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		user := User{}
		if err = cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, *user.toDomain())
	}

	err = cur.Close(context.Background())
	if err != nil {
		return nil, err
	}

	return &users, nil
}
