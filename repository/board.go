package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

// Board methods
func (r *Repository) CreateBoard(board *entity.Board) error {
	return r.db.Create(board).Error
}

func (r *Repository) GetBoardByID(id string) (*entity.Board, error) {
	var board entity.Board
	err := r.db.Where("id = ?", id).First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (r *Repository) GetBoardsByUserID(userID string) ([]entity.Board, error) {
	var boards []entity.Board
	err := r.db.Where("(owner_id = ? OR id IN (SELECT board_id FROM members WHERE user_id = ?)) AND archived = false", userID, userID).
		Find(&boards).Error
	return boards, err
}

func (r *Repository) UpdateBoard(board *entity.Board) error {
	return r.db.Save(board).Error
}

func (r *Repository) ArchiveBoard(boardID string) error {
	return r.db.Model(&entity.Board{}).Where("id = ?", boardID).Update("archived", true).Error
}

func (r *Repository) GetBoardWithColumnsAndTickets(boardID string) (*entity.Board, []entity.Column, error) {
	// Get board first
	board, err := r.GetBoardByID(boardID)
	if err != nil {
		return nil, nil, err
	}

	// Get columns with tickets
	columns, err := r.GetColumnsWithTicketsByBoardID(boardID)
	if err != nil {
		return board, nil, err
	}

	return board, columns, nil
}

func (r *Repository) DeleteBoard(boardID string) error {
	return r.db.Delete(&entity.Board{}, "id = ?", boardID).Error
}
