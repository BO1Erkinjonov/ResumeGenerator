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

type UserTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *repo.UserRepo
}

func (s *UserTestSuite) SetupTest() {
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
func (s *UserTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestUserCrud() {
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
	resp, err := s.Repository.CreateUser(ctx, user)
	s.Suite.NoError(err)
	s.Suite.NotNil(resp)
	s.Suite.Equal(resp.ID, user.ID)
	s.Suite.Equal(resp.FirstName, user.FirstName)
	s.Suite.Equal(resp.LastName, user.LastName)
	s.Suite.Equal(resp.Email, user.Email)
	s.Suite.Equal(resp.Password, user.Password)
	s.Suite.Equal(resp.Username, user.Username)
	s.Suite.Equal(resp.ImageUrl, user.ImageUrl)

	respGet, err := s.Repository.GetUserById(ctx, &entity.FieldValueReq{
		Field: "id",
		Value: user.ID,
	})

	s.Suite.NoError(err)
	s.Suite.NotNil(respGet)
	s.Suite.Equal(respGet.ID, user.ID)
	s.Suite.Equal(respGet.FirstName, user.FirstName)
	s.Suite.Equal(respGet.LastName, user.LastName)
	s.Suite.Equal(respGet.Email, user.Email)
	s.Suite.Equal(respGet.Password, user.Password)
	s.Suite.Equal(respGet.Username, user.Username)
	s.Suite.Equal(respGet.ImageUrl, user.ImageUrl)

	all, err := s.Repository.GetAllUsers(ctx, &entity.GetAllReq{
		Field:  "",
		Values: "",
		Limit:  10,
		Offset: 0,
	})

	s.Suite.NoError(err)
	s.Suite.NotNil(all)

	newFirstName := "Deyorbek"
	newLastName := "Orifjonv"
	newPassword := "+_+diyor2005+_+"
	newUserName := "D1YORTOP4EEK"

	result, err := s.Repository.UpdateUserById(ctx, &entity.UpdateUserReq{
		UserId:    user.ID,
		FirstName: newFirstName,
		LastName:  newLastName,
		Password:  newPassword,
		UserName:  newUserName,
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(result)
	s.Suite.Equal(result.ID, user.ID)
	s.Suite.Equal(result.FirstName, newFirstName)
	s.Suite.Equal(result.LastName, newLastName)
	s.Suite.Equal(result.Email, user.Email)
	s.Suite.Equal(result.Password, newPassword)
	s.Suite.Equal(result.Username, newUserName)
	s.Suite.Equal(result.ImageUrl, user.ImageUrl)

	resultDel, err := s.Repository.DeleteUserById(ctx, &entity.DeleteReq{ID: resp.ID})
	s.Suite.NoError(err)
	s.Suite.NotNil(resultDel)
}
