package job

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/service"
	"time"
)

type RankingJob struct {
	rankingSvc service.RankingService
}

func NewRankingJob(rankingSvc service.RankingService) *RankingJob {
	return &RankingJob{rankingSvc: rankingSvc}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return r.rankingSvc.TopN(ctx)
}
