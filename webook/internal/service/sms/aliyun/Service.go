package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	ali "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
	"github.com/lyydsheep/generic_tools"
)

type Service struct {
	client   *ali.Client
	signName string
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms2.NameData, numbers ...string) error {
	m := make(map[string]any, len(args))
	for i := range args {
		m[args[i].Name] = args[i].Data
	}
	bCode, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req := &ali.SendSmsRequest{
		SignName:      generic_tools.ToPtr[string](s.signName),
		TemplateCode:  generic_tools.ToPtr[string](tpl),
		TemplateParam: generic_tools.ToPtr[string](string(bCode)),
	}
	for _, v := range numbers {
		req.PhoneNumbers = generic_tools.ToPtr[string](v)
		resp, err := s.client.SendSms(req)
		if err != nil {
			return err
		}
		if *resp.Body.Code == "OK" {
			fmt.Println("succeed")
		}
	}
	return nil
}

func NewService(signName string, client *ali.Client) *Service {
	return &Service{
		client:   client,
		signName: signName,
	}
}
