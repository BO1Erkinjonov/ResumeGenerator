package api

import (
	"resume-generator/api/middleware/casbin"
	//_ "resume-generator/api/docs"
	"resume-generator/internal/usecase"
	"time"

	v1 "resume-generator/api/handlers/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"resume-generator/internal/pkg/config"
)

type RouteOption struct {
	ContextTimeout time.Duration
	Logger         *zap.Logger
	Config         *config.Config
	User           usecase.User
}

// NewRoute
// @title Generate resume
// @version 1.7
// // @host 18.158.24.26:9050
// @host localhost:9050
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		User:           option.User,

		//BrokerProducer: option.BrokerProducer,
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Замените на порт вашего локального фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(casbin.NewAuthorizer())

	api := router.Group("/v1")

	// auth
	auth := api.Group("/auth")
	auth.POST("/register/", HandlerV1.Register)
	auth.POST("/verification/", HandlerV1.Verification)
	auth.POST("/login/", HandlerV1.LogIn)

	// user
	user := api.Group("/user")
	user.GET("/get/", HandlerV1.GetUser)
	user.GET("/all/", HandlerV1.GetAllUsers)
	user.PUT("/update/", HandlerV1.UpdateUser)
	user.DELETE("/delete/", HandlerV1.DeleteUser)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
