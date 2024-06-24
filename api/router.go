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
// @host :9050
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	router.Use(casbin.NewAuthorizer())

	api := router.Group("/v1")

	auth := api.Group("/auth")
	auth.POST("/register/", HandlerV1.Register)
	auth.POST("/verification/", HandlerV1.Verification)
	auth.POST("/login/", HandlerV1.LogIn)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
