package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Member repository methods
func (r *Repository) CreateMember(member *entity.Member) error {
	return r.db.Create(member).Error
}

func (r *Repository) GetAllMembers() ([]entity.Member, error) {
	var members []entity.Member
	err := r.db.Order("created_at DESC").Find(&members).Error
	return members, err
}

func (r *Repository) GetMemberByID(id string) (*entity.Member, error) {
	var member entity.Member
	err := r.db.Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *Repository) GetMembersByBoardID(boardID string) ([]entity.Member, error) {
	var members []entity.Member
	err := r.db.Where("board_id = ?", boardID).Order("full_name ASC").Find(&members).Error
	return members, err
}

func (r *Repository) UpdateMember(member *entity.Member) error {
	return r.db.Save(member).Error
}

func (r *Repository) DeleteMember(id string) error {
	return r.db.Delete(&entity.Member{}, "id = ?", id).Error
}

func (r *Repository) DeleteMembersByBoardID(boardID string) error {
	return r.db.Where("board_id = ?", boardID).Delete(&entity.Member{}).Error
}
