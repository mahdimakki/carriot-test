package ports

import "context"

type WarningQueue interface {
	PublishAllLogs(ctx context.Context, value string) error
	PublishWarningLogs(ctx context.Context, value string) error
	Subscribe(ctx context.Context)
	SetHandler(func(context.Context, []byte) error)
}
