package service

import (
	"context"
	"fmt"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository"
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
	"math/rand"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, input string) error
	generateCode(biz string) string
	getKey(biz, phone string) string
}

type BasicCodeService struct {
	repo repository.CodeRepository
	sms  sms2.SMS
}

var (
	ErrTooFrequent = repository.ErrTooFrequent
	ErrNotMatch    = repository.ErrNotMatch
	ErrExceed      = repository.ErrExceed
)

func NewCodeService(cr repository.CodeRepository, sms sms2.SMS) CodeService {
	return &BasicCodeService{repo: cr, sms: sms}
}

// Send 发送验证码，业务类型，目标手机号
func (cs *BasicCodeService) Send(ctx context.Context, biz, phone string) error {
	//生成一个验证码
	code := cs.generateCode(biz)
	//将验证码放入Redis
	key := cs.getKey(biz, phone)
	err := cs.repo.Store(ctx, key, code)
	if err != nil {
		return err
	}
	//发送验证码
	err = cs.sms.Send(ctx, "这是一个tpl", []sms2.NameData{
		{Name: "code", Data: code},
	}, phone)
	return err
}

// Verify 校验验证码
func (cs *BasicCodeService) Verify(ctx context.Context, biz, phone, input string) error {
	key := cs.getKey(biz, phone)
	err := cs.repo.Check(ctx, key, input)
	return err
}

func (cs *BasicCodeService) generateCode(biz string) string {
	switch biz {
	case "login":
		return fmt.Sprintf("%06d", rand.Intn(1000000))
	default:
		return fmt.Sprintf("%06d", rand.Intn(1000000))
	}
}

func (cs *BasicCodeService) getKey(biz, phone string) string {
	return fmt.Sprintf("code:%s:%s", biz, phone)
}
