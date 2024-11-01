package database

import (
	"context"

	"gorm.io/gorm"
)

type contextKey string

const txKey = contextKey("tx")

func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db
}

type Service[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, entity *T) error
	Get(ctx context.Context, entity *T) (*T, error)
}

type TransactionalService[T any] struct {
	db        *gorm.DB
	decorated Service[T]
}

// NewTransactionalRepository constructor for the transactional decorator
func NewTransactionalRepository[T any](db *gorm.DB, decorated Service[T]) *TransactionalService[T] {
	return &TransactionalService[T]{
		db:        db,
		decorated: decorated,
	}
}

func (t *TransactionalService[T]) Create(ctx context.Context, entity *T) (*T, error) {
	tx := t.db.Begin()
	ctxWithTx := SetTx(ctx, tx)

	res, err := t.decorated.Create(ctxWithTx, entity)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, tx.Commit().Error

}

func (t *TransactionalService[T]) Update(ctx context.Context, entity *T) error {
	tx := t.db.Begin()
	ctxWithTx := SetTx(ctx, tx)
	err := t.decorated.Update(ctxWithTx, entity)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (t *TransactionalService[T]) Delete(ctx context.Context, entity *T) error {
	tx := t.db.Begin()
	ctxWithTx := SetTx(ctx, tx)
	err := t.decorated.Delete(ctxWithTx, entity)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (t *TransactionalService[T]) Get(ctx context.Context, entity *T) (*T, error) {
	return t.decorated.Get(ctx, entity)
}
