package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"resume-generator/api/models"
	"resume-generator/internal/entity"
	"time"
)

// GetAllUsers ...
// @Summary GetAllUsers
// @Description GetAllUsers - Api for get users
// @Tags User
// @Accept json
// @Produce json
// @Param field query int false "field"
// @Param value query int false "value"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} models.ListUsers
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/all/ [get]
func (h *HandlerV1) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	field := c.Query("field")
	value := c.Query("value")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	offset := cast.ToUint64(limit) * (cast.ToUint64(page) - 1)
	resp, err := h.user.GetAllUsers(ctx, &entity.GetAllReq{
		Field:  field,
		Values: value,
		Limit:  cast.ToUint64(limit),
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Get all users error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetUser ...
// @Summary GetUser
// @Description GetUser - Api for get users
// @Tags User
// @Accept json
// @Produce json
// @Param field query string false "field"
// @Param value query string false "value"
// @Success 200 {object} models.UserBody
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/get/ [get]
func (h *HandlerV1) GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	field := c.Query("field")
	value := c.Query("value")

	resp, err := h.user.GetUserById(ctx, &entity.FieldValueReq{
		Field: field,
		Value: value,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Get user error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateUser ...
// @Summary UpdateUser
// @Description UpdateUser - Api for get users
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUser body models.UpdateUserReq true "updateUserModel"
// @Success 200 {object} models.UserBody
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/update/ [put]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var body models.UserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Update user error", zap.Error(err))
		return
	}

	resp, err := h.user.UpdateUserById(ctx, &entity.UpdateUserReq{
		UserId:    body.ID,
		UserName:  body.Username,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Password:  body.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Update user error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteUser ...
// @Summary DeleteUser
// @Description DeleteUser - Api for get users
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param DeleteUser query models.DeleteReq true "deleteUserModel"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/delete/ [delete]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	id := c.Query("id")

	resp, err := h.user.DeleteUserById(ctx, &entity.DeleteReq{
		ID: id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Update user error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}
