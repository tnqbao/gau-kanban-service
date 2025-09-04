package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// AddMemberToBoard thêm member vào board
func (ctrl *Controller) AddMemberToBoard(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Add Member] Add member to board request received")

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Add Member] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra user có tồn tại trong bảng users không
	user, err := ctrl.Repository.GetUserByID(req.UserID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Add Member] User not found: %s", req.UserID)
		utils.JSON404(c, "User not found")
		return
	}

	// Kiểm tra board có tồn tại không
	board, err := ctrl.Repository.GetBoardByID(req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Add Member] Board not found: %s", req.BoardID)
		utils.JSON404(c, "Board not found")
		return
	}

	// Kiểm tra user đã là member của board chưa
	exists, err := ctrl.Repository.IsMemberExists(req.UserID, req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Add Member] Failed to check member existence")
		utils.JSON500(c, "Failed to check member existence")
		return
	}

	if exists {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Add Member] User %s is already a member of board %s", req.UserID, req.BoardID)
		utils.JSON400(c, "User is already a member of this board")
		return
	}

	// Tạo member mới
	member := &entity.Member{
		UserID:  req.UserID,
		BoardID: req.BoardID,
		Role:    "member", // default role
	}

	if err := ctrl.Repository.CreateMember(member); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Add Member] Failed to add member")
		utils.JSON500(c, "Failed to add member to board")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Add Member] Successfully added user %s to board %s", user.FullName, board.Title)

	// Response với thông tin member đã thêm
	utils.JSON200(c, gin.H{
		"message": "Member added to board successfully",
		"data": gin.H{
			"id":        member.ID,
			"user_id":   member.UserID,
			"board_id":  member.BoardID,
			"role":      member.Role,
			"full_name": user.FullName,
		},
	})
}
