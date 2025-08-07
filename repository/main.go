package repository

import (
	"github.com/tnqbao/gau-kanban-service/infra"
	"gorm.io/gorm"
)

// Repositories holds all repository interfaces
type Repositories struct {
	Column         ColumnRepositoryInterface
	Ticket         TicketRepositoryInterface
	Label          LabelRepositoryInterface
	TaskAssignment TaskAssignmentRepositoryInterface
	TicketComment  TicketCommentRepositoryInterface
	TicketLabel    TicketLabelRepositoryInterface
}

// NewRepositories creates and initializes all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Column:         NewColumnRepository(db),
		Ticket:         NewTicketRepository(db),
		Label:          NewLabelRepository(db),
		TaskAssignment: NewTaskAssignmentRepository(db),
		TicketComment:  NewTicketCommentRepository(db),
		TicketLabel:    NewTicketLabelRepository(db),
	}
}

// RepositoryInterface defines a unified interface for all repositories
type RepositoryInterface interface {
	GetRepositories() *Repositories
}

// RepositoryManager implements RepositoryInterface
type RepositoryManager struct {
	repos *Repositories
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) RepositoryInterface {
	return &RepositoryManager{
		repos: NewRepositories(db),
	}
}

// GetRepositories returns all repositories
func (r *RepositoryManager) GetRepositories() *Repositories {
	return r.repos
}

type Repository struct {
	db *gorm.DB
	//cacheDb *redis.Client
}

var repository *Repository

func InitRepository(infra *infra.Infra) *Repository {
	repository = &Repository{
		db: infra.Postgres.DB,
		//cacheDb: infra.Redis.Client,
	}
	if repository.db == nil {
		panic("database connection is nil")
	}
	return repository
}

func GetRepository() *Repository {
	if repository == nil {
		panic("repository not initialized")
	}
	return repository
}
