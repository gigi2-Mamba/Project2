package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"project0/internal/domain"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Save(ctx context.Context, art domain.Article) (int64,error)
	Update(ctx context.Context, article domain.Article) error

	//Repository 层 同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context,id int64, uid int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article,error)
}

type CacheArticleRepository struct {
	dao dao.ArticleDao


    // 分离制作库和线上库 V2
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
	// 强耦合，没法摆脱依赖。 严格来说repository抽象存储层不用考虑具体的db实现
	db  *gorm.DB
}

func (c *CacheArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {


	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return  nil,err
	}
	res := slice.Map[dao.Article,domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return  c.toDomain(src)
	})
	//也可以异步处理缓存回写
    //go func(){
	//
	//}

   return  res,nil

}

func NewCacheArticleRepository(dao dao.ArticleDao,cache cache.ArticleCache) ArticleRepository {

	return &CacheArticleRepository{
		dao: dao,
	    }
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context,id int64, uid int64, status domain.ArticleStatus) error{
    // 怎么这个方法不直接用article嫌碍事吗  因为只是修改一个状态
    return c.dao.SyncStatus(ctx,id,uid,status.ToUint8())
}

// 在Resposioty 来分离制作库和线上库 叫做v2 对于service层处理
func NewCacheArticleRepositoryV2(authorDAO dao.ArticleAuthorDAO,readerDAO dao.ArticleReaderDAO,) *CacheArticleRepository {
	return &CacheArticleRepository{readerDAO: readerDAO, authorDAO: authorDAO}
}
// 在dao层处理制作库和线上库 解决方案v3
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
    return  c.dao.Sync(ctx,c.toEntity(art))
}
// 同库不同表用事务实现,在repository层面要加db字段    reposioty的v2
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	dbtx := c.db.WithContext(ctx).Begin()
	if dbtx.Error != nil {
		return 0,dbtx.Error
	}
    defer c.db.Rollback()

   authorDAO := dao.NewArticleGORMAuthorDAO(dbtx)
   readerDAO := dao.NewArticleReaderGORMDAO(dbtx)

	artEntity := c.toEntity(art)
	// 这么使用时有一个技巧的
	var ( id = art.Id
		err error)
	if id > 0 {
		err = authorDAO.Update(ctx, artEntity)
	} else {
		id, err = authorDAO.Create(ctx, artEntity)
	}
	if err != nil {
		return 0, err
	}
	artEntity.Id = id
	err = readerDAO.Upsert(ctx, artEntity)
	if err != nil {
		return 0, err
	}
	dbtx.Commit()
	return id, nil
}
// 两个库就可以用这个实现,而且是非事务处理
func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artEntity := c.toEntity(art)
	// 这么使用时有一个技巧的，为了给if这种局部区域用
	var ( id = art.Id
		err error)
	if id > 0 {
		err = c.authorDAO.Update(ctx, artEntity)
	} else {
		id, err = c.authorDAO.Create(ctx, artEntity)
	}
	if err != nil {
		return 0, err
	}
	// 是为了上面的else分支做补全
	artEntity.Id = id
	err = c.readerDAO.Upsert(ctx, artEntity)
	return id, err
}



func (c *CacheArticleRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx,c.toEntity(art))
}
func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) (error) {
	return c.dao.UpdateById(ctx,c.toEntity(art))
}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id: art.Id,
		Title:  art.Title,
		Content: art.Content,
		AuthorId: art.Author.Id,
		Status: art.Status.ToUint8(),
		//Status: int(art.Status),
		//Fucker: art.Fucker,
	}
}

func (c *CacheArticleRepository) toDomain(art dao.Article) domain.Article {
	 return domain.Article{
		 Id: art.Id,
		 Title: art.Title,
		 Content: art.Content,
		 Status: domain.ArticleStatus(art.Status),
		 Author: domain.Author{
			 Id: art.Id,
		 },
		 Ctime: time.UnixMilli(art.Ctime),
		 Utime: time.UnixMilli(art.Utime),
	 }
}

