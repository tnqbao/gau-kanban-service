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
	}

	return response, nil
}
