package dao

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)
type ArticleS3DAO struct {
	ArticleGROMDAO
	oss *s3.S3
}

func NewArticleS3DAO(db ArticleGROMDAO, oss *s3.S3) *ArticleS3DAO {
	return &ArticleS3DAO{ArticleGROMDAO: db, oss: oss}
}


func (a *ArticleS3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 方法外定义更简洁？
	id := art.Id
	//log.Println("should 2 published : ",art.Status)
	// 闭包处理事务
	err :=a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//创造了本体，使用了事务实现
		dao := NewArticleGROMDAO(tx)
		//使用var 方便给if-else使用
		var (
			err error
		)
		if id > 0 {
			err = dao.UpdateById(ctx, art)
		} else {
			id, err = dao.Insert(ctx, art)
		}

		if err != nil {
			return  err
		}
		now := time.Now().UnixMilli()
		art.Id = id
		publishArt := PublishedArticleV2{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Ctime:    now,
			Utime:    now,
			Status:   art.Status,
		}
		// 确保代码逻辑分支生成的返回值要传到需要的实体上
		art.Id = id
		err = tx.Clauses(clause.OnConflict{
			// 对mysql不起效，但是可以兼容别的方言
			Columns: []clause.Column{{Name: "id"}},
			//做更新，mysql就支持这个
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title": publishArt.Title,

				"utime": now,
				"status":publishArt.Status,
			}),
		}).Create(&publishArt).Error
		return err
	})

	if err != nil {
		return 0, err
	}
	_, err = a.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("webook-1314583317"),
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return id, err
}


func (a *ArticleS3DAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	now := time.Now().UnixMilli()
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? and author_id = ?", uid, id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return  errors.New("ID 不对或者创作者不对")
		}
		return tx.Model(&PublishedArticleV2{}).
			Where("id = ?", uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
	if err != nil {
		return err
	}
	const statusPrivate = 3
	if status == statusPrivate {
		_, err = a.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: ekit.ToPtr[string]("webook-1314583317"),
			Key:    ekit.ToPtr[string](strconv.FormatInt(id, 10)),
		})
	}
	return err
}

//  oss
type PublishedArticleV2 struct {
	Id      int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	//更新时间
	Utime int64 `bson:"utime,omitempty"`
	Status uint8 `bson:"status,omitempty"`
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	//Content string `gorm:"type=BLOB"`


}