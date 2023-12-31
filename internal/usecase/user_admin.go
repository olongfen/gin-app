package usecase

import (
	"context"
	"time"

	"gin-app/internal/domain"

	"github.com/google/uuid"
	gormgenerics "github.com/olongfen/gorm-generics"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// UserAdminConfig 用户管理员配置
type UserAdminConfig struct {
	Repo           domain.UserRepo
	ContextTimeout time.Duration
}

// userAdminUsecase 用户管理员用例
type userAdminUsecase struct {
	cfg UserAdminConfig
}

// NewUserAdminUsecase 新建用户管理员用例
func NewUserAdminUsecase(cfg UserAdminConfig) domain.UserAdminUsecase {
	return &userAdminUsecase{
		cfg: cfg,
	}
}

// List 获取用户列表
func (u *userAdminUsecase) List(ctx context.Context, req *domain.UserAdminListReq) (*domain.UserAdminListResp, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	data, total, err := u.cfg.Repo.Find(ctx, &gormgenerics.Limit{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Count:    true,
	}, clause.OrderBy{
		Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "created_at"}, Desc: true},
		},
	})
	if err != nil {
		return nil, err
	}
	ret := &domain.UserAdminListResp{}
	ret.Pagination = &domain.Pagination{
		Total:    total,
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
	}
	for _, v := range data {
		ret.List = append(ret.List, &domain.UserListInfo{
			ID:     v.ID,
			Uuid:   v.Uuid,
			Status: v.Status,
			UserInfo: domain.UserInfo{
				Gender:    v.Gender,
				Email:     v.Email,
				Username:  v.Username,
				Phone:     v.Phone,
				CreatedAt: v.CreatedAt,
			},
		})
	}
	return ret, nil
}

// Add 添加用户
func (u *userAdminUsecase) Add(ctx context.Context, req *domain.UserAdminAddReq) error {
	ctx, cancelFunc := context.WithTimeout(ctx, u.cfg.ContextTimeout)
	defer cancelFunc()
	user := &domain.User{
		Uuid:     uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Gender:   req.Gender,
		Status:   req.Status,
		Phone:    req.Phone,
	}
	password, err := bcrypt.GenerateFromPassword([]byte(req.Phone), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	if err = u.cfg.Repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}
