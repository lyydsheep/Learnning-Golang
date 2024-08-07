package ioc

import (
	sms2 "github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/service/sms/memory"
)

func InitSMS() sms2.SMS {
	return memory.NewService()
}
