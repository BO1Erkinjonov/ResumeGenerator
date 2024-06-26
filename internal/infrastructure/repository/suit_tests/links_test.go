package suit_tests

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"log"
	"resume-generator/internal/entity"
	repo "resume-generator/internal/infrastructure/repository/postgresql"
	"resume-generator/internal/pkg/config"
	db "resume-generator/internal/pkg/postgres"
	"testing"
	"time"
)

type LinkTestSuite struct {
	suite.Suite
	CleanUpFunc      func()
	Repository       *repo.LinkRepository
	ResumeRepository *repo.ResumeRepository
	RepoUser         *repo.UserRepo
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
	s.Repository = repo.NewLinRepository(pgPool)
	s.ResumeRepository = repo.NewResumeRepo(pgPool)
	s.RepoUser = repo.NewUserRepo(pgPool)
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

	user := &entity.User{
		ID:        uuid.NewString(),
		FirstName: "Test first name",
		LastName:  "Test last name",
		Email:     "Test email",
		Password:  "Test password",
		Username:  "Test username",
		ImageUrl:  "Test image",
	}

	userResponse, err := s.RepoUser.CreateUser(ctx, user)
	s.NotNil(userResponse)
	s.NoError(err)

	resme := &entity.Resume{
		ID:          uuid.NewString(),
		UserID:      userResponse.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Category:    "PROGRAMING",
		BirthData:   time.Now(),
		Salary:      "12344555",
		Description: "Test description",
		WorkType:    "offline",
	}
	result, err := s.ResumeRepository.CreateResume(ctx, resme)
	s.NotNil(result)
	s.NoError(err)

	link := &entity.Link{
		ID:       uuid.NewString(),
		ResumeID: result.ID,
		LinkName: "Test link name",
		LinkURL:  "Test link url",
	}
	resultLink, err := s.Repository.CreateLink(ctx, link)
	s.NotNil(resultLink)
	s.NoError(err)
	fmt.Println(resultLink, "wegwegewgewg")
	s.Suite.Equal(link.LinkURL, resultLink.LinkURL)
	s.Suite.Equal(link.LinkName, resultLink.LinkName)
	s.Suite.Equal(link.ResumeID, resultLink.ResumeID)

	resultLink, err = s.Repository.UpdateLinkById(ctx, &entity.LinksUpdateReq{
		LinID:    resultLink.ID,
		LinkName: "qwert",
		LinkURL:  "qwerl",
	})

	s.NotNil(resultLink)
	s.NoError(err)

	all, err := s.Repository.GetAllLinks(ctx, &entity.GetAllReq{
		Field:  "id",
		Values: resultLink.ID,
		Limit:  20,
	})
	s.NotNil(all)
	s.NoError(err)
	resultLink, err = s.Repository.GetLinkById(ctx, &entity.FieldValueReq{
		Field: "id",
		Value: resultLink.ID,
	})
	s.NotNil(resultLink)
	s.NoError(err)

	res, err := s.Repository.DeleteLink(ctx, &entity.DeleteReq{ID: resultLink.ID})
	s.NoError(err)
	s.NotNil(res)
	res, err = s.ResumeRepository.DeleteResume(ctx, &entity.DeleteReq{ID: resme.ID})
	s.NotNil(res)
	s.NoError(err)

	res, err = s.RepoUser.DeleteUserById(ctx, &entity.DeleteReq{
		IsHardDeleted: true,
		ID:            userResponse.ID,
	})
	s.NotNil(res)
	s.NoError(err)
}
