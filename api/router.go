package api

import (
	"fmt"
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	api := router.Group("/v1")

	fmt.Println(HandlerV1)
	//animalType := api.Group("/animal-type")
	//animalType.POST("", HandlerV1.CreateAnimalTypes)
	//animalType.GET("/get", HandlerV1.GetAnimalTypes)
	//animalType.GET("", HandlerV1.ListAnimalTypes)
	//animalType.PUT("", HandlerV1.UpdateAnimalTypes)
	//animalType.DELETE("", HandlerV1.DeleteAnimalTypes)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
