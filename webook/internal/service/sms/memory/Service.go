package memory

import (
	"context"
	"fmt"
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tpl string, args []sms2.NameData, numbers ...string) error {
	fmt.Printf("验证码 %s 发送成功", args[0].Data)
	return nil
}
