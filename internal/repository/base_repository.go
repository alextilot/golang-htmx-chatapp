package repository

import (
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create a new record
func (r *BaseRepository[T]) Create(entity *T) (*T, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// FindByID retrieves one record by ID
func (r *BaseRepository[T]) FindByID(id string) (*T, error) {
	var entity T
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// FindAll retrieves all records
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Update updates existing fields (partial update)
func (r *BaseRepository[T]) Update(entity *T, id string) (*T, error) {
	if err := r.db.Model(entity).Where("id = ?", id).Updates(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *BaseRepository[T]) delete(id string, hard bool) error {
	query := r.db
	if hard {
		query = query.Unscoped()
	}
	res := query.Delete(new(T), "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete performs a soft delete
func (r *BaseRepository[T]) Delete(id string) error {
	return r.delete(id, false)
}

// DeleteHard permanently deletes a record
func (r *BaseRepository[T]) DeleteHard(id string) error {
	return r.delete(id, true)
}
