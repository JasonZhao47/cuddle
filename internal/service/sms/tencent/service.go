package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	sms2 "github.com/jasonzhao47/cuddle/internal/service/sms"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type SmsService struct {
	client   *sms.Client
	appId    *string
	signName *string
}

func NewSMSService(client *sms.Client, appId string, signName string) sms2.Service {
	return &SmsService{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *SmsService) Send(ctx context.Context, tplId string, args []string, phoneNums []string) error {
	request := sms.NewSendSmsRequest()

	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = common.StringPtr(tplId)
	request.TemplateParamSet = common.StringPtrs(args)
	request.PhoneNumberSet = common.StringPtrs(phoneNums)

	response, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	for _, status := range response.Response.SendStatusSet {
		if status.Code != nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败, code: %s, 原因: %s", *status.Code, *status.Message)
		}
	}
	// log here
	return nil
}
