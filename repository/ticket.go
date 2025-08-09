package repository

import (
	"fmt"
	"github.com/tnqbao/gau-kanban-service/entity"
)

func (r *Repository) generateTicketNumber() (string, error) {
	var count int64
	err := r.db.Model(&entity.Ticket{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("TASK-%04d", count+1), nil
}

func (r *Repository) CreateTicket(ticket *entity.Ticket) error {
	// Generate ticket number nếu chưa có
	if ticket.TicketNo == "" {
		ticketNo, err := r.generateTicketNumber()
		if err != nil {
			return err
		}
		ticket.TicketNo = ticketNo
	}

	// Nếu position chưa được set, đặt ticket ở cuối column
	if ticket.Position == 0 {
		maxPosition, err := r.getMaxPositionInColumn(ticket.ColumnID)
		if err != nil {
			return err
		}
		ticket.Position = maxPosition + 1
	}

	return r.db.Create(ticket).Error
}

// getMaxPositionInColumn lấy position cao nhất trong column
func (r *Repository) getMaxPositionInColumn(columnID string) (int, error) {
	var maxPosition int
	err := r.db.Model(&entity.Ticket{}).
		Where("column_id = ?", columnID).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

func (r *Repository) GetAllTicket() ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *Repository) GetTicketByID(id string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *Repository) GetTicketByColumnID(columnID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Where("column_id = ?", columnID).Order("position ASC, created_at ASC").Find(&tickets).Error
	return tickets, err
}

func (r *Repository) UpdateTicket(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *Repository) DeleteTicket(id string) error {
	return r.db.Delete(&entity.Ticket{}, "id = ?", id).Error
}

func (r *Repository) MoveTicketToColumn(ticketID, columnID string) error {
	return r.db.Model(&entity.Ticket{}).Where("id = ?", ticketID).Update("column_id", columnID).Error
}

// MoveTicketToColumnWithPosition di chuyển ticket sang column khác với position cụ thể
func (r *Repository) MoveTicketToColumnWithPosition(ticketID, newColumnID string, newPosition int) error {
	// Lấy thông tin ticket hiện tại
	ticket, err := r.GetTicketByID(ticketID)
	if err != nil {
		return err
	}

	oldColumnID := ticket.ColumnID
	oldPosition := ticket.Position

	// Bắt đầu transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if oldColumnID == newColumnID {
		// Nếu di chuyển trong cùng column, dùng logic UpdateTicketPosition
		if newPosition < oldPosition {
			err = tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position >= ? AND position < ? AND id != ?",
					newColumnID, newPosition, oldPosition, ticketID).
				Update("position", tx.Raw("position + 1")).Error
		} else if newPosition > oldPosition {
			err = tx.Model(&entity.Ticket{}).
				Where("column_id = ? AND position > ? AND position <= ? AND id != ?",
					newColumnID, oldPosition, newPosition, ticketID).
				Update("position", tx.Raw("position - 1")).Error
		}
	} else {
		// Di chuyển giữa các column khác nhau

		// 1. Cập nhật position trong column cũ (đóng khoảng trống)
		err = tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position > ?", oldColumnID, oldPosition).
			Update("position", tx.Raw("position - 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 2. Cập nhật position trong column mới (tạo chỗ trống)
		err = tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position >= ?", newColumnID, newPosition).
			Update("position", tx.Raw("position + 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	// 3. Cập nhật ticket với column và position mới
	err = tx.Model(&entity.Ticket{}).
		Where("id = ?", ticketID).
		Updates(map[string]interface{}{
			"column_id": newColumnID,
			"position":  newPosition,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) UpdateTicketPosition(ticketID string, newPosition int) error {
	// Lấy thông tin ticket hiện tại
	ticket, err := r.GetTicketByID(ticketID)
	if err != nil {
		return err
	}

	oldPosition := ticket.Position
	columnID := ticket.ColumnID

	// Bắt đầu transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Nếu di chuyển lên (newPosition < oldPosition)
	if newPosition < oldPosition {
		// Các tickets có position từ newPosition đến oldPosition-1 sẽ tăng position lên 1
		err = tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position >= ? AND position < ? AND id != ?",
				columnID, newPosition, oldPosition, ticketID).
			Update("position", tx.Raw("position + 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if newPosition > oldPosition {
		// Nếu di chuyển xuống (newPosition > oldPosition)
		// Các tickets có position từ oldPosition+1 đến newPosition sẽ giảm position xuống 1
		err = tx.Model(&entity.Ticket{}).
			Where("column_id = ? AND position > ? AND position <= ? AND id != ?",
				columnID, oldPosition, newPosition, ticketID).
			Update("position", tx.Raw("position - 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Cập nhật position của ticket được di chuyển
	err = tx.Model(&entity.Ticket{}).
		Where("id = ?", ticketID).
		Update("position", newPosition).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *Repository) GetTicketWithAssignments(ticketID string) (*entity.Ticket, []entity.TaskAssignment, error) {
	var ticket entity.Ticket
	var assignments []entity.TaskAssignment

	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Where("ticket_id = ?", ticketID).Find(&assignments).Error
	if err != nil {
		return &ticket, nil, err
	}

	return &ticket, assignments, nil
}

func (r *Repository) GetTicketWithLabels(ticketID string) (*entity.Ticket, []entity.Label, error) {
	var ticket entity.Ticket
	var labels []entity.Label

	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, nil, err
	}

	err = r.db.Table("labels").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticketID).
		Find(&labels).Error

	if err != nil {
		return &ticket, nil, err
	}

	return &ticket, labels, nil
}

func (r *Repository) SearchTicket(query string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	searchPattern := "%" + query + "%"
	err := r.db.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").Find(&tickets).Error
	return tickets, err
}

func (r *Repository) GetTicketDetail(ticketID string) (*TicketDetailDTO, error) {
	var ticket entity.Ticket
	var column entity.Column

	// Lấy thông tin ticket
	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, err
	}

	// Lấy thông tin column
	err = r.db.Where("id = ?", ticket.ColumnID).First(&column).Error
	if err != nil {
		return nil, err
	}

	// Lấy labels chi tiết
	var labels []LabelDTO
	err = r.db.Table("labels").
		Select("labels.id, labels.name, labels.color").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticket.ID).
		Scan(&labels).Error
	if err != nil {
		labels = []LabelDTO{}
	}

	// Lấy assignees chi tiết
	var assignees []AssigneeDTO
	err = r.db.Table("task_assignments").
		Select("task_assignments.id, task_assignments.user_id, task_assignments.user_full_name as user_full_name, task_assignments.assigned_at").
		Where("ticket_id = ?", ticket.ID).
		Scan(&assignees).Error
	if err != nil {
		assignees = []AssigneeDTO{}
	}

	// Lấy comments chi tiết
	var comments []CommentDTO
	err = r.db.Table("ticket_comments").
		Select("id, user_id, content, created_at").
		Where("ticket_id = ?", ticket.ID).
		Order("created_at ASC").
		Scan(&comments).Error
	if err != nil {
		comments = []CommentDTO{}
	}

	// Xác định completed
	completed := column.Name == "DONE" || column.Name == "COMPLETED"

	ticketDTO := TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		TicketID:    ticket.ID,
		Labels:      labels,
		Assignees:   assignees,
		Comments:    comments,
		Completed:   completed,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	columnDTO := ColumnDTO{
		ID:        column.ID,
		Name:      column.Name,
		Position:  column.Position,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
	}

	result := &TicketDetailDTO{
		Ticket:    ticketDTO,
		Column:    columnDTO,
		Labels:    labels,
		Assignees: assignees,
		Comments:  comments,
	}

	return result, nil
}

func (r *Repository) GetTicketWithAllRelations(ticketID string) (*TicketDTO, error) {
	var ticket entity.Ticket
	var column entity.Column

	// Lấy thông tin ticket
	err := r.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		return nil, err
	}

	// Lấy thông tin column để xác định completed
	err = r.db.Where("id = ?", ticket.ColumnID).First(&column).Error
	if err != nil {
		return nil, err
	}

	// Lấy labels chi tiết
	var labels []LabelDTO
	err = r.db.Table("labels").
		Select("labels.id, labels.name, labels.color").
		Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
		Where("ticket_labels.ticket_id = ?", ticket.ID).
		Scan(&labels).Error
	if err != nil {
		labels = []LabelDTO{}
	}

	// Lấy assignees chi tiết
	var assignees []AssigneeDTO
	err = r.db.Table("task_assignments").
		Select("task_assignments.id, task_assignments.user_id, task_assignments.user_full_name as user_full_name, task_assignments.assigned_at").
		Where("ticket_id = ?", ticket.ID).
		Scan(&assignees).Error
	if err != nil {
		assignees = []AssigneeDTO{}
	}

	// Lấy comments chi tiết
	var comments []CommentDTO
	err = r.db.Table("ticket_comments").
		Select("id, user_id, content, created_at").
		Where("ticket_id = ?", ticket.ID).
		Order("created_at ASC").
		Scan(&comments).Error
	if err != nil {
		comments = []CommentDTO{}
	}

	// Xác định completed
	completed := column.Name == "DONE" || column.Name == "COMPLETED"

	result := &TicketDTO{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		TicketID:    ticket.ID,
		Labels:      labels,
		Assignees:   assignees,
		Comments:    comments,
		Completed:   completed,
		DueDate:     ticket.DueDate,
		Priority:    ticket.Priority,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
	}

	return result, nil
}
