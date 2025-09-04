package data

import (
	"context"
	"universal/app/system/internal/biz"
	"universal/app/system/internal/data/model"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Create(ctx context.Context, u *biz.User, password string) (*biz.User, error) {
	// 对密码进行哈希加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.log.WithContext(ctx).Errorf("密码加密失败: %v", err)
		return nil, err
	}

	// 转换为数据库模型，ID由GORM Hook自动生成
	dbUser := &model.User{
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: string(hashedPassword),
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		Status:   u.Status,
	}

	// 保存到数据库
	if err := r.data.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		r.log.WithContext(ctx).Errorf("创建用户失败: %v", err)
		return nil, err
	}

	// 转换回业务模型
	return &biz.User{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Phone:    dbUser.Phone,
		Nickname: dbUser.Nickname,
		Avatar:   dbUser.Avatar,
		Status:   dbUser.Status,
	}, nil
}

func (r *userRepo) GetByID(context.Context, int64) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) Update(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}
func (r *userRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func (r *userRepo) List(ctx context.Context, req *biz.ListUserRequest) (*biz.ListUserResponse, error) {
	var dbUsers []*model.User
	var total int64

	// 构建查询条件
	query := r.data.db.WithContext(ctx).Model(&model.User{})

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 状态过滤
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("查询用户总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Limit(int(req.PageSize)).Offset(int(offset)).Find(&dbUsers).Error; err != nil {
		r.log.WithContext(ctx).Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	// 转换为业务模型
	var bizUsers []*biz.User
	for _, dbUser := range dbUsers {
		bizUsers = append(bizUsers, &biz.User{
			ID:       dbUser.ID,
			Username: dbUser.Username,
			Email:    dbUser.Email,
			Phone:    dbUser.Phone,
			Nickname: dbUser.Nickname,
			Avatar:   dbUser.Avatar,
			Status:   dbUser.Status,
		})
	}

	return &biz.ListUserResponse{
		Users:    bizUsers,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (r *userRepo) GetByUsername(context.Context, string) (*biz.User, error) {
	return nil, nil
} // 根据用户名获取
func (r *userRepo) GetByEmail(context.Context, string) (*biz.User, error) {
	return nil, nil
} // 根据邮箱获取
func (r *userRepo) ExistsByUsername(context.Context, string) (bool, error) {
	return false, nil
} // 检查用户名是否存在
func (r *userRepo) ExistsByEmail(context.Context, string) (bool, error) {
	return false, nil
} // 检查邮箱是否存在
