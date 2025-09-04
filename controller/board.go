package controller

import (
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

	// Get or create user in database
	user, err := ctrl.Repository.GetOrCreateUser(userIDStr, req.FullName)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Failed to get or create user")
		utils.JSON500(c, "Failed to process user information")
		return
	}

	// Create board
	board := &entity.Board{
		Title:       req.Title,
		Description: req.Description,
		OwnerID:     user.ID,
	}

	if err := ctrl.Repository.CreateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Failed to create board")
		utils.JSON500(c, "Failed to create board")
		return
	}

	// Add owner as member with owner role
	member := &entity.Member{
		UserID:  user.ID,
		BoardID: board.ID,
		Role:    "owner",
	}

	if err := ctrl.Repository.CreateMember(member); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Board] Failed to add owner as member")
		// Log error but don't fail the board creation
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create Board] Board created but owner not added as member: %s", board.ID)
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Board] Board created successfully: %s", board.ID)
	utils.JSON200(c, gin.H{
		"message": "Board created successfully",
		"data":    board,
	})
}

// GetBoards lấy danh sách boards mà user hiện tại là member
func (ctrl *Controller) GetBoards(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Boards] Get boards request received")

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Get Boards] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get boards where user is a member
	boards, err := ctrl.Repository.GetBoardsByUserID(userIDStr)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Boards] Failed to get boards for user: %s", userIDStr)
		utils.JSON500(c, "Failed to get boards")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Boards] Successfully retrieved %d boards for user: %s", len(boards), userIDStr)
	utils.JSON200(c, gin.H{
		"message": "Boards retrieved successfully",
		"data":    boards,
	})
}

// GetBoardByID lấy thông tin board với columns và tickets
func (ctrl *Controller) GetBoardByID(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Board By ID] Get board by ID request received for ID: %s", boardID)

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Get Board By ID] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Check if user has access to this board (is member or owner)
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember {
		// Also check if user is owner
		board, err := ctrl.Repository.GetBoardByID(boardID)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Board not found: %s", boardID)
			utils.JSON404(c, "Board not found")
			return
		}

		if board.OwnerID != userIDStr {
			ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Get Board By ID] User %s does not have access to board %s", userIDStr, boardID)
			utils.JSON403(c, "Access denied: You are not a member of this board")
			return
		}
	}

	// Get board with columns and tickets
	board, columns, err := ctrl.Repository.GetBoardWithColumnsAndTickets(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Failed to get board with columns and tickets: %s", boardID)
		utils.JSON500(c, "Failed to get board details")
		return
	}

	// Get members of the board
	members, err := ctrl.Repository.GetMembersByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Failed to get board members: %s", boardID)
		utils.JSON500(c, "Failed to get board members")
		return
	}

	// Get labels of the board
	labels, err := ctrl.Repository.GetLabelsByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Board By ID] Failed to get board labels: %s", boardID)
		utils.JSON500(c, "Failed to get board labels")
		return
	}

	// If no columns, return empty array
	if columns == nil {
		columns = []entity.Column{}
	}

	// If no members, return empty array
	if members == nil {
		members = []entity.Member{}
	}

	// If no labels, return empty array
	if labels == nil {
		labels = []entity.Label{}
	}

	// Ensure each column has tickets array (even if empty)
	for i := range columns {
		if columns[i].Tickets == nil {
			columns[i].Tickets = []entity.Ticket{}
		}
	}

	// Set the relationships directly on the board object
	board.Columns = columns
	board.Members = members
	board.Labels = labels

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Board By ID] Successfully retrieved board %s with %d columns, %d members, and %d labels", boardID, len(columns), len(members), len(labels))
	utils.JSON200(c, gin.H{
		"message": "Board retrieved successfully",
		"data":    board,
	})
}

// UpdateBoard cập nhật thông tin board
func (ctrl *Controller) UpdateBoard(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Board] Update board request received for ID: %s", boardID)

	var req UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Update Board] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get board and check ownership/membership
	board, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Check if user is owner or member
	isOwner := board.OwnerID == userIDStr
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isOwner && !isMember {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Update Board] User %s does not have access to board %s", userIDStr, boardID)
		utils.JSON403(c, "Access denied: You are not a member of this board")
		return
	}

	// Update board fields if provided
	if req.Title != "" {
		board.Title = req.Title
	}
	if req.Description != "" {
		board.Description = req.Description
	}

	if err := ctrl.Repository.UpdateBoard(board); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Board] Failed to update board: %s", boardID)
		utils.JSON500(c, "Failed to update board")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Board] Board updated successfully: %s", boardID)
	utils.JSON200(c, gin.H{
		"message": "Board updated successfully",
		"data":    board,
	})
}

// DeleteBoard xóa board (chỉ owner mới được xóa)
func (ctrl *Controller) DeleteBoard(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Board] Delete board request received for ID: %s", boardID)

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Delete Board] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get board and check ownership
	board, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Board] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Only owner can delete board
	if board.OwnerID != userIDStr {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Delete Board] User %s is not owner of board %s", userIDStr, boardID)
		utils.JSON403(c, "Access denied: Only board owner can delete the board")
		return
	}

	if err := ctrl.Repository.DeleteBoard(boardID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Board] Failed to delete board: %s", boardID)
		utils.JSON500(c, "Failed to delete board")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Board] Board deleted successfully: %s", boardID)
	utils.JSON200(c, gin.H{
		"message": "Board deleted successfully",
	})
}

// ArchiveBoard archive board (chỉ owner mới được archive)
func (ctrl *Controller) ArchiveBoard(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Archive Board] Archive board request received for ID: %s", boardID)

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Archive Board] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get board and check ownership
	board, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Archive Board] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Only owner can archive board
	if board.OwnerID != userIDStr {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Archive Board] User %s is not owner of board %s", userIDStr, boardID)
		utils.JSON403(c, "Access denied: Only board owner can archive the board")
		return
	}

	// Check if already archived
	if board.Archived {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Archive Board] Board %s is already archived", boardID)
		utils.JSON400(c, "Board is already archived")
		return
	}

	if err := ctrl.Repository.ArchiveBoard(boardID); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Archive Board] Failed to archive board: %s", boardID)
		utils.JSON500(c, "Failed to archive board")
		return
	}

	// Get updated board for response
	board.Archived = true

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Archive Board] Board archived successfully: %s", boardID)
	utils.JSON200(c, gin.H{
		"message": "Board archived successfully",
		"data":    board,
	})
}

// SearchTickets searches for tickets by title within a board
func (ctrl *Controller) SearchTickets(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Search Tickets] Search tickets request received")

	boardID := c.Param("id")
	if boardID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Search Tickets] Board ID is required")
		utils.JSON400(c, "Board ID is required")
		return
	}

	var req SearchTicketsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Invalid query parameters")
		utils.JSON400(c, "Invalid query parameters: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Search Tickets] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get board and verify it exists
	board, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Check if user has access to this board
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember && board.OwnerID != userIDStr {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Search Tickets] User %s does not have access to board %s", userIDStr, boardID)
		utils.JSON403(c, "Access denied: You are not a member of this board")
		return
	}

	// Get all columns for this board
	columns, err := ctrl.Repository.GetColumnsByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Failed to get columns for board: %s", boardID)
		utils.JSON500(c, "Failed to get board columns")
		return
	}

	// For each column, get tickets that match the search criteria
	for i := range columns {
		tickets, err := ctrl.Repository.SearchTicketsByColumnAndTitle(columns[i].ID, req.Search)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Failed to search tickets for column: %s", columns[i].ID)
			// Don't fail the request, just set empty array
			tickets = []entity.Ticket{}
		}
		columns[i].Tickets = tickets
	}

	// Get board members
	members, err := ctrl.Repository.GetMembersByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Failed to get members")
		members = []entity.Member{}
	}

	// Get board labels
	labels, err := ctrl.Repository.GetLabelsByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Search Tickets] Failed to get labels")
		labels = []entity.Label{}
	}

	// Set relationships
	board.Columns = columns
	board.Members = members
	board.Labels = labels

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Search Tickets] Search completed successfully for board: %s with search term: %s", boardID, req.Search)
	utils.JSON200(c, gin.H{
		"message": "Tickets searched successfully",
		"data":    board,
	})
}

// FilterTickets filters tickets by assignee, label, or status within a board
func (ctrl *Controller) FilterTickets(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Filter Tickets] Filter tickets request received")

	boardID := c.Param("id")
	if boardID == "" {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Filter Tickets] Board ID is required")
		utils.JSON400(c, "Board ID is required")
		return
	}

	var req FilterTicketsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Invalid query parameters")
		utils.JSON400(c, "Invalid query parameters: "+err.Error())
		return
	}

	// Extract user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, nil, "[Filter Tickets] User ID not found in context")
		utils.JSON401(c, "Authentication required")
		return
	}

	userIDStr := userID.(uuid.UUID).String()

	// Get board and verify it exists
	board, err := ctrl.Repository.GetBoardByID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Board not found: %s", boardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Check if user has access to this board
	isMember, err := ctrl.Repository.IsMemberOfBoard(userIDStr, boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Failed to check member access")
		utils.JSON500(c, "Failed to check board access")
		return
	}

	if !isMember && board.OwnerID != userIDStr {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Filter Tickets] User %s does not have access to board %s", userIDStr, boardID)
		utils.JSON403(c, "Access denied: You are not a member of this board")
		return
	}

	// Get all columns for this board
	columns, err := ctrl.Repository.GetColumnsByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Failed to get columns for board: %s", boardID)
		utils.JSON500(c, "Failed to get board columns")
		return
	}

	// For each column, get tickets that match the filter criteria
	for i := range columns {
		tickets, err := ctrl.Repository.FilterTicketsByColumn(columns[i].ID, req.Assignee, req.Label, req.Status)
		if err != nil {
			ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Failed to filter tickets for column: %s", columns[i].ID)
			// Don't fail the request, just set empty array
			tickets = []entity.Ticket{}
		}
		columns[i].Tickets = tickets
	}

	// Get board members
	members, err := ctrl.Repository.GetMembersByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Failed to get members")
		members = []entity.Member{}
	}

	// Get board labels
	labels, err := ctrl.Repository.GetLabelsByBoardID(boardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Filter Tickets] Failed to get labels")
		labels = []entity.Label{}
	}

	// Set relationships
	board.Columns = columns
	board.Members = members
	board.Labels = labels

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Filter Tickets] Filter completed successfully for board: %s", boardID)
	utils.JSON200(c, gin.H{
		"message": "Tickets filtered successfully",
		"data":    board,
	})
}
