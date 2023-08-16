package telegram

import "context"

type User struct {
	ctx    context.Context
	cancel context.CancelFunc
}
