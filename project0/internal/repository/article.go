package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"log"
	"project0/internal/domain"
	"project0/internal/repository/cache"
	"project0/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	//Repository 层 同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type CacheArticleRepository struct {
	dao   dao.ArticleDao
	cache cache.ArticleCache
	//缺少领域服务信息，注入repo
	userDao dao.UserDao
	// 分离制作库和线上库 V2
	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
	// 强耦合，没法摆脱依赖。 严格来说repository抽象存储层不用考虑具体的db实现
	db *gorm.DB
}

func (c *CacheArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.cache.GetPub(ctx, id)
	if err == nil {
		//log.Println("获取到了缓存？？")
		return res, nil
	} else {
		//日志记录
	}
    log.Println("没有该发布帖子的缓存")
	pubArt, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err

	}
	res = c.toDomain(dao.Article(pubArt))
	author, err := c.userDao.Profile(ctx, res.Author.Id)
	if err != nil {
		return res, err
	}
	res.Author.Name = author.NickName
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := c.cache.SetPub(ctx, res)
		if er != nil {
			//记录日志
		}
	}()
	return res, nil
}

func (c *CacheArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	// 先走缓存,新奇缓存方案，通过相关业务预加载生成的缓存直接获取
	res, err := c.cache.Get(ctx, id)
	if err == nil {
		return res, err
	}

	art, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res = c.toDomain(art)
	// 异步回写缓存
	go func() {
		er := c.cache.Set(ctx, res)
		if er != nil {
			//记录日志
		}
	}()
	return res, nil
}

func (c *CacheArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 引入缓存,是否访问缓存
	// 不严格也可以< 100取出，但是注意处理取出的多余数据
	if offset == 0 && limit == 100 { //先写上层，再写下层，先不关心底层实现
		res, err := c.cache.GetFirstPage(ctx, uid)
		// 对err = nil 做优先处理，不同的err不同处理
		if err == nil {
			return res, err
		} else {
			//没有出错，但是可能未命中缓存
			// 做一些日志处理,那应该做什么呢 tbd
		}
	}

	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	//也可以异步处理缓存回写
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if offset == 0 && limit == 100 {
			err = c.cache.SetFirstPage(ctx, uid, res)
			if err != nil {
				// 记录日志
			}
		}
	}()
	// 相关业务缓存预加载

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		c.preCache(ctx, res)
	}()

	return res, nil
}

func NewCacheArticleRepository(dao dao.ArticleDao, cache cache.ArticleCache,user dao.UserDao) ArticleRepository {
	return &CacheArticleRepository{
		dao:   dao,
		cache: cache,
		userDao: user,
	}
}

// 同步status不要搞缓存吗
func (c *CacheArticleRepository) SyncStatus(ctx context.Context, id int64, uid int64, status domain.ArticleStatus) error {
	err := c.dao.SyncStatus(ctx, id, uid, status.ToUint8())
	if err == nil {
		er := c.cache.DeleteFirstPage(ctx, uid)
		if er != nil {
			//记录日志
		}
	}
	return err
}

// 在Resposioty 来分离制作库和线上库 叫做v2 对于service层处理
func NewCacheArticleRepositoryV2(authorDAO dao.ArticleAuthorDAO, readerDAO dao.ArticleReaderDAO) *CacheArticleRepository {
	return &CacheArticleRepository{readerDAO: readerDAO, authorDAO: authorDAO}
}

// 在dao层处理制作库和线上库 解决方案v3
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	// 操作err == nil 一般都是当前方法是返回错误的技巧
	if err == nil {
		er := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if er != nil {
			//记录日志
		}
	}
	// 预设置缓存，刚发布有人会看
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(),time.Second)
		defer cancel()
		log.Println("debug author_id ",art.Author.Id)
		profile ,er := c.userDao.Profile(ctx,art.Author.Id)
		if er !=nil {
			//记录日志
		}
		// 灵活设置过期时间
		art.Author = domain.Author{
			Id: profile.Id,
			Name:  profile.NickName,
		}
		er = c.cache.SetPub(ctx,art)
		if er != nil {
			//记录日志
		}
	}()

	return id, err
}

// 同库不同表用事务实现,在repository层面要加db字段    reposioty的v2
func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	dbtx := c.db.WithContext(ctx).Begin()
	if dbtx.Error != nil {
		return 0, dbtx.Error
	}
	defer c.db.Rollback()

	authorDAO := dao.NewArticleGORMAuthorDAO(dbtx)
	readerDAO := dao.NewArticleReaderGORMDAO(dbtx)

	artEntity := c.toEntity(art)
	// 这么使用时有一个技巧的
	var (
		id  = art.Id
		err error
	)
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
	var (
		id  = art.Id
		err error
	)
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

// 上层是service.Edit
func (c *CacheArticleRepository) Save(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Insert(ctx, c.toEntity(art))
	if err == nil {
		er := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if er != nil {
			//记录日志
		}
	}
	return id, err
}
func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	err := c.dao.UpdateById(ctx, c.toEntity(art))
	if err == nil {
		er := c.cache.DeleteFirstPage(ctx, art.Author.Id)
		if er != nil {
			//记录日志
		}
	}
	return err

}

func (c *CacheArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   art.Status.ToUint8(),
	}
}

func (c *CacheArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id: art.Id,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

// 相关业务缓存预加载
func (c *CacheArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 大文件都不缓存
	const contentSizeThreshold = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) < contentSizeThreshold {
		err := c.cache.Set(ctx, arts[0])
		if err != nil {
			//记录日志
			log.Println("相关业务缓存预设置失败: ", err)
		}
	}

}
