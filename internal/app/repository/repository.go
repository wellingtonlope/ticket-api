package repository

type AllRepositories struct {
	UserRepository   UserRepository
	TicketRepository TicketRepository
}

type Repositories interface {
	GetRepositories() (*AllRepositories, error)
}
