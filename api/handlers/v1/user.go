package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"resume-generator/internal/entity"
	"time"
)

// GetAllUsers ...
// @Summary GetAllUsers
// @Description GetAllUsers - Api for get users
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} models.ListUsers
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/all/ [get]
func (h *HandlerV1) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	offset := cast.ToUint64(limit) * (cast.ToUint64(page) - 1)
	resp, err := h.user.GetAllUsers(ctx, &entity.GetAllUserReq{
		Field:  "",
		Values: "",
		Limit:  0,
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

//// GetUser ...
//// @Summary GetUser
//// @Description GetUser - Api for get users
//// @Tags User
//// @Accept json
//// @Produce json
//// @Param field query int false "field"
//// @Param value query int false "value"
//// @Success 200 {object} models.UserBody
//// @Failure 400 {object} models.StandardErrorModel
//// @Failure 500 {object} models.StandardErrorModel
//// @Router /v1/user/all/ [get]
//func (h *HandlerV1) GetUser(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
//	defer cancel()
//
//}
