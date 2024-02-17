package service

import (
	"context"
	"errors"
	"log"
	"project0/internal/domain"
	"project0/internal/events/article"
	"project0/internal/repository"
	"project0/pkg/loggerDefine"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid int64, id int64) error
	GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article,error)
	GetById(ctx context.Context, id int64) (domain.Article,error)
	GetPubById(ctx context.Context, id,uid int64) (domain.Article,error)

}

type articleService struct {
	repo repository.ArticleRepository
	producer article.Producer
	// service层分发制作库和线上库  v1写法
	readerRepo repository.ArticleReaderRepository
	authorRepo repository.ArticleAuthorRepository
	l          loggerDefine.LoggerV1
}

func NewArticleService(repo repository.ArticleRepository,producer article.Producer) ArticleService {
	return &articleService{
		repo: repo,
		producer: producer}
}
func (a *articleService) GetPubById(ctx context.Context, id,uid int64) (domain.Article, error) {
	// 如果是微服务版本，可以直接调用其他服务来补全前端缺失的领域信息
	//log.Println("(a *articleService) GetPubById ",id)
    // 这里只是获取整个帖子的详情
	art, err := a.repo.GetPubById(ctx, id)
	// 在这里决定，发送一条消息给kafka.  这么一看没什么用？
	if err == nil {
        //log.Println("发送消息给kafka ")
		go func() {
			//log.Println("接下来发送信息")
			er := a.producer.ProduceReadEvent(article.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er != nil {
				a.l.Error("发送readEvent失败",
					loggerDefine.Int64("aid",id),
					loggerDefine.Int64("uid",uid))
			}
		}()
	}
	return art,nil
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {

	return  a.repo.GetById(ctx, id)
}

func (a *articleService) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {

	return  a.repo.GetByAuthor(ctx, uid, offset, limit)
}



func (a *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return a.repo.SyncStatus(ctx, id, uid, domain.ArticleStatusPrivate)
}

func NewArticleServiceV1(authorRepo repository.ArticleAuthorRepository, readerRepo repository.ArticleReaderRepository,
	l loggerDefine.LoggerV1) *articleService {
	return &articleService{readerRepo: readerRepo, authorRepo: authorRepo, l: l}
}
// 发表状态是属于业务逻辑，所以在服务层就改变了状态
func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}
// service层面分离制作库和线上库 v1   各个层面的实现分离制作库和线上库 service层面 是第一个版本
// 靠重试
func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	//return 0,nil
	// 在tdd测试时，直接使用一个panic在这里，逻辑上静态语法检查就会miss
	// tdd模式下写完的方法需要经常改
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.authorRepo.Update(ctx, art)
	} else {
		id, err = a.authorRepo.Create(ctx, art)
	}

	if err != nil {
		log.Println("制作出问题")
		return 0, err
	}
	// 制作库和线上库同用一个id保持一致。 延伸神道： 近似但是不同用处的功能共用同样的东西，可以让另一个方法省个返回值。
	//节约空间
	art.Id = id
	//--实践版本
	//线上库可能已经有了
	//可能也没有
	//err = a.readerRepo.Save(ctx, art)
	// 制作库保存成功，线上库保存失败,实践中这样就可以。 上线观察
	//if err != nil {
	//   a.l.Error("保存到制作库成功但是线上库失败",loggerDefine.Int64("aid",art.Id),
	//	   loggerDefine.Error(err))
	//   return 0, err
	//}
	// --实践版本
	// 引入重试
	for i := 0; i < 3; i++ {
		err = a.readerRepo.Save(ctx, art)
		if err != nil {
			// if err != nil 做日志记录，进行下一步逻辑
			a.l.Error("保存到制作库成功但是线上库失败", loggerDefine.Int64("aid", art.Id))
			loggerDefine.Error(err)
		} else {
			return id, nil
		}
	}
	a.l.Error("保存到制作库成功但是线上库失败,重试耗尽", loggerDefine.Int64("aid", art.Id))
	loggerDefine.Error(err)

	return id, errors.New("保存到制作库成功但是线上库失败,重试耗尽")
}
// 共用同一个接口实现修改和新建
func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnPublished
	//article.Fucker = 222
	//log.Println("理应设置成功了")
	//log.Println("理应设置成功了",article.Status)
	// 不复杂就直接分发id
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}
	return a.repo.Save(ctx, article)
}
