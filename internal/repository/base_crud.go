package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// CREATE
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// READ
func (r *BaseRepository[T]) FindByID(ctx context.Context, id string, out *T) error {
	err := r.db.WithContext(ctx).
		First(out, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}

// Note: there is deliberately no generic Update(ctx, id, updates any) here.
// An untyped "updates any" parameter lets callers pass arbitrary field
// maps, which hides domain operations (like soft-deleting a group) behind
// stringly-typed keys instead of an explicit method signature. Add a
// purpose-built, typed method to the concrete repository instead (see
// GroupRepository.Deactivate / UpdateDetails for the pattern).

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
