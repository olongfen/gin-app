package repository

import (
	"gin-app/internal/domain"

	gormgenerics "github.com/olongfen/gorm-generics"
)

// userRepo 学生存储库实现
type userRepo struct {
	gormgenerics.BasicRepo[domain.User]
	db gormgenerics.Database
}

var _ domain.UserRepo = (*userRepo)(nil)

// NewUserRepo 新建学生存储库
func NewUserRepo(data gormgenerics.Database) domain.UserRepo {
	d := gormgenerics.NewBasicRepository[domain.User](data)
	stu := &userRepo{
		d,
		data,
	}
	return stu
}
