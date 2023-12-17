package async

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/domain"
	"github.com/jasonzhao47/cuddle/internal/repository"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"time"
)

type SmsService struct {
	repo repository.AsyncSmsRepository
}

func NewService(repo repository.AsyncSmsRepository) sms.Service {
	svc := &SmsService{repo: repo}
	go func() {
		// start Async
		svc.AsyncSend()
	}()
	return svc
}

func (s *SmsService) AsyncSend() {
	// 抢占调度
	// preemptive
	// 抢到了直接cancel
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptiveGetSms(ctx)
	cancel()
	switch err {
	case nil:
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := s.Send(ctx, as.TplId, as.Args, as.PhoneNums)
		if er != nil {
			// 记录发送错误
		}
		res := er == nil
		er = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if er != nil {
			// 发送短信成功
			// 记录数据库存储错误
		}
	case repository.ErrSmsNotFound:
		time.Sleep(time.Second)
	default:
		// 规避网络抖动问题
		// 记录一下抢占失败了
		time.Sleep(time.Second)
	}
}

func (s *SmsService) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	if s.needAsync() {
		// 转储到本地数据库
		return s.repo.Add(ctx, domain.AsyncSms{
			Id:        1,
			TplId:     tplId,
			RetryMax:  0,
			PhoneNums: phoneNums,
			Args:      args,
		})
	}
	return s.Send(ctx, tplId, args, phoneNums)
}

func (s *SmsService) needAsync() bool {
	return true
}
