package failover

import (
	"context"
	"github.com/jasonzhao47/cuddle/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailOverSmsService struct {
	svcs      []sms.Service
	idx       uint64
	threshold uint64
	cnt       int32
}

func NewTimeoutFailOverSmsService(svcs []sms.Service) *TimeoutFailOverSmsService {
	return &TimeoutFailOverSmsService{svcs: svcs}
}

func (t *TimeoutFailOverSmsService) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	// 把当前的index调出来
	idx := atomic.LoadUint64(&t.idx)
	// 把当前的超时次数调出来
	cnt := atomic.LoadInt32(&t.cnt)
	// 超过了阈值，需要更换一个服务
	if uint64(cnt) > t.threshold {
		newIdx := idx % uint64(len(t.svcs))
		// 如果存储成功
		if atomic.CompareAndSwapUint64(&t.idx, idx, newIdx) {
			// 重新开始计算
			atomic.StoreUint64(&t.idx, 0)
		}
	}
	// 真正去用一个服务商发送消息
	err := t.svcs[t.idx].Send(ctx, tplId, args, phoneNums)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	default:
	}
	return nil
}
