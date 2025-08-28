package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateBoard tạo board mới
func (ctrl *Controller) CreateBoard(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Board] Create new board request received")

	var req CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Create Board] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	board := &entity.Board{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Failed to create board")
		utils.JSON500(c, err.Error())
		return
	}

	// Automatically add the creator as a member of the board
	member := &entity.Member{
		BoardID:   board.ID,
		MemberID:  userIDStr,
		FullName:  req.FullName, // Use fullname from request instead of "Board Creator"
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateMember(member); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Failed to add creator as member")
		// Log error but don't fail the board creation
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create Board] Board created but creator not added as member: %s", board.ID)
	} else {
		ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Board] Creator added as member: %s", member.ID)
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Board] Board created successfully: %s", board.ID)
	utils.JSON200(c, gin.H{
		"message": "Board created successfully",
		"data":    board,
	})
}

// GetBoards lấy tất cả boards
func (ctrl *Controller) GetBoards(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Boards] Get all boards request received")

	boards, err := ctrl.Repository.GetAllBoards()
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Boards] Failed to get boards")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Boards] Successfully retrieved %d boards", len(boards))
	utils.JSON200(c, gin.H{
		"message": "Boards retrieved successfully",
		"data":    boards,
	})
}

// GetBoardByID lấy board theo ID kèm theo tất cả members
func (ctrl *Controller) GetBoardByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Board By ID] Get board by ID request received for ID: %s", id)

	board, err := ctrl.Repository.GetBoardByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Board not found: %s", id)
		utils.JSON404(c, "Board not found")
		return
	}

	// Lấy tất cả members của board
	members, err := ctrl.Repository.GetMembersByBoardID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Failed to get members for board: %s", id)
		utils.JSON500(c, "Failed to get board members")
		return
	}

	// Convert members to DTO
	memberDTOs := make([]MemberDTO, len(members))
	for i, member := range members {
		memberDTOs[i] = MemberDTO{
			ID:        member.ID,
			BoardID:   member.BoardID,
			MemberID:  member.MemberID,
			FullName:  member.FullName,
			CreatedAt: member.CreatedAt,
			UpdatedAt: member.UpdatedAt,
		}
	}

	// Tạo response với board và members
	response := BoardWithMembersResponse{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		Archived:    board.Archived,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
		Members:     memberDTOs,
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Board By ID] Board and %d members retrieved successfully: %s", len(members), id)
	utils.JSON200(c, gin.H{
		"message": "Board and members retrieved successfully",
		"data":    response,
	})
}

// UpdateBoard cập nhật board
func (ctrl *Controller) UpdateBoard(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Board] Update board request received for ID: %s", id)

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	board, err := ctrl.Repository.GetBoardByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Board not found: %s", id)
		utils.JSON404(c, "Board not found")
		return
	}

	// Cập nhật các field nếu có trong request
	if req.Title != "" {
		board.Title = req.Title
	}
	if req.Description != "" {
		board.Description = req.Description
	}

	board.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Failed to update board: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Board] Board updated successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Board updated successfully",
		"data":    board,
	})
}

// DeleteBoard xóa board
func (ctrl *Controller) DeleteBoard(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Board] Delete board request received for ID: %s", id)

	// Kiểm tra board có tồn tại không
	_, err := ctrl.Repository.GetBoardByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Board] Board not found: %s", id)
		utils.JSON404(c, "Board not found")
		return
	}

	if err := ctrl.Repository.DeleteBoard(id); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Board] Failed to delete board: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Board] Board deleted successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Board deleted successfully",
	})
}

// ArchiveBoard archives a board (soft delete)
func (ctrl *Controller) ArchiveBoard(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Archive Board] Archive board request received for ID: %s", id)

	board, err := ctrl.Repository.GetBoardByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Archive Board] Board not found: %s", id)
		utils.JSON404(c, "Board not found")
		return
	}

	board.Archived = true
	board.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Archive Board] Failed to archive board: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Archive Board] Board archived successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Board archived successfully",
		"data":    board,
	})
}

// RestoreBoard restores an archived board
func (ctrl *Controller) RestoreBoard(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Restore Board] Restore board request received for ID: %s", id)

	board, err := ctrl.Repository.GetBoardByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Restore Board] Board not found: %s", id)
		utils.JSON404(c, "Board not found")
		return
	}

	board.Archived = false
	board.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Restore Board] Failed to restore board: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Restore Board] Board restored successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Board restored successfully",
		"data":    board,
	})
}
