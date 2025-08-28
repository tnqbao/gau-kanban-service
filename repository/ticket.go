package repository

import (
	"fmt"
	"github.com/tnqbao/gau-kanban-service/entity"
)

// GenerateTicketNumber tạo ticket number theo format TASK-XXXX
func (r *Repository) GenerateTicketNumber() (string, error) {
	var nextVal int
	err := r.db.Raw("SELECT nextval('ticket_number_seq')").Scan(&nextVal).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TASK-%04d", nextVal), nil
}

func (r *Repository) CreateTicket(ticket *entity.Ticket) error {
	// Generate ticket number nếu chưa có
	if ticket.TicketNo == "" {
		ticketNo, err := r.GenerateTicketNumber()
		if err != nil {
			return err
		}
		ticket.TicketNo = ticketNo
	}

	// Nếu position chưa được set, đặt ticket ở cuối column
	if ticket.Position == 0 {
		maxPosition, err := r.GetMaxTicketPositionInColumn(ticket.ColumnID)
		if err != nil {
			return err
		}
		ticket.Position = maxPosition + 1
	}

	return r.db.Create(ticket).Error
}

// GetMaxTicketPositionInColumn lấy position cao nhất trong column
func (r *Repository) GetMaxTicketPositionInColumn(columnID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *Repository) GetAllTickets() ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Order("position ASC, created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *Repository) GetTicketsByColumnID(columnID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Where("column_id = ?", columnID).
		Order("position ASC, created_at DESC").
		Find(&tickets).Error
	return tickets, err
}

func (r *Repository) GetTicketByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.First(&ticket, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) UpdateTicket(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *Repository) DeleteTicket(id string) error {
	// Xóa các assignments và checklists liên quan trước
	if err := r.DeleteAssignmentsByTicketID(id); err != nil {
		return err
	}
	if err := r.DeleteChecklistsByTicketID(id); err != nil {
		return err
	}

	return r.db.Delete(&entity.Ticket{}, "id = ?", id).Error
}

// UpdateTicketPosition cập nhật vị trí ticket trong column
func (r *Repository) UpdateTicketPosition(ticketID, columnID string, newPosition int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Lấy ticket hiện tại
	var currentTicket entity.Ticket
	if err := tx.First(&currentTicket, "id = ?", ticketID).Error; err != nil {
		tx.Rollback()
		return err
	}

	oldPosition := currentTicket.Position
	oldColumnID := currentTicket.ColumnID

	// Nếu di chuyển trong cùng column
	if oldColumnID == columnID {
		if oldPosition == newPosition {
			tx.Commit()
			return nil // Không có thay đổi
		}

		if oldPosition < newPosition {
			// Di chuyển xuống: giảm position của các ticket từ oldPosition+1 đến newPosition
			if err := tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position > ? AND position <= ?", columnID, oldPosition, newPosition).
				Update("position", r.db.Raw("position - 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Di chuyển lên: tăng position của các ticket từ newPosition đến oldPosition-1
			if err := tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position >= ? AND position < ?", columnID, newPosition, oldPosition).
				Update("position", r.db.Raw("position + 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		// Di chuyển sang column khác
		// Giảm position của các ticket sau vị trí cũ trong column cũ
		if err := tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position > ?", oldColumnID, oldPosition).
			Update("position", r.db.Raw("position - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Tăng position của các ticket từ newPosition trở đi trong column mới
		if err := tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position >= ?", columnID, newPosition).
			Update("position", r.db.Raw("position + 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Cập nhật ticket hiện tại
	if err := tx.Model(&currentTicket).Updates(map[string]interface{}{
		"column_id": columnID,
		"position":  newPosition,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// MoveTicketToColumn di chuyển ticket sang column khác (đặt ở cuối)
func (r *Repository) MoveTicketToColumn(ticketID, columnID string) error {
	maxPosition, err := r.GetMaxTicketPositionInColumn(columnID)
	if err != nil {
		return err
	}
	return r.UpdateTicketPosition(ticketID, columnID, maxPosition+1)
}

// MoveTicketToColumnWithPosition di chuyển ticket sang column khác với position cụ thể
func (r *Repository) MoveTicketToColumnWithPosition(ticketID, columnID string, position int) error {
	return r.UpdateTicketPosition(ticketID, columnID, position)
}

// GetTicketWithDetails lấy ticket kèm assignments và checklists
func (r *Repository) GetTicketWithDetails(ticketID string) (*TicketWithDetailsResponse, error) {
	// Lấy ticket
	ticket, err := r.GetTicketByID(ticketID)
	if err != nil {
		return nil, err
	}

	// Lấy assignments
	assignments, err := r.GetAssignmentsByTicketID(ticketID)
	if err != nil {
		return nil, err
	}

	// Lấy checklists
	checklists, err := r.GetChecklistsByTicketID(ticketID)
	if err != nil {
		return nil, err
	}

	// Lấy comments
	var comments []entity.TicketComment
	err = r.db.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// Chuyển đổi sang DTOs
	var assignmentDTOs []AssignmentDTO
	for _, assignment := range assignments {
		assignmentDTOs = append(assignmentDTOs, AssignmentDTO{
			ID:           assignment.ID,
			TicketID:     assignment.TicketID,
			UserID:       assignment.UserID,
			UserFullName: assignment.UserFullName,
			AssignedAt:   assignment.AssignedAt,
		})
	}

	var checklistDTOs []ChecklistDTO
	for _, checklist := range checklists {
		checklistDTOs = append(checklistDTOs, ChecklistDTO{
			ID:        checklist.ID,
			TicketID:  checklist.TicketID,
			Title:     checklist.Title,
			Completed: checklist.Completed,
			Position:  checklist.Position,
			CreatedAt: checklist.CreatedAt,
			UpdatedAt: checklist.UpdatedAt,
		})
	}

	var commentDTOs []CommentDTO
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, CommentDTO{
			ID:        comment.ID,
			TicketID:  comment.TicketID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: "", // TicketComment entity doesn't have UpdatedAt field
		})
	}

	response := &TicketWithDetailsResponse{
		ID:          ticket.ID,
		TicketNo:    ticket.TicketNo,
		ColumnID:    ticket.ColumnID,
		Title:       ticket.Title,
		Description: ticket.Description,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		Position:    ticket.Position,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		Assignments: assignmentDTOs,
		Checklists:  checklistDTOs,
		Comments:    commentDTOs,
	}

	return response, nil
}

// ChangeTicketPosition changes ticket position with advanced handling (within same column or between columns)
func (r *Repository) ChangeTicketPosition(ticketID string, newColumnID string, newPosition int) error {
	// Get current ticket info
	var currentTicket entity.Ticket
	if err := r.db.Where("id = ?", ticketID).First(&currentTicket).Error; err != nil {
		return err
	}

	currentColumnID := currentTicket.ColumnID
	currentPosition := currentTicket.Position

	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	if currentColumnID == newColumnID {
		// Moving within the same column
		if currentPosition == newPosition {
			// No change needed
			return tx.Commit().Error
		}

		if newPosition < currentPosition {
			// Moving up (decrease position number)
			// Shift other tickets down (increase their position)
			if err := tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position >= ? AND position < ? AND id != ?",
					newColumnID, newPosition, currentPosition, ticketID).
				Update("position", r.db.Raw("position + 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Moving down (increase position number)
			// Shift other tickets up (decrease their position)
			if err := tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position > ? AND position <= ? AND id != ?",
					newColumnID, currentPosition, newPosition, ticketID).
				Update("position", r.db.Raw("position - 1")).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// Update current ticket position
		if err := tx.Model(&entity.Ticket{}).
			Where("id = ?", ticketID).
			Update("position", newPosition).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Moving between different columns

		// 1. Shift tickets in old column to fill the gap
		if err := tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position > ?", currentColumnID, currentPosition).
			Update("position", r.db.Raw("position - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 2. Make space in new column by shifting tickets down
		if err := tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position >= ?", newColumnID, newPosition).
			Update("position", r.db.Raw("position + 1")).Error; err != nil {
			tx.Rollback()
			return err
		}

		// 3. Update current ticket with new column and position
		if err := tx.Model(&entity.Ticket{}).
			Where("id = ?", ticketID).
			Updates(map[string]interface{}{
				"column_id": newColumnID,
				"position":  newPosition,
			}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// SearchTickets performs fuzzy search on tickets by title and ticket number
func (r *Repository) SearchTickets(query string, boardID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	// Build the search query
	searchQuery := r.db.Where("title ILIKE ? OR ticket_no ILIKE ?", "%"+query+"%", "%"+query+"%")

	// Add board filter if provided
	if boardID != "" {
		searchQuery = searchQuery.Joins("JOIN columns ON tickets.column_id = columns.id").
			Where("columns.board_id = ?", boardID)
	}

	err := searchQuery.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

// FilterTickets filters tickets by label, assignee, and status
func (r *Repository) FilterTickets(labelID, assigneeID, status, boardID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket

	query := r.db.Table("tickets")

	// Join with columns for board filtering
	if boardID != "" {
		query = query.Joins("JOIN columns ON tickets.column_id = columns.id").
			Where("columns.board_id = ?", boardID)
	}

	// Filter by label
	if labelID != "" {
		query = query.Joins("JOIN ticket_labels ON tickets.id = ticket_labels.ticket_id").
			Where("ticket_labels.label_id = ?", labelID)
	}

	// Filter by assignee
	if assigneeID != "" {
		query = query.Joins("JOIN task_assignments ON tickets.id = task_assignments.ticket_id").
			Where("task_assignments.user_id = ?", assigneeID)
	}

	// Filter by status (column title)
	if status != "" {
		if boardID == "" {
			query = query.Joins("JOIN columns ON tickets.column_id = columns.id")
		}
		query = query.Where("LOWER(columns.title) = LOWER(?)", status)
	}

	err := query.Order("tickets.created_at DESC").Find(&tickets).Error
	return tickets, err
}
