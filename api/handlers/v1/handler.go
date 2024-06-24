package v1

import (
	"resume-generator/internal/pkg/config"
	token "resume-generator/internal/pkg/tokens"
	"resume-generator/internal/usecase"
	"time"

	"go.uber.org/zap"
)

type HandlerV1 struct {
	ContextTimeout time.Duration
	log            *zap.Logger
	cfg            *config.Config
	jwthandler     token.JWTHandler
	user           usecase.User
}

// HandlerV1Config ...
type HandlerV1Config struct {
	ContextTimeout time.Duration
	Logger         *zap.Logger
	Config         *config.Config
	User           usecase.User
	JWThandler     token.JWTHandler
}

// New ...
func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		ContextTimeout: c.ContextTimeout,
		log:            c.Logger,
		cfg:            c.Config,
		user:           c.User,
		jwthandler:     c.JWThandler,
	}
}
