package repository

import (
	"fmt"

	"github.com/tnqbao/gau-kanban-service/entity"
)

// Column methods
func (r *Repository) GetColumnsByBoardID(boardID string) ([]entity.Column, error) {
	var columns []entity.Column
	err := r.db.Where("board_id = ?", boardID).Order("position ASC").Find(&columns).Error
	if err != nil {
		return nil, err
	}

	// Initialize empty tickets array for each column
	for i := range columns {
		columns[i].Tickets = []entity.Ticket{}
	}

	return columns, nil
}

func (r *Repository) GetColumnsWithTicketsByBoardID(boardID string) ([]entity.Column, error) {
	var columns []entity.Column
	err := r.db.Preload("Tickets").Where("board_id = ?", boardID).Order("position ASC").Find(&columns).Error
	return columns, err
}

func (r *Repository) CreateColumn(column *entity.Column) error {
	return r.db.Create(column).Error
}

func (r *Repository) GetColumnByID(id string) (*entity.Column, error) {
	var column entity.Column
	err := r.db.Where("id = ?", id).First(&column).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *Repository) GetColumnWithTicketsByID(id string) (*entity.Column, error) {
	var column entity.Column
	err := r.db.Preload("Tickets").Where("id = ?", id).First(&column).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *Repository) GetMaxOrderByBoardID(boardID string) (int, error) {
	var maxOrder int
	err := r.db.Model(&entity.Column{}).Where("board_id = ?", boardID).Select("COALESCE(MAX(position), 0)").Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *Repository) GetMinOrderByBoardID(boardID string) (int, error) {
	var minOrder int
	err := r.db.Model(&entity.Column{}).Where("board_id = ?", boardID).Select("COALESCE(MIN(position), 0)").Scan(&minOrder).Error
	return minOrder, err
}

func (r *Repository) UpdateColumn(column *entity.Column) error {
	return r.db.Save(column).Error
}

func (r *Repository) DeleteColumn(columnID string) error {
	return r.db.Delete(&entity.Column{}, "id = ?", columnID).Error
}

// Reorder methods - Adaptive Spacing Strategy (Optimal)
func (r *Repository) GetColumnOrderPosition(boardID string, position string) (int, error) {
	var order int

	switch {
	case position == "first":
		// Get minimum order and use adaptive spacing
		var minOrder int
		err := r.db.Model(&entity.Column{}).
			Where("board_id = ?", boardID).
			Select("COALESCE(MIN(position), 1000)").
			Scan(&minOrder).Error
		if err != nil {
			return 0, err
		}

		// Adaptive spacing: use larger gap if there's room, smaller if tight
		if minOrder >= 100 {
			order = minOrder - 100 // Large gap for plenty of fractional space
		} else {
			order = minOrder - 10 // Smaller gap but still workable
		}
		return order, nil

	case position == "last":
		// Get maximum order and use adaptive spacing
		var maxOrder int
		err := r.db.Model(&entity.Column{}).
			Where("board_id = ?", boardID).
			Select("COALESCE(MAX(position), 0)").
			Scan(&maxOrder).Error
		if err != nil {
			return 0, err
		}

		// Adaptive spacing: maintain consistency with first logic
		if maxOrder >= 1000 {
			order = maxOrder + 100 // Large gap for plenty of fractional space
		} else {
			order = maxOrder + 10 // Smaller gap but still workable
		}
		return order, nil

	case len(position) > 6 && position[:6] == "after:":
		targetID := position[6:]
		return r.getOrderAfterColumn(boardID, targetID)

	case len(position) > 7 && position[:7] == "before:":
		targetID := position[7:]
		return r.getOrderBeforeColumn(boardID, targetID)

	default:
		return 0, fmt.Errorf("invalid position format: %s", position)
	}
}

func (r *Repository) getOrderAfterColumn(boardID string, targetID string) (int, error) {
	var targetOrder, nextOrder int

	// Get target column order
	err := r.db.Model(&entity.Column{}).
		Where("id = ? AND board_id = ?", targetID, boardID).
		Select("position").
		Scan(&targetOrder).Error
	if err != nil {
		return 0, err
	}

	// Get next column order (or use max + 1000 if it's the last)
	err = r.db.Model(&entity.Column{}).
		Where("board_id = ? AND position > ?", boardID, targetOrder).
		Select("COALESCE(MIN(position), ?)", targetOrder+1000).
		Scan(&nextOrder).Error
	if err != nil {
		return 0, err
	}

	// Calculate middle position
	newOrder := (targetOrder + nextOrder) / 2
	if newOrder == targetOrder {
		newOrder = targetOrder + 500 // Ensure different order
	}

	return newOrder, nil
}

func (r *Repository) getOrderBeforeColumn(boardID string, targetID string) (int, error) {
	var targetOrder, prevOrder int

	// Get target column order
	err := r.db.Model(&entity.Column{}).
		Where("id = ? AND board_id = ?", targetID, boardID).
		Select("position").
		Scan(&targetOrder).Error
	if err != nil {
		return 0, err
	}

	// Get previous column order (or use min - 1000 if it's the first)
	err = r.db.Model(&entity.Column{}).
		Where("board_id = ? AND position < ?", boardID, targetOrder).
		Select("COALESCE(MAX(position), ?)", targetOrder-1000).
		Scan(&prevOrder).Error
	if err != nil {
		return 0, err
	}

	// Calculate middle position
	newOrder := (prevOrder + targetOrder) / 2
	if newOrder == targetOrder {
		newOrder = targetOrder - 500 // Ensure different order
	}

	return newOrder, nil
}

func (r *Repository) UpdateColumnOrder(columnID string, newOrder int) error {
	return r.db.Model(&entity.Column{}).
		Where("id = ?", columnID).
		Update("position", newOrder).Error
}
