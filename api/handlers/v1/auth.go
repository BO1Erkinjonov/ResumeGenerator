package v1

import (
	_ "resume-generator/api/docs"
	"resume-generator/api/models"
	"resume-generator/internal/entity"
	l "resume-generator/internal/pkg/logger"
	token "resume-generator/internal/pkg/tokens"

	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Register ...
// @Summary Register
// @Description Register - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param Register body models.User true "createRegisterModel"
// @Success 200 {object} bool
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/register/ [post]
func (h *HandlerV1) Register(c *gin.Context) {
	var (
		body models.UserBody
	)

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json", l.Error(err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	exists_email, err := h.user.CheckUniques(ctx, &entity.FieldValueReq{
		Field: "email",
		Value: body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniques")
		return
	}

	if exists_email.IsExists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email address",
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	body.ID = uuid.New().String()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	byteDate, err := json.Marshal(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to marshaling error")
		return
	}

	err = rdb.SetEx(context.Background(), "email_"+body.Email, byteDate, time.Minute*1).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to redis error")
		return
	}
	code := rand.Int() % 1000000
	err = rdb.SetEx(context.Background(), body.Email, code, time.Minute*1).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to redis error")
		return
	}
	codeResp := strconv.Itoa(code)

	auth := smtp.PlainAuth("", "boburerkinzonov@gmail.com", "llqmgbilccvhltfd", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "boburerkinzonov@gmail.com", []string{body.Email}, []byte(codeResp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to email errored", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, true)
}

// Verification ...
// @Summary Verification
// @Description Verification - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param code query int true "Code"
// @Param email query string true "Email"
// @Success 200 {object} models.AccessToken
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/verification/ [post]
func (h *HandlerV1) Verification(c *gin.Context) {
	codeRegis := c.Query("code")
	emailRegis := c.Query("email")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	respCode, err := rdb.Get(context.Background(), emailRegis).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to redis error")
		return
	}
	var code int
	if err := json.Unmarshal([]byte(respCode), &code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to unmarshal error")
		return
	}

	code1, err := strconv.Atoi(codeRegis)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to otp error")
		return
	}

	if code != code1 {
		c.JSON(http.StatusBadRequest, false)
	} else {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
		defer cancel()
		err = rdb.Del(context.Background(), emailRegis).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to redis error")
			return
		}
		var regis models.UserBody
		respUser, err := rdb.Get(context.Background(), "email_"+emailRegis).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to redis error")
			return
		}
		err = rdb.Del(context.Background(), "email_"+emailRegis).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to redis error")
			return
		}
		if err := json.Unmarshal([]byte(respUser), &regis); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to unmarshal error")
			return
		}

		h.jwthandler = token.JWTHandler{
			Sub:     regis.ID,
			Role:    "user",
			SignKey: h.cfg.Token.Secret,
			Timout:  h.cfg.Token.AccessTTL,
		}

		access, _, err := h.jwthandler.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to token error")
			return
		}
		_, err = h.user.CreateUser(ctx, &entity.User{
			ID:        regis.ID,
			FirstName: regis.FirstName,
			LastName:  regis.LastName,
			Email:     regis.Email,
			Password:  regis.Password,
			Username:  regis.Username,
			ImageUrl:  regis.ImageUrl,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed to created error")
			return
		}

		c.JSON(http.StatusOK, access)
	}
}

// LogIn ...
// @Summary LogIn
// @Description LogIn - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param password query string true "Password"
// @Param email query string true "Email"
// @Success 200 {object} models.AccessToken
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/login/ [post]
func (h *HandlerV1) LogIn(c *gin.Context) {
	password := c.Query("password")
	email := c.Query("email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
	defer cancel()

	user, err := h.user.GetUserById(ctx, &entity.FieldValueReq{
		Field: "email",
		Value: email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.DeletedAt.String() != "0001-01-01 00:00:00 +0000 UTC" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "such user does not exist",
		})
		h.log.Error("the user has already been deleted", l.Error(err))
		return
	}
	if password != user.Password {
		c.JSON(http.StatusBadRequest, "password error")
		return
	}
	h.jwthandler = token.JWTHandler{
		Sub:     user.ID,
		Role:    "user",
		SignKey: h.cfg.Token.Secret,
		Timout:  h.cfg.Token.AccessTTL,
	}
	access, _, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to token error")
		return
	}

	c.JSON(http.StatusOK, access)
}
