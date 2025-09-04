package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Member methods
func (r *Repository) CreateMember(member *entity.Member) error {
	return r.db.Create(member).Error
}

func (r *Repository) GetMembersByBoardID(boardID string) ([]entity.Member, error) {
	var members []entity.Member
	err := r.db.Where("board_id = ?", boardID).Find(&members).Error
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

func (r *Repository) IsMemberOfBoard(userID, boardID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Member{}).Where("user_id = ? AND board_id = ?", userID, boardID).Count(&count).Error
	return count > 0, err
}

func (r *Repository) DeleteMember(id string) error {
	return r.db.Delete(&entity.Member{}, "id = ?", id).Error
}

func (r *Repository) IsMemberExists(userID, boardID string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.Member{}).Where("user_id = ? AND board_id = ?", userID, boardID).Count(&count).Error
	return count > 0, err
}
