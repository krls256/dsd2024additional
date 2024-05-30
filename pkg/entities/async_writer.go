package entities

import "github.com/samber/mo"

type AsyncWriter[S any] interface {
	AsyncWrite(result mo.Either[S, error])
}

func NewLeftWithAction(action string, payload any) LeftWithAction {
	return LeftWithAction{
		Action:  action,
		Payload: payload,
	}
}

type LeftWithAction struct {
	Action  string
	Payload any
}
