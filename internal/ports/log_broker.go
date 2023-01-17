package ports

import (
	"context"
)

type LogBroker interface {
	Subscribe(context.Context)
	SetHandler(func(context.Context, []byte) error)
}
