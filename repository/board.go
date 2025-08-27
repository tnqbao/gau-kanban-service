package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Board repository methods
func (r *Repository) CreateBoard(board *entity.Board) error {
	return r.db.Create(board).Error
}

func (r *Repository) GetAllBoards() ([]entity.Board, error) {
	var boards []entity.Board
	err := r.db.Order("created_at DESC").Find(&boards).Error
	return boards, err
}

func (r *Repository) GetBoardByID(id string) (*entity.Board, error) {
	var board entity.Board
	err := r.db.Where("id = ?", id).First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *Repository) UpdateBoard(board *entity.Board) error {
	return r.db.Save(board).Error
}

func (r *Repository) DeleteBoard(id string) error {
	// This will cascade delete all related columns, tickets, etc.
	return r.db.Delete(&entity.Board{}, "id = ?", id).Error
}
