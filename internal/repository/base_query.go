package repository

import (
	"context"

	"gorm.io/gorm"
)

type QueryOption func(*gorm.DB) *gorm.DB

func (r *BaseRepository[T]) Find(
	ctx context.Context,
	out *[]T,
	opts ...QueryOption,
) error {

	db := r.db.WithContext(ctx)

	for _, opt := range opts {
		db = opt(db)
	}

	return db.Find(out).Error
}

func (r *BaseRepository[T]) First(
	ctx context.Context,
	out *T,
	opts ...QueryOption,
) error {

	db := r.db.WithContext(ctx)

	for _, opt := range opts {
		db = opt(db)
	}

	return db.First(out).Error
}
