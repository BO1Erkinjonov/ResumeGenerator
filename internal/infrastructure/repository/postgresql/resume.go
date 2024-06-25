package postgresql

import (
	"resume-generator/internal/pkg/postgres"
)

//TODO:: type Resume interface {
//	CreateResume(ctx context.Context, resume *entity.Resume) (*entity.Resume, error)
//	GetResumeById(ctx context.Context, req *entity.FieldValueReq) (*entity.Resume, error)
//	CheckUniques(ctx context.Context, req *entity.FieldValueReq) (*entity.Result, error)
//	GetAllResumes(ctx context.Context, req *entity.GetAllReq) ([]*entity.Resume, error)
//	DeleteResume(ctx context.Context, req *entity.DeleteReq) (*entity.Result, error)
//	UpdateResumeById(ctx context.Context, req *entity.UpdateResumeReq) (*entity.Result, error)
//}

type ResumeRepository struct {
	db        *postgres.PostgresDB
	tableName string
}

func NewResumeRepo(pg *postgres.PostgresDB) *ResumeRepository {
	return &ResumeRepository{
		db:        pg,
		tableName: "resumes",
	}
}
