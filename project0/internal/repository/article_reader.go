package repository

import (
	"context"
	"project0/internal/domain"
)

type ArticleReaderRepository interface {
	//Create(ctx context.Context, art domain.Article) (error)
	// insert or update
	Save(ctx context.Context, art domain.Article) error
}
