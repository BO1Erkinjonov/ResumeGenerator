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

type ResumeTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *repo.ResumeRepository
	RepoUser    *repo.UserRepo
}

func (s *ResumeTestSuite) SetupTest() {
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
	s.Repository = repo.NewResumeRepo(pgPool)
	s.RepoUser = repo.NewUserRepo(pgPool)
	s.CleanUpFunc = pgPool.Close

}

func (s *ResumeTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func (r *ResumeTestSuite) TestResumeCrud() {
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

	user, err := r.RepoUser.CreateUser(ctx, user)
	r.Suite.NoError(err)
	r.Suite.NotNil(user)

	resme := &entity.Resume{
		ID:          uuid.NewString(),
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Category:    "PROGRAMING",
		BirthData:   time.Now(),
		Salary:      "12344555",
		Description: "Test description",
		WorkType:    "offline",
	}

	result, err := r.Repository.CreateResume(ctx, resme)
	r.Suite.NoError(err)
	r.Suite.NotNil(result)
	r.Suite.Equal(result.ID, resme.ID)
	r.Suite.Equal(result.UserID, resme.UserID)
	r.Suite.Equal(result.FirstName, resme.FirstName)
	r.Suite.Equal(result.LastName, resme.LastName)
	r.Suite.Equal(result.Category, resme.Category)
	r.Suite.Equal(result.Description, resme.Description)
	r.Suite.Equal(result.WorkType, resme.WorkType)
	r.Suite.Equal(result.Salary, resme.Salary)

	resume, err := r.Repository.GetResumeById(ctx, &entity.FieldValueReq{Field: "id", Value: resme.ID})
	r.Suite.NoError(err)
	r.Suite.Equal(resume.ID, resme.ID)
	r.Suite.Equal(resume.FirstName, resme.FirstName)
	r.Suite.Equal(resume.LastName, resme.LastName)
	r.Suite.Equal(resume.Category, resme.Category)
	r.Suite.Equal(resume.Description, resme.Description)
	r.Suite.Equal(resume.WorkType, resme.WorkType)
	r.Suite.Equal(resume.Salary, resme.Salary)
	r.Suite.Equal(resume.BirthData, resme.BirthData)
	r.Suite.Equal(resume.UserID, resume.UserID)

	all, err := r.Repository.GetAllResumes(ctx, &entity.GetAllReq{
		Field:  "",
		Values: "",
		Limit:  10,
		Offset: 0,
	})
	r.Suite.NoError(err)
	r.Suite.NotNil(all)

	r2, err := r.Repository.UpdateResumeById(ctx, &entity.UpdateResumeReq{
		ResumeID:  resume.ID,
		FirstName: "diyorbek",
		LastName:  "orijonov",
		Category:  "dizayn",
		WorkType:  "online",
	})
	r.Suite.NoError(err)
	r.Suite.NotNil(r2)

	resul, err := r.Repository.DeleteResume(ctx, &entity.DeleteReq{
		ID: resume.ID,
	})
	r.Suite.NoError(err)
	r.Suite.NotNil(resul)

}

func TestResumeTestSuite(t *testing.T) {
	suite.Run(t, new(ResumeTestSuite))
}
