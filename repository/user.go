package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// User methods
func (r *Repository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetOrCreateUser(userID string, fullName string) (*entity.User, error) {
	// Try to get existing user first
	user, err := r.GetUserByID(userID)
	if err == nil {
		return user, nil
	}

	// If user doesn't exist, create new one
	newUser := &entity.User{
		ID:       userID,
		FullName: fullName,
	}

	if err := r.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *Repository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) GetAllUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}
