package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

func InitTables(db *gorm.DB) error {

	return db.AutoMigrate(&User{}, &Article{}, &PublishedArticle{})
}

func InitCollection(mdb *mongo.Database) error {
	// 让你自己来做，真的很多东西你不知道怎么搞，何来创造性？  第一批的程序员是怎么工作的
	// think before action
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	col := mdb.Collection("articles")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	if err != nil {
		return err
	}
	livecol := mdb.Collection("published_articles")
	_, err = livecol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})

	return err
}
