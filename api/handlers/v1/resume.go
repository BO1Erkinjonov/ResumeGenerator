package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
	"resume-generator/api/models"
	"resume-generator/internal/entity"
	"resume-generator/internal/pkg/config"
	l "resume-generator/internal/pkg/logger"
	jwt "resume-generator/internal/pkg/tokens"
	"time"
)

// CreateResume ...
// @Summary CreateResume
// @Description CreateResume - Api for create resume
// @Security ApiKeyAuth
// @Tags Resume
// @Accept json
// @Produce json
// @Param CreateResume body models.ReqResume false "createResumeModel"
// @Param work_type query string true "work_type" Enums(offline, online, does not matter) "work_type"
// @Success 200 {object} models.Resume
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/resume/create/ [post]
func (h *HandlerV1) CreateResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var r models.Resume
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		h.log.Error("ShouldBindJSON", zap.Error(err))
		return
	}

	workType := c.Query("work_type")

	tok := c.GetHeader("Authorization")
	claims, err := jwt.ExtractClaim(tok, []byte(config.Token()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	birthDate, err := time.Parse("2006-01-02", r.BirthData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parse birthdate", l.Error(err))
		return
	}

	resp, err := h.resume.CreateResume(ctx, &entity.Resume{
		ID:          uuid.NewString(),
		UserID:      cast.ToString(claims["sub"]),
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		Category:    r.Category,
		BirthData:   birthDate,
		Salary:      r.Salary,
		Description: r.Description,
		WorkType:    workType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("CreateResume", zap.Error(err))
	}
	c.JSON(http.StatusCreated, resp)
}

// GetAllResume ...
// @Summary GetAllResume
// @Description GetAllResume - Api for get resume
// @Tags Resume
// @Accept json
// @Produce json
// @Param field query int false "field"
// @Param value query int false "value"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} models.ListResume
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/resume/all/ [get]
func (h *HandlerV1) GetAllResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()
	field := c.Query("field")
	value := c.Query("value")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	offset := cast.ToUint64(limit) * (cast.ToUint64(page) - 1)
	resp, err := h.resume.GetAllResumes(ctx, &entity.GetAllReq{
		Field:  field,
		Values: value,
		Limit:  cast.ToUint64(limit),
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Get all resumes error", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, resp)
}

// GetResume ...
// @Summary GetResume
// @Description GetResume - Api for get users
// @Tags Resume
// @Accept json
// @Produce json
// @Param field query string false "field"
// @Param value query string false "value"
// @Success 200 {object} models.Resume
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/resume/get/ [get]
func (h *HandlerV1) GetResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	field := c.Query("field")
	value := c.Query("value")

	resp, err := h.resume.GetResumeById(ctx, &entity.FieldValueReq{
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

// UpdateResume ...
// @Summary UpdateResume
// @Description UpdateResume - Api for get users
// @Security ApiKeyAuth
// @Tags Resume
// @Accept json
// @Produce json
// @Param UpdateResume body models.UpdateResume true "updateUserModel"
// @Param work_type query string true "work_type" Enums(offline, online, does not matter) "work_type"
// @Success 200 {object} models.UserBody
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/resume/update/ [put]
func (h *HandlerV1) UpdateResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	var body models.UpdateResume
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		h.log.Error("Update user error", zap.Error(err))
		return
	}
	workType := c.Query("work_type")
	birthDate, err := time.Parse("2006-01-02", body.BirthData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to parse birthdate", l.Error(err))
		return
	}

	resp, err := h.resume.UpdateResumeById(ctx, &entity.UpdateResumeReq{
		ResumeID:    body.ID,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Category:    body.Category,
		BirthDate:   birthDate,
		Salary:      body.Salary,
		Description: body.Description,
		WorkType:    workType,
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

// DeleteResume ...
// @Summary DeleteResume
// @Description DeleteResume - Api for get users
// @Security ApiKeyAuth
// @Tags Resume
// @Accept json
// @Produce json
// @Param DeleteResume query models.DeleteReq true "deleteUserModel"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/resume/delete/ [delete]
func (h *HandlerV1) DeleteResume(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	id := c.Query("id")

	resp, err := h.resume.DeleteResume(ctx, &entity.DeleteReq{
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
