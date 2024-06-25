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

type AnimalTypesTestSuite struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *repo.UserRepo
}

func (s *AnimalTypesTestSuite) SetupTest() {
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

func (s *AnimalTypesTestSuite) TestAnimalTypesCrud() {
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

	all, err := s.Repository.GetAllUsers(ctx, &entity.GetAllUserReq{})

	s.Suite.NoError(err)
	s.Suite.NotNil(all)

}

func (s *AnimalTypesTestSuite) TestDeleteUser() {
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

	result, err := s.Repository.DeleteUserById(ctx, &entity.DeleteUserReq{ID: resp.ID})
	s.Suite.NoError(err)
	s.Suite.NotNil(result)
}

func (s *AnimalTypesTestSuite) TestUpdateUser() {
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

	result, err := s.Repository.UpdateUserById(ctx, &entity.UpdateUserReq{
		UserId:    user.ID,
		FirstName: "Diyorbek",
		LastName:  "Orifjonv",
		Password:  "+_+diyor2005+_+",
		UserName:  "D1YORTOP4EEK",
	})
	s.Suite.NoError(err)
	s.Suite.NotNil(result)
}
func (s *AnimalTypesTestSuite) TearDownTest() {
	s.CleanUpFunc()
}

func TestAnimalTypesTestSuite(t *testing.T) {
	suite.Run(t, new(AnimalTypesTestSuite))
}
