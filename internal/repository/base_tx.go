package repository

import (
	"context"

	"gorm.io/gorm"
)

func (r *BaseRepository[T]) Transaction(
	ctx context.Context,
	fn func(tx *BaseRepository[T]) error,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&BaseRepository[T]{db: tx})
	})
}
