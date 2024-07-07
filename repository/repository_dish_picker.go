package repository

import (
	"context"
)

type DishID string

type PickedDishByUser struct {
	Dish            DishID
	DishDescription string
	UserID          string
}

type RepositoryDishPicker interface {
	PutDishDescription(ctx context.Context, dishDescription string) (DishID, error)
	GetDishDescription(ctx context.Context, dishID DishID) (string, error)
	ListDishVariants(ctx context.Context) ([]DishID, error)

	PickDishForUser(ctx context.Context, userID int64, dishID DishID) error
	GetPickedDishForUser(ctx context.Context, userID int64) (DishID, error)
	ListPickedDishesForUsers(ctx context.Context) (PickedDishByUser, error)
}
