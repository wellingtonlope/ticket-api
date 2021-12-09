package local

import (
	"github.com/wellingtonlope/ticket-api/application/myerrors"
	"github.com/wellingtonlope/ticket-api/domain"
)

type UserRepositoryLocal struct {
	Users []domain.User
}

func (r *UserRepositoryLocal) Insert(user *domain.User) (*domain.User, *myerrors.Error) {
	userGot, err := r.GetById(user.ID)
	if userGot != nil && err == nil {
		return nil, myerrors.NewErrorMessage("user already exist", myerrors.REGISTER_ALREADY_EXISTS)
	}
	if err != nil && err.Type != myerrors.REGISTER_NOT_FOUND {
		return nil, myerrors.NewError(err, myerrors.REPOSITORY)
	}

	r.Users = append(r.Users, *user)
	return user, nil
}

func (r *UserRepositoryLocal) GetById(id string) (*domain.User, *myerrors.Error) {
	for _, user := range r.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, myerrors.NewErrorMessage("user not found", myerrors.REGISTER_NOT_FOUND)
}

func (r *UserRepositoryLocal) GetByEmail(email string) (*domain.User, *myerrors.Error) {
	for _, user := range r.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, myerrors.NewErrorMessage("user not found", myerrors.REGISTER_NOT_FOUND)
}

func (r *UserRepositoryLocal) GetAllOperator() (*[]domain.User, *myerrors.Error) {
	operators := []domain.User{}

	for _, user := range r.Users {
		if user.Profile == domain.PROFILE_OPERATOR {
			operators = append(operators, user)
		}
	}

	return &operators, nil
}

func (r *UserRepositoryLocal) Delete(id string) *myerrors.Error {
	_, err := r.GetById(id)
	if err != nil {
		return myerrors.NewErrorMessage("user not found", myerrors.REGISTER_NOT_FOUND)
	}

	for index, user := range r.Users {
		if user.ID == id {
			r.Users = removeIndexUser(r.Users, index)
		}
	}

	return nil
}

func removeIndexUser(s []domain.User, index int) []domain.User {
	return append(s[:index], s[index+1:]...)
}
