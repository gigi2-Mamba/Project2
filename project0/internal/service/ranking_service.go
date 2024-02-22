package service

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"project0/internal/repository"

	//"github.com/ecodeclub/ekit/queue"

	//queue	"github.com/ecodeclub/ekit/internal/queue"

	"project0/internal/domain"
	"time"
)

/*
Created by society-programmer on 2024/2/20.
*/

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
    interSvc InteractiveService
	artSvc  ArticleService
	batchSize int   // 批次大小，一批取多大也就是limit
	scoreFunc func(likeCnt int64,utime time.Time) float64
	n int   //热榜，前n个，控制传多少个

	// 加入缓存
	repo   repository.RankingRepository
}

func NewBatchRankingService(inter InteractiveService, article ArticleService) RankingService {
	return &BatchRankingService{interSvc: inter, 
		artSvc: article,
	    batchSize: 100,
		n: 100,
		//在New的时候直接初始化func的手段,应该是在new直接写死了。 内部锁死。
	    scoreFunc: func(likeCnt int64, utime time.Time) float64 {
           duration := time.Since(utime).Seconds()
		   return   float64(likeCnt - 1) / math.Pow(duration+2,1.5)
		},
	}
}

func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return b.repo.GetTopN(ctx)
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	arts,err := b.topN(ctx)
	//_,err := b.topN(ctx)
	if err != nil {
		return err
	}
	// 最终要存到缓存里面的
	return b.repo.ReplaceTopN(ctx,arts)
}

// 起初这么写，是绕开了缓存，更加好测。  单元测试先测这里
// 获取数据库前n条的数据及点赞总数
func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
     //在这里找出热榜数据
	//获取前n条数据
	offset := 0
	start:= time.Now()
	ddl := start.Add(-7 * 24 * time.Hour)

	type Score struct {
		score float64
		art domain.Article
	}
	// 得到一个优先队列.  得到优先队列吗，没什么好说的用ekit，后期也可以自己写自己的ekit
    // 为什么这个叫做topn?
	topN :=  queue.NewPriorityQueue(b.n, func(src Score, dst Score) int {
           if   src .score > dst.score {
			   return 1
		   } else if src.score == dst.score {
			   return 0
		   } else {
			   return -1
		   }
	})
       //为什么这里要用for?，有多批，以下代码是一次处理流程
	for  {
		//取数据，取批量的article
		arts ,err := b.artSvc.ListPub(ctx,start,offset,b.batchSize)
		if err != nil {
			return nil, err
		}
		//if len(arts) == 0 {
		//	break
		//}
        //ekit的slice.Map将一个切片转换成另外一个切片。泛型。计算原切片的长度，然后提供元素转换func，然后for range逐个转换
		//利用到的只是。append.
		ids := slice.Map(arts, func(idx int,art domain.Article) int64 {
			return art.Id
		})

		intrMap, err := b.interSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		for _,art := range arts{
			intr := intrMap[art.Id]
			score := b.scoreFunc(intr.LikeCnt,art.Utime)
			ele := Score{
				score: score,
				art: art,
			}
			//入队有检测是否淘汰替换
			err = topN.Enqueue(ele)
			if err ==  queue.ErrOutOfCapacity  {
				//拿出最小元素
				minEle,_ := topN.Dequeue()
				if minEle.score < score {
					topN.Enqueue(ele)
				} else {
					//topN.Enqueue(minEle)
					topN.Enqueue(ele)
				}
			}
		}
		// 这里决定了批量
		offset = offset + len(arts)
		//没有取够一批，我们就直接中断执行
		if len(arts) < b.batchSize ||
			arts[len(arts)-1].Utime.Before(ddl){
			break
		}
	}

	res := make([]domain.Article,topN.Len())
	// 因为小顶堆的特性，反着来。
	for i := topN.Len()-1; i >= 0 ; i-- {
		ele,_ := topN.Dequeue()
		res[i] = ele.art
	}
  return res,nil
}




