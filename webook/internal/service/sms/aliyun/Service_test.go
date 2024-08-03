package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ali "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
	"github.com/lyydsheep/generic_tools"
	"testing"
)

func TestService_Send(t *testing.T) {
	config := &openapi.Config{
		AccessKeyId:     generic_tools.ToPtr[string](""),
		AccessKeySecret: generic_tools.ToPtr[string](""),
		Endpoint:        generic_tools.ToPtr[string]("dysmsapi.aliyuncs.com"),
	}
	client := &ali.Client{}
	client, err := ali.NewClient(config)
	if err != nil {
		panic(err)
	}
	s := NewService("webook后端短信接口", client)
	err = s.Send(nil, "SMS_471400091", []sms2.NameData{
		{Name: "code", Data: "4321"},
	}, "17870189780")
	if err != nil {
		panic(err)
	}
}
