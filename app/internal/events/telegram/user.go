package telegram

import "context"

type User struct {
	cancelGateMode context.CancelFunc
}

func NewUser(cancelGateMode context.CancelFunc) User {
	return User{cancelGateMode: cancelGateMode}
}
