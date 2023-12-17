package sms

import (
	"context"
)

//go:generate mockgen -source=internal/service/sms/sms_service.go -destination=internal/service/mocks/sms_service.mock.go -package=svcmock
type Service interface {
	Send(ctx context.Context, tplId string, args []string, phoneNums []string) error
}
