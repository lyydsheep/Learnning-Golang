package sms

import "context"

type SMS interface {
	Send(ctx context.Context, tpl string, args []NameData, numbers ...string) error
}

type NameData struct {
	Name string
	Data string
}
