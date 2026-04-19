package repository

import (
	"context"
)

// CREATE
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// READ
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string, out *T) error {
	return r.db.WithContext(ctx).
		First(out, "id = ?", id).Error
}

// UPDATE
func (r *BaseRepository[T]) Update(ctx context.Context, id string, updates any) error {
	return r.db.WithContext(ctx).
		Model(new(T)).
		Where("id = ?", id).
		Updates(updates).Error
}

// DELETE (soft if model supports it)
func (r *BaseRepository[T]) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Delete(new(T), "id = ?", id).Error
}

// HARD DELETE
func (r *BaseRepository[T]) DeleteHard(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Delete(new(T), "id = ?", id).Error
}
