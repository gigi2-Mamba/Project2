package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleReaderDAO interface {
	// insert  or update
	Upsert(ctx context.Context, art Article) error
}

type ArticleGORMDAO struct {
	db *gorm.DB
}

func (a *ArticleGORMDAO) Upsert(ctx context.Context, art Article) error {
	//TODO implement me
	panic("implement me")
}

func NewArticleReaderGORMDAO(db *gorm.DB) ArticleReaderDAO {
	return &ArticleGORMDAO{db: db}
}
