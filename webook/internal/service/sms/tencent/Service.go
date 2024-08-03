package tencent

import (
	"context"
	"fmt"
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
	"github.com/lyydsheep/generic_tools/Slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   sms.Client
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms2.NameData, numbers ...string) error {
	request := sms.NewSendSmsRequest()
	//配置信息
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = &tpl
	request.SetContext(ctx)
	request.PhoneNumberSet = stringToPtrSlice(numbers)
	strings := make([]string, len(args))
	for i := range strings {
		strings[i] = args[i].Data
	}
	request.TemplateParamSet = stringToPtrSlice(strings)
	resp, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	for _, v := range resp.Response.SendStatusSet {
		if v.Code == nil || *(v.Code) != "Ok" {
			return fmt.Errorf("发送失败，code：%s，reason：%s", *v.Code, *v.Message)
		}
	}
	return nil
}

func stringToPtrSlice(s []string) []*string {
	return Slice.Map[string, *string](s, func(s string) *string {
		return &s
	})
}

func NewService(appId string, signName string, client sms.Client) *Service {
	return &Service{
		appId:    &appId,
		signName: &signName,
		client:   client,
	}
}
