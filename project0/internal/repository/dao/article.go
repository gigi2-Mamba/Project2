package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type Article struct {
	Id       int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	Ctime    int64 `bson:"ctime,omitempty"`
	//更新时间
	Utime  int64  `bson:"utime,omitempty"`
	Status uint8  `bson:"status,omitempty"`
	Title  string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	//Content string `gorm:"type=BLOB"`
	Content string `bson:"content,omitempty"`
	//Fucker int `bson:"fucker,omitempty"`
}
type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, entity Article) (int64, error) //  这里原本有三个实现干嘛
	SyncStatus(ctx context.Context, id int64, uid int64, status uint8) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error)
	GetById(ctx context.Context, id int64) (Article, error)
	GetPubById(ctx context.Context, id int64) (PublishedArticle, error)
	ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedArticle, error)
}

// 同库不同表,高度相似可以做类型延伸
type PublishedArticle Article

// 使用组合也可以，或许这个扩展性更好
type PublishArticleV1 struct {
	Article
}
type ArticleGROMDAO struct {
	db *gorm.DB
}

func (a *ArticleGROMDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedArticle, error) {
	// 小微书第一次批量查询？
	//dao层不直接引用domain的方法
	const ArticlePublishStatus = 2
	var res []PublishedArticle
	err := a.db.WithContext(ctx).Where("utime < ? and status = ?", start.UnixMilli(), ArticlePublishStatus).
		Offset(offset).Limit(limit).First(&res).Error
	//if err != nil {
	//	return nil, err
	//}   要返回的err已经是最后一个且是返回值。直接return作为方法最后一行.有其他参数的话就构造err但是不需要对if err !=nil 做处理
	return res, err
}

func NewArticleGROMDAO(db *gorm.DB) ArticleDao {
	return &ArticleGROMDAO{db: db}
}
func (a *ArticleGROMDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var art PublishedArticle
	// first默认应该是按什么顺序查找呢
	err := a.db.Where(" id = ? ", id).First(&art).Error

	return art, err
}

func (a *ArticleGROMDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := a.db.Where("id = ?", id).Find(&art).Error
	return art, err
}

func (a *ArticleGROMDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {

	var arts []Article
	err := a.db.Where("author_id = ?", uid).Offset(offset).Limit(limit).
		Order("utime DESC").Find(&arts).Error
	return arts, err
}

func (a *ArticleGROMDAO) SyncStatus(ctx context.Context, id int64, uid int64, status uint8) error {
	// 同步状态需要同步时间
	now := time.Now().UnixMilli()
	// 闭包封装处理事务
	return a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? and author_id= ?", id, uid).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			})

		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("更新失败,ID不对或者作者不对")
		}

		return tx.Model(&PublishedArticle{}).Where("id = ?", id).
			Updates(map[string]any{
				"utime":  now,
				"status": status,
			}).Error
	})
}

//  current use 闭包实现事务。  属于v3dao层分离制作库和线上库    同库不同表
func (a *ArticleGROMDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 方法外定义更简洁？
	id := art.Id
	log.Println("happen here1")
	//log.Println("should 2 published : ",art.Status)
	// 闭包处理事务
	err := a.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//创造了本体，使用了事务实现
		dao := NewArticleGROMDAO(tx)
		//使用var 方便给if-else使用
		var (
			err error
		)
		if id > 0 {
			log.Println("here gorm dao, update: ", art.Id)

			err = dao.UpdateById(ctx, art)
		} else {
			log.Println("here gorm dao, create: ", art.Id)
			id, err = dao.Insert(ctx, art)
		}

		if err != nil {
			return err
		}
		publishArt := PublishedArticle(art)
		now := time.Now().UnixMilli()
		publishArt.Ctime = now
		publishArt.Utime = now
		// 确保代码逻辑分支生成的返回值要传到需要的实体上
		art.Id = id
		err = tx.Clauses(clause.OnConflict{
			// 对mysql不起效，但是可以兼容别的方言
			Columns: []clause.Column{{Name: "id"}},
			//做更新，mysql就支持这个
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   publishArt.Title,
				"content": publishArt.Content,
				"utime":   now,
				"status":  publishArt.Status,
			}),
		}).Create(&publishArt).Error
		return err
	})
	return id, err
}

// dao层做事务处理，dao层面同步数据的v1 手动操作事务
func (a *ArticleGROMDAO) SyncV1(ctx context.Context, art Article) (int64, error) {
	// 开启事务
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 防止后面业务panic
	defer tx.Rollback()
	//创造了本体，使用了事务实现
	dao := NewArticleGROMDAO(tx)
	//使用var 方便给if-else使用
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	publishArt := PublishedArticle(art)
	now := time.Now().UnixMilli()
	publishArt.Ctime = now
	publishArt.Utime = now
	// 确保代码逻辑分支生成的返回值要传到需要的实体上
	art.Id = id
	err = tx.Clauses(clause.OnConflict{
		// 对mysql不起效，但是可以兼容别的方言
		Columns: []clause.Column{{Name: "id"}},
		//做更新，mysql就支持这个
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   publishArt.Title,
			"content": publishArt.Content,
			"utime":   now,
			//tbdxx
			"status": publishArt.Status,
		}),
	}).Create(&publishArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, err
}

func (a *ArticleGROMDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	res := a.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
			"status":  art.Status,
		})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errors.New("更新失败,ID不对或者作者不对")
	}

	return nil
}

func (a *ArticleGROMDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	//log.Printf("id: %v,title: %v,status:%v,auhtorid: %v,ctime:%v,utime:%v /n",art.Id,art.Title,art.Status,art.AuthorId,art.Ctime,art.Utime)
	err := a.db.Table("articles").WithContext(ctx).Create(&art).Error
	return art.Id, err
}
