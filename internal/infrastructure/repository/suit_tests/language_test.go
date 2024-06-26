package suit_tests

import (
	"context"
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

type LanguageTestSuite struct {
	suite.Suite
	CleanUpFunc      func()
	Repository       *repo.LanguageRepository
	ResumeRepository *repo.ResumeRepository
	RepoUser         *repo.UserRepo
}

func (s *LanguageTestSuite) SetupTest() {
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
	s.Repository = repo.NewLanguageRepository(pgPool)
	s.ResumeRepository = repo.NewResumeRepo(pgPool)
	s.RepoUser = repo.NewUserRepo(pgPool)
	s.CleanUpFunc = pgPool.Close
}
func (s *LanguageTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestLanguageTestSuite(t *testing.T) {
	suite.Run(t, new(LanguageTestSuite))
}

func (s *LanguageTestSuite) Test_NewLanguageRepository() {
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

	language := &entity.Language{
		ID:       uuid.NewString(),
		Name:     "Test language name",
		Level:    "A1",
		ResumeID: result.ID,
	}

	resL, err := s.Repository.CreateLanguage(ctx, language)
	s.NotNil(resL)
	s.NoError(err)
	s.Equal(language.ID, resL.ID)
	s.Equal(language.Name, resL.Name)
	s.Equal(language.Level, resL.Level)
	s.Equal(language.ResumeID, resL.ResumeID)

	resL, err = s.Repository.UpdateLanguageById(ctx, &entity.LanguagesUpdateReq{
		ID:    resL.ID,
		Name:  "Test language name",
		Level: "A2",
	})
	s.NotNil(resL)
	s.NoError(err)
	all, err := s.Repository.GetAllLanguage(ctx, &entity.GetAllReq{})
	s.NotNil(all)
	s.NoError(err)
	resL, err = s.Repository.GetLanguageById(ctx, &entity.FieldValueReq{
		Field: "id",
		Value: resL.ID,
	})
	s.NotNil(resL)
	s.NoError(err)
	reslDe, err := s.Repository.DeleteLanguage(ctx, &entity.DeleteReq{
		ID: resL.ID,
	})
	s.NotNil(reslDe)
	s.NoError(err)

	res, err := s.ResumeRepository.DeleteResume(ctx, &entity.DeleteReq{ID: resme.ID})
	s.NotNil(res)
	s.NoError(err)

	res, err = s.RepoUser.DeleteUserById(ctx, &entity.DeleteReq{
		IsHardDeleted: true,
		ID:            userResponse.ID,
	})
	s.NotNil(res)
	s.NoError(err)

}
