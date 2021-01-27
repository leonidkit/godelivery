package storage

import "context"

type Item struct {
	ID         int64
	FormatType string
	Format     []byte
}

type ItemRepository interface {
	// Add item to repository
	Create(ctx context.Context, item Item) error

	// Remove item from repository by ID and FormatType
	Delete(ctx context.Context, id int64, formatType string) error

	// Select item from repository by ID and FormatType
	Get(ctx context.Context, id int64, formatType string) (Item, error)
}
