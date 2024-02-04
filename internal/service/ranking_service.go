package service

import (
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

func NewBatchRankingService(batchSize int, artSvc ArticleService, usrActSvc UserActivityService, repo repository.RankingRepository) *BatchRankingService {
	return &BatchRankingService{batchSize: batchSize, artSvc: artSvc, usrActSvc: usrActSvc, repo: repo}
}

// 装饰器实现log
// 封装为Job
// 用cron job定时调用计算
// TDD

func (b *BatchRankingService) TopN(ctx context.Context) error {
	// 计算接口
	// 用队列计算前N个阅读数最多的
	// 新建一个队列大小为N，把阅读数装进去
	// 每次从数据库里拉一批进来，不要全拉进来，内存会炸
	// 用定时任务定时计算，达到更新目的
}

func (b *BatchRankingService) GetTopN(ctx context.Context) ([]domain.Article, error) {
	// 查询接口
	// 本地缓存 + redis缓存
	// 本地缓存兜底，redis缓存是即时数据
	// 因为前N个这个业务并不是强实时性
	// 如果MySQL和Redis两个都宕机，要给前端告知，让它重试
	// failover策略，高可用常见的套路
}
