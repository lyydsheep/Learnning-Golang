package sms

import "context"

type SMS interface {
	send(ctx context.Context, tpl string, args []string, numbers ...string) error
}
