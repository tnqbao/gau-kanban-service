package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateMember tạo member mới
func (ctrl *Controller) CreateMember(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Member] Create new member request received")

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Member] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	// Kiểm tra board có tồn tại không
	_, err := ctrl.Repository.GetBoardByID(req.BoardID)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Member] Board not found: %s", req.BoardID)
		utils.JSON404(c, "Board not found")
		return
	}

	member := &entity.Member{
		BoardID:   req.BoardID,
		MemberID:  req.MemberID,
		FullName:  req.FullName,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := ctrl.Repository.CreateMember(member); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create Member] Failed to create member")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create Member] Member created successfully: %s", member.ID)
	utils.JSON200(c, gin.H{
		"message": "Member created successfully",
		"data":    member,
	})
}

// GetMembers lấy tất cả members hoặc theo board_id
func (ctrl *Controller) GetMembers(c *gin.Context) {
	ctx := c.Request.Context()
	boardID := c.Query("board_id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Members] Get members request received with board_id: %s", boardID)

	var members []entity.Member
	var err error

	if boardID != "" {
		members, err = ctrl.Repository.GetMembersByBoardID(boardID)
	} else {
		members, err = ctrl.Repository.GetAllMembers()
	}

	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Members] Failed to get members")
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Members] Successfully retrieved %d members", len(members))
	utils.JSON200(c, gin.H{
		"message": "Members retrieved successfully",
		"data":    members,
	})
}

// GetMemberByID lấy member theo ID
func (ctrl *Controller) GetMemberByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Member By ID] Get member by ID request received for ID: %s", id)

	member, err := ctrl.Repository.GetMemberByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Get Member By ID] Member not found: %s", id)
		utils.JSON404(c, "Member not found")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Get Member By ID] Member retrieved successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Member retrieved successfully",
		"data":    member,
	})
}

// UpdateMember cập nhật member
func (ctrl *Controller) UpdateMember(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Member] Update member request received for ID: %s", id)

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Member] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	member, err := ctrl.Repository.GetMemberByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Member] Member not found: %s", id)
		utils.JSON404(c, "Member not found")
		return
	}

	// Cập nhật full_name nếu có trong request
	if req.FullName != "" {
		member.FullName = req.FullName
	}

	member.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := ctrl.Repository.UpdateMember(member); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Update Member] Failed to update member: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Update Member] Member updated successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Member updated successfully",
		"data":    member,
	})
}

// DeleteMember xóa member
func (ctrl *Controller) DeleteMember(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Member] Delete member request received for ID: %s", id)

	// Kiểm tra member có tồn tại không
	_, err := ctrl.Repository.GetMemberByID(id)
	if err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Member] Member not found: %s", id)
		utils.JSON404(c, "Member not found")
		return
	}

	if err := ctrl.Repository.DeleteMember(id); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Delete Member] Failed to delete member: %s", id)
		utils.JSON500(c, err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Delete Member] Member deleted successfully: %s", id)
	utils.JSON200(c, gin.H{
		"message": "Member deleted successfully",
	})
}
