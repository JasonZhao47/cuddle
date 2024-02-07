package service

import (
	"container/heap"
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"time"
)

type RankingService interface {
	TopN(ctx context.Context) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	n         int
	batchSize int
	artSvc    ArticleService
	usrActSvc UserActivityService
	scoreFn   func(readCnt int64, utime time.Time) float64
	repo      repository.RankingRepository
}

func NewBatchRankingService(n int, batchSize int, artSvc ArticleService, usrActSvc UserActivityService, scoreFn func(readCnt int64, utime time.Time) float64, repo repository.RankingRepository) *BatchRankingService {
	return &BatchRankingService{n: n, batchSize: batchSize, artSvc: artSvc, usrActSvc: usrActSvc, scoreFn: scoreFn, repo: repo}
}

// 装饰器实现log
// 封装为Job
// 用cron job定时调用计算
// TDD

func (b *BatchRankingService) TopN(ctx context.Context) ([]domain.PublishedArticle, error) {
	// 计算接口
	// 用队列计算前N个阅读数最多的
	// 新建一个队列大小为N，把阅读数装进去
	// 每次从数据库里拉一批进来，不要全拉进来，内存会炸
	// 用定时任务定时计算，达到更新目的

	offset := 0
	now := time.Now()
	pq := make(PriorityQueue, 0)
	ddl := now.Add(-7 * 24 * time.Hour)
	heap.Init(&pq)

	for {
		pubArts, err := b.artSvc.ListPub(ctx, now, offset, b.batchSize)
		if err != nil {
			return []domain.PublishedArticle{}, err
		}

		ids := make([]int64, len(pubArts))
		for i := 0; i < len(pubArts); i++ {
			ids[i] = pubArts[i].Id
		}

		userActs, er := b.usrActSvc.GetReadByIds(ctx, "article", ids)
		if er != nil {
			return []domain.PublishedArticle{}, err
		}
		for i := 0; i < len(pubArts); i++ {
			currScore := b.scoreFn(userActs[i].ReadCnt, pubArts[i].Utime)
			if pq.Len() == b.n {
				minEle := heap.Pop(&pq)
				score := minEle.(*Score)
				if score.score < currScore {
					heap.Push(&pq, &Score{
						score: currScore,
						art:   pubArts[i],
					})
				} else {
					heap.Push(&pq, &Score{
						score: score.score,
						art:   score.art,
					})
				}
			} else {
				heap.Push(&pq, &Score{
					score: currScore,
					art:   pubArts[i],
				})
			}
		}
		offset = offset + len(pubArts)
		if len(pubArts) == 0 || len(pubArts) < b.batchSize || pubArts[len(pubArts)-1].Utime.Before(ddl) {
			break
		}
	}
	res := make([]domain.PublishedArticle, pq.Len())

	for i := 0; i < len(res); i++ {
		ele := pq.Pop()
		scr := ele.(*Score)
		res[i] = scr.art
	}
	return res, nil
}

func (b *BatchRankingService) RankTopN(ctx context.Context) error {
	// 查询接口
	// 本地缓存 + redis缓存
	// 本地缓存兜底，redis缓存是即时数据
	// 因为前N个这个业务并不是强实时性
	// 如果MySQL和Redis两个都宕机，要给前端告知，让它重试
	// failover策略，高可用常见的套路
	arts, err := b.TopN(ctx)
	if err != nil {
		return err
	}
	return b.repo.ReplaceTopN(ctx, arts)
}

type Score struct {
	score float64
	art   domain.PublishedArticle
}

type PriorityQueue []*Score

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].score < pq[j].score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	score := x.(*Score)
	*pq = append(*pq, score)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
