package local

import "github.com/wellingtonlope/ticket-api/application/repository"

type Repositories struct{}

func (r *Repositories) GetRepositories() (*repository.AllRepositories, error) {
	return &repository.AllRepositories{
		UserRepository:   &UserRepositoryLocal{},
		TicketRepository: &TicketRepositoryLocal{},
	}, nil
}
