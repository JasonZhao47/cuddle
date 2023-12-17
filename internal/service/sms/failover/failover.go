package failover

import (
	"context"
	"errors"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"sync/atomic"
)

type FailOverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{svcs: svcs}
}

//func (f *FailOverSMSService) SendProto(ctx context.Context, tplId string, args []string, phoneNums []string) error {
//	for _, svc := range f.svcs {
//		err := svc.Send(ctx, tplId, args, phoneNums)
//		if err == nil {
//			return nil
//		}
//		// log.Println(err)
//	}
//	return errors.New("所有第三方服务商都失败了")
//}

func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, phoneNums)
		if err == nil {
			return nil
		}
		if err == context.Canceled || err == context.DeadlineExceeded {
			return err
		}
		// log.Println(err)
	}
	return errors.New("所有第三方服务商都失败了")
}
