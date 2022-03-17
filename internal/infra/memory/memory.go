package memory

import "github.com/wellingtonlope/ticket-api/internal/app/repository"

type Repositories struct{}

func (r *Repositories) GetRepositories() (*repository.AllRepositories, error) {
	return &repository.AllRepositories{
		UserRepository:   &UserRepository{},
		TicketRepository: &TicketRepository{},
	}, nil
}
