package repository

import (
	"github.com/tnqbao/gau-kanban-service/entity"
)

func (r *Repository) CreateColumn(column *entity.Column) error {
	return r.db.Create(column).Error
}

func (r *Repository) GetAllColumn() ([]entity.Column, error) {
	var columns []entity.Column
	err := r.db.Order("position ASC").Find(&columns).Error
	return columns, err
}

func (r *Repository) GetColumnByID(id string) (*entity.Column, error) {
	var column entity.Column
	err := r.db.Where("id = ?", id).First(&column).Error
	if err != nil {
		return nil, err
	}
	return &column, nil
}

func (r *Repository) UpdateColumn(column *entity.Column) error {
	return r.db.Save(column).Error
}

func (r *Repository) DeleteColumn(id string) error {
	return r.db.Delete(&entity.Column{}, "id = ?", id).Error
}

func (r *Repository) UpdateColumnPosition(id string, position int) error {
	return r.db.Model(&entity.Column{}).Where("id = ?", id).Update("position", position).Error
}

func (r *Repository) GetAllColumnWithTickets() ([]ColumnWithTicketsDTO, error) {
	var columns []entity.Column
	var result []ColumnWithTicketsDTO

	// Lấy tất cả columns theo thứ tự position
	err := r.db.Order("position ASC").Find(&columns).Error
	if err != nil {
		return nil, err
	}

	// Duyệt qua từng column và lấy tickets của nó
	for _, column := range columns {
		columnDTO := ColumnWithTicketsDTO{
			ID:      column.ID,
			Title:   column.Title,
			Order:   column.Position,
			Tickets: []TicketDTO{},
		}

		// Lấy tickets của column này
		var tickets []entity.Ticket
		err := r.db.Where("column_id = ?", column.ID).Order("created_at ASC").Find(&tickets).Error
		if err != nil {
			return nil, err
		}

		// Convert tickets thành TicketDTO với thông tin cơ bản
		for _, ticket := range tickets {
			// Lấy labels của ticket
			var labels []LabelDTO
			err := r.db.Table("labels").
				Select("labels.id, labels.name, labels.color").
				Joins("JOIN ticket_labels ON labels.id = ticket_labels.label_id").
				Where("ticket_labels.ticket_id = ?", ticket.ID).
				Scan(&labels).Error
			if err != nil {
				labels = []LabelDTO{} // Nếu có lỗi thì để array rỗng
			}

			// Lấy assignees của ticket
			var assignees []AssigneeDTO
			err = r.db.Table("task_assignments").
				Select("task_assignments.id, task_assignments.user_id, task_assignments.user_full_name as user_full_name, task_assignments.assigned_at").
				Where("ticket_id = ?", ticket.ID).
				Scan(&assignees).Error
			if err != nil {
				assignees = []AssigneeDTO{} // Nếu có lỗi thì để array rỗng
			}

			// Xác định completed dựa trên column name
			completed := column.Title == "DONE" || column.Title == "COMPLETED"

			ticketDTO := TicketDTO{
				ID:          ticket.ID,
				Title:       ticket.Title,
				Description: ticket.Description,
				TicketID:    ticket.ID,
				Labels:      labels,
				Assignees:   assignees,
				Comments:    []CommentDTO{}, // Comments rỗng cho performance
				Completed:   completed,
				DueDate:     ticket.DueDate,
				Priority:    ticket.Priority,
				CreatedAt:   ticket.CreatedAt,
				UpdatedAt:   ticket.UpdatedAt,
			}

			columnDTO.Tickets = append(columnDTO.Tickets, ticketDTO)
		}

		result = append(result, columnDTO)
	}

	return result, nil
}

func (r *Repository) GetAllColumnWithFullTicketDetails() ([]ColumnWithTicketsDTO, error) {
	var columns []entity.Column
	var result []ColumnWithTicketsDTO

	// Lấy tất cả columns theo thứ tự position
	err := r.db.Order("position ASC").Find(&columns).Error
	if err != nil {
		return nil, err
	}

	// Duyệt qua từng column và lấy tickets với đầy đủ thông tin
	for _, column := range columns {
		columnDTO := ColumnWithTicketsDTO{
			ID:      column.ID,
			Title:   column.Title,
			Order:   column.Position,
			Tickets: []TicketDTO{},
		}

		// Lấy tickets của column này
		var tickets []entity.Ticket
		err := r.db.Where("column_id = ?", column.ID).Order("created_at ASC").Find(&tickets).Error
		if err != nil {
			return nil, err
		}

		// Convert tickets thành TicketDTO với đầy đủ thông tin
		for _, ticket := range tickets {
			// Lấy labels chi tiết
			var labels []LabelDTO
			err := r.db.Table("labels").
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
			completed := column.Title == "DONE" || column.Title == "COMPLETED"

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

			columnDTO.Tickets = append(columnDTO.Tickets, ticketDTO)
		}

		result = append(result, columnDTO)
	}

	return result, nil
}
