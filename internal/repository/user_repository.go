package repository

import (
	"context"

	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[model.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository[model.User](db),
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	err := r.DB().
		WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUsernameOrEmail(
	ctx context.Context,
	username, email string,
) (*model.User, error) {

	var user model.User

	err := r.DB().
		WithContext(ctx).
		Where("username = ? OR email = ?", username, email).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Search(ctx context.Context, query string) ([]model.User, error) {
	var users []model.User

	q := r.DB().WithContext(ctx)

	if query != "" {
		q = q.Where("username LIKE ?", "%"+query+"%")
	}

	err := q.Find(&users).Error
	return users, err
}

func (r *UserRepository) Exists(ctx context.Context, username, email string) (bool, error) {
	var count int64

	err := r.DB().
		WithContext(ctx).
		Model(&model.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error

	return count > 0, err
}
