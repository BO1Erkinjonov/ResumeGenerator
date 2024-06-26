package suit_tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"log"
	repo "resume-generator/internal/infrastructure/repository/postgresql"
	"resume-generator/internal/pkg/config"
	db "resume-generator/internal/pkg/postgres"
	"testing"
	"time"
)

type LinkTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *repo.UserRepo
}

func (s *LinkTestSuite) SetupTest() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println(err)
		return
	}
	pgPool, err := db.New(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	s.Repository = repo.NewUserRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}
func (s *LinkTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestLinkTestSuite(t *testing.T) {
	suite.Run(t, new(LinkTestSuite))
}

func (s *LinkTestSuite) Test_AddLink() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(2))
	defer cancel()

}
