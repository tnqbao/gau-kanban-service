package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-kanban-service/entity"
	"github.com/tnqbao/gau-kanban-service/utils"
)

// CreateUser creates a new user with specified user_id and fullname
func (ctrl *Controller) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create User] Create user request received")

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create User] Invalid request body")
		utils.JSON400(c, "Invalid request body: "+err.Error())
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create User] Creating user with ID: %s, FullName: %s", req.UserID, req.FullName)

	// Check if user already exists
	existingUser, err := ctrl.Repository.GetUserByID(req.UserID)
	if err == nil && existingUser != nil {
		ctrl.Provider.LoggerProvider.WarningWithContextf(ctx, "[Create User] User already exists with ID: %s", req.UserID)
		utils.JSON409(c, "User already exists with this ID")
		return
	}

	// Create user with specified ID and fullname
	user := &entity.User{
		ID:       req.UserID,
		FullName: req.FullName,
	}

	if err := ctrl.Repository.CreateUser(user); err != nil {
		ctrl.Provider.LoggerProvider.ErrorWithContextf(ctx, err, "[Create User] Failed to create user")
		utils.JSON500(c, "Failed to create user")
		return
	}

	ctrl.Provider.LoggerProvider.InfoWithContextf(ctx, "[Create User] User created successfully: %s", user.ID)
	utils.JSON201(c, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}
